package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/AlfredDobradi/ledgerlog/internal/cli"
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

	fmt.Printf("Registering %s with the given public key at %s...", cmd.Email, instanceURL)

	raw, fileErr := os.ReadFile(cmd.KeyPath)
	if fileErr != nil {
		cli.Failure()
		return fmt.Errorf("Failed to open %s: %w", cmd.KeyPath, fileErr)
	}

	request := models.RegisterRequest{
		Email:     cmd.Email,
		Name:      cmd.PreferredName,
		PublicKey: strings.Trim(string(raw), "\n"),
	}

	jsonRaw, jsonErr := json.Marshal(request)
	if jsonErr != nil {
		cli.Failure()
		return fmt.Errorf("Failed to marshal request %w", jsonErr)
	}
	body := bytes.NewBuffer(jsonRaw)

	r, requestErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", instanceURL, server.RouteAPIRegister), body)
	if requestErr != nil {
		cli.Failure()
		return fmt.Errorf("Failed to create HTTP request: %w", requestErr)
	}

	client := http.DefaultClient
	response, clientErr := client.Do(r)
	if clientErr != nil {
		cli.Failure()
		return fmt.Errorf("Failed to send HTTP request: %w", clientErr)
	}

	responseBody, readBodyErr := io.ReadAll(response.Body)
	if readBodyErr != nil {
		cli.Failure()
		return fmt.Errorf("Failed to read response body: %w", readBodyErr)
	}

	if response.StatusCode != http.StatusOK {
		cli.Failure()
		return fmt.Errorf("The server returned an error: %s", response.Status)
	}

	var registerResponse models.RegisterResponse
	if err := json.Unmarshal(responseBody, &registerResponse); err != nil {
		cli.Failure()
		return fmt.Errorf("Failed to unmarshal response body: %w", err)
	}

	cli.Success()

	return nil
}
