package ssh

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	AuthHeader string = "Authorization"
)

type AuthData struct {
	Email     string
	Signature *ssh.Signature
}

func ParsePublicKey(raw []byte) (ssh.PublicKey, error) {
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

type Client struct {
	Signer ssh.Signer
	Email  string
}

func NewClient(privKeyPath string, email string) (Client, error) {
	c := Client{
		Email: email,
	}
	der, err := os.ReadFile(privKeyPath)
	if err != nil {
		return c, err
	}

	key, err := ssh.ParseRawPrivateKey(der)
	if err != nil {
		return c, err
	}

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return c, err
	}

	c.Signer = signer

	return c, nil
}

func (c Client) SignRequest(r *http.Request, body []byte) error {
	signature, err := c.Signer.Sign(rand.Reader, body)
	if err != nil {
		log.Panicf("Sign: %v", err)
	}
	rawSig, err := json.Marshal(signature)
	if err != nil {
		log.Panicf("Marshal sig: %v", err)
	}
	sig := base64.StdEncoding.EncodeToString(rawSig)

	r.Header.Set(AuthHeader, fmt.Sprintf("%s:%s", c.Email, sig))
	return nil
}
