package handler

type UserInfoResp struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
	ID   uint64 `json:"id"`
}
