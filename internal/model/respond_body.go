package model

type UserInfoResp struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

type TokenResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint64 `json:"expires_in"`
}

type LoginResp struct {
	UserInfo   UserInfoResp `json:"user_info"`
	TokenClass TokenResp    `json:"token_class"`
}
