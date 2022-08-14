package groups

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	message_service "gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"golang.org/x/sync/errgroup"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

type ServiceImpl struct {
	repo           Repo
	userService    userdetail.Service
	messService    message_service.Service
	contextTimeout time.Duration
}

func NewServiceImpl(r Repo, userService userdetail.Service, time time.Duration, messService message_service.Service) Service {
	return &ServiceImpl{
		repo:           r,
		userService:    userService,
		messService:    messService,
		contextTimeout: time,
	}
}

// get group by name for Role doctor
func (s *ServiceImpl) GetGroupByNameForDoctor(ctx context.Context, user string, keyword string, pag utils.Pagination) (groupDtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	groups, err := s.repo.GetGroupByUser(ctx, user)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		pubGroups, err := s.repo.GetGroupByNameForDoctor(ctx, user, keyword, pag)
		if err != nil {
			return err
		}
		if len(pubGroups) > 0 {
			for _, g := range pubGroups {
				pubDto := g.ConvertToDTO()
				message, _ := s.messService.LoadMessageHistoryService(ctx, int(pubDto.Id))
				if len(message) > 0 {
					pubDto.LastMessage = message[0]
				}
				fmt.Println(pubDto)
				_, check := FindObjectInGroup(groups, pubDto.Id)
				if check {
					pubDto.IsMember = true
				} else {
					pubDto.IsMember = false
				}
				if pubDto.Owner == user {
					pubDto.IsOwer = true
				}
				groupDtos = append(groupDtos, pubDto)
			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	return
}

// get group by name for Role patient
func (s *ServiceImpl) GetGroupByNameForPatient(ctx context.Context, user string, keyword string, pag utils.Pagination) (groupDtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	groups, err := s.repo.GetGroupByUser(ctx, user)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		pubGroups, err := s.repo.GetGroupByNameForPatient(ctx, false, user, keyword, pag)
		if err != nil {
			return err
		}
		if len(pubGroups) > 0 {
			for _, g := range pubGroups {
				pubDto := g.ConvertToDTO()
				if g.Thumbnail != "" {
					pubDto = g.ConvertToDTOHaveThumbnail()
				}
				message, _ := s.messService.LoadMessageHistoryService(ctx, int(pubDto.Id))
				if len(message) > 0 {
					pubDto.LastMessage = message[0]
				}
				fmt.Println(pubDto)
				_, check := FindObjectInGroup(groups, pubDto.Id)
				if check {
					pubDto.IsMember = true
				} else {
					pubDto.IsMember = false
				}
				groupDtos = append(groupDtos, pubDto)
			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	return
}

// tao chat 1 1 neu chua co, neu co r tra lai
func (s *ServiceImpl) GetGroupByOwnerAndUserService(ctx context.Context, groupPayload PayLoad, owner string) (groupDtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	group := groupPayload.ConvertToModel()
	_, find := utils.Find(group.Users, owner)
	if find {
		fmt.Println("lỗi trùng ower và users cùng tên")
		return
	}
	fmt.Println("check")
	groups, err := s.repo.GetGroupByOwnerAndUserAndTypeOne(ctx, owner, group.Users[0])
	fmt.Println(owner)
	fmt.Println(group.Users[0])
	fmt.Println(groups)
	if err != nil {
		return nil, err
	}
	if len(groups) <= 0 {
		fmt.Println("nho hon o")
		group.UserCreate = owner
		model, err := s.repo.AddGroupType(ctx, group)
		if err != nil {
			sentry.CaptureException(err)
			return nil, err
		}

		group.Users = append(group.Users, group.UserCreate)
		err = s.repo.AddGroupUser(ctx, group.Users, int(model.ID))
		if err != nil {
			sentry.CaptureException(err)
			return nil, err
		}

		fmt.Println(model.ID)
		group, _ = s.repo.GetGroupIdGroup(ctx, int(model.ID))
		dto := group.ConvertToDTO()
		if group.Thumbnail != "" {
			dto = group.ConvertToDTOHaveThumbnail()
		}
		groupDtos = append(groupDtos, dto)
		return groupDtos, nil
	} else {
		fmt.Println("hello")
		for _, g := range groups {
			dto := g.ConvertToDTO()
			groupDtos = append(groupDtos, dto)
		}
		return
	}
}
func (s *ServiceImpl) GetGroupByPatientService(ctx context.Context, user string, pag utils.Pagination) (groupDtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	groups, err := s.repo.GetGroupByUser(ctx, user)
	if err != nil {
		return
	}
	g, ctx := errgroup.WithContext(ctx)
	//g.Go(func() (err error) {
	//	groups, err := s.repo.GetGroupByUser(ctx, user)
	//	if err != nil {
	//		return
	//	}
	//	for _, group := range groups {
	//		dto := group.ConvertToDTO()
	//		groupDtos = append(groupDtos, dto)
	//	}
	//	return
	//})
	g.Go(func() (err error) {
		pubGroups, err := s.repo.GetGroupByPrivateAndUser(ctx, false, user, pag)
		if err != nil {
			return err
		}
		if len(pubGroups) > 0 {
			for _, g := range pubGroups {
				pubDto := g.ConvertToDTO()
				if g.Thumbnail != "" {
					pubDto = g.ConvertToDTOHaveThumbnail()
				}
				message, _ := s.messService.LoadMessageHistoryService(ctx, int(pubDto.Id))
				if len(message) > 0 {
					pubDto.LastMessage = message[0]
				}
				fmt.Println(pubDto)
				_, check := FindObjectInGroup(groups, pubDto.Id)
				if check {
					pubDto.IsMember = true
				} else {
					pubDto.IsMember = false
				}
				groupDtos = append(groupDtos, pubDto)

			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	return
}

func (s *ServiceImpl) GetGroupById(ctx context.Context, id int) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		group, err := s.repo.GetGroupIdGroup(ctx, id)
		dto = group.ConvertToDTO()
		if group.Thumbnail != "" {
			dto = group.ConvertToDTOHaveThumbnail()
		}
		return
	})

	if err = g.Wait(); err != nil {
		return
	}
	return

}

func (s *ServiceImpl) GetGroupId(ctx context.Context, id int) (group Groups, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		group, err = s.repo.GetGroupIdGroup(ctx, id)
		return
	})

	if err = g.Wait(); err != nil {
		return
	}
	return

}

func (s *ServiceImpl) GetGroupByDoctorService(ctx context.Context, user string, pag utils.Pagination) (groupDtos []Dto, err error) {
	fmt.Println("chu")
	fmt.Println(user)
	fmt.Println("----------")
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	groups, err := s.repo.GetGroupByUser(ctx, user)
	if err != nil {
		return
	}
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		pubGroups, err := s.repo.GetGroupPublicByDoctor(ctx, user, pag)
		if err != nil {
			return err
		}
		if len(pubGroups) > 0 {
			for _, g := range pubGroups {
				pubDto := g.ConvertToDTO()
				if g.Thumbnail != "" {
					pubDto = g.ConvertToDTOHaveThumbnail()
				}
				message, _ := s.messService.LoadMessageHistoryService(ctx, int(pubDto.Id))
				if len(message) > 0 {
					pubDto.LastMessage = message[0]
				}
				fmt.Println(pubDto)
				_, check := FindObjectInGroup(groups, pubDto.Id)
				if check {
					pubDto.IsMember = true
				} else {
					pubDto.IsMember = false
				}
				fmt.Println(pubDto.Owner)
				fmt.Println(user)
				if pubDto.Owner == user {
					pubDto.IsOwer = true
				}
				groupDtos = append(groupDtos, pubDto)
			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	return
}

func (s *ServiceImpl) AddGroupManyService(ctx context.Context, groupPayLoad PayLoad, owner string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	group := groupPayLoad.ConvertToModel()
	index, find := utils.Find(groupPayLoad.Users, owner)
	fmt.Println(index)
	if find {
		fmt.Println("bị trùng và xóa user trùng khỏi list")
		groupPayLoad.Users = utils.RemoveIndex(groupPayLoad.Users, index)
	}

	fmt.Println("danh sach user")
	fmt.Println(group.Users)
	group.UserCreate = owner
	group, err = s.repo.AddGroupType(ctx, group)
	if err != nil {
		return
	}
	dto = group.ConvertToDTO()

	groupPayLoad.Users = append(groupPayLoad.Users, owner)
	err = s.repo.AddGroupUser(ctx, groupPayLoad.Users, int(group.ID))
	if err != nil {
		return
	}
	return
}
func (s *ServiceImpl) UpdateGroupService(ctx context.Context, groupsPayLoad PayLoad, idGroup int) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	fmt.Println(groupsPayLoad)
	group := groupsPayLoad.ConvertToModel()
	group.ID = uint(idGroup)
	fmt.Println(group.Name)
	if group.Thumbnail != "" {
		group, err = s.repo.UpdateGroupWithThumbnail(ctx, group)
		dto = group.ConvertToDTOHaveThumbnail()
		return
	} else {
		group, err = s.repo.UpdateGroup(ctx, group)
		dto = group.ConvertToDTO()
		return
	}

}
func (s *ServiceImpl) DeleteGroupService(ctx context.Context, idGroup int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	//err = s.repo.DeleteGroup(ctx, idGroup)
	err = s.repo.TempDeleteGroup(ctx, idGroup)
	if err != nil {
		return err
	}
	return
}
func (s *ServiceImpl) CheckRoleOwnerInGroupService(ctx context.Context, owner string, idgroup int) (check bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	check = false
	check, err = s.repo.GetOwnerByGroupAndOwner(ctx, owner, idgroup)
	if err != nil {
		return
	}
	return
}
func (s *ServiceImpl) AddUserInGroupService(ctx context.Context, userIds []string, groupId int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	err = s.repo.AddGroupUser(ctx, userIds, groupId)
	if err != nil {
		return
	}
	return
}
func (s *ServiceImpl) DeleteUserInGroupService(ctx context.Context, userIds []string, groupId int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	err = s.repo.DeleteGroupUser(ctx, userIds, groupId)
	if err != nil {
		return
	}
	return
}
func (s *ServiceImpl) GetListUserByGroupService(ctx context.Context, groupId int) (dtos []userdetail.Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	users, err := s.repo.GetListUserByGroup(ctx, groupId)
	if err != nil {
		return
	}
	for _, u := range users {
		user, _ := u.ConvertToDto()
		dtos = append(dtos, user)
	}
	return
}
func (s *ServiceImpl) GetListUserOnlineAndOffByGroupService(ctx context.Context, groupId int) (dtos []userdetail.Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	mapUser, err := s.repo.GetListUserOnlineAndOfflineByGroup(ctx, groupId)
	if err != nil {
		return
	}

	var onlineStr []string
	for _, online := range mapUser[USERON] {
		onlineStr = append(onlineStr, online.ID)
	}
	userIohON := userdetail.GetListFromUserId(onlineStr)

	if len(userIohON) > 0 && userIohON != nil {
		for _, u := range userIohON {
			u.Status = USERON
			dtos = append(dtos, u)
		}
	}

	var offlineStr []string
	for _, offline := range mapUser[USEROFF] {
		offlineStr = append(offlineStr, offline.ID)
	}

	userIohOff := userdetail.GetListFromUserId(offlineStr)

	if len(userIohOff) > 0 && userIohOff != nil {
		for _, u := range userIohOff {
			u.Status = USEROFF
			dtos = append(dtos, u)
		}
	}
	return
}

// truyền vào danh sách group và id của user của chính mình để lấy tên từ người muốn nhắn
func (s *ServiceImpl) GetNameGroupForGroup11(ctx context.Context, dto []Dto, id string) ([]Dto, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	for i, _ := range dto {
		if dto[i].Type == ONE {
			users, _ := s.repo.GetListUserByGroup(ctx, int(dto[i].Id))

			for j, _ := range users {
				if users[j].ID != id {
					value, ok := userdetail.ListUserGlobal[users[j].ID]
					userDetail, _ := s.userService.GetUserDetailByIDService(ctx, users[j].ID)
					if ok == true {
						dto[i].Thumbnail = userDetail.Avatar
						dto[i].Name = value.Username
						if value.FirstName != "" && value.LastName != "" {
							dto[i].Name = value.FirstName + " " + value.LastName
						}
					} else {
						user := userdetail.GetUserFromKCById(users[j].ID)
						log.Printf("USER_KC: %v", user)
						dto[i].Name = user.Username
						dto[i].Thumbnail = userDetail.Avatar
						if user.FirstName != "" && user.LastName != "" {
							dto[i].Name = user.FirstName + " " + user.LastName
						}

					}
					break
				}
			}
		}
	}
	return dto, nil
}
func (s *ServiceImpl) CheckUserAndGroupExits(ctx context.Context, idGroup int, idUser string) bool {
	groups, err := s.repo.GetGroupByUserAndIdGroup(ctx, idUser, idGroup)
	if err != nil {
		sentry.CaptureException(err)
		return false
	}
	fmt.Println(len(groups))
	//if len(groups) > 0 {
	//	return true
	//} else {
	//	return false
	//}

	return len(groups) > 0
}

func (s *ServiceImpl) DeleteGroupNeedDelete(ctx context.Context) (err error) {
	groups, _ := s.repo.GetGroupNeedDelete(ctx)
	for _, group := range groups {
		fmt.Println(group)
		_ = s.repo.DeleteGroup(ctx, int(group.ID))
	}
	return
}

func (s *ServiceImpl) UpdateThumbnailGroupService(id int, file multipart.File, handler *multipart.FileHeader, ctx context.Context) (newImage string, err error) {
	objectName := time.Now().Format("01-02-2021 15:04:05") + "" + handler.Filename
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return
	}
	exists, err := minioClient.BucketExists(ctx, BucketThumbnail)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, BucketThumbnail, minio.MakeBucketOptions{Region: location})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		log.Printf("Successfully created %s\n", id)
	}
	contentType := CheckTypeFile(file)
	_, err = minioClient.PutObject(ctx, BucketThumbnail, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	newImage = objectName
	return
}

func (s *ServiceImpl) UpdateGroupWhenHaveAction(ctx context.Context, idGroup int) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	group, err := s.repo.UpdateGroupWhenHaveAction(ctx, idGroup)
	if err != nil {
		return
	}
	dto = group.ConvertToDTO()
	return
}

func FindObjectInGroup(slice []Groups, id uint) (int, bool) {
	if slice == nil {
		return -1, false
	}
	for i, item := range slice {
		if item.ID == id {
			return i, true
		}
	}
	return -1, false
}

func CheckTypeFile(file multipart.File) string {
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		fmt.Println("err")
	}
	return http.DetectContentType(fileHeader)
}

func GetFileService(bucketName string, objectName string) (shareLink string, err error) {

	value, ok := FileGlobal[objectName]
	if ok {
		if time.Now().Day()-value.CreateAt.Day() < 2 {
			fmt.Println("vào caching")
			shareLink = value.ShareUrl
			return
		}
	}
	fmt.Println(bucketName)
	if objectName == "" {
		fmt.Println("bị rỗng rồi")
		return
	}
	fmt.Println(objectName)
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
	image, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, time.Second*60*60*60, reqParams)
	if err != nil {
		shareLink = ""
		return
	}
	shareLink = image.String()
	newImage := ModelUpload{
		ShareUrl: shareLink,
		NameFile: objectName,
		CreateAt: time.Now(),
	}
	FileGlobal[objectName] = newImage
	fmt.Println(shareLink)
	return
}
