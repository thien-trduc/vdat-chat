package groups

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	m "gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	message_service "gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"log"
	"os"
	"time"

	"net/http"
	"strconv"
)

type Handler struct {
	service     Service
	userService userdetail.Service
}

func NewHandler(r *mux.Router) {
	timeoutContext := time.Duration(2) * time.Second

	messRepo := message_service.NewRepoImpl(database.DB)
	messService := message_service.NewServiceImpl(messRepo, timeoutContext)
	repo := NewRepoImpl(database.DB)
	userRepo := userdetail.NewRepoImpl(database.DB)
	userService := userdetail.NewServiceImpl(userRepo, timeoutContext)
	service := NewServiceImpl(repo, userService, timeoutContext, messService)
	handler := &Handler{service: service, userService: userService}

	minioClient, _ = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})

	minioEndPortStr := os.Getenv("MINIO_END_PORT")
	minioKeyStr := os.Getenv("MINIO_KEY")
	minioSercetStr := os.Getenv("MINIO_ACCESSES")
	if len(minioEndPortStr) > 0 {
		endpoint = minioEndPortStr
	}
	if len(minioKeyStr) > 0 {
		accessKeyID = minioKeyStr
	}
	if len(minioSercetStr) > 0 {
		secretAccessKey = minioSercetStr
	}

	r.HandleFunc("/api/v1/groups", m.Auth(handler.GetListGroupApi)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/name", m.Auth(handler.FindGroupByName)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/groups", m.Auth(handler.CreateGroupApi)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}", m.Auth(handler.UpdateInfoGroupApi)).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}", m.Auth(handler.DeleteGroupApi)).Methods(http.MethodDelete, http.MethodOptions)

	r.HandleFunc("/api/v1/groups/public/{idGroup}/members", m.Auth(handler.AddUserInGroupPublicApi)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}/members", m.Auth(handler.GetListUserOnlineByGroupApi)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}/members", m.Auth(handler.AddUserInGroupApi)).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}/members", m.Auth(handler.UserOutGroupApi)).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}/members/{userId}", m.Auth(handler.DeleteGroupUserApi)).Methods(http.MethodDelete, http.MethodOptions)

}

//API load danh sách groups theo patient hoac theo doctor

