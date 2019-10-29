package main

import (
	"context"
	"io"
	"log"
	"net"
	"strings"

	"github.com/CyCoreSystems/audiosocket"
	"github.com/fatih/color"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	speechv1 "google.golang.org/genproto/googleapis/cloud/speech/v1"
	texttospeechv1 "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

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
							Phrases: keyPhrases,
						},
					},
				},
			},
		},
	}); err != nil {
		return "", errors.Wrap(err, "failed to send recognition config")
	}

	go pipeFromAsterisk(ctx, cancel, r, svc)

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
				color.Green(alt.Transcript)
				return alt.Transcript, nil
			}
		}
	}
	return "", nil
}

func pipeFromAsterisk(ctx context.Context, cancel context.CancelFunc, in io.Reader, out speechv1.Speech_StreamingRecognizeClient) {
	var err error
	var m audiosocket.Message

	defer out.CloseSend() // nolint: errcheck
	defer cancel()

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
