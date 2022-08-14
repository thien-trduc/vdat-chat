package upload

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/groups/v2"
	message_service "gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	GroupService groups.Service
}

func NewHandler(r *mux.Router) {
	timeoutContext := time.Duration(2) * time.Second
	messRepo := message_service.NewRepoImpl(database.DB)
	messService := message_service.NewServiceImpl(messRepo, timeoutContext)
	repo := groups.NewRepoImpl(database.DB)
	userRepo := userdetail.NewRepoImpl(database.DB)
	userService := userdetail.NewServiceImpl(userRepo, timeoutContext)
	service := groups.NewServiceImpl(repo, userService, timeoutContext, messService)

	minioClient, _ = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})

	minioEndPortStr := os.Getenv("MINIO_END_PORT")
	fmt.Println(minioEndPortStr)
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

	handler := &Handler{GroupService: service}
	r.HandleFunc("/api/v1/files/upload", auth.AuthenMiddleJWT(handler.UploadFileForOtherService)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/files/upload/{idGroup}", auth.AuthenMiddleJWT(handler.Upload)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/files/{idGroup}", auth.AuthenMiddleJWT(handler.GetListFileInGroup)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/files/download/{idGroup}/{objectName}", auth.AuthenMiddleJWT(handler.DownloadFileApi)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/files/avatar/user", auth.AuthenMiddleJWT(handler.UpdateAvatarUser)).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc("/api/v1/files/test/caching", auth.AuthenMiddleJWT(handler.testCachingImage)).Methods(http.MethodPost, http.MethodOptions)
}

// add upload file godoc
// @Summary upload file
// @Description upload file
// @Tags upload
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group upload file"
// @Success 200 {object} ModelUpload
// @Router /api/v1/files/upload/{idGroup} [post]
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	params := mux.Vars(r)
	ctx := r.Context()
	idGroup := params["idGroup"]
	bucketName := "group-" + idGroup
	file, handler, err := r.FormFile("file")
	fmt.Println(r.ContentLength)
	//typeFile := CheckTypeFile(file)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("r.FormFile %+v", err)
		utils.ResErr(w, http.StatusInternalServerError)
		return
	}

	checkExists, err := CheckBucketExits(bucketName)
	if err != nil {
		fmt.Printf("CheckBucketExits %+v", err)
		return
	}
	if checkExists {
		log.Printf("We already own %s\n", bucketName)
	} else {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			fmt.Printf("MakeBucket %+v", err)
		}
	}

	objectName := time.Now().Format("01-02-2021 15:04:05") + "_" + handler.Filename
	contentType := "image/png"
	newImage, err := PutObjectTOMinio(bucketName, objectName, file, contentType, r.ContentLength)
	if err != nil {
		fmt.Printf("PutObjectTOMinio %+v", err)
		sentry.CaptureException(err)
		utils.ResErr(w, http.StatusInternalServerError)
	}
	utils.ResponseWithJSON(w, 200, newImage)
}

