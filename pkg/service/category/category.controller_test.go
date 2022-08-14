package category

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

func TestHandler_GetAllCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/all-category", h.GetAllCategory).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/all-category", nil)
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
	var expected = `[{"id":1,"name":"test","parentId":2,"num":0,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:47:24.372393Z","updateAt":"2020-12-24T09:47:24.372393Z"},{"id":2,"name":"test 2","parentId":-1,"num":1,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T10:03:27.480813Z","updateAt":"2020-12-24T10:03:27.480813Z"}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetCategoryById(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/{idCategory}", h.GetCategoryById).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/category/1", nil)
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
	var expected = `{"id":1,"name":"test","parentId":2,"num":0,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:47:24.372393Z","updateAt":"2020-12-24T09:47:24.372393Z"}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetCategoryByName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/search-name", h.GetCategoryByName).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/category/search-name?name=te", nil)
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
	var expected = `[{"id":2,"name":"test 2","parentId":-1,"num":1,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T10:03:27.480813Z","updateAt":"2020-12-24T10:03:27.480813Z"},{"id":1,"name":"test","parentId":2,"num":0,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:47:24.372393Z","updateAt":"2020-12-24T09:47:24.372393Z"}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetCategoryByCreatedBy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/search-created-by", h.GetCategoryByCreatedBy).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/category/search-created-by?createdBy=0106334e-575e-4fca-8b90-cf5726a36c14", nil)
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
	var expected = `[{"id":1,"name":"test","parentId":2,"num":0,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:47:24.372393Z","updateAt":"2020-12-24T09:47:24.372393Z"},{"id":2,"name":"test 2","parentId":-1,"num":1,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T10:03:27.480813Z","updateAt":"2020-12-24T10:03:27.480813Z"}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetCategoryByUpdateBy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/search-update-by", h.GetCategoryByUpdateBy).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/category/search-update-by?updateBy=0106334e-575e-4fca-8b90-cf5726a36c14", nil)
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
	var expected = `[{"id":1,"name":"test","parentId":2,"num":0,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:47:24.372393Z","updateAt":"2020-12-24T09:47:24.372393Z"},{"id":2,"name":"test 2","parentId":-1,"num":1,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T10:03:27.480813Z","updateAt":"2020-12-24T10:03:27.480813Z"}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetCategoryByParentId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/parent/{idParent}", h.GetCategoryByParentId).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/category/parent/2", nil)
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
	var expected = `[{"id":1,"name":"test","parentId":2,"num":0,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:47:24.372393Z","updateAt":"2020-12-24T09:47:24.372393Z"}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_GetCategoryByUserId(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category", h.GetCategoryByUserId).Methods(http.MethodGet)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/category", nil)
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
	var expected = `[{"id":1,"name":"test","parentId":2,"num":0,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T09:47:24.372393Z","updateAt":"2020-12-24T09:47:24.372393Z"},{"id":2,"name":"test 2","parentId":-1,"num":1,"version":1,"createdBy":"0106334e-575e-4fca-8b90-cf5726a36c14","updateBy":"0106334e-575e-4fca-8b90-cf5726a36c14","createdAt":"2020-12-24T10:03:27.480813Z","updateAt":"2020-12-24T10:03:27.480813Z"}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
func TestHandler_UpdateCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/{idCategory}", h.UpdateCategory).Methods(http.MethodPut)

	p := Payload{
		Name:     "test",
		ParentID: 2,
	}

	js, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/api/v1/category/1", bytes.NewBuffer(js))
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
func TestHandler_StoreCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category", h.StoreCategory).Methods(http.MethodPost)

	p := Payload{
		Name: "test",
	}

	js, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/category", bytes.NewBuffer(js))
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
func TestHandler_StoreChildCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/{idCategory}", h.StoreChildCategory).Methods(http.MethodPost)

	p := Payload{
		Name: "test",
	}
	js, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/category/1", bytes.NewBuffer(js))
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
	p.ParentID = 1
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
func TestHandler_DeleteCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	database.Connect()
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	h := &Handler{service: service}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/category/{idCategory}", h.DeleteCategory).Methods(http.MethodDelete)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/api/v1/category/1", nil)
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
