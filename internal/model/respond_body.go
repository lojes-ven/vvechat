package model

type UserInfoResp struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

type LoginResp struct {
	UserInfo     UserInfoResp `json:"user_info"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    uint64       `json:"expires_in"`
}

type RefreshTokenResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint64 `json:"expires_in"`
}
