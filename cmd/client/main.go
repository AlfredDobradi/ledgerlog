package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/AlfredDobradi/ledgerlog/internal/ssh"
)

var (
	privKeyPath string
	email       string
)

func main() {
	flag.StringVar(&privKeyPath, "private-key", "~/.ssh/id_rsa", "Key to SSH private key for signing")
	flag.StringVar(&email, "email", "", "Email address to auth with")
	flag.Parse()

	if email == "" {
		log.Panicln("empty email")
	}

	sshClient, err := ssh.NewClient(privKeyPath, email)
	if err != nil {
		log.Panicf("New SSH Client: %v", err)
	}

	data := map[string]string{
		"test": "hello",
	}
	raw, err := json.Marshal(data)
	if err != nil {
		log.Panicf("Marshal: %v", err)
	}

	body := bytes.NewBuffer(raw)
	r, err := http.NewRequest(http.MethodPost, "http://localhost:8080/test", body)
	if err != nil {
		log.Panicf("New request: %v", err)
	}
	if err := sshClient.SignRequest(r, body.Bytes()); err != nil {
		log.Panicf("Signing request: %v", err)
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Panicf("Request: %v", err)
	}

	log.Println(res.Status)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panicf("Read body: %v", err)
	}
	log.Println(string(resBody))
}
