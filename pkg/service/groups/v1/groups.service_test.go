package groups

import (
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"log"
	"testing"
)

func TestGetListUserOnlineAndOffByGroupService(t *testing.T) {
	database.Connect()
	dtos, err := GetListUserOnlineAndOffByGroupService(1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dtos)
}
func TestGetGroupByPatientService(t *testing.T) {
	//database.Connect()
	//dtos, err := GetGroupByPatientService("ffb63922-8f99-46ba-9648-d07f3ac14757")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, dto := range dtos {
	//	t.Log(dto.Id)
	//}
}
func TestDeleteGroupService(t *testing.T) {
	database.Connect()
	err := DeleteGroupService(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("sucess")
}
