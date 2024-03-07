package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/vivekkumar-git/k8s-event-handler/config"
)

var (
	configFile     = "config.yaml"
	configFilePath = flag.String("config", "", "config folder path")
)

func main() {
	flag.Parse()
	filepath := filepath.Join(*configFilePath, configFile)

	conf, err := config.New(filepath)
	if err != nil {
		log.Fatalf("Error in loading configuration. Error:%s", err.Error())
	}

	informer, err := informer.NewInformer(conf)
	if err != nil {
		log.Fatalf("Error in loading configuration. Error:%s", err.Error())
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	informer.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSTOP)

	<-sigterm
}
