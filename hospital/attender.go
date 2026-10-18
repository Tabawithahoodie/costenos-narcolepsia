package hospital

// Attender es la interfaz que cumple todo el personal que puede atender
// una emergencia de sueño. El hospital guarda su personal en un
// []Attender y no necesita saber si es Doctor u Orderly.
type Attender interface {
	ID() string
	Name() string
	Attend(p *Patient, location string) (EpisodeRecord, error)
}

// Comprobaciones en tiempo de compilación.
var (
	_ Attender = (*Doctor)(nil)
	_ Attender = (*Orderly)(nil)
)
