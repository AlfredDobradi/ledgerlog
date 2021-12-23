package internal

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRebuildSnapshot(t *testing.T) {
	storeID := uuid.New()
	tests := []struct {
		actions     []*action
		expected    int
		expectedErr error
	}{
		{
			actions: []*action{
				{
					id:      uuid.New(),
					action:  ActionCreate,
					storeid: storeID,
					content: map[string]interface{}{
						"test": 1,
					},
				},
			},
			expected:    1,
			expectedErr: nil,
		},
		{
			actions: []*action{
				{
					id:      uuid.New(),
					action:  ActionCreate,
					storeid: storeID,
					content: map[string]interface{}{
						"test": 1,
					},
				},
				{
					id:      uuid.New(),
					action:  ActionUpdate,
					storeid: storeID,
					content: map[string]interface{}{
						"test": 100,
					},
				},
			},
			expected:    100,
			expectedErr: nil,
		},
		{
			actions: []*action{
				{
					id:      uuid.New(),
					action:  ActionCreate,
					storeid: storeID,
					content: map[string]interface{}{
						"test": 1,
					},
				},
				{
					id:      uuid.New(),
					action:  ActionDelete,
					storeid: storeID,
					content: nil,
				},
			},
			expected:    0,
			expectedErr: ErrKeyNotFound{storeID},
		},
		{
			actions: []*action{
				{
					id:      uuid.New(),
					action:  ActionCreate,
					storeid: storeID,
					content: map[string]interface{}{
						"test": 1,
					},
				},
				{
					id:      uuid.New(),
					action:  ActionDelete,
					storeid: storeID,
					content: nil,
				},
				{
					id:      uuid.New(),
					action:  ActionUpdate,
					storeid: storeID,
					content: map[string]interface{}{
						"test": 100,
					},
				},
			},
			expected:    0,
			expectedErr: ErrKeyNotFound{storeID},
		},
	}

	for i, tt := range tests {
		tf := func(t *testing.T) {
			store := &Store{
				mx:       &sync.Mutex{},
				journal:  make([]*action, len(tt.actions)),
				snapshot: make([]*data, 0),
			}
			var prevAction *action
			for ind, act := range tt.actions {
				if prevAction != nil {
					act.prev = prevAction
				}
				act.timestamp = time.Now()
				store.journal[ind] = act
			}

			err := store.RebuildSnapshot()
			if err != nil {
				t.Fatalf("Fail: Expected no error, got %v", err)
			}

			item, err := store.Get(storeID)
			if tt.expectedErr != nil && err == nil {
				t.Fatalf("Fail: Expected error %v, got none", tt.expectedErr)
			} else if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Fatalf("Fail: Expected error %v, got %v", tt.expectedErr, err)
			} else if tt.expectedErr == nil && err != nil {
				t.Fatalf("Fail: Expected no errors, got %v", err)
			}

			if tt.expectedErr == nil {
				if expected, actual := tt.expected, item.content["test"].(int); expected != actual {
					t.Fatalf("Fail: expected %d, got %d", expected, actual)
				}
			}

			t.Log("Pass")
		}
		t.Run(fmt.Sprintf("test_%d", i+1), tf)
	}
}
