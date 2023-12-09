package models

type Produto struct {
	ID             int64   `json:"ID" datastore:"-"`
	Nome           string  `json:"nome" datastore:"nome"`
	Descricao      string  `json:"descricao" datastore:"descricao"`
	Valor          float64 `json:"valor" datastore:"valor"`
	EmEstoque      bool    `json:"emEstoque" datastore:"emEstoque"`
	ImagemURL      string  `json:"imagemURL" datastore:"imagemURL"`
	CodigoDeBarras string  `json:"codigoDeBarras" datastore:"codigoDeBarras"`
}
