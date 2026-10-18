package hospital

// PatientState indica en qué situación está un paciente.
type PatientState int

const (
	Awake PatientState = iota
	AsleepInHallway
	AsleepInBed
)

// String permite imprimir el estado de forma legible.
func (s PatientState) String() string {
	switch s {
	case Awake:
		return "Awake"
	case AsleepInHallway:
		return "AsleepInHallway"
	case AsleepInBed:
		return "AsleepInBed"
	default:
		return "Unknown"
	}
}

// NarcolepsyLevel es la gravedad de la narcolepsia del paciente.
type NarcolepsyLevel int

const (
	Mild NarcolepsyLevel = iota
	Moderate
	Severe
)

func (l NarcolepsyLevel) String() string {
	switch l {
	case Mild:
		return "Mild"
	case Moderate:
		return "Moderate"
	case Severe:
		return "Severe"
	default:
		return "Unknown"
	}
}

// RoomState indica si una habitación tiene cupo o no.
type RoomState int

const (
	Available RoomState = iota
	Occupied
)

func (s RoomState) String() string {
	switch s {
	case Available:
		return "Available"
	case Occupied:
		return "Occupied"
	default:
		return "Unknown"
	}
}
