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

	now := time.Now()
	entry := models.GlicemiaRecord{
		Person:    person,
		Glucose:   input.Glucose,
		Date:      now.Format("2006-01-02"),
		Time:      parsedTime.Format("15:04"),
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
