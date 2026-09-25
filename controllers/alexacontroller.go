package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/diiineeei/ApiLojinha/database"
	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"
)

const TargetPerson = "mary"

// AlexaWebhook processa as requisições enviadas pela Alexa
func AlexaWebhook(c *gin.Context) {
	var req models.AlexaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Alexa] Erro ao decodificar JSON: %v", err)
		c.JSON(http.StatusOK, models.BuildAlexaResponse("Desculpe, ocorreu um erro ao processar o pedido.", true))
		return
	}

	switch req.Request.Type {
	case "LaunchRequest":
		c.JSON(http.StatusOK, models.BuildAlexaRepromptResponse(
			"Olá! Diga o valor da glicemia da Mary para registrar.",
			"Por favor, diga o valor da glicemia, por exemplo: cento e dez.",
		))
		return

	case "IntentRequest":
		intentName := req.Request.Intent.Name
		log.Printf("[Alexa] Intent recebida: %s", intentName)

		switch intentName {
		case "RegistrarGlicemiaIntent":
			slot, exists := req.Request.Intent.Slots["valor"]
			if !exists || strings.TrimSpace(slot.Value) == "" {
				c.JSON(http.StatusOK, models.BuildAlexaRepromptResponse(
					"Não entendi o valor. Qual é a glicemia da Mary?",
					"Por favor, diga o valor da glicemia.",
				))
				return
			}

			// Suporta vírgula decimal como "110,5"
			valClean := strings.ReplaceAll(slot.Value, ",", ".")
			glucose, err := strconv.ParseFloat(valClean, 64)
			if err != nil || glucose <= 0 {
				c.JSON(http.StatusOK, models.BuildAlexaRepromptResponse(
					"Valor inválido. Qual é o número da glicemia da Mary?",
					"Por favor, diga apenas o número da glicemia.",
				))
				return
			}

			loc, err := time.LoadLocation("America/Sao_Paulo")
			if err != nil {
				loc = time.FixedZone("BRT", -3*3600)
			}
			now := time.Now().In(loc)

			entry := models.GlicemiaRecord{
				Person:    TargetPerson,
				Glucose:   glucose,
				Date:      now.Format("2006-01-02"),
				Time:      now.Format("15:04"),
				CreatedAt: now.UTC(),
			}

			savedEntry, err := database.InsertGlicemiaEntry(entry)
			if err != nil {
				log.Printf("[Alexa] Erro ao salvar glicemia no MongoDB: %v", err)
				c.JSON(http.StatusOK, models.BuildAlexaResponse("Desculpe, não consegui salvar a glicemia no momento.", true))
				return
			}

			log.Printf("[Alexa] Glicemia salva com sucesso: %+v", savedEntry)

			var speech string
			if glucose == float64(int64(glucose)) {
				speech = fmt.Sprintf("Anotado! Glicemia de %.0f da Mary registrada às %s.", glucose, entry.Time)
			} else {
				speech = fmt.Sprintf("Anotado! Glicemia de %.1f da Mary registrada às %s.", glucose, entry.Time)
			}

			c.JSON(http.StatusOK, models.BuildAlexaResponse(speech, true))
			return

		case "ConsultarGlicemiaIntent":
			entries, err := database.ListGlicemiaEntries(TargetPerson)
			if err != nil || len(entries) == 0 {
				c.JSON(http.StatusOK, models.BuildAlexaResponse("Não encontrei nenhum registro de glicemia para a Mary.", true))
				return
			}

			last := entries[0]
			var speech string
			if last.Glucose == float64(int64(last.Glucose)) {
				speech = fmt.Sprintf("A última glicemia da Mary foi %.0f, registrada às %s.", last.Glucose, last.Time)
			} else {
				speech = fmt.Sprintf("A última glicemia da Mary foi %.1f, registrada às %s.", last.Glucose, last.Time)
			}

			c.JSON(http.StatusOK, models.BuildAlexaResponse(speech, true))
			return

		case "AMAZON.HelpIntent":
			c.JSON(http.StatusOK, models.BuildAlexaRepromptResponse(
				"Você pode dizer: anota glicemia 110, ou perguntar: qual a última glicemia da Mary.",
				"Diga o valor da glicemia da Mary para registrar.",
			))
			return

		case "AMAZON.CancelIntent", "AMAZON.StopIntent":
			c.JSON(http.StatusOK, models.BuildAlexaResponse("Até logo!", true))
			return

		default:
			c.JSON(http.StatusOK, models.BuildAlexaResponse("Desculpe, não reconheci esse comando.", true))
			return
		}

	case "SessionEndedRequest":
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return

	default:
		c.JSON(http.StatusOK, models.BuildAlexaResponse("Requisição desconhecida.", true))
		return
	}
}
