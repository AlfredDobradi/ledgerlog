package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/AlfredDobradi/ledgerlog/internal/ssh"
)

type Handler struct {
	DB *badgerdb.DB
}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	auth, err := ssh.GetAuthFromRequest(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	pkey, err := h.DB.GetPublicKey(auth.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if pkey == nil {
		log.Println("Auth failed")
		http.Error(w, "Auth failed", http.StatusForbidden)
		return
	}

	if err := pkey.Verify(data, auth.Signature); err != nil {
		http.Error(w, "Invalid signature", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK")) // nolint
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request models.RegisterRequest
	if err := json.Unmarshal(body, &request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.DB.RegisterUser(request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.RegisterResponse{
		Timestamp: time.Now().Format(time.RFC1123),
		Message:   "OK",
	}
	raw, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(raw) // nolint
}

func (h *Handler) handleDebugKeys(w http.ResponseWriter, r *http.Request) {
	d, err := h.DB.GetKeys()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(d) // nolint
}
