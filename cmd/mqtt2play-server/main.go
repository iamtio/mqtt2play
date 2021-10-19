package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"
	"unicode"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iamtio/mqtt2play"
	"github.com/iamtio/mqtt2play/eclogrus"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Brokers  []string       `required:"true" desc:"MQTT Brokers URIs"`
	Username string         `desc:"MQTT Username"`
	Password string         `desc:"MQTT Password"`
	ClientID string         `desc:"MQTT Client ID"`
	Timeout  time.Duration  `default:"3s" desc:"MQTT connection timeout"`
	Qos      uint           `default:"0" desc:"MQTT QoS"`
	Prefix   string         `default:"mqtt2play/" desc:"MQTT topic prefix"`
	SfxDir   string         `default:"sfx/" desc:"Directory with sound files"`
	LogLevel eclogrus.Level `default:"info" desc:"Log level"`
}

var conf Config

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.TraceLevel)
	log.SetReportCaller(true)
}
func initConfig() {
	err := envconfig.Process("mqtt2play", &conf)
	if err != nil {
		envconfig.Usage("mqtt2play", &conf)
		log.Fatal(err.Error())
	}
}

func main() {
	ctx := context.Background()
	initConfig()
	handlerCtx, cancel := context.WithCancel(ctx)
	mqttHandler(handlerCtx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Printf("Shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
}

func mqttHandler(ctx context.Context) {
	// Connection opts
	co := mqtt.NewClientOptions()

	co.ClientID = conf.ClientID
	for _, b := range conf.Brokers {
		co.AddBroker(b)
	}
	co.Username = conf.Username
	co.Password = conf.Password
	client := mqtt.NewClient(co)

	token := client.Connect()
	token.WaitTimeout(conf.Timeout)
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	if err := publishAvailableSounds(client); err != nil {
		log.Fatal(err)
	}
	subscribePlay(ctx, client).WaitTimeout(conf.Timeout)
}
func subscribePlay(ctx context.Context, client mqtt.Client) mqtt.Token {
	return client.Subscribe(mqttTopic("play"), byte(conf.Qos), func(c mqtt.Client, m mqtt.Message) {
		playLogger := log.WithFields(log.Fields{
			"topic": mqttTopic("play"),
		})
		mqttPayload := string(m.Payload())
		for n, r := range mqttPayload {
			if !unicode.IsPrint(r) {
				playLogger.Warnf("non printable character at: %d", n)
				return
			}
		}
		sfxPath := fmt.Sprintf("%s/%s", conf.SfxDir, mqttPayload)
		if _, err := os.Stat(sfxPath); os.IsNotExist(err) {
			playLogger.Warnf("sfx file doesn't exist: %s", sfxPath)
			return
		}

		if err := mqtt2play.PlaySound(ctx, sfxPath); err != nil {
			log.Fatal(err)
		}
	})
}
func publishAvailableSounds(client mqtt.Client) error {
	matched := mqtt2play.FindSfx(conf.SfxDir)
	encoded, err := json.Marshal(matched)
	if err != nil {
		return err
	}
	token := client.Publish(mqttTopic("available"), byte(conf.Qos), false, encoded)
	token.WaitTimeout(conf.Timeout)
	return token.Error()
}

func mqttTopic(path string) string {
	return fmt.Sprintf("%s%s", conf.Prefix, path)
}
