package main

import (
	"github.com/ivofreitas/chat/internal/auth/application"
	"github.com/ivofreitas/chat/pkg/log"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	log.Init()

	server := application.NewServer()

	wg.Add(1)

	go func() {
		defer wg.Done()
		server.Run()
	}()

	wg.Wait()
}
