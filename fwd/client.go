package fwd

import (
	"context"
	"fmt"
	"log"
	"net"
)

type StartCh struct {
	RemotePort string
	LocalPort  string
}

func (pf *StartCh) GetId() string {
	return fmt.Sprintf("%s->%s", pf.LocalPort, pf.RemotePort)
}

func portAvailable(port string) bool {
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

func GetController(sshUser, sshHost, sshPort, keyPath string) (startCh, cancelCh chan StartCh, exitCh chan struct{}, lports map[string]string, runner func()) {
	startCh = make(chan StartCh)
	cancelCh = make(chan StartCh)
	exitCh = make(chan struct{})

	lports = make(map[string]string)
	return startCh, cancelCh, exitCh, lports, func() {
		ctxs := make(map[string]context.CancelFunc)

		// cf := func() {}

		for {
			select {
			case <-exitCh:
				fmt.Println("Exiting...")
				return
			case i := <-startCh:

				// if lports[i.LocalPort] != false {
				// 	fmt.Println("port already in use", i.LocalPort)
				// 	continue
				// }

				if !portAvailable(i.LocalPort) {
					fmt.Printf("port %s already in use: %s\n", i.LocalPort, lports[i.LocalPort])
					continue
				}

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				// cf = cancel
				pf, err := NewForwarder(i.LocalPort, i.RemotePort, sshUser, sshHost, sshPort, keyPath)
				if err != nil {
					log.Println(err)
					continue
				}

				ctxs[i.GetId()] = cancel
				lports[i.LocalPort] = i.GetId()
				go func() {
					if err := pf.Start(ctx); err != nil {
						fmt.Println(err)
					}
				}()
				fmt.Printf("[+] %s\n", i.GetId())
			case i := <-cancelCh:
				if cancel, exists := ctxs[i.GetId()]; exists {
					cancel()
					delete(ctxs, i.GetId())
					delete(lports, i.LocalPort)
					fmt.Printf("[-] %s\n", i.GetId())
				} else {
					fmt.Printf("No forwarder to cancel with id %s\n", i.GetId())
				}
			}
		}
	}
}
