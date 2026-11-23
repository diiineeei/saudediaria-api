package models

import "testing"

func receitaValida() Prescricao {
	return Prescricao{
		Basal: InsulinaBasal{Name: "NPH", Doses: []DoseFixa{{Moment: "antes_cafe", Units: 10}, {Moment: "22h", Units: 5}}},
		Rapida: InsulinaRapida{Name: "Rápida", MinutesBefore: 15, Scales: map[string][]FaixaDose{
			"antes_cafe": {{Min: 70, Max: 80, Units: 1}, {Min: 81, Max: 100, Units: 2}},
		}},
		Metas: MetasGlicemicas{Pre: 95, Pos: 140},
	}
}

func TestPrescricaoValidate(t *testing.T) {
	ok := receitaValida()
	if err := ok.Validate(); err != nil {
		t.Fatalf("receita válida recusada: %v", err)
	}

	casos := map[string]func(p *Prescricao){
		"momento basal inválido": func(p *Prescricao) { p.Basal.Doses[0].Moment = "lanche" },
		"dose basal negativa":    func(p *Prescricao) { p.Basal.Doses[0].Units = -1 },
		"momento rápida inválido": func(p *Prescricao) {
			p.Rapida.Scales["22h"] = []FaixaDose{{Min: 70, Max: 80, Units: 1}}
		},
		"faixa invertida": func(p *Prescricao) { p.Rapida.Scales["antes_cafe"][0].Max = 60 },
		"faixas sobrepostas": func(p *Prescricao) {
			p.Rapida.Scales["antes_cafe"][1].Min = 80
		},
		"dose rápida exagerada": func(p *Prescricao) { p.Rapida.Scales["antes_cafe"][0].Units = 500 },
	}
	for nome, muda := range casos {
		t.Run(nome, func(t *testing.T) {
			p := receitaValida()
			muda(&p)
			if err := p.Validate(); err == nil {
				t.Fatal("esperava erro de validação")
			}
		})
	}
}
