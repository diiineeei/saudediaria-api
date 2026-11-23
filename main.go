package main

import (
	"github.com/diiineeei/ApiLojinha/controllers"
	"github.com/diiineeei/ApiLojinha/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	mongoConfig := os.Getenv("MONGO_CONFIG")
	database.ConnectMongo(mongoConfig)
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/registro", controllers.Registro)
	router.GET("/api/registros", controllers.ListaProdutos)
	router.POST("/api/glicemia/:person", controllers.CreateGlicemiaEntry)
	router.GET("/api/glicemia/:person/list", controllers.ListGlicemiaEntries)
	router.PUT("/api/glicemia/:person/:id", controllers.UpdateGlicemiaEntry)
	router.DELETE("/api/glicemia/:person/:id", controllers.DeleteGlicemiaEntry)
	router.GET("/api/prescricao/:person", controllers.GetPrescricao)
	router.PUT("/api/prescricao/:person", controllers.SavePrescricao)
	router.POST("/api/insulina/:person", controllers.CreateAplicacao)
	router.GET("/api/insulina/:person/list", controllers.ListAplicacoes)
	router.PUT("/api/insulina/:person/:id", controllers.UpdateAplicacao)
	router.DELETE("/api/insulina/:person/:id", controllers.DeleteAplicacao)
	router.POST("/api/alexa", controllers.AlexaWebhook)
	router.POST("/api/alexa/webhook", controllers.AlexaWebhook)
	router.GET("/api/ping", func(c *gin.Context) {
		c.String(200, "test")
	})
	err := router.Run(":8080")
	if err != nil {
		panic(err)
		return
	}
}
