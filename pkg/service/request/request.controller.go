package request

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	m "gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/groups/v2"
	message_service "gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service      Service
	userService  userdetail.Service
	groupService groups.Service
}

func NewHandler(r *mux.Router) {
	timeoutContext := time.Duration(2) * time.Second
	repo := NewRepoImpl(database.DB)
	messRepo := message_service.NewRepoImpl(database.DB)
	messService := message_service.NewServiceImpl(messRepo, timeoutContext)
	groupRepo := groups.NewRepoImpl(database.DB)
	userRepo := userdetail.NewRepoImpl(database.DB)
	userService := userdetail.NewServiceImpl(userRepo, timeoutContext)
	groupService := groups.NewServiceImpl(groupRepo, userService, timeoutContext, messService)
	requestService := NewServiceImpl(repo, userService, groupService, timeoutContext)
	handler := &Handler{service: requestService, userService: userService, groupService: groupService}

	r.HandleFunc("/api/v1/requests", m.Auth(handler.GetListRequest)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/requests/{idRequest}", m.Auth(handler.GetOneRequest)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/requests/request", m.Auth(handler.CreateRequest)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/requests/request/{idGroup}", m.Auth(handler.GetListRequestGroup)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/requests/request/approve/{idRequest}", m.Auth(handler.ApproveRequest)).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc("/api/v1/requests/reject/{idGroup}", m.Auth(handler.GetListRejectRequest)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/requests/request/reject/{idRequest}", m.Auth(handler.RejectRequest)).Methods(http.MethodPatch, http.MethodOptions)
}

// GetListRequest godoc
// @Summary Get all request
// @Description Get all requests
// @Tags request
// @Accept  json
// @Produce  json
// @Success 200 {array} Dto
// @Router /api/v1/requests [get]
func (h *Handler) GetListRequest(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	listRequest, err := h.service.GetAllRequest(ctx)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}
	w.Write(utils.ResponseWithByte(listRequest))
}

// CreateRequestApi godoc
// @Summary Create a new request
// @Description create a new request
// @Tags request
// @Accept  json
// @Produce  json
// @Param requestPayLoad body Payload true "Create request"
// @Success 200 {object} Dto
// @Router /api/v1/requests/request [post]
func (h *Handler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	//owner := ctx.Value(m.UserKey).(string)
	//payload.CreateBy = owner
	group, _ := h.groupService.GetGroupById(ctx, payload.IdGroup)
	if group.Type == groups.ONE {
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	if group.Private == false {
		fmt.Println("group public")
		User := []string{payload.IdInvite}
		_ = h.groupService.AddUserInGroupService(ctx, User, payload.IdGroup)
		dto := Dto{
			IdInvite: userdetail.GetUserById(payload.IdInvite),
			Status:   APPROVE,
			CreateBy: userdetail.GetUserById(payload.CreateBy),
		}
		w.Write(utils.ResponseWithByte(dto))
		return
	}

	check := h.groupService.CheckUserAndGroupExits(ctx, payload.IdGroup, payload.CreateBy)
	checkRole, _ := h.userService.GetUserDetailByIDService(ctx, payload.CreateBy)
	if !check {
		utils.ResErr(w, http.StatusForbidden)
		return
	} else {
		if checkRole.Role == userdetail.PATIENT {
			utils.ResErr(w, http.StatusForbidden)
			return
		}
		checkRequestExits, _ := h.service.CheckExitsRequest(ctx, payload.IdInvite, payload.IdGroup)
		if checkRequestExits != (Dto{}) {
			if checkRequestExits.Status == REJECT {
				updateRequest, err := h.service.UpdateRequest(ctx, int(checkRequestExits.Id), PENDING, payload.CreateBy)
				if err != nil {
					fmt.Printf("UpdateRequest ")
					utils.ResErr(w, http.StatusBadRequest)
					return
				}
				w.Write(utils.ResponseWithByte(updateRequest))
				return
			} else {
				utils.ResErr(w, http.StatusBadRequest)
				return
			}
		}
		dto, err := h.service.CreateRequest(ctx, payload, payload.CreateBy)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusBadRequest)
			return
		}
		w.Write(utils.ResponseWithByte(dto))
	}
}

