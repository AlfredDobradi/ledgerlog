package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	m.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		d, err := bdb.GetKeys()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(d) // nolint
	})

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: m,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go server.ListenAndServe() // nolint

	running := true
	for running {
		<-sigs
		log.Println("Received signal")
		server.Shutdown(context.Background()) // nolint
		log.Println("Shutdown web service")
		badgerdb.Close() // nolint
		log.Println("Closed database")

		running = false
	}
}
