package controllers

import (
	"github.com/diiineeei/ApiLojinha/database"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"
	"strconv"

	"net/http"
)

const Table = "produtos"

func CadastraProduto(ctx *gin.Context) {
	var produto models.Produto
	if err := ctx.ShouldBindJSON(&produto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	record := database.InstancePostgres.Create(&produto)
	if record.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
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
	var (
		id                   int64
		nome, codigoDeBarras string
		emEstoque            bool
		err                  error
	)

	paramID := ctx.Query("id")
	if paramID != "" {
		id, err = strconv.ParseInt(paramID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"erro ao decodificar parametro ID": err})
			ctx.Abort()
			return
		}
	}
	paramEmEstoque := ctx.Query("emEstoque")
	if paramEmEstoque != "" {
		emEstoque, err = strconv.ParseBool(paramEmEstoque)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"erro ao decodificar parametro emEstoque": err})
			ctx.Abort()
			return
		}
	}
	nome = ctx.Query("nome")
	codigoDeBarras = ctx.Query("codigoDeBarras")

	var produtos []models.Produto

	query := database.InstancePostgres.Table(Table)
	if id > 0 {
		query = query.Where("id = ?", id)
	}
	if nome != "" {
		query = query.Where("nome = ?", "%"+nome+"%")
	}
	if codigoDeBarras != "" {
		query = query.Where("codigo_de_barras = ?", codigoDeBarras)
	}
	if emEstoque {
		query = query.Where("em_estoque = ?", emEstoque)
	}
	query.Find(&produtos)
	ctx.JSON(http.StatusOK, produtos)
}
