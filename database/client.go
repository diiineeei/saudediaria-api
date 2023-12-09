package database

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/diiineeei/ApiLojinha/models"
	//_ "database/sql"
	//_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

// var InstancePostgres *gorm.DB
//var dbError error
//var Client *firestore.Client
//
//func Database(connectionString string) {
//	var err error
//	ctx := context.Background()
//	Client, err = firestore.NewClient(ctx, "testes-rodinei")
//	if err != nil {
//		log.Fatalf("Error initializing Cloud Firestore client: %v", err)
//	}
//
//	//user := "postgres"
//	//password := "l;+t#N8@I:8DC#h5"
//	//dbHost := "testes-rodinei:us-central1:lojinha"
//	//databaseName := "lojinha"
//	//
//	//connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, user, databaseName, password)
//	//
//	//InstancePostgres, dbError = gorm.Open(postgres.New(postgres.Config{
//	//	DriverName: "cloudsqlpostgres",
//	//	DSN:        connString,
//	//}))
//
//	//InstancePostgres, dbError = gorm.Open(postgres.Open("postgres@cloudsql(testes-rodinei:us-central1:lojinha)/lojinha?charset=utf8&parseTime=True&loc=UTC"), &gorm.Config{})
//
//	//InstancePostgres, dbError = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
//	//if dbError != nil {
//	//	log.Fatal(dbError)
//	//	panic("Cannot connect to DB")
//	//}
//
//	log.Println("Connected to Database Postgres!")
//}

/*func Migrate() {
	err := InstancePostgres.AutoMigrate(&models.Produto{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database Migration Completed!")
}
*/

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
	key := datastore.IDKey(KindProdutos, produto.ID, nil)

	err = dsClient.Get(ctx, key, e)
	if err != nil && err.Error() != ErrnoSuchEntity {
		return err
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
			produtos = append(produtos, models.Produto{
				ID:             produtoFilter.ID,
				Nome:           e.Nome,
				EmEstoque:      e.EmEstoque,
				CodigoDeBarras: e.CodigoDeBarras,
				Descricao:      e.Descricao,
				ImagemURL:      e.ImagemURL,
				Valor:          e.Valor,
			})
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

	_, err = dsClient.GetAll(ctx, q, &produtos)
	if err != nil {
		return nil, err
	}

	if err != nil && err.Error() != ErrnoSuchEntity {
		return nil, err
	}

	return produtos, nil
}
