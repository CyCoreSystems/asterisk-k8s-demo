package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/CyCoreSystems/audiosocket"
	"github.com/ericchiang/k8s"
	"github.com/ericchiang/k8s/apis/apps/v1"
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
const redisAddr = "redis:6379"
const languageCode = "en-US"

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

	/*
		if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", googleCreds)
		}
	*/

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
	ctx, cancel := context.WithTimeout(pCtx, MaxCallDuration)

	defer func() {
		cancel()

		if _, err := c.Write(audiosocket.HangupMessage()); err != nil {
			log.Println("failed to send hangup message:", err)
		}

	}()

	id, err := getCallID(c)
	if err != nil {
		log.Println("failed to get call ID:", err)
		return
	}
	log.Printf("processing call %s", id.String())

	resp, err := tts.SynthesizeSpeech(ctx, greeting())
	if err != nil {
		log.Println("failed to synthesize greeting:", err)
		return
	}
	if err = sendAudio(c, resp.GetAudioContent()); err != nil {
		log.Println("failed to send greeting to Asterisk:", err)
	}

	for ctx.Err() == nil {
      log.Println("waiting for command")
		resp, err := processCommand(ctx, c)
		if err != nil {
			log.Println("failed to process command:", err)
		}
		if resp != "" {
			if err = speak(ctx, c, resp); err != nil {
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
   defer svc.CloseSend()

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
								"asterisk",
								"asterisks",
								"bye",
								"goodbye",
								"hello",
								"kamailio",
								"kamailios",
								"proxy",
								"proxies",
								"scale",
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

func greeting() *texttospeechv1.SynthesizeSpeechRequest {
	return ttsRequest("Hello.  How may I help you?")
}

func parting() *texttospeechv1.SynthesizeSpeechRequest {
	return ttsRequest("Good-bye.  Thank you for playing.")
}

func scaleAsterisk(ctx context.Context, count int, w io.Writer) (string, error) {
	if count > 10 {
		return "Sorry, I can only scale to ten Asterisk instances", nil
	}

	if err := scaleDeployment(ctx, "asterisk", "voip", int32(count)); err != nil {
		return "Sorry, I failed to scale Asterisk", err
	}

	if count == 1 {
		return "Scaled asterisk to 1 instance", nil
	}
	return fmt.Sprintf("Asterisk has been scaled to %d instances.", count), nil
}

func scaleKamailio(count int, w io.Writer) (string, error) {
	return "I cannot scale proxies yet", errors.New("not implemented")
	/*
		var req *texttospeechv1.SynthesizeSpeechRequest
		var err error

			if count > 2 {
				return "Sorry, I can only scale to two proxy instances", nil
			}
			if count == 1 {
				return "I have scaled the proxies to 1 instance", nil
			}
			return fmt.Sprintf("Proxies have been scaled to %dinstances", count), nil
			req = ttsRequest(fmt.Sprintf("Scaling proxy to %d instances.", count))
	*/

}

func pipeFromAsterisk(ctx context.Context, in io.Reader, out speechv1.Speech_StreamingRecognizeClient) {
	var err error
	var m audiosocket.Message

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
				log.Println("recognition client closed")
				return
			}
			log.Println("failed to send audio data for recognition:", err)
		}
	}
}

func sendAudio(w io.Writer, data []byte) error {

	var chunks int

	for i := 0; i < len(data); {
		var chunkLen = slinChunkSize
		if i+slinChunkSize > len(data) {
			chunkLen = len(data) - i
		}
		if _, err := w.Write(audiosocket.SlinMessage(data[i : i+chunkLen])); err != nil {
			return errors.Wrap(err, "failed to write chunk to audiosocket")
		}
		chunks++
		i += chunkLen
	}

	return nil
}

func processCommand(ctx context.Context, rw io.ReadWriter) (string, error) {
	cmd, err := recognizeRequest(ctx, rw)
	if err != nil {
		return "Sorry, I failed to listen to you", errors.Wrap(err, "failed to recognize request")
	}

	switch {
	case strings.Contains(cmd, "scale"):
		switch {
		case strings.Contains(cmd, "asterisk") || strings.Contains(cmd, "astris"):
			count, err := parseCount(cmd)
			if err != nil {
				return "Sorry, I could not understand how many Asterisk instances to scale to", errors.Wrapf(err, "failed to parse count in phrase (%s)", cmd)
			}
			return scaleAsterisk(ctx, count, rw)
		case strings.Contains(cmd, "prox") || strings.Contains(cmd, "kamailio"):
			count, err := parseCount(cmd)
			if err != nil {
				return "Sorry, I could not understand how many Asterisk instances to scale to", errors.Wrapf(err, "failed to parse count in phrase (%s)", cmd)
			}
			return scaleKamailio(count, rw)
		}
	case strings.Contains(cmd, "hello"):
		return "Hello.  How may I help you?", nil
	case strings.Contains(cmd, "bye"):
		return "Good bye!", ErrHangup
	}

   log.Println("failed to parse command:", cmd)
	return "Sorry, I don't know how to do that", nil
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

func parseCount(msg string) (int, error) {
	for _, word := range strings.Split(msg, " ") {
		// Try direct number parsing
		if count, err := strconv.Atoi(word); err == nil {
			return int(count), nil
		}

		switch word {
		case "one":
			return 1, nil
		case "two":
			return 2, nil
		case "three":
			return 3, nil
		case "four":
			return 4, nil
		case "five":
			return 5, nil
		case "six":
			return 6, nil
		case "seven":
			return 7, nil
		case "eight":
			return 8, nil
		case "nine":
			return 9, nil
		case "ten":
			return 10, nil
		}
	}
	return 0, errors.New("failed to find count in message")
}

func scaleDeployment(ctx context.Context, name, namespace string, size int32) error {
	k, err := k8s.NewInClusterClient()
	if err != nil {
		return errors.Wrap(err, "failed to get kubernetes client")
	}

	d := new(v1.Deployment)
	if err = k.Get(ctx, namespace, name, d); err != nil {
		return errors.Wrap(err, "failed to retrieve current deployment")
	}

	current := d.GetSpec().GetReplicas()
	if current == size {
		return nil
	}

	d.GetSpec().Replicas = &size
	return k.Update(ctx, d)
}
