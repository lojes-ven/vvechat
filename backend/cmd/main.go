package main

import (
	"log"
	"os"

	"github.com/lojes7/inquire/internal/router"
	"github.com/lojes7/inquire/pkg/infra"
)

func main() {
	infra.Init()
	/*infra.GetDB().AutoMigrate(&model.User{})
	infra.GetDB().AutoMigrate(&model.Friendship{})
	infra.GetDB().AutoMigrate(&model.FriendshipRequest{})
	infra.GetDB().AutoMigrate(&model.Message{})
	infra.GetDB().AutoMigrate(&model.Conversation{})
	infra.GetDB().AutoMigrate(&model.MessageUser{})
	infra.GetDB().AutoMigrate(&model.ConversationUser{})
	infra.GetDB().AutoMigrate(&model.File{})*/
	r := router.Launch()

	address := ":" + os.Getenv("PORT")

	err := r.Run(address)
	if err != nil {
		log.Fatalln("路由器出错")
	}
}
