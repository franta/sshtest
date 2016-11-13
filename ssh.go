package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

const key = `-----BEGIN RSA PRIVATE KEY-----

-----END RSA PRIVATE KEY-----`

func run(server string) bool {

	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		fmt.Println("ERROR:", server, err.Error())
		return false
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", server+":22", config)
	if err != nil {
		fmt.Println("ERROR:", server, err.Error())
		return false
	}
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("ERROR:", server, err.Error())
		return false
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("ls"); err != nil {
		fmt.Println("ERROR:", server, err.Error())
		return false
	}
	fmt.Println("OK:", server)
	return true
}

func runWIthTimeout(server string) bool {
	timeout := time.After(3 * time.Second)
	resp := make(chan bool)
	go func(server string) {
		resp <- run(server)
	}(server)
	select {
	case res := <-resp:
		return res
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
