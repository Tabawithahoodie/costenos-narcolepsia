package hospital

// Person guarda los datos de identidad que comparten pacientes, doctores
// y camilleros. Se incrusta (embedding) en cada uno de ellos para no
// repetir los mismos campos.
type Person struct {
	id   string
	name string
	age  int
}

// ID devuelve el identificador de la persona.
func (p Person) ID() string { return p.id }

// Name devuelve el nombre de la persona.
func (p Person) Name() string { return p.name }

// Age devuelve la edad de la persona.
func (p Person) Age() int { return p.age }
