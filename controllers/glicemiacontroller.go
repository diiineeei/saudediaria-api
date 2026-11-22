package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/diiineeei/ApiLojinha/database"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"
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
