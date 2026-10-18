package hospital

import "errors"

// Errores que el hospital puede prever (nada de panic).
var (
	ErrNoRoom          = errors.New("no hay habitación disponible")
	ErrRoomFull        = errors.New("la habitación ya está llena")
	ErrNotInRoom       = errors.New("el paciente no está en esa habitación")
	ErrNilValue        = errors.New("el valor recibido es nulo")
	ErrNilPatient      = errors.New("el paciente es nulo")
	ErrPatientAwake    = errors.New("el paciente está despierto")
	ErrAlreadyInBed    = errors.New("el paciente ya está en una cama")
	ErrNoStaff         = errors.New("no hay personal disponible para atender")
	ErrDuplicate       = errors.New("ya existe un registro con ese identificador")
	ErrPatientNotKnown = errors.New("el paciente no está admitido en el hospital")
	ErrInvalidRecord   = errors.New("el registro del episodio es inválido")
	ErrAlreadyAwake    = errors.New("el paciente ya estaba despierto")
)
