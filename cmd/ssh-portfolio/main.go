package main

import (
	"log"

	sshserver "github.com/Shbhom/ssh-portfolio/internal/ssh-server"
)

func main() {
	addr := ":23234"
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
