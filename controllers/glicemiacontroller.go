package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/diiineeei/ApiLojinha/database"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateGlicemiaEntry(ctx *gin.Context) {
	person := strings.TrimSpace(ctx.Param("person"))
	if person == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "person parameter is required"})
		return
	}

	var input models.GlicemiaInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedTime, err := time.Parse("15:04", input.Time)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "time must be in HH:MM format"})
		return
	}

	period := strings.TrimSpace(input.Period)
	if !models.IsValidPeriod(period) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "period inválido"})
		return
	}

	// Usa o fuso de São Paulo para que a data não vire o dia seguinte à noite (servidor em UTC).
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		loc = time.FixedZone("BRT", -3*3600)
	}
	now := time.Now().In(loc)
	entry := models.GlicemiaRecord{
		Person:    person,
		Glucose:   input.Glucose,
		Date:      now.Format("2006-01-02"),
		Time:      parsedTime.Format("15:04"),
		Period:    period,
		CreatedAt: now.UTC(),
	}

	savedEntry, err := database.InsertGlicemiaEntry(entry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, savedEntry)
}

func ListGlicemiaEntries(ctx *gin.Context) {
	person := strings.TrimSpace(ctx.Param("person"))
	if person == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "person parameter is required"})
		return
	}

	entries, err := database.ListGlicemiaEntries(person)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

// parsePersonAndID valida os parâmetros :person e :id da rota.
func parsePersonAndID(ctx *gin.Context) (string, primitive.ObjectID, bool) {
	person := strings.TrimSpace(ctx.Param("person"))
	if person == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "person parameter is required"})
		return "", primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(ctx.Param("id")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return "", primitive.NilObjectID, false
	}
	return person, id, true
}

// UpdateGlicemiaEntry edita valor, horário, data e momento de uma medição.
func UpdateGlicemiaEntry(ctx *gin.Context) {
	person, id, ok := parsePersonAndID(ctx)
	if !ok {
		return
	}

	var input models.GlicemiaUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Glucose <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "glucose deve ser maior que zero"})
		return
	}

	parsedTime, err := time.Parse("15:04", strings.TrimSpace(input.Time))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "time must be in HH:MM format"})
		return
	}

	period := strings.TrimSpace(input.Period)
	if !models.IsValidPeriod(period) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "period inválido"})
		return
	}

	set := bson.M{"glucose": input.Glucose, "time": parsedTime.Format("15:04")}
	if date := strings.TrimSpace(input.Date); date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "date must be in YYYY-MM-DD format"})
			return
		}
		set["date"] = parsedDate.Format("2006-01-02")
	}

	unset := bson.M{}
	if period == "" {
		unset["period"] = ""
	} else {
		set["period"] = period
	}

	updated, err := database.UpdateGlicemiaEntry(person, id, set, unset)
	if errors.Is(err, database.ErrGlicemiaNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "medição não encontrada"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

// DeleteGlicemiaEntry apaga uma medição.
func DeleteGlicemiaEntry(ctx *gin.Context) {
	person, id, ok := parsePersonAndID(ctx)
	if !ok {
		return
	}

	err := database.DeleteGlicemiaEntry(person, id)
	if errors.Is(err, database.ErrGlicemiaNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "medição não encontrada"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
