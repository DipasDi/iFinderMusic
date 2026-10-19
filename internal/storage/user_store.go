package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"random/internal/model"
)

// UserStore manages persisting registered Telegram bot users to a JSON file.
type UserStore struct {
	mu       sync.RWMutex
	filePath string
	users    map[int64]model.User
}

// NewUserStore creates a new UserStore and loads existing users from the JSON file if available.
func NewUserStore(filePath string) (*UserStore, error) {
	if filePath == "" {
		filePath = "users.json"
	}

	if fi, err := os.Stat(filePath); err == nil && fi.IsDir() {
		filePath = filepath.Join(filePath, "users.json")
	}

	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for users file: %w", err)
		}
	}

	store := &UserStore{
		filePath: filePath,
		users:    make(map[int64]model.User),
	}

	if err := store.load(); err != nil {
		log.Printf("warning: could not load users from %s: %v", filePath, err)
	}

	return store, nil
}

// TrackUser records or updates a Telegram user in memory and persists the state to the JSON file.
func (s *UserStore) TrackUser(id int64, username, firstName, lastName string) error {
	if s == nil || id == 0 {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC().Truncate(time.Second)
	existing, exists := s.users[id]

	if !exists {
		s.users[id] = model.User{
			ID:        id,
			Username:  username,
			FirstName: firstName,
			LastName:  lastName,
			FirstSeen: now,
			LastSeen:  now,
		}
	} else {
		hasInfoChanged := existing.Username != username || existing.FirstName != firstName || existing.LastName != lastName
		timePassed := now.Sub(existing.LastSeen) >= time.Minute

		if !hasInfoChanged && !timePassed {
			return nil
		}

		existing.Username = username
		existing.FirstName = firstName
		existing.LastName = lastName
		existing.LastSeen = now
		s.users[id] = existing
	}

	return s.saveLocked()
}

// GetUser retrieves a user by their Telegram ID.
func (s *UserStore) GetUser(id int64) (model.User, bool) {
	if s == nil {
		return model.User{}, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[id]
	return u, ok
}

// GetAllUsers returns a slice of all tracked users sorted by FirstSeen timestamp.
func (s *UserStore) GetAllUsers() []model.User {
	if s == nil {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]model.User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].FirstSeen.Before(list[j].FirstSeen)
	})

	return list
}

// load parses the JSON file and populates the in-memory users map.
func (s *UserStore) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var list []model.User
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}

	for _, u := range list {
		if u.ID != 0 {
			s.users[u.ID] = u
		}
	}

	return nil
}

// saveLocked serializes the user map to JSON with indentation and writes it to disk.
func (s *UserStore) saveLocked() error {
	list := make([]model.User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].FirstSeen.Before(list[j].FirstSeen)
	})

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0666)
}
