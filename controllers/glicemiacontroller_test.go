package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Validações que respondem 400 antes de tocar no banco.
func TestUpdateDeleteGlicemiaValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/glicemia/:person/list", ListGlicemiaEntries)
	r.PUT("/api/glicemia/:person/:id", UpdateGlicemiaEntry)
	r.DELETE("/api/glicemia/:person/:id", DeleteGlicemiaEntry)

	const validID = "64b7f0c2a1b2c3d4e5f60718"
	cases := []struct {
		name, method, path, body string
	}{
		{"id inválido no PUT", http.MethodPut, "/api/glicemia/mary/abc", `{"glucose":100,"time":"08:00"}`},
		{"id inválido no DELETE", http.MethodDelete, "/api/glicemia/mary/abc", ""},
		{"sem glicemia", http.MethodPut, "/api/glicemia/mary/" + validID, `{"time":"08:00"}`},
		{"glicemia zero", http.MethodPut, "/api/glicemia/mary/" + validID, `{"glucose":0,"time":"08:00"}`},
		{"horário inválido", http.MethodPut, "/api/glicemia/mary/" + validID, `{"glucose":100,"time":"8h"}`},
		{"data inválida", http.MethodPut, "/api/glicemia/mary/" + validID, `{"glucose":100,"time":"08:00","date":"23/11/2026"}`},
		{"momento inválido", http.MethodPut, "/api/glicemia/mary/" + validID, `{"glucose":100,"time":"08:00","period":"lanche"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("esperava 400, recebeu %d: %s", w.Code, w.Body.String())
			}
		})
	}
}