// GetRequestById godoc
// @Summary Get request by idRequest
// @Description Get the request corresponding to the input idRequest
// @Tags request
// @Accept  json
// @Produce  json
// @Param idRequest path int true "ID of the Request to be find"
// @Success 200 {object} Dto
// @Router /api/v1/requests/{idRequest} [get]
func (h *Handler) GetOneRequest(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idRequest"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	request, err := h.service.GetOneRequest(ctx, id)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	w.Write(utils.ResponseWithByte(request))
}

// GetListRequestGroup godoc
// @Summary Get all request in group
// @Description Get all requests in group
// @Tags request
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the idGroup to be get list"
// @Success 200 {array} Dto
// @Router /api/v1/requests/request/{idGroup} [get]
func (h *Handler) GetListRequestGroup(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	idGroup, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	owner := ctx.Value(m.UserKey).(string)
	check := h.groupService.CheckUserAndGroupExits(ctx, idGroup, owner)
	checkRole, err := h.groupService.CheckRoleOwnerInGroupService(ctx, owner, idGroup)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	if !checkRole || !check {
		utils.ResErr(w, http.StatusForbidden)
		return
	}

	listRequest, err := h.service.GetListRequestInGroup(ctx, idGroup)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	w.Write(utils.ResponseWithByte(listRequest))
}

// ApproveRequestApi godoc
// @Summary Approve request to group
// @Description approve request to group
// @Tags request
// @Accept  json
// @Produce  json
// @Param idRequest path int true "ID of the Request to be Approve"
// @Success 200 {object} Dto
// @Router /api/v1/requests/request/approve/{idRequest} [patch]
func (h *Handler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	idRequest, err := strconv.Atoi(params["idRequest"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	request, err := h.service.GetOneRequest(ctx, idRequest)
	if request == (Dto{}) {
		utils.ResErr(w, http.StatusNotFound)
		return
	}
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	owner := ctx.Value(m.UserKey).(string)
	check := h.groupService.CheckUserAndGroupExits(ctx, request.IdGroup, owner)
	checkRole, err := h.groupService.CheckRoleOwnerInGroupService(ctx, owner, request.IdGroup)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	if !checkRole || !check {
		utils.ResErr(w, http.StatusForbidden)
		return
	} else {
		dto, err := h.service.UpdateRequest(ctx, idRequest, APPROVE, owner)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusBadRequest)
			return
		}
		w.Write(utils.ResponseWithByte(dto))
	}
}

// RejectRequestApi godoc
// @Summary Reject request to group
// @Description Reject request to group
// @Tags request
// @Accept  json
// @Produce  json
// @Param idRequest path int true "ID of the Request to be Reject"
// @Success 200 {object} Dto
// @Router /api/v1/requests/request/reject/{idRequest} [patch]
func (h *Handler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	idRequest, err := strconv.Atoi(params["idRequest"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	request, err := h.service.GetOneRequest(ctx, idRequest)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	if request == (Dto{}) {
		utils.ResErr(w, http.StatusNotFound)
		return
	}
	owner := ctx.Value(m.UserKey).(string)
	check := h.groupService.CheckUserAndGroupExits(ctx, request.IdGroup, owner)
	checkRole, err := h.groupService.CheckRoleOwnerInGroupService(ctx, owner, request.IdGroup)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	if !checkRole || !check {
		utils.ResErr(w, http.StatusForbidden)
		return
	} else {
		dto, err := h.service.UpdateRequest(ctx, idRequest, REJECT, owner)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusBadRequest)
			return
		}
		w.Write(utils.ResponseWithByte(dto))
	}
}

// GetListRejectRequestGroup godoc
// @Summary Get all reject request in group
// @Description Get all reject requests in group
// @Tags request
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the idGroup to be get list reject"
// @Success 200 {array} Dto
// @Router /api/v1/requests/reject/{idGroup} [get]
func (h *Handler) GetListRejectRequest(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	idGroup, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	owner := ctx.Value(m.UserKey).(string)
	check := h.groupService.CheckUserAndGroupExits(ctx, idGroup, owner)
	checkRole, err := h.groupService.CheckRoleOwnerInGroupService(ctx, owner, idGroup)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	if !checkRole || !check {
		utils.ResErr(w, http.StatusForbidden)
		return
	}

	listRequest, err := h.service.GetListRequestReject(ctx, idGroup)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	w.Write(utils.ResponseWithByte(listRequest))
}
