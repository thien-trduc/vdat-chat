package groups

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v1"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"

	"net/http"
	"strconv"
)

func RegisterGroupApi(r *mux.Router) {
	r.HandleFunc("/api/v1/groups", auth.AuthenMiddleJWT(GetListGroupApi)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/groups", auth.AuthenMiddleJWT(CreateGroupApi)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}", auth.AuthenMiddleJWT(UpdateInfoGroupApi)).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}", auth.AuthenMiddleJWT(DeleteGroupApi)).Methods(http.MethodDelete, http.MethodOptions)

	r.HandleFunc("/api/v1/groups/{idGroup}/members", auth.AuthenMiddleJWT(GetListUserOnlineByGroupApi)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}/members", auth.AuthenMiddleJWT(AddUserInGroupApi)).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}/members", auth.AuthenMiddleJWT(UserOutGroupApi)).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/v1/groups/{idGroup}/members/{userId}", auth.AuthenMiddleJWT(DeleteGroupUserApi)).Methods(http.MethodDelete, http.MethodOptions)

}

//API load danh sách groups theo patient hoac theo doctor

// GetListGroupApi godoc
// @Summary Get all groups
// @Description Get all groups
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 200 {array} Dto
// @Router /api/v1/groups [get]
func GetListGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	user := auth.JWTparseOwner(r.Header.Get("Authorization"))
	check, err := userdetail.GetUserDetailByIDService(user)
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	fmt.Println(user)
	// danh sách user đã được lưu
	fmt.Println(len(userdetail.ListUserGlobal))

	if check.Role == userdetail.PATIENT {
		fmt.Println("là bệnh nhân")
		//groups, err := GetGroupByPatientService(user)
		groups, err := GetGroupByPatientService(ctx, user)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
			return
		}

		listGroup, _ := getNameGroupForGroup11(groups, check.ID)
		w.Write(utils.ResponseWithByte(listGroup))
	} else {
		fmt.Println("là bác sĩ")
		groups, err := GetGroupByDoctorService(ctx, user)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
			return
		}

		listGroup, _ := getNameGroupForGroup11(groups, check.ID)
		w.Write(utils.ResponseWithByte(listGroup))
	}

}

// api tao group n - n chi doctor dc tao va tao chat 1 1

// CreateOrder godoc
// @Summary Create a new groups
// @Description create a new groups
// @Tags groups
// @Accept  json
// @Produce  json
// @Param groupPayLoad body PayLoad true "Create groups"
// @Success 200 {object} Dto
// @Router /api/v1/groups [post]
func CreateGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)

	var groupPayLoad PayLoad
	err := json.NewDecoder(r.Body).Decode(&groupPayLoad)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}

	owner := auth.JWTparseOwner(r.Header.Get("Authorization"))

	if groupPayLoad.Type == ONE { //api Tạo hội thoại 1 - 1 (nhóm bí mật) ||
		groupsDto, err := GetGroupByOwnerAndUserService(groupPayLoad, owner)
		if err != nil {
			sentry.CaptureException(err)
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
			return
		}
		groupDto11, _ := getNameGroupForGroup11(groupsDto, owner)
		w.Write(utils.ResponseWithByte(groupDto11))
	} else { //Tạo hội thoại n- n
		check, err := userdetail.GetUserDetailByIDService(owner)
		if check.Role == userdetail.PATIENT {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusForbidden)
			utils.ResponseErr(w, http.StatusForbidden)
			return
		}
		groupDto, err := AddGroupManyService(groupPayLoad, owner)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
			return
		}
		w.Write(utils.ResponseWithByte(groupDto))
	}

}

// api update ten nhom

// Updategroups godoc
// @Summary Update group by groupId
// @Description Update the group corresponding to the input groupId
// @Tags groups
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Param groupPayLoad body PayLoad true "update groups"
// @Success 200 {object} Dto
// @Router /api/v1/groups/{idGroup} [put]
func UpdateInfoGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}

	owner := auth.JWTparseOwner(r.Header.Get("Authorization"))
	check, err := CheckRoleOwnerInGroupService(owner, id)
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if !check {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusForbidden)
		utils.ResponseErr(w, http.StatusForbidden)
		return
	} else {
		var groupPayLoad PayLoad
		err = json.NewDecoder(r.Body).Decode(&groupPayLoad)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseErr(w, http.StatusBadRequest)
			return
		}
		newgroup, err := UpdateGroupService(groupPayLoad, id)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
			return
		}
		w.Write(utils.ResponseWithByte(newgroup))
	}

}

