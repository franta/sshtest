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

func run(server string, ret chan<- error) {

	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		ret <- err
		return
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
		ret <- err
		return
	}
	session, err := client.NewSession()
	if err != nil {
		ret <- err
		return
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("ls"); err != nil {
		ret <- err
		return
	}
	ret <- nil
}

func runWIthTimeout(server string, ret chan<- bool) {
	// EDIT timeout
	timeout := time.After(3 * time.Second)
	resp := make(chan error)
	go run(server, resp)
	select {
	case e := <-resp:
		if e != nil {
			fmt.Println("ERROR:", server, e.Error())
			ret <- false
		} else {
			fmt.Println("OK:", server)
			ret <- true
		}
	case <-timeout:
		fmt.Println("ERROR:", server, "connection timeout")
		ret <- false
	}
}

func main() {
	hosts := os.Args[1:]
	results := make(chan bool)
	for _, hostname := range hosts {
		go runWIthTimeout(hostname, results)
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