// add get list file godoc
// @Summary get list file
// @Description get list file
// @Tags upload
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group get list file"
// @Param type query string false "type"
// @Success 200 {array} ModelUpload
// @Router /api/v1/files/{idGroup} [get]
func (h *Handler) GetListFileInGroup(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	params := mux.Vars(r)
	ctx := r.Context()
	idGroup := params["idGroup"]
	bucketName := "group-" + idGroup
	user := auth.JWTparseOwner(r.Header.Get("Authorization"))
	convertIdGroup, _ := strconv.Atoi(idGroup)
	checkUserAndGroup := h.GroupService.CheckUserAndGroupExits(ctx, convertIdGroup, user)
	if !checkUserAndGroup {
		fmt.Println("khong trung voi group")
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	// Initialize minio client object.
	checkExits, err := CheckBucketExits(bucketName)
	if err != nil {
		fmt.Printf("CheckBucketExits %+v", err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if !checkExits {
		utils.ResponseWithJSON(w, 200, []ModelUpload{})
	}
	if len(r.URL.Query()["type"]) == 0 {
		fmt.Println("rong")
		listImage := GetListFileInBucket(bucketName, all)
		utils.ResponseWithJSON(w, 200, listImage)
	} else {
		role := r.URL.Query()["type"][0]
		fmt.Println(role)
		listImage := GetListFileInBucket(bucketName, role)
		utils.ResponseWithJSON(w, 200, listImage)
	}
}

// DownloadFileApi godoc
// @Summary dowload file
// @Description dowload file
// @Tags upload
// @Accept  json
// @Produce  json
// @Param idGroup path int true "ID of the group to get file"
// @Param objectName path string true "name file"
// @Success 200 {object} ModelUpload
// @Router /api/v1/files/download/{idGroup}/{objectName} [get]
func (h *Handler) DownloadFileApi(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	params := mux.Vars(r)
	idGroup := params["idGroup"]
	bucketName := "group-" + idGroup
	objectName := params["objectName"]
	ctx := r.Context()

	user := auth.JWTparseOwner(r.Header.Get("Authorization"))
	convertIdGroup, _ := strconv.Atoi(idGroup)
	checkUserAndGroup := h.GroupService.CheckUserAndGroupExits(ctx, convertIdGroup, user)
	if !checkUserAndGroup {
		fmt.Println("khong trung voi group")
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
	res, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		log.Println(err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	result := ModelUpload{
		ShareUrl: res.String(),
		Location: res.Host + res.RawPath,
		NameFile: objectName,
	}
	log.Println("Successfully saved my-filename.csv")
	utils.ResponseWithJSON(w, 200, result)

}

// add upload file godoc
// @Summary upload file
// @Description upload file
// @Tags upload
// @Accept  json
// @Produce  json
// @Success 200 {object} ModelUpload
// @Router /api/v1/files/upload [post]
func (h *Handler) UploadFileForOtherService(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	file, handler, err := r.FormFile("file")
	//typeFile := CheckTypeFile(file)

	checkExists, err := CheckBucketExits(BucketThumbnail)
	if err != nil {
		fmt.Printf("CheckBucketExits %+v", err)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	if checkExists {
		log.Printf("We already own %s\n", BucketNameAvatar)
	} else {
		err = minioClient.MakeBucket(context.Background(), BucketThumbnail, minio.MakeBucketOptions{Region: location, ObjectLocking: true})
		log.Printf("Successfully created %s\n", BucketNameAvatar)
	}

	objectName := time.Now().Format("01-02-2021 15:04:05") + "#" + handler.Filename
	contentType := "image/png"
	newImage, err := PutObjectTOMinio(BucketThumbnail, objectName, file, contentType, r.ContentLength)
	if err != nil {
		fmt.Println(err)
		utils.ResErr(w, http.StatusInternalServerError)
	}
	utils.ResponseWithJSON(w, 200, newImage)
}

// UpdateAvatarUser godoc
// @Summary upload file
// @Description upload file avatar for user
// @Tags upload
// @Accept  json
// @Produce  json
// @Success 200 {object} ModelUpload
// @Router /api/v1/files/avatar/user [patch]
func (h *Handler) UpdateAvatarUser(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	ctx := r.Context()
	payload, err := userdetail.JWTparseUser(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		utils.ResponseErr(w, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newImage, err := UpdateUserAvatarService(payload.ID, file, handler, ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ResponseErr(w, http.StatusInternalServerError)
		return
	}
	utils.ResponseWithJSON(w, 200, newImage)
}

func (h *Handler) testCachingImage(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	buket := r.URL.Query()["bucket"]
	nameFile := r.URL.Query()["nameFile"]
	linkshare, _ := GetFileService(buket[0], nameFile[0])
	//_ = RemoveFileService(buket[0], nameFile[0])
	utils.ResponseOk(w, linkshare)
}

//func (h *Handler) UpdateThumbnailForGroup(w http.ResponseWriter, r *http.Request) {
//	cors.SetupResponse(&w, r)
//	ctx := r.Context()
//	params := mux.Vars(r)
//	id, err := strconv.Atoi(params["idGroup"])
//	file, handler, err := r.FormFile("file")
//	if err != nil {
//		fmt.Println(err)
//		utils.ResponseErr(w, http.StatusBadRequest)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	newImage, err := UpdateThumbnailGroupService(id, file, handler, ctx)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		utils.ResponseErr(w, http.StatusInternalServerError)
//		return
//	}
//
//	group, _ := h.GroupService.GetGroupId(ctx, id)
//	group.Thumbnail = newImage.NameFile
//
//	utils.ResponseOk(w, newImage)
//}
