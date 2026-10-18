package hospital

// Patient es un paciente costeño con narcolepsia.
// Todos sus campos son privados: solo se tocan mediante métodos.
type Patient struct {
	Person          // embedding: Patient reutiliza ID(), Name() y Age()
	level           NarcolepsyLevel
	state           PatientState
	currentLocation string
	assignedDoctor  *Doctor
	room            *Room
}

// NewPatient crea un paciente despierto.
func NewPatient(id, name string, age int, level NarcolepsyLevel) *Patient {
	return &Patient{
		Person: Person{id: id, name: name, age: age},
		level:  level,
		state:  Awake,
	}
}

// State devuelve el estado actual del paciente.
func (p *Patient) State() PatientState { return p.state }

// Level devuelve el nivel de narcolepsia.
func (p *Patient) Level() NarcolepsyLevel { return p.level }

// Location devuelve el último lugar donde estuvo el paciente.
func (p *Patient) Location() string { return p.currentLocation }

// AssignedDoctor devuelve el doctor que lo atiende (puede ser nil).
func (p *Patient) AssignedDoctor() *Doctor { return p.assignedDoctor }

// CurrentRoom devuelve la habitación donde está (nil si no tiene).
func (p *Patient) CurrentRoom() *Room { return p.room }

// SufferSleepAttack hace que el paciente se duerma de golpe en un lugar.
// Queda dormido en el pasillo hasta que se le asigne una habitación.
func (p *Patient) SufferSleepAttack(location string) {
	if p.state != Awake {
		return // ya estaba dormido, no se vuelve a dormir
	}
	p.state = AsleepInHallway
	p.currentLocation = location
}

// WakeUp despierta al paciente y libera su habitación si tenía una.
func (p *Patient) WakeUp() error {
	if p.state == Awake {
		return ErrAlreadyAwake
	}
	if p.room != nil {
		if err := p.room.Release(p); err != nil {
			return err
		}
	}
	p.state = Awake
	return nil
}
