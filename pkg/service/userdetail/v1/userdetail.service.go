package userdetail

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var GlobalToken string = Connect()
var ListUserGlobal = make(map[string]User)

func AddUserDetailService(payload Payload) error {
	detail := payload.convertToModel()
	err := NewRepoImpl(database.DB).AddUserDetail(detail)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}
func UpdateUserDetailservice(payload Payload) error {
	detail := payload.convertToModel()
	err := NewRepoImpl(database.DB).UpdateUserDetail(detail)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}
func GetUserDetailByIDService(id string) (Dto, error) {
	var dto Dto
	detail, err := NewRepoImpl(database.DB).GetUserDetailById(id)
	if err != nil {
		sentry.CaptureException(err)
		return dto, err
	}
	dto = detail.ConvertToDto()
	if detail.Avatar != "" {
		linkImage, _ := getFileService(bucketNameAvatar, detail.Avatar)
		dto.Avatar = linkImage
	}
	return dto, nil
}
func GetListUserDetailService() ([]Dto, error) {
	dtos := make([]Dto, 0)
	userdetails, err := NewRepoImpl(database.DB).GetListUser()
	if err != nil {
		sentry.CaptureException(err)
		return dtos, err
	}
	for _, detail := range userdetails {
		if detail.Avatar != "" {
			linkImage, _ := getFileService(bucketNameAvatar, detail.Avatar)
			detail.Avatar = linkImage
		}
		dto := detail.ConvertToDto()
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

//neu user chua co thi add thong tin tu token vao database user
func CheckUserDetailService(payload Payload) (Dto, error) {
	var dto Dto
	userdetail, err := NewRepoImpl(database.DB).GetUserDetailById(payload.ID)
	if err != nil {
		return dto, err
	}
	if userdetail == (UserDetail{}) {
		payload.Role = PATIENT
		err = AddUserDetailService(payload)
		if err != nil {
			fmt.Println("loi 1 ")
			fmt.Println(err)
			sentry.CaptureException(err)
			return dto, err
		}
	} else {
		payload.Role = userdetail.Role
		err = UpdateUserDetailservice(payload)
		if err != nil {
			fmt.Println("loi 2")
			sentry.CaptureException(err)
			return dto, err
		}
	}

	userdetail, err = NewRepoImpl(database.DB).GetUserDetailById(payload.ID)
	if err != nil {
		sentry.CaptureException(err)
		return dto, err
	}
	dto = userdetail.ConvertToDto()
	fmt.Println("day la dto")
	fmt.Println(dto)
	return dto, nil
}

func JWTparseUser(tokenHeader string) (Payload, error) {
	var payload Payload
	splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token

	block, _ := pem.Decode([]byte(auth.Jwtkey))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)

	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	tokenPart := splitted[1]
	tk := &auth.UserClaims{}
	_, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	})
	if err != nil {
		return payload, err
	}
	payload = Payload{
		ID:       tk.Subject,
		FullName: tk.FullName,
		Username: tk.UserName,
		First:    tk.GivenName,
		Last:     tk.FamilyName,
	}
	return payload, nil
}

func getData(keyword string, page string, pageSize string) []Dto {

	if !utils.CheckTokenExp(GlobalToken) {
		GlobalToken = Connect()
	}
	//fmt.Println(time)

	size, _ := strconv.Atoi(pageSize)
	pageInt, _ := strconv.Atoi(page)

	if size <= 0 {
		size = 10
	}
	if pageInt <= 1 {
		pageInt = 1
	}

	num := size * pageInt
	var expectNum int
	if pageInt > 1 {
		expectNum = size * (pageInt - 1)
	} else {
		expectNum = size * pageInt
	}
	if pageInt == 1 {
		expectNum = 0
	}
	var (
		urlHost string = "https://vdat-mcsvc-kc-admin-api-auth-proxy.vdatlab.com/auth/admin/realms/vdatlab.com/users?search="
	)
	URL := fmt.Sprintf(urlHost+"%s"+"&max=%s"+"&first=%s", keyword, strconv.Itoa(num), strconv.Itoa(expectNum))
	var bearer = "Bearer " + GlobalToken

	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var users []User
	json.Unmarshal([]byte(body), &users)
	//fmt.Print(users)
	var userDtos []Dto
	if len(users) == 0 {
		empty := make([]Dto, 0)
		return empty
	}
	for i, _ := range users {
		//fmt.Println(users[i].ID)
		detail, _ := NewRepoImpl(database.DB).GetUserDetailById(users[i].ID)
		if detail == (UserDetail{}) {
			fmt.Println("khong ton tai")
			users[i].Role = PATIENT
			payload := Payload{
				ID:   users[i].ID,
				Role: PATIENT,
			}
			err = AddUserDetailService(payload)
			if err != nil {
				fmt.Println(err)
			}
			Dto := users[i].ConvertUserToDto()
			userDtos = append(userDtos, Dto)

		} else {
			users[i].Role = detail.Role
			Dto := users[i].ConvertUserToDto()
			userDtos = append(userDtos, Dto)
		}
	}
	//fmt.Print(string(body))
	fmt.Println(len(userDtos))
	return userDtos
}

func GetListFromUserId(listUser []string) []Dto {
	var userDtos []Dto

	fmt.Println("danh sách id người dùng")
	for i, _ := range listUser {
		fmt.Println(listUser[i])

		value, ok := ListUserGlobal[listUser[i]]
		if ok == true {
			detail, _ := NewRepoImpl(database.DB).GetUserDetailById(listUser[i])
			if detail == (UserDetail{}) {
				value.Role = ""
			} else {
				value.Role = detail.Role
			}
			dto := value.ConvertUserToDto()
			userDtos = append(userDtos, dto)
		} else {
			user := GetUserFromKCById(listUser[i])
			if !(user == (User{})) {
				ListUserGlobal[listUser[i]] = user

				detail, _ := NewRepoImpl(database.DB).GetUserDetailById(listUser[i])
				if detail == (UserDetail{}) {
					user.Role = ""
				} else {
					user.Role = detail.Role
				}
				dto := user.ConvertUserToDto()
				userDtos = append(userDtos, dto)
			} else {
				fmt.Printf("-- user có id: %s ko tồn tại (ko add vào danh sách trả về)\n", listUser[i])
			}
		}
	}
	//fmt.Println(userDtos)
	fmt.Println(len(userDtos))
	return userDtos
}

func GetUserFromKCById(id string) User {
	if !utils.CheckTokenExp(GlobalToken) {
		GlobalToken = Connect()
	}
	var (
		urlHost string = "https://vdat-mcsvc-kc-admin-api-auth-proxy.vdatlab.com/auth/admin/realms/vdatlab.com/users/"
		bearer         = "Bearer " + GlobalToken
	)

	req, err := http.NewRequest("GET", urlHost+id, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var user User
	json.Unmarshal(body, &user)
	ListUserGlobal[id] = user
	return user
}

func getFileService(bucketName string, objectName string) (shareLink string, err error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Println(err)
		return
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
	image, _ := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, time.Second*60*60*60, reqParams)
	shareLink = image.String()
	return
}