// DeleteOrder godoc
// @Summary Delete group identified by the given idGroup
// @Description Delete the group corresponding to the input idGroup
// @Tags groups
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Success 204 "No Content"
// @Router /api/v1/groups/{idGroup} [delete]
func DeleteGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	owner := auth.JWTparseOwner(r.Header.Get("Authorization"))
	check, err := CheckRoleOwnerInGroupService(owner, id)
	if !check {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusForbidden)
		utils.ResponseErr(w, http.StatusForbidden)
		return
	} else {
		err = DeleteGroupService(id)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
			return
		}
		utils.ResponseOk(w, true)
	}
}

//API thêm thành viên vào 1 nhóm va chi owner moi dc them

// add user to group godoc
// @Summary add user to group
// @Description add user to group
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Param groupPayLoad body PayLoad true "add user to group"
// @Success 200 {object} boolean
// @Router /api/v1/groups/{idGroup}/members [patch]
func AddUserInGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	owner := auth.JWTparseOwner(r.Header.Get("Authorization"))
	check, err := CheckRoleOwnerInGroupService(owner, id)
	if !check {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusForbidden)
		utils.ResponseErr(w, http.StatusForbidden)
		return
	} else {
		var groupPayload PayLoad
		err = json.NewDecoder(r.Body).Decode(&groupPayload)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseErr(w, http.StatusBadRequest)
			return
		}

		err = AddUserInGroupService(groupPayload.Users, id)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
		}
		result := utils.ResponseBool{Result: true}
		w.Write(utils.ResponseWithByte(result))
	}
}

//API xoa thành viên khoi 1 nhóm va chi owner moi dc xoa

// delete user to group by admin godoc
// @Summary delete user to group by admin
// @Description delete user to group by admin
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID group"
// @Param userId path int true "ID user want delete"
// @Success 200
// @Router /api/v1/groups/{idGroup}/members/{userId} [delete]
func DeleteGroupUserApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)

	params := mux.Vars(r)
	groupID, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusBadRequest)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	userid := params["userId"]
	owner := auth.JWTparseOwner(r.Header.Get("Authorization"))
	check, err := CheckRoleOwnerInGroupService(owner, groupID)
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if !check {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusForbidden)
		utils.ResponseErr(w, http.StatusForbidden)
		return
	} else {
		//xoa thanh vien trong nhom
		users := []string{userid}
		err := DeleteUserInGroupService(users, groupID)
		if err != nil {
			sentry.CaptureException(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
			return
		}
		result := utils.ResponseBool{Result: true}
		w.Write(utils.ResponseWithByte(result))
	}
}

//API user outgroup nhung owner ko dc out

// delete user to group godoc
// @Summary delete user to group
// @Description delete user to group
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be add user"
// @Success 200
// @Router /api/v1/groups/{idGroup}/members [delete]
func UserOutGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)

	params := mux.Vars(r)
	groupID, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	owner := auth.JWTparseOwner(r.Header.Get("Authorization"))
	check, err := CheckRoleOwnerInGroupService(owner, groupID)

	if check {
		w.WriteHeader(http.StatusForbidden)
		utils.ResponseErr(w, http.StatusForbidden)
		return
	} else {
		//
		users := []string{owner}
		err := DeleteUserInGroupService(users, groupID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.ResponseErr(w, http.StatusInternalServerError)
			return
		}
		result := utils.ResponseBool{Result: true}
		w.Write(utils.ResponseWithByte(result))
	}

}
func GetListUserByGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	params := mux.Vars(r)
	groupID, err := strconv.Atoi(params["idGroup"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}

	users, err := GetListUserByGroupService(groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	w.Write(utils.ResponseWithByte(users))
}

// GetListMemberGroupApi godoc
// @Summary Get all member groups
// @Description Get all member groups
// @Tags groupUser
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to be updated"
// @Success 200 {array} []userdetail.Dto
// @Router /api/v1/groups/{idGroup}/members [get]
func GetListUserOnlineByGroupApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	params := mux.Vars(r)
	groupID, err := strconv.Atoi(params["idGroup"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.ResponseErr(w, http.StatusBadRequest)
		return
	}
	users, err := GetListUserOnlineAndOffByGroupService(groupID)

	if users == nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		utils.ResponseErr(w, http.StatusGatewayTimeout)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	w.Write(utils.ResponseWithByte(users))
}
