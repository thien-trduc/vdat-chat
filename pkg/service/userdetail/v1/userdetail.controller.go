package userdetail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/useronline/v1"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func RegisterUserApi(r *mux.Router) {

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

	r.HandleFunc("/api/v1/user", auth.AuthenMiddleJWT(GetUserApi)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/user/info", auth.AuthenMiddleJWT(CheckUserDetailApi)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/user/online", auth.AuthenMiddleJWT(UserLogOutApi)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/user/role", auth.AuthenMiddleJWT(ExchangeRole)).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc("/api/v1/user/detail", auth.AuthenMiddleJWT(ExchangeRole)).Methods(http.MethodPatch, http.MethodOptions)
}

//API tìm kiếm người dùng filtter

// find user by keyword godoc
// @Summary find users by keyword
// @Description find user by keyword
// @Tags user
// @Accept  json
// @Produce  json
// @Param keyword query string false "name search by keyword"
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {array} Dto
// @Router /api/v1/user [get]
func GetUserApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	if len(r.URL.Query()["keyword"]) == 0 || len(r.URL.Query()["page"]) == 0 || len(r.URL.Query()["pageSize"]) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	fil := r.URL.Query()["keyword"]
	page := r.URL.Query()["page"]
	pageSize := r.URL.Query()["pageSize"]

	if page[0] == "" {
		page[0] = "1"
	}
	if pageSize[0] == "" {
		pageSize[0] = "10"
	}
	listUser := getData(fil[0], page[0], pageSize[0])
	if len(listUser) == 0 {
		json.NewEncoder(w).Encode(listUser)
	} else {
		w.Write(utils.ResponseWithByte(listUser))
	}

	//a:= []string{"b9018379-8394-4205-9104-2d85d69943db","b767e36c-e4a9-4d8c-886c-181427ec4e2c","f51ae747-9ab1-446b-bc66-62c49ec307df","0bc77b02-ecba-43f5-82b0-f846f721984b","cefcb41e-bf21-4bc0-97d2-4981ab946a2b","9ea08917-bbfc-49f3-bc8f-158f745f0ff7","8ba1f1ca-b72f-4cbc-8562-36e174376694","5614020d-5322-4c9d-b1ce-80297b05f83f","c98f749f-f45c-485c-aba2-256f1bdc7440","425feec0-dd5c-4ce7-b1a7-9ac0223a9b14","9d1cf0d5-8d35-40d5-9cfa-29b0f1a90f15","718f2d59-d841-4f84-845b-b697c8af4a76","91c153ad-2fcf-4a53-9766-c4f4564d94d6","ffb63922-8f99-46ba-9648-d07f3ac14757","d84ffce0-f9a6-43ef-953f-b56debc6bc6f","5a56852a-ba09-4e5f-aae1-4769301688c3","c772f9bf-4128-4898-bdbe-f7edf5fa5b3c"}
	//a := []string{"a", "b767e36c-e4a9-4d8c-886c-181427ec4e2c", "f51ae747-9ab1-446b-bc66-62c49ec307df", "0bc77b02-ecba-43f5-82b0-f846f721984b"}
	//b := []string{"dfa7e1a3-0e98-42d3-b88f-18b1afe07a17", "dc160880-512d-435d-931f-ab9f9187a08f", "afeea15b-b37f-4934-85b4-5e2eecbcc43f", "0bc77b02-ecba-43f5-82b0-f846f721984b", "cefcb41e-bf21-4bc0-97d2-4981ab946a2b", "9ea08917-bbfc-49f3-bc8f-158f745f0ff7", "8ba1f1ca-b72f-4cbc-8562-36e174376694", "5614020d-5322-4c9d-b1ce-80297b05f83f", "c98f749f-f45c-485c-aba2-256f1bdc7440", "425feec0-dd5c-4ce7-b1a7-9ac0223a9b14", "9d1cf0d5-8d35-40d5-9cfa-29b0f1a90f15", "718f2d59-d841-4f84-845b-b697c8af4a76", "91c153ad-2fcf-4a53-9766-c4f4564d94d6", "ffb63922-8f99-46ba-9648-d07f3ac14757", "d84ffce0-f9a6-43ef-953f-b56debc6bc6f", "5a56852a-ba09-4e5f-aae1-4769301688c3", "c772f9bf-4128-4898-bdbe-f7edf5fa5b3c"}
	////user := GetListFromUserId(a)
	//user := GetListFromUserId(b)
	//w.Write(utils.ResponseWithByte(user))
}

