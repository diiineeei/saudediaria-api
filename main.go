package main

import (
	"github.com/diiineeei/ApiLojinha/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//	func init() {
//		database.Database("postgres://admin:admin@localhost:5432/admin")
//		//database.Migrate()
//	}
func main() {
	router := gin.Default()
	router.Use(cors.Default())
	//router.Use(static.Serve("/", static.LocalFile("static_production", false)))
	router.POST("/api/cadastrar/produto", controllers.CadastraProduto)
	router.POST("/api/upload", controllers.UploadProduto)
	router.GET("/api/lista/produtos", controllers.ListaProdutos)
	router.GET("/api/ping", func(c *gin.Context) {
		c.String(200, "test")
	})
	err := router.Run(":8080")
	if err != nil {
		panic(err)
		return
	}
}
