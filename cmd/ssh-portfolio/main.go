package main

import (
	"flag"
	"fmt"
	"log"

	sshserver "github.com/Shbhom/ssh-portfolio/internal/ssh-server"
)

func main() {
	port := flag.Int("port", 22, "port on which wish ssh will run")

	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	hostKeyPath := "ssh_host_ed25519"

	srv, err := sshserver.New(addr, hostKeyPath)
	if err != nil {
		log.Fatalf("failed to create ssh server: %v", err)
	}

	log.Printf("Starting SSH server on %s ...", addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
