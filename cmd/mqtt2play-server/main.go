package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/iamtio/mqtt2play"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Broker string `required:"true" desc:"MQTT Broker URI"`
	Qos    uint   `default:"0" desc:"MQTT QoS"`
	Prefix string `default:"mqtt2play/" desc:"MQTT topic prefix"`
	SfxDir string `default:"sfx" desc:"Directory with sound files"`
}

var conf Config

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}
func initConfig() {
	err := envconfig.Process("mqtt2play", &conf)
	if err != nil {
		envconfig.Usage("mqtt2play", &conf)
		log.Fatal(err.Error())
	}
}

func main() {
	initConfig()
	ctx := context.Background()
	playCtx, stopPlay := context.WithCancel(ctx)
	go func() {
		if err := mqtt2play.PlaySound(playCtx, fmt.Sprintf("%s/evil_laugh.wav", conf.SfxDir)); err != nil {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Printf("Shutting down...")
	stopPlay()
	time.Sleep(3 * time.Second)
}
