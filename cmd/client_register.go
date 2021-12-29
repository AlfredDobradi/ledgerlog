package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
)

type RegisterCmd struct {
	Email   string `help:"Email address to register" required:""`
	KeyPath string `help:"Path to public key" type:"existingfile" default:"~/.ssh/id_rsa.pub"`
}

func (cmd *RegisterCmd) Run(ctx *Context) error {
	raw, fileErr := os.ReadFile(cmd.KeyPath)
	if fileErr != nil {
		return fileErr
	}

	request := models.RegisterRequest{
		Email:     cmd.Email,
		PublicKey: string(raw),
	}

	jsonRaw, jsonErr := json.Marshal(request)
	if jsonErr != nil {
		return jsonErr
	}

	body := bytes.NewBuffer(jsonRaw)

	r, requestErr := http.NewRequest(http.MethodPost, "http://localhost:8080/register", body)
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

	fmt.Printf("Successfully registered %s with the given public key\n", cmd.Email)
	return nil
}
