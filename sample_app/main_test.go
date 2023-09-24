package main

import (
	"github.com/quasi-go/di"
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	// Set logger
	logger := log.New(os.Stdout, "DI: ", 0)
	di.SetLogger(logger)
	di.SetLogLevel(di.LOG_LEVEL_ALL)
}
