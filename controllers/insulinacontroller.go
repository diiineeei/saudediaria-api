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
	"go.mongodb.org/mongo-driver/mongo"
)

func personParam(ctx *gin.Context) (string, bool) {
	person := strings.TrimSpace(ctx.Param("person"))
	if person == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "person parameter is required"})
		return "", false
	}
	return person, true
}

// GetPrescricao devolve a receita de insulina (404 se ainda não houver).
func GetPrescricao(ctx *gin.Context) {
	person, ok := personParam(ctx)
	if !ok {
		return
	}
	p, err := database.GetPrescricao(person)
	if errors.Is(err, mongo.ErrNoDocuments) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "receita não cadastrada"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, p)
}

// SavePrescricao cria ou substitui a receita de insulina.
func SavePrescricao(ctx *gin.Context) {
	person, ok := personParam(ctx)
	if !ok {
		return
	}
	var p models.Prescricao
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.Person = person
	if err := p.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.UpdatedAt = time.Now().UTC()
	saved, err := database.SavePrescricao(p)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, saved)
}

// validaAplicacao confere o corpo e devolve os campos normalizados.
func validaAplicacao(ctx *gin.Context, input models.AplicacaoInput) (date, hhmm string, ok bool) {
	kind := strings.TrimSpace(input.Kind)
	if kind != "basal" && kind != "rapida" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "kind deve ser basal ou rapida"})
		return "", "", false
	}
	if !models.IsValidMomentoInsulina(strings.TrimSpace(input.Moment)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "momento inválido"})
		return "", "", false
	}
	if input.Units <= 0 || input.Units > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "units deve estar entre 0 e 100"})
		return "", "", false
	}
	t, err := time.Parse("15:04", strings.TrimSpace(input.Time))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "time must be in HH:MM format"})
		return "", "", false
	}
	if d := strings.TrimSpace(input.Date); d != "" {
		parsed, err := time.Parse("2006-01-02", d)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "date must be in YYYY-MM-DD format"})
			return "", "", false
		}
		date = parsed.Format("2006-01-02")
	}
	return date, t.Format("15:04"), true
}

// CreateAplicacao registra uma dose de insulina aplicada.
func CreateAplicacao(ctx *gin.Context) {
	person, ok := personParam(ctx)
	if !ok {
		return
	}
	var input models.AplicacaoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	date, hhmm, ok := validaAplicacao(ctx, input)
	if !ok {
		return
	}
	now := time.Now().In(saoPaulo())
	if date == "" {
		date = now.Format("2006-01-02")
	}
	a := models.AplicacaoInsulina{
		Person:    person,
		Kind:      strings.TrimSpace(input.Kind),
		Name:      strings.TrimSpace(input.Name),
		Units:     input.Units,
		Moment:    strings.TrimSpace(input.Moment),
		Date:      date,
		Time:      hhmm,
		Glucose:   input.Glucose,
		Suggested: input.Suggested,
		CreatedAt: now.UTC(),
	}
	saved, err := database.InsertAplicacao(a)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, saved)
}

func ListAplicacoes(ctx *gin.Context) {
	person, ok := personParam(ctx)
	if !ok {
		return
	}
	list, err := database.ListAplicacoes(person)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func UpdateAplicacao(ctx *gin.Context) {
	person, id, ok := parsePersonAndID(ctx)
	if !ok {
		return
	}
	var input models.AplicacaoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	date, hhmm, ok := validaAplicacao(ctx, input)
	if !ok {
		return
	}
	set := bson.M{
		"kind":   strings.TrimSpace(input.Kind),
		"name":   strings.TrimSpace(input.Name),
		"units":  input.Units,
		"moment": strings.TrimSpace(input.Moment),
		"time":   hhmm,
	}
	if date != "" {
		set["date"] = date
	}
	updated, err := database.UpdateAplicacao(person, id, set)
	if errors.Is(err, mongo.ErrNoDocuments) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "aplicação não encontrada"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

func DeleteAplicacao(ctx *gin.Context) {
	person, id, ok := parsePersonAndID(ctx)
	if !ok {
		return
	}
	err := database.DeleteAplicacao(person, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "aplicação não encontrada"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func saoPaulo() *time.Location {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.FixedZone("BRT", -3*3600)
	}
	return loc
}
