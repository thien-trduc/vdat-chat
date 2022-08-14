package article

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
	r.HandleFunc("/api/v1/article", auth.AuthenMiddleJWT(handler.GetArticleByUserId)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/article", auth.AuthenMiddleJWT(handler.StoreArticle)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/article/{idArticle}", auth.AuthenMiddleJWT(handler.GetArticleById)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/article/{idArticle}", auth.AuthenMiddleJWT(handler.UpdateArticle)).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/v1/article/{idArticle}", auth.AuthenMiddleJWT(handler.DeleteArticle)).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/v1/search-article", auth.AuthenMiddleJWT(handler.GetArticleByTitle)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/article/share/{idArticle}", auth.AuthenMiddleJWT(handler.ShareArticle)).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/v1/all-article", auth.AuthenMiddleJWT(handler.GetAllArticle)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/category/{idCategory}/article", auth.AuthenMiddleJWT(handler.GetArticleByCategory)).Methods(http.MethodGet, http.MethodOptions)
}

func (h *Handler) GetAllArticle(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()

	page, err := strconv.Atoi(r.URL.Query()["page"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	maxItemPage, err := strconv.Atoi(r.URL.Query()["maxItemPage"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	pag := utils.Pagination{
		Page:        int32(page),
		MaxItemPage: int32(maxItemPage),
	}

	list, err := h.service.Fetch(ctx, pag)
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

func (h *Handler) GetArticleByUserId(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()

	page, err := strconv.Atoi(r.URL.Query()["page"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	maxItemPage, err := strconv.Atoi(r.URL.Query()["maxItemPage"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	pag := utils.Pagination{
		Page:        int32(page),
		MaxItemPage: int32(maxItemPage),
	}

	user := auth.JWTparseOwner(r.Header.Get("Authorization"))
	list, err := h.service.GetByUserId(ctx, user, pag)
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

func (h *Handler) GetArticleById(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idArticle"])
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

func (h *Handler) GetArticleByTitle(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := r.URL.Query()["title"]
	if params == nil {
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(r.URL.Query()["page"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	maxItemPage, err := strconv.Atoi(r.URL.Query()["maxItemPage"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	pag := utils.Pagination{
		Page:        int32(page),
		MaxItemPage: int32(maxItemPage),
	}

	list, err := h.service.GetByTitle(ctx, params[0], pag)
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

func (h *Handler) GetArticleByCategory(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idCategory"])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(r.URL.Query()["page"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	maxItemPage, err := strconv.Atoi(r.URL.Query()["maxItemPage"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	pag := utils.Pagination{
		Page:        int32(page),
		MaxItemPage: int32(maxItemPage),
	}
	dto, err := h.service.GetByCategory(ctx, int64(id), pag)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	utils.ResponseOk(w, dto)
}

func (h *Handler) StoreArticle(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	user := auth.JWTparseOwner(r.Header.Get("Authorization"))

	id, err := strconv.Atoi(params["idArticle"])
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

func (h *Handler) ShareArticle(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idArticle"])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	dto, err := h.service.UpdateWithNumShare(ctx, int64(id))
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	utils.ResponseOk(w, dto)
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["idArticle"])
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
