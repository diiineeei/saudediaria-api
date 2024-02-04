package database

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/diiineeei/ApiLojinha/models"
)

const ProjectID = "adega-do-vilao"
const KindProdutos = "produtos"
const ErrnoSuchEntity = "datastore: no such entity"

func PutProdutos(produto models.Produto) error {
	ctx := context.Background()

	dsClient, err := datastore.NewClient(ctx, ProjectID)
	if err != nil {
		return err
	}
	defer dsClient.Close()

	e := new(models.Produto)
	var key *datastore.Key

	key = datastore.IDKey(KindProdutos, produto.ID, nil)

	if produto.ID == 0 {
		key = datastore.IncompleteKey(KindProdutos, nil)
	}

	if produto.ID != 0 {
		err = dsClient.Get(ctx, key, e)
		if err != nil && err.Error() != ErrnoSuchEntity {
			return err
		}
	}

	e.Nome = produto.Nome
	e.Descricao = produto.Descricao
	e.Valor = produto.Valor
	e.EmEstoque = produto.EmEstoque
	e.ImagemURL = produto.ImagemURL
	e.CodigoDeBarras = produto.CodigoDeBarras

	if _, err := dsClient.Put(ctx, key, e); err != nil {
		return err
	}

	return nil
}

func GetProdutos(produtoFilter models.Produto) ([]models.Produto, error) {
	var produtos []models.Produto
	ctx := context.Background()

	dsClient, err := datastore.NewClient(ctx, ProjectID)
	if err != nil {
		return nil, err
	}
	defer dsClient.Close()

	key := datastore.IDKey(KindProdutos, produtoFilter.ID, nil)

	if produtoFilter.ID != 0 {
		e := new(models.Produto)
		err = dsClient.Get(ctx, key, e)
		if err != nil && err.Error() != ErrnoSuchEntity {
			return nil, err
		}
		if e != nil {
			e.ID = key.ID
			produtos = append(produtos, *e)
		}
	}

	q := datastore.NewQuery(KindProdutos)
	if produtoFilter.EmEstoque {
		q.FilterField("EmEstoque", "=", produtoFilter.EmEstoque)
	}
	if produtoFilter.CodigoDeBarras != "" {
		q.FilterField("CodigoDeBarras", "=", produtoFilter.CodigoDeBarras)
	}
	if produtoFilter.Nome != "" {
		q.FilterField("Nome", "=", produtoFilter.Nome)
	}

	keys, err := dsClient.GetAll(ctx, q, &produtos)
	if err != nil {
		return nil, err
	}

	for i := range produtos {
		produtos[i].ID = keys[i].ID
	}

	if err != nil && err.Error() != ErrnoSuchEntity {
		return nil, err
	}

	return produtos, nil
}
