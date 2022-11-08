package main

import (
	"flag"
	"log"

	pingdns "github.com/jamespwilliams/dns-over-ping"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logLevel := zap.LevelFlag("log-level", zap.InfoLevel, "one of DEBUG, INFO, WARN, ERROR. defaults to INFO")

	flag.Parse()

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = zap.NewAtomicLevelAt(*logLevel)
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("failed to build logger: %v", err)
	}

	panic(pingdns.NewServer(logger).Serve())
}
