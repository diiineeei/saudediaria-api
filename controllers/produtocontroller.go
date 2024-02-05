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

func UploadProduto(ctx *gin.Context) {
	var produto models.Produto
	err := ctx.Bind(&produto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	file, err := ctx.FormFile("imagem")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileReader.Close()

	fileBin, err := io.ReadAll(fileReader)

	c := context.Background()
	client, err := storage.NewClient(c)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"failed to get default GCS bucket name": err.Error()})
		return
	}
	defer client.Close()

	// Nome do bucket e caminho onde o arquivo será salvo
	bucketName := "adega-produtos"
	uploadPath := fmt.Sprintf("%s.png", time.Now().Format("20060102040405"))
	//uploadPath := "adega-do-vilao/" + filename

	// Iniciando o upload
	wc := client.Bucket(bucketName).Object(uploadPath).NewWriter(c)

	wc.ContentType = "image/png"
	_, err = wc.Write(fileBin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "wc.Write failed: " + err.Error()})
		return
	}

	if err := wc.Close(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "GCS Writer.Close failed: " + err.Error()})
		return
	}

	//downloadURL, err := client.Bucket(bucketName).SignedURL(uploadPath, &storage.SignedURLOptions{
	//	Expires: time.Now().AddDate(0, 0, 7),
	//	Scheme:  storage.SigningSchemeV4,
	//	Method:  "GET",
	//})
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting download URL: " + err.Error()})
	//	return
	//}

	produto.ImagemURL = fmt.Sprintf("https://storage.cloud.google.com/adega-produtos/%s", uploadPath)
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
