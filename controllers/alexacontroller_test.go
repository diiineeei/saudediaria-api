package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diiineeei/ApiLojinha/models"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/alexa/webhook", AlexaWebhook)
	return r
}

func TestAlexaLaunchRequest(t *testing.T) {
	r := setupTestRouter()

	reqPayload := models.AlexaRequest{
		Version: "1.0",
		Request: models.AlexaPayload{
			Type: "LaunchRequest",
		},
	}

	body, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/alexa/webhook", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperava status 200, recebeu %d", w.Code)
	}

	var resp models.AlexaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Falha ao deserializar resposta: %v", err)
	}

	if resp.Response.OutputSpeech == nil || resp.Response.OutputSpeech.Text == "" {
		t.Fatalf("Esperava fala na resposta da Alexa")
	}

	if resp.Response.ShouldEndSession {
		t.Fatalf("No LaunchRequest, a sessão não deve ser encerrada imediatamente")
	}
}

func TestAlexaHelpIntent(t *testing.T) {
	r := setupTestRouter()

	reqPayload := models.AlexaRequest{
		Version: "1.0",
		Request: models.AlexaPayload{
			Type: "IntentRequest",
			Intent: models.AlexaIntent{
				Name: "AMAZON.HelpIntent",
			},
		},
	}

	body, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/alexa/webhook", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperava status 200, recebeu %d", w.Code)
	}

	var resp models.AlexaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Falha ao deserializar resposta: %v", err)
	}

	if resp.Response.OutputSpeech == nil || resp.Response.OutputSpeech.Text == "" {
		t.Fatalf("Esperava fala de ajuda")
	}
}
