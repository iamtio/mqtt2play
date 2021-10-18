package mqtt2play

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/h2non/filetype"
	log "github.com/sirupsen/logrus"
)

func PlaySound(ctx context.Context, filepath string) error {
	fileType, err := GetAudioFileType(filepath)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	var streamer beep.StreamSeekCloser
	var format beep.Format
	switch fileType {
	case "mp3":
		streamer, format, err = mp3.Decode(f)
	case "wav":
		streamer, format, err = wav.Decode(f)
	default:
		return fmt.Errorf("file type - %s is not supported", fileType)
	}
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	for {
		select {
		case <-done:
			log.Info("Stop playing due to end of audio")
			return nil
		case <-ctx.Done():
			log.Info("Stop playing due to context done")
			speaker.Close()
			return nil
		}
	}
}

func GetAudioFileType(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return "", err
	}

	if !filetype.IsAudio(head) {
		return "", fmt.Errorf("file is not audio format")
	}

	kind, err := filetype.Match(head)
	if err != nil {
		return "", err
	}
	return kind.Extension, nil
}
