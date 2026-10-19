# Hospital de los Costeños con Narcolepsia

Trabajo de Programación Orientada a Objetos — Universidad EIA (Go 1.21+, solo librería estándar).

**Integrante:** _Samuel Tabares Gil_ — grupo de 1 estudiante.

## Cómo correrlo

```bash
go run .          # corre el escenario de demostración
go test ./...     # corre los tests
gofmt -l .        # no debe imprimir nada
go vet ./...      # no debe reportar nada
go run ./cmd/web   # interfaz gráfica (bono): abre http://localhost:8080
```

## Estructura

```
costenos-narcolepsia/
  main.go               // solo arma el escenario e imprime
  hospital/
    hospital.go         // Hospital: personal, pacientes, habitaciones, consultas
    person.go           // Person (se incrusta en Patient, Doctor y Orderly)
    patient.go          // Patient
    doctor.go           // Doctor (cumple Attender)
    orderly.go          // Orderly / camillero (cumple Attender)
    attender.go         // interfaz Attender
    room.go             // Room
    episode.go          // EpisodeRecord
    states.go           // estados como tipos con constantes (iota)
    errors.go           // errores previstos del hospital
    hospital_test.go    // tests
```

## Decisiones de diseño

- **Encapsulación:** todos los campos van en minúscula (no exportados), así que
  fuera del paquete `hospital` solo se tocan con métodos. `main` no puede hacer
  `p.state = Awake`: el compilador no lo deja.
- **Composición:** `Person` (id, name, age) se incrusta **por valor** en
  `Patient`, `Doctor` y `Orderly`. Así reutilizan `ID()`, `Name()` y `Age()` sin
  repetir campos.
- **Interfaz `Attender`:** la cumplen `*Doctor` y `*Orderly`. El hospital guarda
  el personal en un `[]Attender` y en `dispatch` los recorre por turnos sin saber
  qué tipo concreto es. Si uno no puede atender (devuelve error) pasa al siguiente.
- **Estados como tipos:** `PatientState`, `NarcolepsyLevel` y `RoomState` son
  `int` con constantes `iota` y método `String()`, no strings sueltos.
- **Errores, no panic:** lo previsible (no hay cama, paciente despierto, id
  repetido) devuelve un `error`. Se usan errores centinela (`ErrNoRoom`...) y
  `%w` para poder comprobarlos con `errors.Is`.
- **`MyEpisodes()` sin recibir el hospital:** el `Doctor` guarda sus propios
  episodios en `Attend`, por eso la relación vive dentro de la entidad.
- **`SevereReport` devuelve `map[*Patient]int`:** el mapa cruza directamente
  cada paciente grave con su conteo. Como los mapas no tienen orden, `main`
  ordena por id antes de imprimir.
- **`ReportSleepAttack`:** junta el flujo completo (dormirse, asignar cama,
  despachar personal, registrar). Si no hay cama, el episodio igual se registra.
- **Limitación:** el hospital no es seguro para goroutines (no se hizo el bono de
  concurrencia). Si dos goroutines asignaran la misma habitación a la vez habría
  una condición de carrera; la solución sería un `sync.Mutex` en `Hospital`.
