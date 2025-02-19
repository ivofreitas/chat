package main

import (
	"github.com/ivofreitas/chat/internal/bot/application"
	"github.com/ivofreitas/chat/pkg/log"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	log.Init()

	receiver := application.NewReceiver()

	wg.Add(1)

	go func() {
		defer wg.Done()
		receiver.Run()
	}()

	wg.Wait()
}
