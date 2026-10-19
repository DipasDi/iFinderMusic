package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestNewUserStore_CreatesFileAndDir verifies that NewUserStore creates required directories and initializes correctly.
func TestNewUserStore_CreatesFileAndDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "userstore_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	filePath := filepath.Join(tempDir, "subfolder", "users.json")
	store, err := NewUserStore(filePath)
	if err != nil {
		t.Fatalf("NewUserStore failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected non-nil store")
	}

	users := store.GetAllUsers()
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

// TestUserStore_TrackUser_NewAndExisting tests registering new users and updating existing user records.
func TestUserStore_TrackUser_NewAndExisting(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "userstore_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	filePath := filepath.Join(tempDir, "users.json")
	store, err := NewUserStore(filePath)
	if err != nil {
		t.Fatalf("NewUserStore failed: %v", err)
	}

	err = store.TrackUser(1001, "first_user", "John", "Doe")
	if err != nil {
		t.Fatalf("TrackUser failed: %v", err)
	}

	u, found := store.GetUser(1001)
	if !found {
		t.Fatal("expected user 1001 to be found")
	}
	if u.Username != "first_user" || u.FirstName != "John" || u.LastName != "Doe" {
		t.Errorf("unexpected user data: %+v", u)
	}

	// Update username
	err = store.TrackUser(1001, "new_username", "John", "Smith")
	if err != nil {
		t.Fatalf("TrackUser update failed: %v", err)
	}

	uUpdated, found := store.GetUser(1001)
	if !found {
		t.Fatal("expected user 1001 to still exist")
	}
	if uUpdated.Username != "new_username" || uUpdated.LastName != "Smith" {
		t.Errorf("expected updated username and lastname, got %+v", uUpdated)
	}
	if !uUpdated.FirstSeen.Equal(u.FirstSeen) {
		t.Errorf("expected FirstSeen to remain unchanged, got %v vs %v", uUpdated.FirstSeen, u.FirstSeen)
	}
}

// TestUserStore_PersistenceRoundTrip tests reloading users from an existing JSON file after restart.
func TestUserStore_PersistenceRoundTrip(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "userstore_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	filePath := filepath.Join(tempDir, "users.json")
	store1, err := NewUserStore(filePath)
	if err != nil {
		t.Fatalf("store1 init failed: %v", err)
	}

	_ = store1.TrackUser(1, "alice", "Alice", "")
	_ = store1.TrackUser(2, "bob", "Bob", "")

	// Open second store on same file to simulate application restart
	store2, err := NewUserStore(filePath)
	if err != nil {
		t.Fatalf("store2 init failed: %v", err)
	}

	all := store2.GetAllUsers()
	if len(all) != 2 {
		t.Fatalf("expected 2 users loaded from file, got %d", len(all))
	}

	uAlice, okAlice := store2.GetUser(1)
	if !okAlice || uAlice.Username != "alice" {
		t.Errorf("alice mismatch in store2: %+v", uAlice)
	}

	uBob, okBob := store2.GetUser(2)
	if !okBob || uBob.Username != "bob" {
		t.Errorf("bob mismatch in store2: %+v", uBob)
	}
}

// TestUserStore_ConcurrentAccess verifies safe concurrent reads and writes across multiple goroutines.
func TestUserStore_ConcurrentAccess(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "userstore_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	filePath := filepath.Join(tempDir, "users.json")
	store, err := NewUserStore(filePath)
	if err != nil {
		t.Fatalf("NewUserStore failed: %v", err)
	}

	var wg sync.WaitGroup
	workers := 20
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			uid := int64(id + 1)
			_ = store.TrackUser(uid, fmt.Sprintf("user_%d", uid), "Name", "Surname")
			_, _ = store.GetUser(uid)
			_ = store.GetAllUsers()
		}(i)
	}
	wg.Wait()

	all := store.GetAllUsers()
	if len(all) != workers {
		t.Errorf("expected %d users recorded concurrently, got %d", workers, len(all))
	}
}

// TestUserStore_NilStore verifies that methods on a nil store instance execute safely without panicking.
func TestUserStore_NilStore(t *testing.T) {
	var store *UserStore

	if err := store.TrackUser(1, "test", "fn", "ln"); err != nil {
		t.Errorf("expected nil error on nil store, got %v", err)
	}
	if _, ok := store.GetUser(1); ok {
		t.Errorf("expected ok=false on nil store")
	}
	if users := store.GetAllUsers(); users != nil {
		t.Errorf("expected nil slice from GetAllUsers on nil store, got %+v", users)
	}
}
