package main

import (
	"log"
	"vvechat/internal/infra"
	"vvechat/internal/router"
	"vvechat/pkg/utils"
)

func main() {
	err := utils.InitSnowflake()
	if err != nil {
		log.Println(err)
		return
	}

	err = infra.InitConfig()
	if err != nil {
		log.Println(err)
		return
	}
	
	db, err := infra.InitDatabase()
	if err != nil {
		log.Println(err)
		return
	}

	r := router.Launch(db)
	r.Run(":8080")
}
