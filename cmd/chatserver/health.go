package main

import (
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"net/http"
)

func CheckHelthHandlr(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	utils.ResponseOk(w, http.StatusOK)
}
