package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Períodos aceitos para a medição de glicemia. O período é opcional.
const (
	PeriodoAntesCafe   = "antes_cafe"
	PeriodoPosCafe     = "pos_cafe"
	PeriodoAntesAlmoco = "antes_almoco"
	PeriodoPosAlmoco   = "pos_almoco"
	PeriodoAntesJantar = "antes_jantar"
	PeriodoPosJantar   = "pos_jantar"
)

var periodosValidos = map[string]bool{
	PeriodoAntesCafe:   true,
	PeriodoPosCafe:     true,
	PeriodoAntesAlmoco: true,
	PeriodoPosAlmoco:   true,
	PeriodoAntesJantar: true,
	PeriodoPosJantar:   true,
}

// IsValidPeriod informa se o período é vazio (não informado) ou um dos valores aceitos.
func IsValidPeriod(period string) bool {
	return period == "" || periodosValidos[period]
}

type GlicemiaInput struct {
	Glucose float64 `json:"glucose" binding:"required"`
	Time    string  `json:"time" binding:"required"`
	Period  string  `json:"period"`
}

// GlicemiaUpdateInput é o corpo do PUT /api/glicemia/:person/:id.
// Date é opcional (mantém a data atual se vier vazia); Period vazio limpa o momento.
type GlicemiaUpdateInput struct {
	Glucose float64 `json:"glucose" binding:"required"`
	Time    string  `json:"time" binding:"required"`
	Date    string  `json:"date"`
	Period  string  `json:"period"`
}

type GlicemiaRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Person    string             `bson:"person" json:"person"`
	Glucose   float64            `bson:"glucose" json:"glucose"`
	Date      string             `bson:"date" json:"date"`
	Time      string             `bson:"time" json:"time"`
	Period    string             `bson:"period,omitempty" json:"period,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
