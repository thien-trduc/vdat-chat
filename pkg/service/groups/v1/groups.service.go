package groups

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/message/v1"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v1"
	"golang.org/x/sync/errgroup"
)

// tao chat 1 1 neu chua co, neu co r tra lai
func GetGroupByOwnerAndUserService(groupPayload PayLoad, owner string) ([]Dto, error) {
	groupdtos := make([]Dto, 0)

	group := groupPayload.ConvertToModel()
	groups, err := NewRepoImpl(database.DB).GetGroupByOwnerAndUserAndTypeOne(owner, group.Users[0])
	if err != nil {
		return nil, err
	} else {
		if len(groups) <= 0 {
			group.UserCreate = owner
			group, err := NewRepoImpl(database.DB).AddGroupType(group)
			if err != nil {
				sentry.CaptureException(err)
				return nil, err
			}
			group.Users = append(group.Users, group.UserCreate)
			err = NewRepoImpl(database.DB).AddGroupUser(group.Users, int(group.ID))
			if err != nil {
				sentry.CaptureException(err)
				return nil, err
			}
			groups, err = NewRepoImpl(database.DB).GetGroupByOwnerAndUserAndTypeOne(group.UserCreate, group.Users[0])
			if err != nil {
				sentry.CaptureException(err)
				return nil, err
			}
			for _, group := range groups {
				groupdto := group.ConvertToDTO()

				groupdtos = append(groupdtos, groupdto)
			}
			return groupdtos, nil
		} else {
			for _, group := range groups {
				groupdto := group.ConvertToDTO()
				groupdtos = append(groupdtos, groupdto)
			}
			return groupdtos, nil
		}
	}
}

