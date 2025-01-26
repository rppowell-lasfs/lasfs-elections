package main

import (
	"election/api"
	"election/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// fmt.Printf("config:\n%v\n", config.Envs)
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)
	r.Run()
}
