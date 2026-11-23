package models

import (
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Momentos de aplicação de insulina.
const (
	MomentoAntesCafe   = "antes_cafe"
	MomentoAntesAlmoco = "antes_almoco"
	MomentoAntesJantar = "antes_jantar"
	MomentoNoite22h    = "22h"
)

var momentosInsulina = map[string]bool{
	MomentoAntesCafe:   true,
	MomentoAntesAlmoco: true,
	MomentoAntesJantar: true,
	MomentoNoite22h:    true,
}

// Momentos em que a insulina rápida é calculada pela glicemia (dextro).
var momentosEscala = map[string]bool{
	MomentoAntesCafe:   true,
	MomentoAntesAlmoco: true,
	MomentoAntesJantar: true,
}

// IsValidMomentoInsulina informa se o momento de aplicação é aceito.
func IsValidMomentoInsulina(m string) bool { return momentosInsulina[m] }

// DoseFixa é uma dose que não depende da glicemia (ex.: insulina NPH).
type DoseFixa struct {
	Moment string  `bson:"moment" json:"moment"`
	Units  float64 `bson:"units" json:"units"`
}

// FaixaDose: se a glicemia estiver entre Min e Max (inclusive), aplicar Units.
type FaixaDose struct {
	Min   float64 `bson:"min" json:"min"`
	Max   float64 `bson:"max" json:"max"`
	Units float64 `bson:"units" json:"units"`
}

type InsulinaBasal struct {
	Name  string     `bson:"name" json:"name"`
	Doses []DoseFixa `bson:"doses" json:"doses"`
}

type InsulinaRapida struct {
	Name          string                 `bson:"name" json:"name"`
	MinutesBefore int                    `bson:"minutesBefore" json:"minutesBefore"`
	Scales        map[string][]FaixaDose `bson:"scales" json:"scales"`
}

// MetasGlicemicas: alvo antes das refeições (pre) e 1h depois (pos), em mg/dL.
type MetasGlicemicas struct {
	Pre float64 `bson:"pre" json:"pre"`
	Pos float64 `bson:"pos" json:"pos"`
}

// Prescricao guarda a receita de insulina de uma pessoa (um documento por pessoa).
type Prescricao struct {
	Person    string          `bson:"person" json:"person"`
	Basal     InsulinaBasal   `bson:"basal" json:"basal"`
	Rapida    InsulinaRapida  `bson:"rapida" json:"rapida"`
	Metas     MetasGlicemicas `bson:"metas" json:"metas"`
	Doctor    string          `bson:"doctor,omitempty" json:"doctor,omitempty"`
	IssuedAt  string          `bson:"issuedAt,omitempty" json:"issuedAt,omitempty"`
	Notes     string          `bson:"notes,omitempty" json:"notes,omitempty"`
	UpdatedAt time.Time       `bson:"updatedAt" json:"updatedAt"`
}

// Validate confere se a receita é coerente antes de gravar.
func (p *Prescricao) Validate() error {
	for _, d := range p.Basal.Doses {
		if !IsValidMomentoInsulina(d.Moment) {
			return fmt.Errorf("momento inválido na insulina basal: %q", d.Moment)
		}
		if d.Units < 0 || d.Units > 100 {
			return fmt.Errorf("dose basal fora do intervalo 0–100: %v", d.Units)
		}
	}
	for moment, faixas := range p.Rapida.Scales {
		if !momentosEscala[moment] {
			return fmt.Errorf("momento inválido na insulina rápida: %q", moment)
		}
		ordenadas := append([]FaixaDose(nil), faixas...)
		sort.Slice(ordenadas, func(i, j int) bool { return ordenadas[i].Min < ordenadas[j].Min })
		for i, f := range ordenadas {
			if f.Min <= 0 || f.Max < f.Min {
				return fmt.Errorf("faixa inválida em %s: %v–%v", moment, f.Min, f.Max)
			}
			if f.Units < 0 || f.Units > 100 {
				return fmt.Errorf("dose fora do intervalo 0–100 em %s: %v", moment, f.Units)
			}
			if i > 0 && f.Min <= ordenadas[i-1].Max {
				return fmt.Errorf("faixas sobrepostas em %s", moment)
			}
		}
	}
	if p.Metas.Pre < 0 || p.Metas.Pos < 0 {
		return fmt.Errorf("metas inválidas")
	}
	return nil
}

// AplicacaoInsulina registra uma dose aplicada.
type AplicacaoInsulina struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Person    string             `bson:"person" json:"person"`
	Kind      string             `bson:"kind" json:"kind"` // "basal" ou "rapida"
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	Units     float64            `bson:"units" json:"units"`
	Moment    string             `bson:"moment" json:"moment"`
	Date      string             `bson:"date" json:"date"`
	Time      string             `bson:"time" json:"time"`
	Glucose   *float64           `bson:"glucose,omitempty" json:"glucose,omitempty"`     // dextro usado no cálculo
	Suggested *float64           `bson:"suggested,omitempty" json:"suggested,omitempty"` // dose sugerida pela receita
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// AplicacaoInput é o corpo do POST/PUT de aplicação de insulina.
type AplicacaoInput struct {
	Kind      string   `json:"kind" binding:"required"`
	Name      string   `json:"name"`
	Units     float64  `json:"units"`
	Moment    string   `json:"moment" binding:"required"`
	Date      string   `json:"date"`
	Time      string   `json:"time" binding:"required"`
	Glucose   *float64 `json:"glucose"`
	Suggested *float64 `json:"suggested"`
}
