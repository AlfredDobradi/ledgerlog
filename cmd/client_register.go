package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/server"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
)

type RegisterCmd struct {
	Email       string `help:"Email address to register" required:"" env:"LEDGER_EMAIL"`
	KeyPath     string `help:"Path to public key" type:"existingfile" default:"~/.ssh/id_rsa.pub" env:"LEDGER_PUB_KEY"`
	InstanceURL string `help:"URL of the instance you want to send the post to" required:"" env:"LEDGER_URL"`
}

func (cmd *RegisterCmd) Run(ctx *Context) error {
	keyPath := config.GetSettings().User.PublicKeyPath
	if cmd.KeyPath != "" {
		keyPath = cmd.KeyPath
	}

	raw, fileErr := os.ReadFile(keyPath)
	if fileErr != nil {
		return fileErr
	}

	email := config.GetSettings().User.Email
	if cmd.Email != "" {
		email = cmd.Email
	}

	request := models.RegisterRequest{
		Email:     email,
		PublicKey: string(raw),
	}

	jsonRaw, jsonErr := json.Marshal(request)
	if jsonErr != nil {
		return jsonErr
	}
	body := bytes.NewBuffer(jsonRaw)

	instanceURL := config.GetSettings().Instance.URL
	if cmd.InstanceURL != "" {
		instanceURL = cmd.InstanceURL
	}
	r, requestErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", instanceURL, server.RouteAPIRegister), body)
	if requestErr != nil {
		return requestErr
	}

	client := http.DefaultClient
	response, clientErr := client.Do(r)
	if clientErr != nil {
		return clientErr
	}

	responseBody, readBodyErr := io.ReadAll(response.Body)
	if readBodyErr != nil {
		return readBodyErr
	}

	if response.StatusCode != http.StatusOK {
		output := response.Status
		if ctx.Debug {
			output = fmt.Sprintf("%s - %s", output, responseBody)
		}
		return fmt.Errorf("%s", output)
	}

	var registerResponse models.RegisterResponse
	if err := json.Unmarshal(responseBody, &registerResponse); err != nil {
		return err
	}

	fmt.Printf("Successfully registered %s with the given public key\n", email)
	return nil
}
