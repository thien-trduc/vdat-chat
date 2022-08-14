package category

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service Service
}

func NewHandler(r *mux.Router) {
	timeoutContext := time.Duration(2) * time.Second
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	handler := &Handler{service: service}
	r.HandleFunc("/api/v1/category", auth.AuthenMiddleJWT(handler.GetCategoryByUserId)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/category", auth.AuthenMiddleJWT(handler.StoreCategory)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/category/{idCategory}", auth.AuthenMiddleJWT(handler.DeleteCategory)).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/v1/category/{idCategory}", auth.AuthenMiddleJWT(handler.StoreChildCategory)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/category/{idCategory}", auth.AuthenMiddleJWT(handler.UpdateCategory)).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/v1/category/{idCategory}", auth.AuthenMiddleJWT(handler.GetCategoryById)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/category/search-name", auth.AuthenMiddleJWT(handler.GetCategoryByName)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/category/search-created-by", auth.AuthenMiddleJWT(handler.GetCategoryByCreatedBy)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/category/search-updated-by", auth.AuthenMiddleJWT(handler.GetCategoryByUpdateBy)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/category/parent/{idParent}", auth.AuthenMiddleJWT(handler.GetCategoryByParentId)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/all-category", auth.AuthenMiddleJWT(handler.GetAllCategory)).Methods(http.MethodGet, http.MethodOptions)
}

func (h *Handler) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()

	list, err := h.service.Fetch(ctx)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if len(list) > 0 {
		utils.ResponseOk(w, list)
	} else {
		results := make([]Dto, 0)
		utils.ResponseOk(w, results)
	}
}

func (h *Handler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idCategory"])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	dto, err := h.service.GetByID(ctx, int64(id))
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	utils.ResponseOk(w, dto)

}

func (h *Handler) GetCategoryByUserId(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()

	user := auth.JWTparseOwner(r.Header.Get("Authorization"))
	list, err := h.service.GetByCreatedBy(ctx, user)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if len(list) > 0 {
		utils.ResponseOk(w, list)
	} else {
		results := make([]Dto, 0)
		utils.ResponseOk(w, results)
	}

}

func (h *Handler) GetCategoryByName(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := r.URL.Query()["name"]
	if params == nil {
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	list, err := h.service.GetByName(ctx, params[0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if len(list) > 0 {
		utils.ResponseOk(w, list)
	} else {
		results := make([]Dto, 0)
		utils.ResponseOk(w, results)
	}
}

func (h *Handler) GetCategoryByCreatedBy(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()

	params := r.URL.Query()["createdBy"]
	if params == nil {
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	list, err := h.service.GetByCreatedBy(ctx, params[0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if len(list) > 0 {
		utils.ResponseOk(w, list)
	} else {
		results := make([]Dto, 0)
		utils.ResponseOk(w, results)
	}

}

func (h *Handler) GetCategoryByUpdateBy(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()

	params := r.URL.Query()["updateBy"]
	if params == nil {
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	list, err := h.service.GetByUpdateBy(ctx, params[0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if len(list) > 0 {
		utils.ResponseOk(w, list)
	} else {
		results := make([]Dto, 0)
		utils.ResponseOk(w, results)
	}
}

func (h *Handler) GetCategoryByParentId(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	params := mux.Vars(r)
	ctx := r.Context()

	id, err := strconv.Atoi(params["idParent"])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	list, err := h.service.GetByParentId(ctx, int64(id))
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if len(list) > 0 {
		utils.ResponseOk(w, list)
	} else {
		results := make([]Dto, 0)
		utils.ResponseOk(w, results)
	}
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	user := auth.JWTparseOwner(r.Header.Get("Authorization"))

	id, err := strconv.Atoi(params["idCategory"])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}

	var p *Payload
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	dto, err := h.service.Update(ctx, p, int64(id), user)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	utils.ResponseOk(w, dto)
}

func (h *Handler) StoreCategory(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()

	user := auth.JWTparseOwner(r.Header.Get("Authorization"))
	var p *Payload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}

	dto, err := h.service.Store(ctx, p, user)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	utils.ResponseOk(w, dto)
}

func (h *Handler) StoreChildCategory(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idCategory"])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}

	user := auth.JWTparseOwner(r.Header.Get("Authorization"))
	var p *Payload
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	dto, err := h.service.StoreChild(ctx, p, int64(id), user)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	utils.ResponseOk(w, dto)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["idCategory"])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	err = h.service.Delete(ctx, int64(id))
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	result := utils.ResponseBool{Result: true}
	utils.ResponseOk(w, result)
}
