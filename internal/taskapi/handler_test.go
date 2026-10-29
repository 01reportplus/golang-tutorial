package taskapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func request(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestCreateListUpdateDelete(t *testing.T) {
	h := NewHandler(NewStore())
	w := request(h, "GET", "/api/tasks", "")
	if w.Code != 200 || strings.TrimSpace(w.Body.String()) != "[]" {
		t.Fatalf("empty list: %d %s", w.Code, w.Body)
	}
	w = request(h, "POST", "/api/tasks", `{"title":" 学习 Go "}`)
	if w.Code != 201 || w.Header().Get("Location") != "/api/tasks/1" {
		t.Fatalf("create: %d %s", w.Code, w.Body)
	}
	var task Task
	if err := json.Unmarshal(w.Body.Bytes(), &task); err != nil {
		t.Fatal(err)
	}
	if task.ID != 1 || task.Title != "学习 Go" || task.Done || task.CreatedAt.IsZero() {
		t.Fatalf("unexpected task: %+v", task)
	}
	w = request(h, "PATCH", "/api/tasks/1", `{"done":true}`)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"done":true`) {
		t.Fatalf("patch: %d %s", w.Code, w.Body)
	}
	w = request(h, "PATCH", "/api/tasks/1", `{"done":false}`)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"done":false`) {
		t.Fatalf("explicit false: %d %s", w.Code, w.Body)
	}
	w = request(h, "GET", "/api/tasks/1", "")
	if w.Code != 200 {
		t.Fatalf("get: %d %s", w.Code, w.Body)
	}
	w = request(h, "DELETE", "/api/tasks/1", "")
	if w.Code != 204 || w.Body.Len() != 0 {
		t.Fatalf("delete: %d %s", w.Code, w.Body)
	}
	if w = request(h, "GET", "/api/tasks/1", ""); w.Code != 404 {
		t.Fatalf("deleted task: %d", w.Code)
	}
}

func TestInvalidRequests(t *testing.T) {
	h := NewHandler(NewStore())
	for _, tt := range []struct {
		name, method, path, body string
		status                   int
	}{
		{"empty title", "POST", "/api/tasks", `{"title":" "}`, 400},
		{"long title", "POST", "/api/tasks", `{"title":"` + strings.Repeat("中", 201) + `"}`, 400},
		{"unknown field", "POST", "/api/tasks", `{"title":"go","other":1}`, 400},
		{"trailing value", "POST", "/api/tasks", `{"title":"go"} {}`, 400},
		{"malformed", "POST", "/api/tasks", `{`, 400},
		{"null object", "POST", "/api/tasks", `null`, 400},
		{"empty body", "POST", "/api/tasks", ``, 400},
		{"oversize", "POST", "/api/tasks", `{"title":"` + strings.Repeat("x", maxBody) + `"}`, 413},
		{"bad id", "GET", "/api/tasks/abc", "", 400},
		{"negative id", "GET", "/api/tasks/-1", "", 400},
		{"missing", "GET", "/api/tasks/999", "", 404},
		{"missing patch", "PATCH", "/api/tasks/999", `{"done":true}`, 404},
		{"required done", "PATCH", "/api/tasks/1", `{}`, 400},
		{"null done", "PATCH", "/api/tasks/1", `{"done":null}`, 400},
		{"bad done type", "PATCH", "/api/tasks/1", `{"done":1}`, 400},
		{"wrong method", "PUT", "/api/tasks", `{}`, 405},
		{"missing delete", "DELETE", "/api/tasks/999", "", 404},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := request(h, tt.method, tt.path, tt.body)
			if w.Code != tt.status {
				t.Fatalf("got %d %s, want %d", w.Code, w.Body, tt.status)
			}
		})
	}
}

func TestMediaTypeAndHealth(t *testing.T) {
	h := NewHandler(NewStore())
	req := httptest.NewRequest("POST", "/api/tasks", strings.NewReader(`{"title":"go"}`))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 415 {
		t.Fatalf("missing content type: %d", w.Code)
	}
	req = httptest.NewRequest("POST", "/api/tasks", strings.NewReader(`{"title":"go"}`))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 201 || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("JSON charset: %d %s", w.Code, w.Body)
	}
	if got := request(h, "GET", "/healthz", ""); got.Code != 200 {
		t.Fatalf("health: %d", got.Code)
	}
}

func TestConcurrentStore(t *testing.T) {
	store := NewStore()
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			task, err := store.Create(fmt.Sprintf("task %d", i))
			if err != nil {
				t.Error(err)
				return
			}
			if _, err := store.SetDone(task.ID, true); err != nil {
				t.Error(err)
			}
			_ = store.List()
		})
	}
	wg.Wait()
	tasks := store.List()
	if len(tasks) != 100 {
		t.Fatalf("got %d tasks", len(tasks))
	}
	for i, task := range tasks {
		if task.ID != int64(i+1) || !task.Done {
			t.Fatalf("bad task: %+v", task)
		}
	}
	tasks[0].Title = "changed outside"
	stored, _ := store.Get(tasks[0].ID)
	if stored.Title == tasks[0].Title {
		t.Fatal("list leaked mutable state")
	}
}

func TestHTTPRoundTrip(t *testing.T) {
	server := httptest.NewServer(NewHandler(NewStore()))
	defer server.Close()
	resp, err := server.Client().Post(server.URL+"/api/tasks", "application/json", strings.NewReader(`{"title":"round trip"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status=%d", resp.StatusCode)
	}
}