//func GetGroupByPatientService(user string) ([]Dto, error) {
//	groupdtos := make([]Dto, 0)
//	groups, err := NewRepoImpl(database.DB).GetGroupByUser(user)
//	if err != nil {
//		sentry.CaptureException(err)
//		return nil, err
//	}
//	pubGroups, err := NewRepoImpl(database.DB).GetGroupByPrivateAndUser(false, user)
//	if err != nil {
//		sentry.CaptureException(err)
//		return nil, err
//	}
//	if len(pubGroups) > 0 {
//		for _, g := range pubGroups {
//			pubDto := g.ConvertToDTO()
//			pubDto.IsMember = false
//			groupdtos = append(groupdtos, pubDto)
//		}
//	}
//	for _, group := range groups {
//		groupdto := group.ConvertToDTO()
//		groupdtos = append(groupdtos, groupdto)
//	}
//	return groupdtos, nil
//}
func GetGroupByPatientService(ctx context.Context, user string) (groupDtos []Dto, err error) {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		groups, err := NewRepoImpl(database.DB).GetGroupByUser(user)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		for _, group := range groups {
			dto := group.ConvertToDTO()
			groupDtos = append(groupDtos, dto)
		}
		return
	})
	g.Go(func() (err error) {
		pubGroups, err := NewRepoImpl(database.DB).GetGroupByPrivateAndUser(false, user)
		if err != nil {
			sentry.CaptureException(err)
			return err
		}
		if len(pubGroups) > 0 {
			for _, g := range pubGroups {
				pubDto := g.ConvertToDTO()
				pubDto.IsMember = false
				groupDtos = append(groupDtos, pubDto)
			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		sentry.CaptureException(err)
		return
	}
	return
}

//func GetGroupByDoctorService(user string) ([]Dto, error) {
//	groupdtos := make([]Dto, 0)
//	groups, err := NewRepoImpl(database.DB).GetGroupByUser(user)
//	if err != nil {
//		sentry.CaptureException(err)
//		return nil, err
//	}
//	pubGroups, err := NewRepoImpl(database.DB).GetGroupPublicByDoctor(user)
//	if err != nil {
//		sentry.CaptureException(err)
//		return nil, err
//	}
//	if len(pubGroups) > 0 {
//		for _, g := range pubGroups {
//			pubDto := g.ConvertToDTO()
//			pubDto.IsMember = false
//			groupdtos = append(groupdtos, pubDto)
//		}
//	}
//	for _, group := range groups {
//		groupdto := group.ConvertToDTO()
//		groupdtos = append(groupdtos, groupdto)
//	}
//	return groupdtos, nil
//}
func GetGroupByDoctorService(ctx context.Context, user string) (groupDtos []Dto, err error) {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		groups, err := NewRepoImpl(database.DB).GetGroupByUser(user)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		for _, group := range groups {
			dto := group.ConvertToDTO()
			groupDtos = append(groupDtos, dto)
		}
		return
	})
	g.Go(func() (err error) {
		pubGroups, err := NewRepoImpl(database.DB).GetGroupPublicByDoctor(user)
		if err != nil {
			sentry.CaptureException(err)
			return err
		}
		if len(pubGroups) > 0 {
			for _, g := range pubGroups {
				pubDto := g.ConvertToDTO()
				pubDto.IsMember = false
				groupDtos = append(groupDtos, pubDto)
			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		sentry.CaptureException(err)
		return
	}
	return
}
func AddGroupManyService(groupPayLoad PayLoad, owner string) (Dto, error) {
	var groupdto Dto
	group := groupPayLoad.ConvertToModel()
	group.UserCreate = owner
	group.Users = append(group.Users, owner)
	group, err := NewRepoImpl(database.DB).AddGroupType(group)
	if err != nil {
		sentry.CaptureException(err)
		return groupdto, err
	}
	groupdto = group.ConvertToDTO()
	err = NewRepoImpl(database.DB).AddGroupUser(group.Users, int(group.ID))
	if err != nil {
		sentry.CaptureException(err)
		return groupdto, err
	}
	return groupdto, nil
}

func UpdateGroupService(groupsPayLoad PayLoad, idGroup int) (Dto, error) {
	var groupdto Dto
	group := groupsPayLoad.ConvertToModel()
	group.ID = uint(idGroup)
	group, err := NewRepoImpl(database.DB).UpdateGroup(group)
	if err != nil {
		sentry.CaptureException(err)
		return groupdto, err
	}
	groupdto = group.ConvertToDTO()
	return groupdto, err
}
func DeleteGroupService(idgroup int) error {
	//ch := make(chan error)
	//go func(ch1 chan error) {
	//
	//	err := message.DeleteMessageService(idgroup)
	//
	//	ch1 <- err
	//
	//	close(ch1)
	//
	//}(ch)
	//for {
	//	select {
	//	case kq:= <- ch:
	//		if kq == nil{
	//			fmt.Println("debug")
	//			return NewRepoImpl(database.DB).DeleteGroup(idgroup)
	//		}
	//	default:
	//		close(ch)
	//	}
	//}
	err := message.DeleteMessageService(idgroup)
	if err != nil {
		return err
	}
	return NewRepoImpl(database.DB).DeleteGroup(idgroup)
}
func CheckRoleOwnerInGroupService(owner string, idgroup int) (bool, error) {
	return NewRepoImpl(database.DB).GetOwnerByGroupAndOwner(owner, idgroup)
}
func AddUserInGroupService(userIds []string, groupId int) error {
	return NewRepoImpl(database.DB).AddGroupUser(userIds, groupId)
}
func DeleteUserInGroupService(userIds []string, groupId int) error {
	return NewRepoImpl(database.DB).DeleteGroupUser(userIds, groupId)
}
func GetListUserByGroupService(groupId int) ([]userdetail.Dto, error) {
	dtos := make([]userdetail.Dto, 0)
	users, err := NewRepoImpl(database.DB).GetListUserByGroup(groupId)
	if err != nil {
		sentry.CaptureException(err)
		return dtos, err
	}
	for _, u := range users {
		user := u.ConvertToDto()
		dtos = append(dtos, user)
	}
	return dtos, nil
}

func GetListUserOnlineAndOffByGroupService(groupId int) ([]userdetail.Dto, error) {
	dtos := make([]userdetail.Dto, 0)
	mapUser, err := NewRepoImpl(database.DB).GetListUserOnlineAndOfflineByGroup(groupId)
	if err != nil {
		sentry.CaptureException(err)
		return dtos, err
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
	return dtos, nil
}

// truyền vào danh sách group và id của user của chính mình để lấy tên từ người muốn nhắn
func getNameGroupForGroup11(dto []Dto, id string) ([]Dto, error) {
	for i, _ := range dto {
		if dto[i].Type == ONE {
			users, _ := NewRepoImpl(database.DB).GetListUserByGroup(int(dto[i].Id))

			for j, _ := range users {
				if users[j].ID != id {
					value, ok := userdetail.ListUserGlobal[users[j].ID]
					userDetail, _ := userdetail.GetUserDetailByIDService(users[j].ID)
					if ok == true {
						dto[i].Thumbnail = userDetail.Avatar
						dto[i].Name = value.Username
						if value.FirstName != "" && value.LastName != "" {
							dto[i].Name = value.FirstName + " " + value.LastName
						}
					} else {
						user := userdetail.GetUserFromKCById(users[j].ID)
						fmt.Println(user)
						fmt.Println(user.Username)
						fmt.Println(user.LastName)
						fmt.Println(user.FirstName)
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

func CheckUserAndGroupExits(idGroup int, idUser string) bool {
	groups, _ := NewRepoImpl(database.DB).GetGroupByUserAndIdGroup(idUser, idGroup)
	fmt.Println(len(groups))
	if len(groups) > 0 {
		return true
	} else {
		return false
	}
}
