package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/kloudlite/docker-ssh-proxy/fwd"
)

func main() {
	// Configuration
	// remoteHost := "localhost:8080"
	sshUser := "kl"
	sshHost := "localhost"
	sshPort := "61100"
	keyPath := filepath.Join(xdg.Home, ".ssh", "id_rsa")

	// Run controller in a goroutine
	startCh, cancelCh, exitCh, lports, runner := fwd.GetController(sshUser, sshHost, sshPort, keyPath)
	go runner()

	// Example usage
	for i := 0; i < 100; i++ {
		startCh <- fwd.StartCh{
			RemotePort: "8080",
			LocalPort:  fmt.Sprint(8000 + i),
		}
		time.Sleep(10 * time.Millisecond)
	}

	for i := 0; i < 100; i++ {
		fmt.Println(lports[fmt.Sprint(8000+i)])

		startCh <- fwd.StartCh{
			RemotePort: "8080",
			LocalPort:  fmt.Sprint(8000 + i),
		}
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(10 * time.Second)

	for i := 0; i < 100; i++ {
		cancelCh <- fwd.StartCh{
			RemotePort: "8080",
			LocalPort:  fmt.Sprint(8000 + i),
		}
		time.Sleep(10 * time.Millisecond)
	}

	exitCh <- struct{}{}
	time.Sleep(2 * time.Second)
}

// Controller function to handle forwarding processes
