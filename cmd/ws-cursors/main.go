package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/error0x001/ws-cursors/internal"
)

func main() {
	config := internal.NewConfig()
	exit := make(chan bool)
	sender := internal.NewSender()
	go sender.Listen(exit)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		internal.Index(config.TemplatePath, config.GetWSPath(), writer, request)
	})
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		internal.WS(sender, writer, request)
	})
	go func() {
		err := http.ListenAndServe(config.GetOnlyPort(), nil)
		if err != nil {
			log.Fatalf("run error: %s", err.Error())
		}
	}()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	exit <- true
	time.Sleep(config.ShutDownTime)
}
