package handler

type RequestUidAndPassword struct {
	Uid      string `json:"name" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}