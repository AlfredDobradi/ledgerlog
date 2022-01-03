package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/database"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type PostUpdate struct {
	Posts    []models.PostDisplay
	TotalNum int
}

func (h *Handler) handlePostsSocket(w http.ResponseWriter, r *http.Request) {
	connID := uuid.New()
	connTime := time.Now()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Printf("[%s] New connection established", connID.String())

	firstUpdate, err := h.checkPosts(connID, 30, time.Unix(0, 0))
	if err != nil {
		log.Println(err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())) // nolint
		return
	}
	lastUpdate := time.Unix(0, 0)
	if len(firstUpdate.Posts) > 0 {
		lastUpdate = firstUpdate.Posts[0].Timestamp
	}
	if err := conn.WriteJSON(firstUpdate); err != nil {
		log.Println(err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())) // nolint
		return
	}

	ticker := time.NewTicker(5 * time.Second)

	stop := make(chan struct{})
	conn.SetCloseHandler(func(code int, text string) error {
		dur := time.Since(connTime)
		stop <- struct{}{}
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed")) // nolint
		log.Printf("[%s] Closed connection after %s", connID, dur)
		return nil
	})

	go func() {
		for {
			select {
			case <-ticker.C:
				update, err := h.checkPosts(connID, 30, lastUpdate)
				if err != nil {
					log.Println(err)
					conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())) // nolint
					return
				}
				if len(update.Posts) > 0 {
					if err := conn.WriteJSON(update); err != nil {
						log.Println(err)
						conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())) // nolint
						return
					}
					lastUpdate = update.Posts[0].Timestamp
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil && !errors.Is(err, &websocket.CloseError{}) {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error())) // nolint
			return
		}
	}
}

func (h *Handler) checkPosts(connID uuid.UUID, maxnum int, since time.Time) (PostUpdate, error) {
	// log.Printf("[%s] Checking for posts since %s...", connID.String(), since.Format(time.RFC3339))
	conn, err := database.GetConnection(context.TODO())
	defer func() {
		if err := conn.Close(context.TODO()); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()
	if err != nil {
		return PostUpdate{}, err
	}

	posts, allPosts, err := conn.GetPostsSince(maxnum, since)
	if err != nil {
		return PostUpdate{}, err
	}

	return PostUpdate{
		Posts:    posts,
		TotalNum: allPosts,
	}, err
}
