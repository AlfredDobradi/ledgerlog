package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/server"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/AlfredDobradi/ledgerlog/internal/ssh"
)

type SendCmd struct {
	PrivKeyPath string `help:"Path to your private key" type:"existingfile" default:"~/.ssh/id_rsa" env:"LEDGER_PRIV_KEY"`
	Email       string `help:"Your registered email" required:"" env:"LEDGER_EMAIL"`
	InstanceURL string `help:"URL of the instance you want to send the post to" required:"" env:"LEDGER_URL"`

	Message string `arg:"" help:"The content of your post" required:""`
}

func (cmd *SendCmd) Run(ctx *Context) error {
	keyPath := config.GetSettings().User.PrivateKeyPath
	if cmd.PrivKeyPath != "" {
		keyPath = cmd.PrivKeyPath
	}
	email := config.GetSettings().User.Email
	if cmd.Email != "" {
		email = cmd.Email
	}
	sshClient, err := ssh.NewClient(keyPath, email)
	if err != nil {
		return err
	}

	postRequest := models.SendPostRequest{
		Message: cmd.Message,
	}
	raw, err := json.Marshal(postRequest)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(raw)
	instanceURL := config.GetSettings().Instance.URL
	if cmd.InstanceURL != "" {
		instanceURL = cmd.InstanceURL
	}

	url := fmt.Sprintf("%s%s", instanceURL, server.RouteAPISend)
	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	if err := sshClient.SignRequest(r, body.Bytes()); err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	log.Println(res.Status)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println(string(resBody))

	return nil
}

// var (
// 	privKeyPath string
// 	email       string
// )

// func main() {
// 	flag.StringVar(&privKeyPath, "private-key", "~/.ssh/id_rsa", "Key to SSH private key for signing")
// 	flag.StringVar(&email, "email", "", "Email address to auth with")
// 	flag.Parse()

// 	if email == "" {
// 		log.Panicln("empty email")
// 	}

// 	sshClient, err := ssh.NewClient(privKeyPath, email)
// 	if err != nil {
// 		log.Panicf("New SSH Client: %v", err)
// 	}

// 	data := map[string]string{
// 		"test": "hello",
// 	}
// 	raw, err := json.Marshal(data)
// 	if err != nil {
// 		log.Panicf("Marshal: %v", err)
// 	}

// 	body := bytes.NewBuffer(raw)
// 	r, err := http.NewRequest(http.MethodPost, "http://localhost:8080/test", body)
// 	if err != nil {
// 		log.Panicf("New request: %v", err)
// 	}
// 	if err := sshClient.SignRequest(r, body.Bytes()); err != nil {
// 		log.Panicf("Signing request: %v", err)
// 	}

// 	res, err := http.DefaultClient.Do(r)
// 	if err != nil {
// 		log.Panicf("Request: %v", err)
// 	}

// 	log.Println(res.Status)
// 	resBody, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		log.Panicf("Read body: %v", err)
// 	}
// 	log.Println(string(resBody))
// }
