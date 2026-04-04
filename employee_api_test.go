package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type stubUserRepo struct{}

func (stubUserRepo) Create(ctx context.Context, name string) (*User, error) {
	return &User{ID: 1, Name: name, CreatedAt: time.Unix(1, 0).UTC()}, nil
}

func (stubUserRepo) List(ctx context.Context) ([]User, error) {
	return nil, nil
}

type stubEmployeeRepo struct {
	nextID int64
	items  []Employee
}

func (s *stubEmployeeRepo) CreateEmployee(ctx context.Context, name, email string) (*Employee, error) {
	s.nextID++
	e := Employee{ID: s.nextID, Name: name, Email: email, CreatedAt: time.Unix(1, 0).UTC()}
	s.items = append(s.items, e)
	return &e, nil
}

func (s *stubEmployeeRepo) ListEmployees(ctx context.Context) ([]Employee, error) {
	return s.items, nil
}

func TestEmployeeAPI_CreateAndList(t *testing.T) {
	er := &stubEmployeeRepo{}
	h := &Handler{userRepo: stubUserRepo{}, employeeRepo: er}
	r := chi.NewRouter()
	r.Post("/employees", h.CreateEmployee)
	r.Get("/employees", h.ListEmployees)

	body := `{"name":"Ada","email":"ada@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(body))
	req = req.WithContext(context.Background())
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var created Employee
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.Name != "Ada" || created.Email != "ada@example.com" {
		t.Fatalf("unexpected created: %+v", created)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/employees", nil)
	req2 = req2.WithContext(context.Background())
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("GET status = %d", rec2.Code)
	}
	var listed []Employee
	if err := json.Unmarshal(rec2.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0].Name != "Ada" {
		t.Fatalf("unexpected list: %+v", listed)
	}
}

type duplicateEmailEmployeeRepo struct{}

func (duplicateEmailEmployeeRepo) CreateEmployee(ctx context.Context, name, email string) (*Employee, error) {
	return nil, fmt.Errorf("insert employee: %w", &pgconn.PgError{Code: "23505"})
}

func (duplicateEmailEmployeeRepo) ListEmployees(ctx context.Context) ([]Employee, error) {
	return nil, nil
}

func TestEmployeeAPI_DuplicateEmail(t *testing.T) {
	h := &Handler{userRepo: stubUserRepo{}, employeeRepo: duplicateEmailEmployeeRepo{}}
	r := chi.NewRouter()
	r.Post("/employees", h.CreateEmployee)

	req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(`{"name":"Bob","email":"bob@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("POST status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestEmployeeAPI_Validation(t *testing.T) {
	h := &Handler{userRepo: stubUserRepo{}, employeeRepo: &stubEmployeeRepo{}}
	r := chi.NewRouter()
	r.Post("/employees", h.CreateEmployee)

	for _, tc := range []struct {
		body string
		want int
	}{
		{`{}`, http.StatusBadRequest},
		{`{"name":"x"}`, http.StatusBadRequest},
		{`{"email":"a@b.com"}`, http.StatusBadRequest},
	} {
		req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(tc.body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != tc.want {
			t.Fatalf("body %q: status = %d want %d", tc.body, rec.Code, tc.want)
		}
	}
}
