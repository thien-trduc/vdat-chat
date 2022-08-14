package groups

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"mime/multipart"
	"time"
)

var FileGlobal = make(map[string]ModelUpload)
var minioClient *minio.Client

type AbstractModel struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
type Groups struct {
	AbstractModel
	UserCreate  string   `json:"userId"`
	Name        string   `json:"nameGroup"`
	Type        string   `json:"type"`
	Private     bool     `json:"private"`
	Thumbnail   string   `json:"thumbnail"`
	Description string   `json:"description"`
	Users       []string `json:"users"`
}
type GroupsUsers struct {
	ID         uint       `json:"id"`
	CreatedAt  *time.Time `json:"joinAt"`
	UserIDJoin string     `json:"userId"`
}

func (g *Groups) ConvertToDTO() Dto {
	//if g.Thumbnail != "" {
	//	g.Thumbnail, _ = GetFileService(BucketThumbnail, g.Thumbnail)
	//}
	dto := Dto{
		Id:          g.ID,
		Name:        g.Name,
		Type:        g.Type,
		Private:     g.Private,
		Thumbnail:   g.Thumbnail,
		Owner:       g.UserCreate,
		Description: g.Description,
		IsMember:    true,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
	return dto
}

type ModelUpload struct {
	ShareUrl string
	Location string
	NameFile string
	Type     string
	CreateAt time.Time
}

func (g *Groups) ConvertToDTOHaveThumbnail() Dto {
	shareImage, err := GetFileService(BucketThumbnail, g.Thumbnail)
	if err != nil {
		fmt.Println(err)
		shareImage = ""
	}

	dto := Dto{
		Id:          g.ID,
		Name:        g.Name,
		Type:        g.Type,
		Private:     g.Private,
		Thumbnail:   shareImage,
		Owner:       g.UserCreate,
		Description: g.Description,
		IsMember:    true,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
	return dto
}

type Repo interface {
	GetGroupIdGroup(ctx context.Context, idGroup int) (group Groups, err error)
	GetGroupByOwnerAndUserAndTypeOne(ctx context.Context, owner string, user string) ([]Groups, error)
	GetGroupByUser(ctx context.Context, user string) ([]Groups, error)
	GetGroupByPrivateAndUser(ctx context.Context, private bool, user string, pag utils.Pagination) ([]Groups, error)
	GetGroupByType(ctx context.Context, typeGroup string, user string) ([]Groups, error)
	GetOwnerByGroupAndOwner(ctx context.Context, owner string, groupId int) (bool, error)
	GetListUserByGroup(ctx context.Context, idGourp int) ([]userdetail.UserDetail, error)
	GetListUserOnlineAndOfflineByGroup(ctx context.Context, idGroup int) (map[string][]userdetail.UserDetail, error)
	GetGroupPublicByDoctor(ctx context.Context, user string, pag utils.Pagination) ([]Groups, error)
	AddGroupType(ctx context.Context, group Groups) (Groups, error)
	AddGroupUser(ctx context.Context, users []string, idgroup int) error
	UpdateGroup(ctx context.Context, group Groups) (Groups, error)
	UpdateGroupWithThumbnail(ctx context.Context, group Groups) (Groups, error)
	DeleteGroup(ctx context.Context, idGourp int) error
	DeleteGroupUser(ctx context.Context, users []string, idgroup int) error
	GetGroupByUserAndIdGroup(ctx context.Context, user string, idGroup int) ([]Groups, error)
	TempDeleteGroup(ctx context.Context, idGroup int) error
	GetGroupNeedDelete(ctx context.Context) ([]Groups, error)
	GetGroupByNameForDoctor(ctx context.Context, user string, keyword string, pag utils.Pagination) ([]Groups, error)
	GetGroupByNameForPatient(ctx context.Context, private bool, user string, keyword string, pag utils.Pagination) ([]Groups, error)
	UpdateGroupWhenHaveAction(ctx context.Context, idGroup int) (Groups, error)
}

type Service interface {
	GetGroupById(ctx context.Context, id int) (dto Dto, err error)
	GetGroupByOwnerAndUserService(ctx context.Context, groupPayload PayLoad, owner string) ([]Dto, error)
	GetGroupByPatientService(ctx context.Context, user string, pag utils.Pagination) ([]Dto, error)
	GetGroupByDoctorService(ctx context.Context, user string, pag utils.Pagination) ([]Dto, error)
	AddGroupManyService(ctx context.Context, groupPayLoad PayLoad, owner string) (Dto, error)
	UpdateGroupService(ctx context.Context, groupsPayLoad PayLoad, idGroup int) (Dto, error)
	DeleteGroupService(ctx context.Context, idGroup int) error
	CheckRoleOwnerInGroupService(ctx context.Context, owner string, idgroup int) (bool, error)
	AddUserInGroupService(ctx context.Context, userIds []string, groupId int) error
	DeleteUserInGroupService(ctx context.Context, userIds []string, groupId int) error
	GetListUserByGroupService(ctx context.Context, groupId int) ([]userdetail.Dto, error)
	GetListUserOnlineAndOffByGroupService(ctx context.Context, groupId int) ([]userdetail.Dto, error)
	GetNameGroupForGroup11(ctx context.Context, dto []Dto, id string) ([]Dto, error)
	CheckUserAndGroupExits(ctx context.Context, idGroup int, idUser string) bool
	DeleteGroupNeedDelete(ctx context.Context) error
	GetGroupByNameForDoctor(ctx context.Context, user string, keyword string, pag utils.Pagination) (groupDtos []Dto, err error)
	GetGroupByNameForPatient(ctx context.Context, user string, keyword string, pag utils.Pagination) (groupDtos []Dto, err error)
	UpdateGroupWhenHaveAction(ctx context.Context, idGroup int) (Dto, error)
	GetGroupId(ctx context.Context, id int) (group Groups, err error)
	UpdateThumbnailGroupService(id int, file multipart.File, handler *multipart.FileHeader, ctx context.Context) (newImage string, err error)
}