// GetListGroupApi godoc
// @Summary Get all groups
// @Description Get all groups
// @Tags groups
// @Accept  json
// @Produce  json
// @Param page query string true "current page"
// @Param keyword query string false "name search by keyword"
// @Param pageSize query string true "max item in page"
// @Success 200 {array} Dto
// @Router /api/v1/groups [get]
func (h *Handler) GetListGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	user := ctx.Value(m.UserKey).(string)
	check, err := h.userService.GetUserDetailByIDService(ctx, user)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}

	if len(r.URL.Query()["page"]) <= 0 || len(r.URL.Query()["pageSize"]) <= 0 {
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(r.URL.Query()["page"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	maxItemPage, err := strconv.Atoi(r.URL.Query()["pageSize"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	pag := utils.Pagination{
		Page:        int32(page),
		MaxItemPage: int32(maxItemPage),
	}

	keyword := r.URL.Query()["keyword"]

	if check.Role == userdetail.PATIENT {
		if len(keyword) > 0 {
			groups, err := h.service.GetGroupByNameForPatient(ctx, user, keyword[0], pag)
			if err != nil {
				sentry.CaptureException(err)
				log.Printf("Exception : %s", err)
				utils.ResErr(w, http.StatusInternalServerError)
				return
			}
			if len(groups) <= 0 {
				utils.ResponseOk(w, groups)
				return
			}
			listGroup, _ := h.service.GetNameGroupForGroup11(ctx, groups, check.ID)
			w.Write(utils.ResponseWithByte(listGroup))
			return
		}
		groups, err := h.service.GetGroupByPatientService(ctx, user, pag)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		if len(groups) <= 0 {
			utils.ResponseOk(w, groups)
			return
		}
		listGroup, _ := h.service.GetNameGroupForGroup11(ctx, groups, check.ID)
		w.Write(utils.ResponseWithByte(listGroup))
	} else {
		if len(keyword) > 0 {
			groups, err := h.service.GetGroupByNameForDoctor(ctx, user, keyword[0], pag)
			if err != nil {
				sentry.CaptureException(err)
				log.Printf("Exception : %s", err)
				utils.ResErr(w, http.StatusInternalServerError)
				return
			}
			if len(groups) <= 0 {
				utils.ResponseOk(w, groups)
				return
			}
			listGroup, _ := h.service.GetNameGroupForGroup11(ctx, groups, check.ID)
			w.Write(utils.ResponseWithByte(listGroup))
			return
		}

		groups, err := h.service.GetGroupByDoctorService(ctx, user, pag)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		if len(groups) <= 0 {
			fmt.Println("rỗng")
			fmt.Println(groups)
			utils.ResponseOk(w, groups)
			return
		}
		listGroup, _ := h.service.GetNameGroupForGroup11(ctx, groups, check.ID)
		w.Write(utils.ResponseWithByte(listGroup))
	}

}

func (h *Handler) GetListGroupByIdApi(w http.ResponseWriter, r *http.Request) {

}

// FindGroupByName godoc
// @Summary Find Group By Name
// @Description Find Group By Name
// @Tags groups
// @Accept  json
// @Produce  json
// @Param page query string true "current page"
// @Param pageSize query string true "max item in page"
// @Param keyword query string false "name search by keyword"
// @Success 200 {array} Dto
// @Router /api/v1/groups/name [get]
func (h *Handler) FindGroupByName(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	user := ctx.Value(m.UserKey).(string)
	check, err := h.userService.GetUserDetailByIDService(ctx, user)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}

	if len(r.URL.Query()["keyword"]) <= 0 || len(r.URL.Query()["page"]) <= 0 || len(r.URL.Query()["pageSize"]) <= 0 {
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}

	page, err := strconv.Atoi(r.URL.Query()["page"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	maxItemPage, err := strconv.Atoi(r.URL.Query()["pageSize"][0])
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	pag := utils.Pagination{
		Page:        int32(page),
		MaxItemPage: int32(maxItemPage),
	}
	keyword := r.URL.Query()["keyword"]

	if check.Role == userdetail.PATIENT {
		groups, err := h.service.GetGroupByNameForPatient(ctx, user, keyword[0], pag)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		if len(groups) <= 0 {
			utils.ResponseOk(w, groups)
			return
		}
		listGroup, _ := h.service.GetNameGroupForGroup11(ctx, groups, check.ID)
		w.Write(utils.ResponseWithByte(listGroup))
	} else {
		groups, err := h.service.GetGroupByNameForDoctor(ctx, user, keyword[0], pag)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		if len(groups) <= 0 {
			utils.ResponseOk(w, groups)
			return
		}
		listGroup, _ := h.service.GetNameGroupForGroup11(ctx, groups, check.ID)
		w.Write(utils.ResponseWithByte(listGroup))
	}
}

// api tao group n - n chi doctor dc tao va tao chat 1 1

// CreateGroupApi godoc
// @Summary Create a new groups
// @Description create a new groups
// @Tags groups
// @Accept  json
// @Produce  json
// @Param groupPayLoad body PayLoad true "Create groups"
// @Success 200 {object} Dto
// @Router /api/v1/groups [post]
func (h *Handler) CreateGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	var groupPayLoad PayLoad
	err := json.NewDecoder(r.Body).Decode(&groupPayLoad)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}

	owner := ctx.Value(m.UserKey).(string)

	if groupPayLoad.Type == ONE { //api Tạo hội thoại 1 - 1 (nhóm bí mật) ||
		groupsDto, err := h.service.GetGroupByOwnerAndUserService(ctx, groupPayLoad, owner)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		if len(groupsDto) <= 0 {
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		groupDto11, err := h.service.GetNameGroupForGroup11(ctx, groupsDto, owner)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		log.Printf("Group: %v", groupDto11)
		w.Write(utils.ResponseWithByte(groupDto11))
	} else { //Tạo hội thoại n- n
		check, err := h.userService.GetUserDetailByIDService(ctx, owner)
		if check.Role == userdetail.PATIENT {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		groupDto, err := h.service.AddGroupManyService(ctx, groupPayLoad, owner)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		w.Write(utils.ResponseWithByte(groupDto))
	}

}

// api update ten nhom

// UpdateInfoGroupApi godoc
// @Summary Update group by groupId
// @Description Update the group corresponding to the input groupId
// @Tags groups
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Param groupPayLoad body PayLoad true "update groups"
// @Success 200 {object} Dto
// @Router /api/v1/groups/{idGroup} [put]
func (h *Handler) UpdateInfoGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}

	owner := ctx.Value(m.UserKey).(string)
	check, err := h.service.CheckRoleOwnerInGroupService(ctx, owner, id)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}
	if !check {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusForbidden)
		return
	} else {
		var groupPayLoad PayLoad
		err = json.NewDecoder(r.Body).Decode(&groupPayLoad)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusBadRequest)
			return
		}
		newGroup, err := h.service.UpdateGroupService(ctx, groupPayLoad, id)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		w.Write(utils.ResponseWithByte(newGroup))
	}

}

// DeleteGroupApi godoc
// @Summary Delete group identified by the given idGroup
// @Description Delete the group corresponding to the input idGroup
// @Tags groups
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Success 200 {object} utils.ResponseBool
// @Router /api/v1/groups/{idGroup} [delete]
func (h *Handler) DeleteGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	owner := ctx.Value(m.UserKey).(string)
	check, err := h.service.CheckRoleOwnerInGroupService(ctx, owner, id)
	if !check {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusForbidden)
		return
	} else {
		err = h.service.DeleteGroupService(ctx, id)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		result := utils.ResponseBool{Result: true}
		w.Write(utils.ResponseWithByte(result))
	}
}

