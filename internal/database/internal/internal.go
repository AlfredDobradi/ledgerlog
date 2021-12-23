package internal

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type DataAction int

const (
	ActionCreate DataAction = iota
	ActionDelete
	ActionUpdate
)

type action struct {
	prev      *action
	timestamp time.Time
	id        uuid.UUID
	storeid   uuid.UUID
	action    DataAction
	content   map[string]interface{}
}

type data struct {
	id      uuid.UUID
	content map[string]interface{}
}

type Store struct {
	mx           *sync.Mutex
	lastSnapshot time.Time
	journal      []*action
	snapshot     []*data
}

type ErrKeyNotFound struct {
	key uuid.UUID
}

func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf("key %s not found", e.key)
}

func (s *Store) Get(key uuid.UUID) (*data, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	for _, item := range s.snapshot {
		if item.id == key {
			return item, nil
		}
	}
	return nil, ErrKeyNotFound{key}
}

func (s *Store) RebuildSnapshot() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	for index, task := range s.journal {
		log.Printf("Applying task at index %d", index)
		switch task.action {
		case ActionCreate:
			if uuid, err := s.snapshotCreate(task.storeid, task.content); err != nil {
				return err
			} else if task.storeid != uuid {
				task.storeid = uuid
			}
		case ActionUpdate:
			if _, err := s.snapshotUpdate(task.storeid, task.content); err != nil {
				return err
			}
		case ActionDelete:
			if err := s.snapshotDelete(task.storeid); err != nil {
				return err
			}
		}
	}

	s.lastSnapshot = time.Now()
	return nil
}

func (s *Store) snapshotCreate(id uuid.UUID, content map[string]interface{}) (uuid.UUID, error) {
	if id == uuid.Nil {
		id = uuid.New()
	} else {
		for _, item := range s.snapshot {
			if item.id == id {
				return uuid.Nil, fmt.Errorf("key %s already exists and will not be updated by a create task", id.String())
			}
		}
	}
	log.Printf("\tCREATE: %s", id.String())
	s.snapshot = append(s.snapshot, &data{
		id:      id,
		content: content,
	})
	return id, nil
}

func (s *Store) snapshotUpdate(id uuid.UUID, newContent map[string]interface{}) (uuid.UUID, error) {
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("trying to update record with nil UUID")
	}

	for _, item := range s.snapshot {
		if item.id == id {
			log.Printf("\tUPDATE: %s", id.String())
			for key, value := range newContent {
				item.content[key] = value
			}
			return item.id, nil
		}
	}
	return uuid.Nil, ErrKeyNotFound{key: id}
}

func (s *Store) snapshotDelete(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("trying to delete record with nil UUID")
	}

	var itemIndex = -1
	for i, item := range s.snapshot {
		if item.id == id {
			itemIndex = i
		}
	}
	if itemIndex == -1 {
		return ErrKeyNotFound{key: id}
	}
	log.Printf("\tDELETE: %s", id.String())
	snapshot := append(s.snapshot[:itemIndex], s.snapshot[itemIndex+1:]...)
	s.snapshot = snapshot
	return nil
}
