package models

type Produto struct {
	ID             string  `json:"id,omitempty"` // Firestore armazena IDs como string, lidamos com isso na conversão
	Nome           string  `json:"nome"`
	Descricao      string  `json:"descricao"`
	Valor          float64 `json:"valor"`
	EmEstoque      bool    `json:"emEstoque"`
	ImagemURL      string  `json:"imagemURL"`
	CodigoDeBarras string  `json:"codigoDeBarras"`
}
