package hospital

import (
	"fmt"
	"sync/atomic"
	"time"
)

// episodeCounter genera los números de los episodios (E-001, E-002...).
var episodeCounter int64

// EpisodeRecord guarda un episodio de sueño repentino.
type EpisodeRecord struct {
	id       string
	dateTime time.Time
	patient  *Patient
	attender Attender // quien intervino: puede ser Doctor u Orderly
	location string
	room     *Room // puede ser nil si no hubo habitación
}

// newEpisodeRecord crea el registro con la hora actual. La habitación
// es la que el paciente tenga en ese momento (nil si sigue en el pasillo).
func newEpisodeRecord(p *Patient, a Attender, location string) EpisodeRecord {
	n := atomic.AddInt64(&episodeCounter, 1)
	return EpisodeRecord{
		id:       fmt.Sprintf("E-%03d", n),
		dateTime: time.Now(),
		patient:  p,
		attender: a,
		location: location,
		room:     p.room,
	}
}

func (e EpisodeRecord) ID() string          { return e.id }
func (e EpisodeRecord) DateTime() time.Time { return e.dateTime }
func (e EpisodeRecord) Patient() *Patient   { return e.patient }
func (e EpisodeRecord) Attender() Attender  { return e.attender }
func (e EpisodeRecord) Location() string    { return e.location }
func (e EpisodeRecord) Room() *Room         { return e.room }

// Summary devuelve una línea lista para imprimir.
func (e EpisodeRecord) Summary() string {
	destino := "stays in hallway"
	if e.room != nil {
		destino = fmt.Sprintf("room %d", e.room.Number())
	}
	return fmt.Sprintf("%s  %s  %-16s -> %s",
		e.dateTime.Format("2006-01-02 15:04"), e.patient.ID(), e.location, destino)
}
