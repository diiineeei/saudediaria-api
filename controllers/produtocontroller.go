package controllers

import (
	"github.com/diiineeei/ApiLojinha/database"
	"github.com/diiineeei/ApiLojinha/log"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CadastraProduto(ctx *gin.Context) {
	var produto models.Produto
	err := ctx.ShouldBindJSON(&produto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	err = database.PutProdutos(produto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated,
		gin.H{
			"ID":             produto.ID,
			"Nome":           produto.Nome,
			"Descricao":      produto.Descricao,
			"Valor":          produto.Valor,
			"EmEstoque":      produto.EmEstoque,
			"ImagemURL":      produto.ImagemURL,
			"CodigoDeBarras": produto.CodigoDeBarras,
		},
	)
}

func ListaProdutos(ctx *gin.Context) {
	var filterProduto models.Produto
	var err error
	log.Infof(ctx.Request, "Teste")
	paramID := ctx.Query("id")
	if paramID != "" {
		filterProduto.ID, err = strconv.ParseInt(paramID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"erro ao decodificar parametro ID": err})
			ctx.Abort()
			return
		}
	}
	paramEmEstoque := ctx.Query("emEstoque")
	if paramEmEstoque != "" {
		filterProduto.EmEstoque, err = strconv.ParseBool(paramEmEstoque)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"erro ao decodificar parametro emEstoque": err})
			ctx.Abort()
			return
		}
	}
	filterProduto.Nome = ctx.Query("nome")
	filterProduto.CodigoDeBarras = ctx.Query("codigoDeBarras")

	produtos, err := database.GetProdutos(filterProduto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"erro GetProdutos": err})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, produtos)
}
