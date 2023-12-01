package controllers

import (
	"github.com/diiineeei/ApiLojinha/database"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"

	"net/http"
)

func RegisterContacts(context *gin.Context) {
	var produto models.Produto
	if err := context.ShouldBindJSON(&produto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record := database.InstancePostgres.Create(&produto)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated,
		gin.H{
			"ID":        produto.ID,
			"Nome":      produto.Nome,
			"Descricao": produto.Descricao,
			"Valor":     produto.Valor,
			"EmEstoque": produto.EmEstoque,
			"ImagemURL": produto.ImagemURL,
		},
	)
}
