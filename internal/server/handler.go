package server

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/database"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/AlfredDobradi/ledgerlog/internal/ssh"
)

const (
	RouteAPIRegister string = "/api/register"
	RouteAPISend     string = "/api/send"
	RouteAPIWS       string = "/api/stream"
	RouteAPIPosts    string = "/api"
	RouteIndex       string = "/"

	PostsPerPage int = 30

	FetchMaxErrors int = 3
)

type Handler struct{}

func NewHandler() *Handler {
	h := &Handler{}

	return h
}

func (h *Handler) handleAPISend(w http.ResponseWriter, r *http.Request) {
	db, err := database.GetDB()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close(context.TODO())

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

	user, err := db.FindUser(map[string]string{"email": auth.Email})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.PublicKey == nil {
		log.Println("Auth failed")
		http.Error(w, "Auth failed", http.StatusForbidden)
		return
	}

	if err := user.PublicKey.Verify(data, auth.Signature); err != nil {
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

	postRequest.Owner = user.ID
	if err := db.AddPost(postRequest); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK")) // nolint
}

func (h *Handler) handleAPIRegister(w http.ResponseWriter, r *http.Request) {
	db, err := database.GetDB()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close(context.TODO())

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

	if _, err := ssh.ParsePublicKey([]byte(request.PublicKey)); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.RegisterUser(request); err != nil {
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

func (h *Handler) handleAPIPosts(w http.ResponseWriter, r *http.Request) {
	db, err := database.GetDB()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close(context.TODO())

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := r.Form.Get("page")
	if page == "" {
		page = "1"
	}

	pageNum, err := strconv.ParseInt(page, 10, 32)
	if err != nil {
		pageNum = 1
	}
	d, err := db.GetPosts(int(pageNum), PostsPerPage)
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

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	indexTemplate, err := os.ReadFile(fmt.Sprintf("%s/templates/index.gohtml", PublicPath()))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl, err := template.New("index").Parse(string(indexTemplate))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	siteConfig := config.GetSettings().Site

	pageData := struct {
		SiteConfig config.SiteSettings
	}{
		SiteConfig: siteConfig,
	}

	output := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(output, pageData); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(output.Bytes()) // nolint
}
