package article

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	timeoutContext = time.Duration(3) * time.Second
)

func TestHandler_GetAllArticle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/all-article", h.GetAllArticle).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/all-article", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var expected = `[{"id":6,"content":"test","title":"test","thumbnail":"test","version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:56:44.6872Z","updateAt":"2020-12-24T09:56:44.6872Z","slug":"test","idCategory":1},{"id":7,"content":"test 2","title":"test 2","thumbnail":"test 2","version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:56:44.6872Z","updateAt":"2020-12-24T09:56:44.6872Z","slug":"test","idCategory":2}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetArticleById(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/article/{idArticle}", h.GetArticleById).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/article/6", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var expected = `{"id":6,"content":"test","title":"test","thumbnail":"test","version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:56:44.6872Z","updateAt":"2020-12-24T09:56:44.6872Z","slug":"test","idCategory":1}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetArticleByTitle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/search-article", h.GetArticleByTitle).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/search-article?title=te", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var expected = `[{"id":6,"content":"test","title":"test","thumbnail":"test","version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:56:44.6872Z","updateAt":"2020-12-24T09:56:44.6872Z","slug":"test","idCategory":1},{"id":7,"content":"test 2","title":"test 2","thumbnail":"test 2","version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:56:44.6872Z","updateAt":"2020-12-24T09:56:44.6872Z","slug":"test","idCategory":2}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetArticleByUserId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/article", h.GetArticleByUserId).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/article", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Connect())
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var expected = `[{"id":6,"content":"test","title":"test","thumbnail":"test","version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:56:44.6872Z","updateAt":"2020-12-24T09:56:44.6872Z","slug":"test","idCategory":1},{"id":7,"content":"test 2","title":"test 2","thumbnail":"test 2","version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:56:44.6872Z","updateAt":"2020-12-24T09:56:44.6872Z","slug":"test","idCategory":2}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetArticleByCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/{idCategory}/article", h.GetArticleByCategory).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/category/1/article", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var expected = `[{"id":6,"content":"test","title":"test","thumbnail":"test","version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:56:44.6872Z","updateAt":"2020-12-24T09:56:44.6872Z","slug":"test","idCategory":1}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_StoreArticle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/article", h.StoreArticle).Methods(http.MethodPost)

	p := Payload{
		Content:    "test",
		Title:      "test",
		Thumbnail:  "test",
		IdCategory: 1,
	}

	js, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/article", bytes.NewBuffer(js))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Connect())
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expected Payload
	err = json.NewDecoder(res.Body).Decode(&expected)
	if err != nil {
		t.Errorf("handler returned wrong body : %v", res.Body.String())
		return
	}

	if expected != p {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
}
func TestHandler_UpdateArticle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/article/{idArticle}", h.UpdateArticle).Methods(http.MethodPut)

	p := Payload{
		Content:    "test",
		Title:      "test",
		Thumbnail:  "test",
		IdCategory: 1,
	}

	js, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/api/v1/article/6", bytes.NewBuffer(js))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Connect())
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expected Payload
	err = json.NewDecoder(res.Body).Decode(&expected)
	if err != nil {
		t.Errorf("handler returned wrong body : %v", res.Body.String())
		return
	}

	if expected != p {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
}
func TestHandler_DeleteArticle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/article/{idArticle}", h.DeleteArticle).Methods(http.MethodDelete)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/api/v1/article/6", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var expected = `{"result":true}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
