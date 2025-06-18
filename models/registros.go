package models

type Registro struct {
	PressureSys   string `json:"pressureSys"`
	PressureDia   string `json:"pressureDia"`
	PressurePulse string `json:"pressurePulse"`
	Glucose       string `json:"glucose"`
	Date          string `json:"date"`
	Time          string `json:"time"`
}
