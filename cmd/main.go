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
		log.Fatalln(err)
	}

	err = infra.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = infra.InitDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	r := router.Launch()
	r.Run(":8080")
}
