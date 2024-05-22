package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/kloudlite/docker-ssh-proxy/fwd"
)

func main() {
	// Configuration
	// localPort := "8080"
	remoteAddr := "127.0.0.1:8080"
	sshUser := "kl"
	sshHost := "localhost"
	sshPort := "61100"
	keyPath := filepath.Join(xdg.Home, ".ssh", "id_rsa")

	ctxs := make([]context.CancelFunc, 100)

	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		pf, err := fwd.NewForwarder(fmt.Sprint(8000+i), remoteAddr, sshUser, sshHost, sshPort, keyPath)
		if err != nil {
			log.Println(err)
		}

		ctxs[i] = cancel
		go pf.Start(ctx)
	}

	time.Sleep(10 * time.Second)

	for i, cancel := range ctxs {
		fmt.Printf("canceling %d", i+1)
		time.Sleep(time.Second * 2)
		cancel()
	}

	select {}
}
