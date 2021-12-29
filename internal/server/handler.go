package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/AlfredDobradi/ledgerlog/internal/ssh"
)

const (
	RouteAPIRegister string = "/api/register"
	RouteAPISend     string = "/api/send"
	RouteAPIPosts    string = "/api"
	RouteDebugKeys   string = "/debug/keys"
)

type Handler struct {
	DB *badgerdb.DB
}

func (h *Handler) handleSend(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		http.Error(w, "Invalid signature", http.StatusInternalServerError)
		return
	}

	var postRequest models.SendPostRequest
	if err := json.Unmarshal(data, &postRequest); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.DB.AddPost(auth.Email, postRequest); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (h *Handler) handlePosts(w http.ResponseWriter, r *http.Request) {
	d, err := h.DB.GetPosts()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.Slice(d, func(i, j int) bool {
		return d[i].Timestamp.After(d[j].Timestamp)
	})

	data, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data) // nolint
}