//API thêm thành viên vào 1 nhóm va chi owner moi dc them

// AddUserInGroupApi godoc
// @Summary add user to group
// @Description add user to group
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Param groupPayLoad body PayLoad true "add user to group"
// @Success 200 {object} utils.ResponseBool
// @Router /api/v1/groups/{idGroup}/members [patch]
func (h *Handler) AddUserInGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	owner := ctx.Value(m.UserKey).(string)
	check, err := h.service.CheckRoleOwnerInGroupService(ctx, owner, id)
	if !check {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusForbidden)
		return
	} else {
		var groupPayload PayLoad
		err = json.NewDecoder(r.Body).Decode(&groupPayload)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusBadRequest)
			return
		}

		err = h.service.AddUserInGroupService(ctx, groupPayload.Users, id)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		result := utils.ResponseBool{Result: true}
		w.Write(utils.ResponseWithByte(result))
	}
}

//API Them thanh vien vao nhom cong khai

// AddUserInGroupPublicApi godoc
// @Summary add user to group public for QR CODE
// @Description add user to group for QR CODE
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Success 200 {object} utils.ResponseBool
// @Router /api/v1/groups/public/{idGroup}/members [post]
func (h *Handler) AddUserInGroupPublicApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idGroup"])

	user := ctx.Value(m.UserKey).(string)
	fmt.Println("id user")
	fmt.Println(user)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}

	var users []string
	users = append(users, user)
	fmt.Println(users)

	group, err := h.service.GetGroupById(ctx, id)
	group.IsMember = false
	fmt.Println(group)
	fmt.Println(Dto{})
	if group == (Dto{}) {
		result := utils.ResponseBool{Result: false}
		w.Write(utils.ResponseWithByte(result))
		return
	}
	if group.Private == false && group.Type == MANY {
		err = h.service.AddUserInGroupService(ctx, users, id)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
	} else {
		result := utils.ResponseBool{Result: false}
		w.Write(utils.ResponseWithByte(result))
		return
	}

	result := utils.ResponseBool{Result: true}
	w.Write(utils.ResponseWithByte(result))
}

//API xoa thành viên khoi 1 nhóm va chi owner moi dc xoa

// DeleteGroupUserApi godoc
// @Summary delete user to group by admin
// @Description delete user to group by admin
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID group"
// @Param userId path int true "ID user want delete"
// @Success 200 {object} utils.ResponseBool
// @Router /api/v1/groups/{idGroup}/members/{userId} [delete]
func (h *Handler) DeleteGroupUserApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	groupID, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	userid := params["userId"]
	owner := ctx.Value(m.UserKey).(string)
	check, err := h.service.CheckRoleOwnerInGroupService(ctx, owner, groupID)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}
	if !check {
		utils.ResErr(w, http.StatusForbidden)
		return
	} else {
		//xoa thanh vien trong nhom
		users := []string{userid}
		err := h.service.DeleteUserInGroupService(ctx, users, groupID)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		result := utils.ResponseBool{Result: true}
		w.Write(utils.ResponseWithByte(result))
	}
}

//API user outgroup nhung owner ko dc out

// UserOutGroupApi godoc
// @Summary delete user to group
// @Description delete user to group
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be add user"
// @Success 200 {object} utils.ResponseBool
// @Router /api/v1/groups/{idGroup}/members [delete]
func (h *Handler) UserOutGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	groupID, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	owner := ctx.Value(m.UserKey).(string)
	check, err := h.service.CheckRoleOwnerInGroupService(ctx, owner, groupID)

	if check {
		utils.ResErr(w, http.StatusForbidden)
		return
	} else {
		users := []string{owner}
		err := h.service.DeleteUserInGroupService(ctx, users, groupID)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			utils.ResErr(w, http.StatusInternalServerError)
			return
		}
		result := utils.ResponseBool{Result: true}
		w.Write(utils.ResponseWithByte(result))
	}

}

// GetListUserOnlineByGroupApi godoc
// @Summary Get all member groups
// @Description Get all member groups
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Success 200 {array} []userdetail.Dto
// @Router /api/v1/groups/{idGroup}/members [get]
func (h *Handler) GetListUserOnlineByGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	params := mux.Vars(r)
	groupID, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusBadRequest)
		return
	}
	users, err := h.service.GetListUserOnlineAndOffByGroupService(ctx, groupID)

	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}

	if users == nil {
		utils.ResErr(w, http.StatusGatewayTimeout)
		return
	}

	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}
	w.Write(utils.ResponseWithByte(users))
}
