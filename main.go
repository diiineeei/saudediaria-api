package main

import (
	"github.com/diiineeei/ApiLojinha/controllers"
	"github.com/diiineeei/ApiLojinha/database"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	database.ConnectPostgres("postgres://admin:admin@localhost:5432/admin")
	database.Migrate()
}
func main() {
	router := gin.Default()
	router.POST("/api/cadastrar/produto", controllers.CadastraProduto)
	router.GET("/api/lista/produtos", controllers.ListaProdutos)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
