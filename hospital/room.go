package hospital

// Room es una habitación con cierta capacidad de camas.
type Room struct {
	number    int
	capacity  int
	state     RoomState
	occupants []*Patient
}

// NewRoom crea una habitación disponible.
func NewRoom(number, capacity int) *Room {
	return &Room{number: number, capacity: capacity, state: Available}
}

// Number devuelve el número de la habitación.
func (r *Room) Number() int { return r.number }

// Capacity devuelve cuántos pacientes caben.
func (r *Room) Capacity() int { return r.capacity }

// State devuelve el estado de la habitación.
func (r *Room) State() RoomState { return r.state }

// IsAvailable dice si todavía cabe otro paciente.
func (r *Room) IsAvailable() bool { return r.state == Available }

// Occupy mete al paciente en la habitación y lo deja AsleepInBed.
func (r *Room) Occupy(p *Patient) error {
	if p == nil {
		return ErrNilPatient
	}
	if !r.IsAvailable() {
		return ErrRoomFull
	}
	r.occupants = append(r.occupants, p)
	if len(r.occupants) >= r.capacity {
		r.state = Occupied
	}
	p.room = r
	p.state = AsleepInBed
	return nil
}

// Release saca al paciente de la habitación.
func (r *Room) Release(p *Patient) error {
	if p == nil {
		return ErrNilPatient
	}
	for i, o := range r.occupants {
		if o == p {
			r.occupants = append(r.occupants[:i], r.occupants[i+1:]...)
			r.state = Available
			p.room = nil
			return nil
		}
	}
	return ErrNotInRoom
}
