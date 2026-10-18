package hospital

import "fmt"

// Orderly (camillero) mueve pacientes dormidos pero NO puede diagnosticar,
// por eso no tiene DiagnosePatient.
type Orderly struct {
	Person
}

// NewOrderly crea un camillero.
func NewOrderly(id, name string, age int) *Orderly {
	return &Orderly{Person: Person{id: id, name: name, age: age}}
}

// Attend cumple la interfaz Attender: registra que movió al paciente.
func (o *Orderly) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if p == nil {
		return EpisodeRecord{}, ErrNilPatient
	}
	if p.state == Awake {
		return EpisodeRecord{}, fmt.Errorf("paciente %s: %w", p.ID(), ErrPatientAwake)
	}
	return newEpisodeRecord(p, o, location), nil
}
