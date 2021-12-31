package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/server"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
)

type RegisterCmd struct {
	Email         string `help:"Email address to register" required:"" env:"LEDGER_EMAIL"`
	PreferredName string `help:"Your preferred name on the site" required:"" env:"LEDGER_PREFERRED_NAME"`
	KeyPath       string `help:"Path to public key" type:"existingfile" required:"" default:"~/.ssh/id_rsa.pub" env:"LEDGER_PUB_KEY"`
	InstanceURL   string `help:"URL of the instance you want to send the post to" required:"" default:"http://localhost:8080" env:"LEDGER_INSTANCE_URL"`
}

func (cmd *RegisterCmd) Run(ctx *Context) error {
	instanceURL := config.GetSettings().Instance.URL
	if cmd.InstanceURL != "" {
		instanceURL = cmd.InstanceURL
	}

	raw, fileErr := os.ReadFile(cmd.KeyPath)
	if fileErr != nil {
		return fileErr
	}

	request := models.RegisterRequest{
		Email:     cmd.Email,
		Name:      cmd.PreferredName,
		PublicKey: strings.Trim(string(raw), "\n"),
	}

	jsonRaw, jsonErr := json.Marshal(request)
	if jsonErr != nil {
		return jsonErr
	}
	body := bytes.NewBuffer(jsonRaw)

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
		if ctx.Debug || config.GetSettings().Debug {
			output = fmt.Sprintf("%s - %s", output, responseBody)
		}
		return fmt.Errorf("%s", output)
	}

	var registerResponse models.RegisterResponse
	if err := json.Unmarshal(responseBody, &registerResponse); err != nil {
		return err
	}

	fmt.Printf("Successfully registered %s with the given public key\n", cmd.Email)
	return nil
}
