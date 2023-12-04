package main

import (
	"github.com/diiineeei/ApiLojinha/controllers"
	"github.com/diiineeei/ApiLojinha/database"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	// Initialize Database
	database.ConnectPostgres("postgres://admin:admin@localhost:5432/admin")
	database.Migrate()
}
func main() {
	// Initialize Router
	router := gin.Default()
	router.POST("/api/cadastrar/produto", controllers.CadastraProduto)

	//api := router.Group("/api")

	//{
	//	//api.POST("/token", controllers.GenerateToken)
	//	api.POST("/user/register", controllers.RegisterUser)
	//	secured := api.Group("/secured").Use(middlewares.Auth())
	//	{
	//		secured.POST("/contact/register", controllers.CadastraProduto)
	//		secured.POST("/contact/register/list", controllers.RegisterContactsList)
	//	}
	//}
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
