package models

import (
	"gorm.io/gorm"
)

type Produto struct {
	gorm.Model
	Nome      string  `json:"nome"`
	Descricao string  `json:"descricao"`
	Valor     float64 `json:"valor"`
	EmEstoque bool    `json:"emEstoque"`
	ImagemURL string  `json:"imagemURL"`
}
