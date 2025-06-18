package controllers

import (
	"cloud.google.com/go/storage"
	"github.com/diiineeei/ApiLojinha/database"
	"github.com/diiineeei/ApiLojinha/log"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Registro(ctx *gin.Context) {
	var registro models.Registro
	err := ctx.ShouldBindJSON(&registro)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	err = database.PutRegistro(registro)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Created"})
}

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

func ListaProdutos(ctx *gin.Context) {
	var filter models.Registro
	var err error
	log.Infof(ctx.Request, "Teste")

	produtos, err := database.GetRegistros(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"erro GetProdutos": err})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, produtos)
}
