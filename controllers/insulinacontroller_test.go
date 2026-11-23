package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Validações que respondem 400 antes de tocar no banco.
func TestInsulinaValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/api/prescricao/:person", SavePrescricao)
	r.POST("/api/insulina/:person", CreateAplicacao)
	r.PUT("/api/insulina/:person/:id", UpdateAplicacao)
	r.DELETE("/api/insulina/:person/:id", DeleteAplicacao)

	const id = "64b7f0c2a1b2c3d4e5f60718"
	cases := []struct{ name, method, path, body string }{
		{"receita com faixa sobreposta", http.MethodPut, "/api/prescricao/mary",
			`{"rapida":{"scales":{"antes_cafe":[{"min":70,"max":90,"units":1},{"min":85,"max":100,"units":2}]}}}`},
		{"tipo inválido", http.MethodPost, "/api/insulina/mary", `{"kind":"x","moment":"antes_cafe","units":5,"time":"08:00"}`},
		{"momento inválido", http.MethodPost, "/api/insulina/mary", `{"kind":"rapida","moment":"lanche","units":5,"time":"08:00"}`},
		{"sem unidades", http.MethodPost, "/api/insulina/mary", `{"kind":"rapida","moment":"antes_cafe","time":"08:00"}`},
		{"unidades demais", http.MethodPost, "/api/insulina/mary", `{"kind":"basal","moment":"22h","units":250,"time":"22:00"}`},
		{"horário inválido", http.MethodPost, "/api/insulina/mary", `{"kind":"basal","moment":"22h","units":17,"time":"22h"}`},
		{"data inválida", http.MethodPost, "/api/insulina/mary", `{"kind":"basal","moment":"22h","units":17,"time":"22:00","date":"23/11/2026"}`},
		{"id inválido no PUT", http.MethodPut, "/api/insulina/mary/abc", `{"kind":"basal","moment":"22h","units":17,"time":"22:00"}`},
		{"id inválido no DELETE", http.MethodDelete, "/api/insulina/mary/abc", ""},
		{"PUT com momento inválido", http.MethodPut, "/api/insulina/mary/" + id, `{"kind":"basal","moment":"x","units":17,"time":"22:00"}`},
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
