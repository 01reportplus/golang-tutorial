// Package taskapi implements a small in-memory teaching API, without persistence or authentication.
package taskapi

import (
	"cmp"
	"errors"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var (
	ErrNotFound = errors.New("task not found")
	ErrTitle    = errors.New("title must contain 1 to 200 Unicode code points after trimming")
)

// Task is returned by value; it contains no mutable reference fields.
type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// Store serializes all access to tasks and nextID. A Store must not be copied.
type Store struct {
	mu     sync.Mutex
	nextID int64
	tasks  map[int64]Task
}

func NewStore() *Store { return &Store{tasks: make(map[int64]Task)} }

func (s *Store) Create(title string) (Task, error) {
	title = strings.TrimSpace(title)
	if !utf8.ValidString(title) || title == "" || utf8.RuneCountInString(title) > 200 {
		return Task{}, ErrTitle
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	task := Task{ID: s.nextID, Title: title, CreatedAt: time.Now().UTC()}
	s.tasks[task.ID] = task
	return task, nil
}

func (s *Store) List() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		out = append(out, task)
	}
	slices.SortFunc(out, func(a, b Task) int { return cmp.Compare(a.ID, b.ID) })
	return out
}

func (s *Store) Get(id int64) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return task, nil
}

func (s *Store) SetDone(id int64, done bool) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	task.Done = done
	s.tasks[id] = task
	return task, nil
}

func (s *Store) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
