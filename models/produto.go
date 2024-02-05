package models

type Produto struct {
	ID             int64   `form:"ID" json:"ID" datastore:"-"`
	Nome           string  `form:"nome" json:"nome" datastore:"Nome"`
	Descricao      string  `form:"descricao" json:"descricao" datastore:"Descricao"`
	Valor          float64 `form:"valor" json:"valor" datastore:"Valor"`
	EmEstoque      bool    `form:"emEstoque" json:"emEstoque" datastore:"EmEstoque"`
	ImagemURL      string  `form:"imagemURL" json:"imagemURL" datastore:"ImagemURL"`
	CodigoDeBarras string  `form:"codigoDeBarras" json:"codigoDeBarras" datastore:"CodigoDeBarras"`
}
