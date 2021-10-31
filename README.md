# What is it?
This project is a CLI appliction what subscribes to MQTT topic and plays sounds from files stored on disk

At now, application supports .wav and .mp3 file formats
# Usage
1. Place your sound files into `sfx/` directory
1. Start application
1. Send filename to `mqtt2play/play` topic as text
1. Enjoy

# Contributing
This project is written in Golang. If you want to contribute code:

1. Ensure you are running golang version 1.16 or greater for go module support
1. Check-out the project: `git clone https://github.com/iamtio/mqtt2play && cd mqtt2play`
1. Make sure `libasound2-dev` installed on your system (Debian, Ubuntu)
1. Make changes to the code
1. Build the project, e.g. via `go build ./cmd/mqtt2play-server`
1. Evaluate and test your changes `MQTT2PLAY_BROKERS=localhost:1883; ./mqtt2play-server`
1. Make a pull request