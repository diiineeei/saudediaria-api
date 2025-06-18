package database

import (
	"context"
	"github.com/diiineeei/ApiLojinha/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

const (
	CollectionRegistros = "registros"
	DatabaseName        = "saude"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
	once       sync.Once
)

func ConnectMongo(uri string) {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(uri)

		var err error
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Erro ao conectar MongoDB: %v", err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Erro no ping do MongoDB: %v", err)
		}

		collection = client.Database(DatabaseName).Collection(CollectionRegistros)
	})
}

// PutRegistro insere um registro de saúde diária
func PutRegistro(registro models.Registro) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, registro)
	return err
}

// GetRegistros busca registros conforme filtro (exemplo básico por data)
func GetRegistros(filtro models.Registro) ([]models.Registro, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	// Filtra por data se fornecida
	if filtro.Date != "" {
		filter["date"] = filtro.Date
	}

	// Se quiser pode filtrar por horário também, glicemia etc, ex:
	if filtro.Time != "" {
		filter["time"] = filtro.Time
	}

	findOptions := options.Find()
	findOptions.SetLimit(10)
	findOptions.SetSort(bson.D{{Key: "date", Value: -1}}) // ordena por data decrescente

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultados []models.Registro
	for cursor.Next(ctx) {
		var registro models.Registro
		if err := cursor.Decode(&registro); err != nil {
			return nil, err
		}
		resultados = append(resultados, registro)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return resultados, nil
}
