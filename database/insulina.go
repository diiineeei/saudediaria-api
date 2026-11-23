package database

import (
	"context"
	"time"

	"github.com/diiineeei/ApiLojinha/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPrescricao devolve a receita da pessoa ou mongo.ErrNoDocuments se não houver.
func GetPrescricao(person string) (models.Prescricao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p models.Prescricao
	err := prescricaoCollection.FindOne(ctx, bson.M{"person": person}).Decode(&p)
	return p, err
}

// SavePrescricao grava (substitui) a receita da pessoa.
func SavePrescricao(p models.Prescricao) (models.Prescricao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := prescricaoCollection.ReplaceOne(ctx, bson.M{"person": p.Person}, p, options.Replace().SetUpsert(true))
	return p, err
}

func InsertAplicacao(a models.AplicacaoInsulina) (models.AplicacaoInsulina, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := insulinaCollection.InsertOne(ctx, a)
	if err != nil {
		return a, err
	}
	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		a.ID = id
	}
	return a, nil
}

func ListAplicacoes(person string) ([]models.AplicacaoInsulina, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := insulinaCollection.Find(ctx, bson.M{"person": person},
		options.Find().SetSort(bson.D{{Key: "date", Value: -1}, {Key: "time", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := []models.AplicacaoInsulina{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateAplicacao altera uma aplicação e devolve o documento atualizado.
func UpdateAplicacao(person string, id primitive.ObjectID, set bson.M) (models.AplicacaoInsulina, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var updated models.AplicacaoInsulina
	err := insulinaCollection.FindOneAndUpdate(ctx, bson.M{"_id": id, "person": person}, bson.M{"$set": set},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updated)
	return updated, err
}

func DeleteAplicacao(person string, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := insulinaCollection.DeleteOne(ctx, bson.M{"_id": id, "person": person})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
