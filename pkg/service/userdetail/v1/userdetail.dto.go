package userdetail

type Dto struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Username string `json:"userName"`
	First    string `json:"first"`
	Last     string `json:"last"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	HostName string `json:"hostName"`
	SocketID string `json:"socketId"`
	Avatar   string `json:"avatar"`
}
