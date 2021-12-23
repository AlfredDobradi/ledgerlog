package ssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type AuthData struct {
	Email     string
	Signature *ssh.Signature
}

func ParsePublicKey(path string) (ssh.PublicKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading public key file: %w", err)
	}

	pk, _, _, _, err := ssh.ParseAuthorizedKey(raw)
	if err != nil {
		return nil, fmt.Errorf("Error parsing key: %w", err)
	}

	pkey, err := ssh.ParsePublicKey(pk.Marshal())
	if err != nil {
		return nil, fmt.Errorf("Error parsing public key: %w", err)
	}

	return pkey, nil
}

func GetAuthFromRequest(r *http.Request) (*AuthData, error) {
	auth := r.Header.Get("Authorization")

	authParts := strings.Split(auth, ":")

	sigJSON, err := base64.StdEncoding.DecodeString(strings.Trim(authParts[1], " "))
	if err != nil {
		return nil, fmt.Errorf("Error parsing signature: %w", err)
	}

	var sig ssh.Signature
	if err := json.Unmarshal(sigJSON, &sig); err != nil {
		return nil, fmt.Errorf("Error unmarshaling signature: %w", err)
	}

	return &AuthData{Email: authParts[0], Signature: &sig}, nil
}
