package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/AlfredDobradi/ledgerlog/internal/ssh"
	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
)

func main() {
	opts := badger.DefaultOptions("./tmp")
	bdb, err := badgerdb.GetConnection(opts)
	if err != nil {
		log.Panicln(err)
	}
	m := mux.NewRouter()

	m.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
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

		pkey, err := bdb.GetPublicKey(auth.Email)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := pkey.Verify(data, auth.Signature); err != nil {
			http.Error(w, "Invalid signature", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK")) // nolint
	})

	m.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
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

		if err := bdb.RegisterUser(request); err != nil {
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
	})

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: m,
	}

	log.Panicln(server.ListenAndServe())
}
