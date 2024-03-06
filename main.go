package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"supernet.tools/tcp-proxy-server/config"
	"supernet.tools/tcp-proxy-server/network"
	"supernet.tools/tcp-proxy-server/service"
)

var conf *config.AppConf

func init() {
	// Initialize logger
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		CallerWithSkipFrameCount(2). // good for debugging, but will be slow on production
		Logger()

	// Prepare configuration
	var confPath string
	flag.StringVar(&confPath, "config", "", "Path to config file (see more in config/example.yaml)")
	flag.Parse()

	conf = config.Init(confPath)
}

func main() {
	// Initialize encryptor
	encryptor := service.NewDummyEncryptor(conf)

	// Initialize proxy server
	server := network.NewProxy(conf, encryptor)

	// Start listener
	server.Start()
}
