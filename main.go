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
	//router.Use(static.Serve("/", static.LocalFile("static_production", false)))
	router.POST("/api/cadastrar/produto", controllers.CadastraProduto)
	router.GET("/api/lista/produtos", controllers.ListaProdutos)
	router.GET("/api/ping", func(c *gin.Context) {
		c.String(200, "test")
	})
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
