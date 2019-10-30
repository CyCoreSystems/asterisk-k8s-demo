package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/CyCoreSystems/audiosocket"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/fatih/color"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// MaxCallDuration is the maximum amount of time to allow a call to be up before it is terminated.
const MaxCallDuration = 5 * time.Minute

// MaxRecognitionDuration is the maximum amount of time to allow for a single voice recognition session to complete
const MaxRecognitionDuration = time.Minute

const listenAddr = ":8080"
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
		color.Magenta("ending call %s", id.String())

		// Tell caller good-bye
		speak(ctx, c, partingMessage) // nolint: errcheck

		// Tell AudioSocket to shut down, if it is still up
		c.Write(audiosocket.HangupMessage()) // nolint: errcheck

		cancel()
	}()

	id, err = getCallID(c)
	if err != nil {
		log.Println("failed to get call ID:", err)
		return
	}
	color.Magenta("processing call %s", id.String())

	a := &App{
		c:  c,
		id: id,
	}
	if err := a.Run(ctx); err != nil {
		if err == ErrHangup {
			return
		}
		if ctx.Err() != nil {
			return
		}
		log.Println(err)
	}
	return
}
