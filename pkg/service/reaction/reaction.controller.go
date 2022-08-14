package reaction

import (
	"encoding/json"
	"fmt"
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
	r.HandleFunc("/api/v1/reaction", auth.AuthenMiddleJWT(handler.StoreReaction)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/reaction/{idArticle}", auth.AuthenMiddleJWT(handler.GetReactionByArticle)).Methods(http.MethodGet, http.MethodOptions)
}

func (h *Handler) GetReactionByArticle(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idArticle"])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	list, err := h.service.GetReactionByArticle(ctx, int64(id))
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

func (h *Handler) StoreReaction(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) UpdateReaction(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	fmt.Println(ctx)
}
