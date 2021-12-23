package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	privKeyPath string
)

func main() {
	flag.StringVar(&privKeyPath, "private-key", "~/.ssh/id_rsa", "Key to SSH private key for signing")
	flag.Parse()

	der, err := os.ReadFile(privKeyPath)
	if err != nil {
		log.Panicf("Readfile: %v", err)
	}

	key, err := ssh.ParseRawPrivateKey(der)
	if err != nil {
		log.Panicf("Parse: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		log.Panicf("New signer: %v", err)
	}

	data := map[string]string{
		"test": "hello",
	}
	raw, err := json.Marshal(data)
	if err != nil {
		log.Panicf("Marshal: %v", err)
	}

	signature, err := signer.Sign(rand.Reader, raw)
	if err != nil {
		log.Panicf("Sign: %v", err)
	}
	rawSig, err := json.Marshal(signature)
	if err != nil {
		log.Panicf("Marshal sig: %v", err)
	}
	sig := base64.StdEncoding.EncodeToString(rawSig)

	body := bytes.NewBuffer(raw)
	r, err := http.NewRequest(http.MethodPost, "http://localhost:8080/test", body)
	if err != nil {
		log.Panicf("New request: %v", err)
	}
	r.Header.Add("X-Signature", sig)

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
