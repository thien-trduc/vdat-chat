package userdetail

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	_ "github.com/getsentry/sentry-go"
	"github.com/minio/minio-go/v7"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
)

var GlobalToken string = Connect()
var ListUserGlobal = make(map[string]User)
var minioClient *minio.Client

type ServiceImpl struct {
	repo           Repo
	contextTimeout time.Duration
}

func NewServiceImpl(r Repo, time time.Duration) Service {
	return &ServiceImpl{
		repo:           r,
		contextTimeout: time,
	}
}

func (s *ServiceImpl) AddUserDetailService(ctx context.Context, payload Payload) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	detail := payload.convertToModel()
	err := s.repo.AddUserDetail(ctx, detail)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}

func (s *ServiceImpl) UpdateUserDetailservice(ctx context.Context, payload Payload) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	detail := payload.convertToModel()
	err := s.repo.UpdateUserDetail(ctx, detail)
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
		return err
	}
	return nil
}

func (s *ServiceImpl) GetUserDetailByIDService(ctx context.Context, id string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	detail, err := s.repo.GetUserDetailById(ctx, id)
	if err != nil {
		sentry.CaptureException(err)
		return dto, err
	}
	dto, _ = detail.ConvertToDto()
	return dto, nil
}

func (s *ServiceImpl) GetListUserDetailService(ctx context.Context) (dtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	userdetails, err := s.repo.GetListUser(ctx)
	if err != nil {
		sentry.CaptureException(err)
		return dtos, err
	}
	for _, detail := range userdetails {
		dto, _ := detail.ConvertToDto()
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (s *ServiceImpl) CheckUserDetailService(ctx context.Context, payload Payload) (Dto, error) {
	ctx, _ = context.WithTimeout(ctx, s.contextTimeout)
	var dto Dto
	userdetail, err := NewRepoImpl(database.DB).GetUserDetailById(ctx, payload.ID)
	if err != nil {
		return dto, err
	}
	if userdetail == (UserDetail{}) {
		payload.Role = PATIENT
		err = s.AddUserDetailService(ctx, payload)
		if err != nil {
			sentry.CaptureException(err)
			return dto, err
		}
	} else {
		payload.Role = userdetail.Role
		err = s.UpdateUserDetailservice(ctx, payload)
		if err != nil {
			sentry.CaptureException(err)
			return dto, err
		}
	}

	userdetail, err = NewRepoImpl(database.DB).GetUserDetailById(ctx, payload.ID)
	if err != nil {
		sentry.CaptureException(err)
		return dto, err
	}
	dto, _ = userdetail.ConvertToDto()
	return dto, nil
}

//func (s *ServiceImpl) GetListFromUserIdv2(ctx context.Context, listUser []string) (userDtos []Dto, err error) {
//	g, ctx := errgroup.WithContext(ctx)
//	for i, _ := range listUser {
//		g.Go(func() error {
//			value, ok := ListUserGlobal[listUser[i]]
//			if ok == true {
//				detail, err2 := s.repo.GetUserDetailById(ctx, listUser[i])
//				if err2 != nil {
//					return err2
//				}
//				if detail == (UserDetail{}) {
//					value.Role = ""
//				} else {
//					value.Role = detail.Role
//				}
//				dto := value.ConvertUserToDto()
//				userDtos = append(userDtos, dto)
//			} else {
//				user := GetUserFromKCById(listUser[i])
//				if !(user == (User{})) {
//					ListUserGlobal[listUser[i]] = user
//
//					detail, err2 := s.repo.GetUserDetailById(ctx, listUser[i])
//					if err2 != nil {
//						return err2
//					}
//					if detail == (UserDetail{}) {
//						user.Role = ""
//					} else {
//						user.Role = detail.Role
//					}
//					dto := user.ConvertUserToDto()
//					userDtos = append(userDtos, dto)
//				} else {
//					fmt.Printf("-- user có id: %s ko tồn tại (ko add vào danh sách trả về)\n", listUser[i])
//				}
//			}
//			return nil
//		})
//	}
//	if err = g.Wait();err != nil {
//		return
//	}
//	return
//}

func JWTparseUser(tokenHeader string) (payload Payload, err error) {
	splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token

	block, _ := pem.Decode([]byte(auth.Jwtkey))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)

	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	tokenPart := splitted[1]
	tk := &auth.UserClaims{}
	_, err = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
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
	timeoutContext := time.Duration(2) * time.Second
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	for i, _ := range users {
		//fmt.Println(users[i].ID)
		detail, _ := NewRepoImpl(database.DB).GetUserDetailById(context.Background(), users[i].ID)
		if detail == (UserDetail{}) {
			fmt.Println("khong ton tai")
			users[i].Role = PATIENT
			payload := Payload{
				ID:   users[i].ID,
				Role: PATIENT,
			}
			err = service.AddUserDetailService(context.Background(), payload)
			if err != nil {
				fmt.Println(err)
			}
			Dto := users[i].ConvertUserToDto()
			userDtos = append(userDtos, Dto)

		} else {
			users[i].Role = detail.Role
			Dto := users[i].ConvertUserToDto()
			if detail.Avatar != "" {
				Dto.Avatar, err = getFileService(bucketNameAvatar, detail.Avatar)
				if err != nil {
					fmt.Println(err)
					Dto.Avatar = ""
				}
			}
			userDtos = append(userDtos, Dto)
		}
	}
	//fmt.Print(string(body))
	fmt.Println(len(userDtos))
	return userDtos
}

func GetUserByEmail(email string, page string, pageSize string) []Dto {

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
		urlHost string = "https://vdat-mcsvc-kc-admin-api-auth-proxy.vdatlab.com/auth/admin/realms/vdatlab.com/users?email="
	)
	URL := fmt.Sprintf(urlHost+"%s"+"&max=%s"+"&first=%s", email, strconv.Itoa(num), strconv.Itoa(expectNum))
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
	timeoutContext := time.Duration(2) * time.Second
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	for i, _ := range users {
		//fmt.Println(users[i].ID)
		detail, _ := NewRepoImpl(database.DB).GetUserDetailById(context.Background(), users[i].ID)
		if detail == (UserDetail{}) {
			fmt.Println("khong ton tai")
			users[i].Role = PATIENT
			payload := Payload{
				ID:   users[i].ID,
				Role: PATIENT,
			}
			err = service.AddUserDetailService(context.Background(), payload)
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

func GetUserByUsername(username string, page string, pageSize string) []Dto {

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
		urlHost string = "https://vdat-mcsvc-kc-admin-api-auth-proxy.vdatlab.com/auth/admin/realms/vdatlab.com/users?username="
	)
	URL := fmt.Sprintf(urlHost+"%s"+"&max=%s"+"&first=%s", username, strconv.Itoa(num), strconv.Itoa(expectNum))
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
	timeoutContext := time.Duration(2) * time.Second
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	for i, _ := range users {
		//fmt.Println(users[i].ID)
		detail, _ := NewRepoImpl(database.DB).GetUserDetailById(context.Background(), users[i].ID)
		if detail == (UserDetail{}) {
			fmt.Println("khong ton tai")
			users[i].Role = PATIENT
			payload := Payload{
				ID:   users[i].ID,
				Role: PATIENT,
			}
			err = service.AddUserDetailService(context.Background(), payload)
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

	//fmt.Println("danh sách id người dùng")
	for i, _ := range listUser {
		//fmt.Println(listUser[i])

		value, ok := ListUserGlobal[listUser[i]]
		if ok == true {
			detail, _ := NewRepoImpl(database.DB).GetUserDetailById(context.Background(), listUser[i])
			if detail == (UserDetail{}) {
				value.Role = ""
			} else {
				value.Role = detail.Role
			}
			dto := value.ConvertUserToDto()
			if detail.Avatar != "" {
				dto.Avatar, _ = getFileService(bucketNameAvatar, detail.Avatar)
			}
			userDtos = append(userDtos, dto)
		} else {
			user := GetUserFromKCById(listUser[i])
			if !(user == (User{})) {
				ListUserGlobal[listUser[i]] = user

				detail, _ := NewRepoImpl(database.DB).GetUserDetailById(context.Background(), listUser[i])
				if detail == (UserDetail{}) {
					user.Role = ""
				} else {
					user.Role = detail.Role
				}
				dto := user.ConvertUserToDto()
				if detail.Avatar != "" {
					dto.Avatar, _ = getFileService(bucketNameAvatar, detail.Avatar)
				}
				userDtos = append(userDtos, dto)
			} else {
				fmt.Printf("-- user có id: %s ko tồn tại (ko add vào danh sách trả về)\n", listUser[i])
			}
		}
	}
	//fmt.Println(userDtos)
	//fmt.Println(len(userDtos))
	return userDtos
}

func GetUserFromKCById(id string) User {
	value, ok := userdetail.ListUserGlobal[id]
	if ok {
		return User(value)
	}

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
		sentry.CaptureException(err)
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var user User
	json.Unmarshal(body, &user)

	ListUserGlobal[id] = user
	return user
}

func GetUserById(id string) Dto {
	timeoutContext := time.Duration(2) * time.Second
	repo := NewRepoImpl(database.DB)
	service := NewServiceImpl(repo, timeoutContext)
	value, ok := ListUserGlobal[id]
	fmt.Println()
	if ok {
		fmt.Println("co ton tai")
		userDto, _ := service.GetUserDetailByIDService(context.Background(), id)
		getUser := User(value)
		userDto.Username = getUser.Username
		userDto.First = getUser.FirstName
		userDto.Last = getUser.LastName
		userDto.FullName = getUser.FirstName + " " + getUser.LastName
		return userDto
	}

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
		sentry.CaptureException(err)
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var user User
	json.Unmarshal(body, &user)

	fmt.Println(user)
	payload := Payload{
		ID:   id,
		Role: PATIENT,
	}
	userDto, err := service.CheckUserDetailService(context.Background(), payload)
	if err != nil {
		fmt.Println(err)
		return Dto{}
	}
	userDto.FullName = user.FirstName + " " + user.LastName
	userDto.Username = user.Username
	userDto.First = user.FirstName
	userDto.Last = user.LastName
	ListUserGlobal[id] = user
	return userDto
}

func getFileService(bucketName string, objectName string) (shareLink string, err error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
	image, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, time.Second*60*60*60, reqParams)
	if err != nil {
		return "", err
	}
	shareLink = image.String()
	return
}
