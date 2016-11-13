package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

// ADD your key here
const key = `-----BEGIN RSA PRIVATE KEY-----

-----END RSA PRIVATE KEY-----`

func run(server string) error {

	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		// EDIT username
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	// EDIT port
	client, err := ssh.Dial("tcp", server+":22", config)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("ls"); err != nil {
		return err
	}
	return nil
}

func runWIthTimeout(server string) bool {
	// EDIT timeout
	timeout := time.After(3 * time.Second)
	resp := make(chan error)
	go func(server string) {
		resp <- run(server)
	}(server)
	select {
	case e := <-resp:
		if e != nil {
			fmt.Println("ERROR:", server, e.Error())
			return false
		} else {
			fmt.Println("OK:", server)
			return true
		}
	case <-timeout:
		fmt.Println("ERROR:", server, "connection timeout")
		return false
	}
}

func main() {
	hosts := os.Args[1:]
	results := make(chan bool, 10)
	for _, hostname := range hosts {
		go func(hostname string) {
			results <- runWIthTimeout(hostname)
		}(hostname)
	}
	ok := true
	for i := 0; i < len(hosts); i++ {
		res := <-results
		if !res {
			ok = false
		}
	}
	if !ok {
		os.Exit(1)
	}
}
