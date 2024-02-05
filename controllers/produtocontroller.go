package controllers

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/diiineeei/ApiLojinha/database"
	"github.com/diiineeei/ApiLojinha/log"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
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

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader

func UploadProduto(ctx *gin.Context) {
	var produto models.Produto
	err := ctx.ShouldBind(&produto)
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

	file, err := ctx.FormFile("imagem")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blobFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("%s.png", time.Now().Format("20060102040405"))

	c := context.Background()
	client, err := storage.NewClient(c)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"failed to get default GCS bucket name": err.Error()})
		return
	}
	defer client.Close()

	uploader := &ClientUploader{
		cl:         client,
		bucketName: "adega-produtos",
		projectID:  "adega-do-vilao",
		uploadPath: "adega-do-vilao/",
	}

	wc := uploader.cl.Bucket(uploader.bucketName).Object(uploader.uploadPath + filename).NewWriter(ctx)
	if _, err := io.Copy(wc, blobFile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"io.Copy": err.Error()})
		return
	}
	if err := wc.Close(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Writer.Close": err.Error()})
		return
	}

	produto.ImagemURL = filename

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
