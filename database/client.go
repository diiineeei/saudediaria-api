package database

import (
	"context"
	"github.com/diiineeei/ApiLojinha/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

const (
	CollectionRegistros = "registros"
	CollectionGlicemia  = "glicemia"
	DatabaseName        = "saude"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
	glicemiaCollection *mongo.Collection
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
		glicemiaCollection = client.Database(DatabaseName).Collection(CollectionGlicemia)
	})
}

// PutRegistro insere um registro de saúde diária
func PutRegistro(registro models.Registro) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, registro)
	return err
}

func InsertGlicemiaEntry(entry models.GlicemiaRecord) (models.GlicemiaRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := glicemiaCollection.InsertOne(ctx, entry)
	if err != nil {
		return entry, err
	}

	if insertedID, ok := result.InsertedID.(primitive.ObjectID); ok {
		entry.ID = insertedID
	}

	return entry, nil
}

func ListGlicemiaEntries(person string) ([]models.GlicemiaRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if person != "" {
		filter["person"] = person
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := glicemiaCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.GlicemiaRecord
	for cursor.Next(ctx) {
		var entry models.GlicemiaRecord
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		results = append(results, entry)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// ErrGlicemiaNotFound indica que a medição não existe para essa pessoa.
var ErrGlicemiaNotFound = mongo.ErrNoDocuments

// UpdateGlicemiaEntry altera uma medição e devolve o documento já atualizado.
func UpdateGlicemiaEntry(person string, id primitive.ObjectID, set bson.M, unset bson.M) (models.GlicemiaRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var updated models.GlicemiaRecord
	err := glicemiaCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id, "person": person},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)
	return updated, err
}

// DeleteGlicemiaEntry remove uma medição. Devolve ErrGlicemiaNotFound se não existir.
func DeleteGlicemiaEntry(person string, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := glicemiaCollection.DeleteOne(ctx, bson.M{"_id": id, "person": person})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrGlicemiaNotFound
	}
	return nil
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
