package database

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/diiineeei/ApiLojinha/utils"
)

const CollectionProdutos = "produtos"
const DatabaseProdutos = "produtos"

func PutProdutos(produto *models.Produto) error {
	ctx := context.Background()

	//opt := option.WithCredentialsFile("rodinei-9f69e511f308.json")
	//fsClient, err := firestore.NewClientWithDatabase(ctx, utils.AppEngineGetProjectID(), DatabaseProdutos, opt)
	//if err != nil {
	//	return err
	//}

	fsClient, err := firestore.NewClientWithDatabase(ctx, utils.AppEngineGetProjectID(), DatabaseProdutos)
	if err != nil {
		return err
	}
	defer fsClient.Close()

	collection := fsClient.Collection(CollectionProdutos)

	var docRef *firestore.DocumentRef

	if produto.ID != "" {
		docRef = collection.NewDoc() // Gera um ID automático
		produto.ID = docRef.ID       // Converte ID gerado para int, se necessário
	} else {
		docRef = collection.Doc(produto.ID) // Usa o ID fornecido
	}

	_, err = docRef.Set(ctx, produto)
	if err != nil {
		return err
	}

	return nil
}

// GetProdutos Buscar produtos com filtro
func GetProdutos(produtoFilter models.Produto) ([]models.Produto, error) {
	ctx := context.Background()
	fsClient, err := firestore.NewClientWithDatabase(ctx, utils.AppEngineGetProjectID(), DatabaseProdutos)
	if err != nil {
		return nil, err
	}
	defer fsClient.Close()

	collection := fsClient.Collection(CollectionProdutos)
	var produtos []models.Produto

	if produtoFilter.ID != "" {
		docRef := collection.Doc(produtoFilter.ID)
		docSnap, err := docRef.Get(ctx)
		if err != nil {
			return nil, err
		}
		var produto models.Produto
		err = docSnap.DataTo(&produto)
		if err != nil {
			return nil, err
		}
		produto.ID = produtoFilter.ID
		return []models.Produto{produto}, nil
	}

	query := collection.Query

	if produtoFilter.EmEstoque {
		query = query.Where("EmEstoque", "==", produtoFilter.EmEstoque)
	}
	if produtoFilter.CodigoDeBarras != "" {
		query = query.Where("CodigoDeBarras", "==", produtoFilter.CodigoDeBarras)
	}
	if produtoFilter.Nome != "" {
		query = query.Where("Nome", "==", produtoFilter.Nome)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var produto models.Produto
		err := doc.DataTo(&produto)
		if err != nil {
			return nil, err
		}
		produto.ID = doc.Ref.ID
		produtos = append(produtos, produto)
	}

	return produtos, nil
}
