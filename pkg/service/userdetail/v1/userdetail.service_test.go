package userdetail

import (
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"testing"
)

func TestAddUserDetailService(t *testing.T) {
	database.Connect()
	p := Payload{
		ID:   "test",
		Role: ADMIN,
	}
	err := AddUserDetailService(p)
	if err != nil {
		t.Fatal(err)
		return
	}
}
func TestGetUserDetailByIDService(t *testing.T) {
	database.Connect()
	dto, err := GetUserDetailByIDService("893a4692-63bb-4919-80d9-aece678c0422")
	if err != nil {
		t.Fatal(err)
		return
	}
	if dto != (Dto{}) {
		t.Log(dto)
	}
}
func TestGetListUserDetailService(t *testing.T) {
	database.Connect()
	dtos, err := GetListUserDetailService()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(dtos)
}
func TestUpdateUserDetailservice(t *testing.T) {
	database.Connect()

	payload := Payload{
		ID:   "ffb63922-8f99-46ba-9648-d07f3ac14757",
		Role: DOCTOR,
	}
	err := UpdateUserDetailservice(payload)
	if err != nil {
		t.Fatal(err)
		return
	}

}

//func TestCheckUserDetailService(t *testing.T) {
//	database.Connect()
//	//dto,err := CheckUserDetailService("893a4692-63bb-4919-80d9-aece678c0422")
//	dto,err := CheckUserDetailService("893a46")
//	if err!= nil{
//		t.Fatal(err)
//	}
//	t.Log(dto)
//}
