package hospital

import "fmt"

// Doctor es el personal médico: puede diagnosticar.
type Doctor struct {
	Person
	specialty string
	patients  []*Patient      // pacientes bajo su cuidado
	episodes  []EpisodeRecord // episodios que ha atendido
}

// NewDoctor crea un doctor.
func NewDoctor(id, name string, age int, specialty string) *Doctor {
	return &Doctor{
		Person:    Person{id: id, name: name, age: age},
		specialty: specialty,
	}
}

// Specialty devuelve la especialidad.
func (d *Doctor) Specialty() string { return d.specialty }

// DiagnosePatient deja al paciente bajo el cuidado del doctor.
func (d *Doctor) DiagnosePatient(p *Patient) error {
	if p == nil {
		return ErrNilPatient
	}
	p.assignedDoctor = d
	for _, q := range d.patients {
		if q == p {
			return nil // ya lo tenía
		}
	}
	d.patients = append(d.patients, p)
	return nil
}

// AttendSleepEmergency atiende la emergencia: el paciente debe estar dormido.
func (d *Doctor) AttendSleepEmergency(p *Patient) error {
	if p == nil {
		return ErrNilPatient
	}
	if p.state == Awake {
		return fmt.Errorf("paciente %s: %w", p.ID(), ErrPatientAwake)
	}
	return d.DiagnosePatient(p)
}

// Attend cumple la interfaz Attender: atiende y guarda el episodio.
func (d *Doctor) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if err := d.AttendSleepEmergency(p); err != nil {
		return EpisodeRecord{}, err
	}
	rec := newEpisodeRecord(p, d, location)
	d.episodes = append(d.episodes, rec)
	return rec, nil
}

// MyEpisodes devuelve (una copia de) los episodios en los que intervino.
func (d *Doctor) MyEpisodes() []EpisodeRecord {
	copia := make([]EpisodeRecord, len(d.episodes))
	copy(copia, d.episodes)
	return copia
}
