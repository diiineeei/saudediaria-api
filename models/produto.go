package models

type Produto struct {
	ID             int64   `json:"ID" datastore:"-"`
	Nome           string  `json:"nome" datastore:"Nome"`
	Descricao      string  `json:"descricao" datastore:"Descricao"`
	Valor          float64 `json:"valor" datastore:"Valor"`
	EmEstoque      bool    `json:"emEstoque" datastore:"EmEstoque"`
	ImagemURL      string  `json:"imagemURL" datastore:"ImagemURL"`
	CodigoDeBarras string  `json:"codigoDeBarras" datastore:"CodigoDeBarras"`
}
