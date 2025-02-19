package main

import (
	"github.com/ivofreitas/chat/internal/chat/application"
	"github.com/ivofreitas/chat/pkg/log"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	log.Init()

	server := application.NewServer()
	receiver := application.NewReceiver()

	wg.Add(2)

	go func() {
		defer wg.Done()
		server.Run()
	}()

	go func() {
		defer wg.Done()
		receiver.Run()
	}()

	wg.Wait()
}
