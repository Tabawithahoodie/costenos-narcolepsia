package hospital

import (
	"fmt"
	"time"
)

// Hospital es el sistema central: guarda personal, pacientes, habitaciones
// y el historial de episodios.
type Hospital struct {
	name      string
	doctors   []*Doctor
	staff     []Attender // doctores y camilleros juntos (polimorfismo)
	patients  []*Patient
	rooms     []*Room
	history   []EpisodeRecord
	nextStaff int // por dónde sigue el turno de despacho
}

// NewHospital crea un hospital vacío.
func NewHospital(name string) *Hospital {
	return &Hospital{name: name}
}

// Name devuelve el nombre del hospital.
func (h *Hospital) Name() string { return h.name }

// History devuelve una copia de todos los episodios.
func (h *Hospital) History() []EpisodeRecord {
	copia := make([]EpisodeRecord, len(h.history))
	copy(copia, h.history)
	return copia
}

// hasStaff dice si ya hay alguien del personal con ese id.
func (h *Hospital) hasStaff(id string) bool {
	for _, a := range h.staff {
		if a.ID() == id {
			return true
		}
	}
	return false
}

// hasPatient dice si el paciente está admitido.
func (h *Hospital) hasPatient(p *Patient) bool {
	for _, q := range h.patients {
		if q == p {
			return true
		}
	}
	return false
}

// HireDoctor contrata un doctor.
func (h *Hospital) HireDoctor(d *Doctor) error {
	if d == nil {
		return ErrNilValue
	}
	if h.hasStaff(d.ID()) {
		return fmt.Errorf("personal %s: %w", d.ID(), ErrDuplicate)
	}
	h.doctors = append(h.doctors, d)
	h.staff = append(h.staff, d)
	return nil
}

// HireOrderly contrata un camillero.
func (h *Hospital) HireOrderly(o *Orderly) error {
	if o == nil {
		return ErrNilValue
	}
	if h.hasStaff(o.ID()) {
		return fmt.Errorf("personal %s: %w", o.ID(), ErrDuplicate)
	}
	h.staff = append(h.staff, o)
	return nil
}

// AddRoom registra una habitación.
func (h *Hospital) AddRoom(r *Room) error {
	if r == nil {
		return ErrNilValue
	}
	for _, x := range h.rooms {
		if x.Number() == r.Number() {
			return fmt.Errorf("habitación %d: %w", r.Number(), ErrDuplicate)
		}
	}
	h.rooms = append(h.rooms, r)
	return nil
}

// AdmitPatient admite un paciente.
func (h *Hospital) AdmitPatient(p *Patient) error {
	if p == nil {
		return ErrNilPatient
	}
	for _, q := range h.patients {
		if q.ID() == p.ID() {
			return fmt.Errorf("paciente %s: %w", p.ID(), ErrDuplicate)
		}
	}
	h.patients = append(h.patients, p)
	return nil
}

// AssignRoom busca la primera habitación disponible y se la da al paciente.
// Si no hay ninguna devuelve error y el paciente sigue en el pasillo.
func (h *Hospital) AssignRoom(p *Patient) (*Room, error) {
	if p == nil {
		return nil, ErrNilPatient
	}
	switch p.State() {
	case Awake:
		return nil, fmt.Errorf("paciente %s: %w", p.ID(), ErrPatientAwake)
	case AsleepInBed:
		return nil, fmt.Errorf("paciente %s: %w", p.ID(), ErrAlreadyInBed)
	}
	for _, r := range h.rooms {
		if r.IsAvailable() {
			if err := r.Occupy(p); err != nil {
				return nil, err
			}
			return r, nil
		}
	}
	return nil, fmt.Errorf("paciente %s se queda en el pasillo: %w", p.ID(), ErrNoRoom)
}

// RegisterEpisode guarda un episodio en el historial del hospital.
func (h *Hospital) RegisterEpisode(rec EpisodeRecord) error {
	if rec.patient == nil || rec.attender == nil {
		return ErrInvalidRecord
	}
	h.history = append(h.history, rec)
	return nil
}

// dispatch le pasa el paciente al personal por turnos. No sabe si quien
// atiende es Doctor u Orderly: solo usa la interfaz Attender.
func (h *Hospital) dispatch(p *Patient, location string) (EpisodeRecord, error) {
	for i := 0; i < len(h.staff); i++ {
		idx := (h.nextStaff + i) % len(h.staff)
		rec, err := h.staff[idx].Attend(p, location)
		if err == nil {
			h.nextStaff = (idx + 1) % len(h.staff)
			return rec, nil
		}
	}
	return EpisodeRecord{}, ErrNoStaff
}

// ReportSleepAttack junta todo el flujo: el paciente se duerme, se intenta
// darle habitación, se despacha personal y se registra el episodio.
// Si no hay habitación, el episodio igual queda registrado y se devuelve
// el error de la habitación.
func (h *Hospital) ReportSleepAttack(p *Patient, location string) (*Room, error) {
	if p == nil {
		return nil, ErrNilPatient
	}
	if !h.hasPatient(p) {
		return nil, fmt.Errorf("paciente %s: %w", p.ID(), ErrPatientNotKnown)
	}
	p.SufferSleepAttack(location)
	room, roomErr := h.AssignRoom(p)

	rec, err := h.dispatch(p, location)
	if err != nil {
		return room, err
	}
	if err := h.RegisterEpisode(rec); err != nil {
		return room, err
	}
	return room, roomErr
}

// PatientsInHallway (consulta 5.1) devuelve los pacientes dormidos en pasillo.
func (h *Hospital) PatientsInHallway() []*Patient {
	var result []*Patient
	for _, p := range h.patients {
		if p.State() == AsleepInHallway {
			result = append(result, p)
		}
	}
	return result
}

// SevereReport (consulta 5.4) cruza los pacientes Severe con sus episodios
// del día de hoy.
func (h *Hospital) SevereReport() map[*Patient]int {
	report := make(map[*Patient]int)
	for _, p := range h.patients {
		if p.Level() == Severe {
			report[p] = 0
		}
	}
	y, m, d := time.Now().Date()
	for _, rec := range h.history {
		ry, rm, rd := rec.DateTime().Date()
		if ry != y || rm != m || rd != d {
			continue
		}
		if _, esSevero := report[rec.Patient()]; esSevero {
			report[rec.Patient()]++
		}
	}
	return report
}
