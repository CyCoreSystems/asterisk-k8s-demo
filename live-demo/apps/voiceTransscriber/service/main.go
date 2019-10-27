package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/CyCoreSystems/audiosocket"
	"github.com/fatih/color"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	speechv1 "google.golang.org/genproto/googleapis/cloud/speech/v1"
	texttospeechv1 "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

// MaxCallDuration is the maximum amount of time to allow a call to be up before it is terminated.
const MaxCallDuration = 2 * time.Minute

// MaxRecognitionDuration is the maximum amount of time to allow for a single voice recognition session to complete
const MaxRecognitionDuration = time.Minute

const listenAddr = ":8080"
const languageCode = "en-US"

const greetingMessage = "Hello. Speak, and I will transscribe you."
const partingMessage = "Good bye. Thanks for calling."

// slinChunkSize is the number of bytes which should be sent per Slin
// audiosocket message.  Larger data will be chunked into this size for
// transmission of the AudioSocket.
//
// This is based on 8kHz, 20ms, 16-bit signed linear.
const slinChunkSize = 320 // 8000Hz * 20ms * 2 bytes

// ErrHangup indicates that the call should be terminated or has been terminated
var ErrHangup = errors.New("Hangup")

var recog *speech.Client
var tts *texttospeech.Client
var googleCreds = "/var/secrets/google/google.json"

func main() {
	var err error

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", googleCreds)
	}

	ctx := context.Background()

	if recog, err = speech.NewClient(ctx); err != nil {
		log.Fatalln("failed to connect to Google Speech API service:", err)
	}
	defer recog.Close()
	if tts, err = texttospeech.NewClient(ctx); err != nil {
		log.Fatalln("failed to connect to Google Text-to-Speech API service:", err)
	}
	defer tts.Close()

	if err = Listen(ctx); err != nil {
		log.Fatalln("listen failure:", err)
	}
	log.Println("exiting")
}

// Listen listens for and responds to Audiosocket connections
func Listen(ctx context.Context) error {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to bind listener to socket %s", listenAddr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("failed to accept new connection:", err)
			continue
		}

		go Handle(ctx, conn)
	}
}

// Handle processes a call
func Handle(pCtx context.Context, c net.Conn) {
	var err error
	var id uuid.UUID

	ctx, cancel := context.WithTimeout(pCtx, MaxCallDuration)

	defer func() {
		if err := speak(ctx, c, partingMessage); err != nil {
			log.Println("failed to speak closing:", err)
		}

		cancel()

		color.Magenta("ending call %s", id.String())
		if _, err := c.Write(audiosocket.HangupMessage()); err != nil {
			log.Println("failed to send hangup message:", err)
		}
	}()

	id, err = getCallID(c)
	if err != nil {
		log.Println("failed to get call ID:", err)
		return
	}
	color.Magenta("processing call %s", id.String())

	if err = speak(ctx, c, greetingMessage); err != nil {
		log.Println("failed to send greeting to Asterisk:", err)
	}

	for ctx.Err() == nil {
		resp, err := processCommand(ctx, c)
		if err == ErrHangup {
			return
		}
		if err != nil {
			log.Println("failed to process command:", err)
		}
		if resp != "" {
			if err := speak(ctx, c, resp); err != nil {
				log.Println("failed to speak response:", err)
			}
		}
		if err != nil {
			return
		}
	}
}

func getCallID(c net.Conn) (uuid.UUID, error) {
	m, err := audiosocket.NextMessage(c)
	if err != nil {
		return uuid.Nil, err
	}

	if m.Kind() != audiosocket.KindID {
		return uuid.Nil, errors.Errorf("invalid message type %d getting CallID", m.Kind())
	}

	return uuid.FromBytes(m.Payload())
}

func ttsRequest(msg string) *texttospeechv1.SynthesizeSpeechRequest {
	return &texttospeechv1.SynthesizeSpeechRequest{
		Input: &texttospeechv1.SynthesisInput{
			InputSource: &texttospeechv1.SynthesisInput_Text{
				Text: msg,
			},
		},
		Voice: &texttospeechv1.VoiceSelectionParams{
			LanguageCode: languageCode,
		},
		AudioConfig: &texttospeechv1.AudioConfig{
			AudioEncoding:   texttospeechv1.AudioEncoding_LINEAR16,
			SampleRateHertz: 8000,
		},
	}
}