// checkUser godoc
// @Summary check user api
// @Description check user api
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} Dto
// @Router /api/v1/user/info [get]
func CheckUserDetailApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)

	if (*r).Method == "OPTIONS" {
		return
	}

	payload, err := JWTparseUser(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dto, err := CheckUserDetailService(payload)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	//dto.HostName = utils.GetLocalIP()
	dto.HostName = payload.ID
	//dto.SocketID = utils.ArraySocketId[0]
	//utils.ArraySocketId = utils.DeleteItemInArray(utils.ArraySocketId)
	//utils.WriteLines(utils.ArraySocketId, "socketid.data")

	uo := useronline.Payload{
		HostName: dto.HostName,
		SocketID: dto.SocketID,
		UserID:   payload.ID,
	}
	err = useronline.AddUserOnlineService(uo)
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// utils.ResponseErr(w, http.StatusOK)
	w.Write(utils.ResponseWithByte(dto))

	// check user he thong neu login chua ton tai thong tin trong he thong thi ghi vao database

}

// user logout godoc
// @Summary user logout
// @Description user logout api
// @Tags user
// @Accept  json
// @Produce  json
// @Param hostName query string false "hostName"
// @Param socketId query string false "socketId"
// @Success 200 {object} boolean
// @Router /api/v1/user/online [post]
func UserLogOutApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	if len(r.URL.Query()["hostName"]) == 0 || len(r.URL.Query()["socketId"]) == 0 {
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	hostname := r.URL.Query()["hostName"][0]
	socketID := r.URL.Query()["socketId"][0]
	err := useronline.DeleteUserOnlineService(socketID, hostname)
	if err != nil {
		sentry.CaptureException(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	//utils.ArraySocketId = utils.RestoreItemArray(utils.ArraySocketId, socketID)
	//utils.WriteLines(utils.ArraySocketId, "socketid.data")
	w.Write(utils.ResponseWithByte(true))
}

// user exchange role godoc
// @Summary user exchange role
// @Description user exchange role api
// @Tags user
// @Accept  json
// @Produce  json
// @Param role query string false "role"
// @Param idUser query string false "idUser"
// @Success 200 {object} boolean
// @Router /api/v1/user/role [patch]
func ExchangeRole(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	if len(r.URL.Query()["role"]) == 0 || len(r.URL.Query()["idUser"]) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	role := r.URL.Query()["role"][0]
	idUser := r.URL.Query()["idUser"][0]
	user := auth.JWTparseOwner(r.Header.Get("Authorization"))
	check, err := GetUserDetailByIDService(user)
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if check.Role != ADMIN {
		fmt.Println("khong la admin")
		w.WriteHeader(http.StatusForbidden)
		utils.ResponseErr(w, http.StatusForbidden)
	} else {
		fmt.Println("la admin")
		payload := Payload{
			ID:   idUser,
			Role: role,
		}
		err = UpdateUserDetailservice(payload)
		utils.ResponseOk(w, true)
	}
}

func Connect() string {
	const (
		clientSecret string = "7161982e-cabe-44d3-ade1-324698d2f5d8"
		clientId     string = "chat.services.vdatlab.com"
		urlHost      string = "https://accounts.vdatlab.com/auth/realms/vdatlab.com/protocol/openid-connect/token"
	)

	client := &http.Client{}
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Add("client_secret", clientSecret)
	data.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", urlHost, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
	}
	resp.Body.Close()
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	var token Token
	json.Unmarshal(f, &token)
	//fmt.Print(token.AccessToken)
	//fmt.Println(string(f))

	return token.AccessToken
}