func recognizeRequest(pCtx context.Context, r io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(pCtx, MaxRecognitionDuration)
	defer cancel()

	svc, err := recog.StreamingRecognize(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to start streaming recognition")
	}

	if err := svc.Send(&speechv1.StreamingRecognizeRequest{
		StreamingRequest: &speechv1.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechv1.StreamingRecognitionConfig{
				Config: &speechv1.RecognitionConfig{
					Encoding:        speechv1.RecognitionConfig_LINEAR16,
					SampleRateHertz: 8000,
					LanguageCode:    languageCode,
					Model:           "command_and_search",
					UseEnhanced:     true,
					SpeechContexts: []*speechv1.SpeechContext{
						&speechv1.SpeechContext{
							Phrases: []string{
								"bye",
								"goodbye",
								"hangup",
							},
						},
					},
				},
			},
		},
	}); err != nil {
		return "", errors.Wrap(err, "failed to send recognition config")
	}

	go pipeFromAsterisk(ctx, r, svc)

	resp, err := svc.Recv()
	if err == io.EOF {
		return "", nil
	}
	if err != nil {
		return "", errors.Wrap(err, "recognition failed")
	}
	if err := resp.Error; err != nil {
		if err.Code == 3 || err.Code == 11 {
			log.Println("recognition exceeded 60-second limit")
		}
		return "", errors.New(err.String())
	}

	for _, result := range resp.Results {
		for _, alt := range result.GetAlternatives() {
			if alt.Transcript != "" {
				return alt.Transcript, nil
			}
		}
	}
	return "", nil
}

func pipeFromAsterisk(ctx context.Context, in io.Reader, out speechv1.Speech_StreamingRecognizeClient) {
	var err error
	var m audiosocket.Message

	defer out.CloseSend() // nolint: errcheck

	for ctx.Err() == nil {
		m, err = audiosocket.NextMessage(in)
		if errors.Cause(err) == io.EOF {
			log.Println("audiosocket closed")
			return
		}
		if m.Kind() == audiosocket.KindHangup {
			log.Println("audiosocket received hangup command")
			return
		}
		if m.Kind() == audiosocket.KindError {
			log.Println("error from audiosocket")
			continue
		}
		if m.Kind() != audiosocket.KindSlin {
			log.Println("ignoring non-slin message", m.Kind())
			continue
		}
		if m.ContentLength() < 1 {
			log.Println("no content")
			continue
		}
		if err = out.Send(&speechv1.StreamingRecognizeRequest{
			StreamingRequest: &speechv1.StreamingRecognizeRequest_AudioContent{
				AudioContent: m.Payload(),
			},
		}); err != nil {
			if err == io.EOF {
				return
			}
			log.Println("failed to send audio data for recognition:", err)
		}
	}
}

func sendAudio(w io.Writer, data []byte) error {
	return audiosocket.SendSlinChunks(w, slinChunkSize, data)
}

func processCommand(ctx context.Context, rw io.ReadWriter) (string, error) {
	cmd, err := recognizeRequest(ctx, rw)
	if err != nil {
		return "Sorry, I failed to listen to you", errors.Wrap(err, "failed to recognize request")
	}
	color.Green(cmd)

	if containsAny(cmd, "laugh", "joke") {
		return "ha ha ha", nil
	}
	if containsAny(cmd, "bye", "goodbye", "hangup", "hang up") {
		return "", ErrHangup
	}

	return "", nil
}

func speak(ctx context.Context, rw io.ReadWriter, msg string) error {
	resp, err := tts.SynthesizeSpeech(ctx, ttsRequest(msg))
	if err != nil {
		return errors.Wrap(err, "failed to synthesize speech")
	}
	if err = sendAudio(rw, resp.GetAudioContent()); err != nil {
		return errors.Wrap(err, "failed to send speech to Asterisk")
	}
	return nil
}

func containsAny(in string, refs ...string) bool {
	inLower := strings.ToLower(in)
	for _, r := range refs {
		if strings.Contains(inLower, r) {
			return true
		}
	}
	return false
}
