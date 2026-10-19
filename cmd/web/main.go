// Command web es la interfaz gráfica del hospital: una API REST con
// net/http y una página HTML/JS. El paquete hospital NO se modifica:
// aquí solo se usan sus métodos públicos.
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"

	"example.com/costenos-narcolepsia/hospital"
)

//go:embed index.html
var indexHTML []byte

// server guarda el hospital y las referencias que necesita la interfaz.
// El hospital no es seguro para varias goroutines, y net/http atiende
// cada petición en su propia goroutine, así que el mutex vive aquí
// (en la capa de la vista) y no dentro del paquete hospital.
type server struct {
	mu       sync.Mutex
	h        *hospital.Hospital
	doctors  []*hospital.Doctor
	staff    []staffDTO
	patients []*hospital.Patient
	rooms    []*hospital.Room
}

// ---- DTOs: la forma en que los datos viajan como JSON ----

type staffDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type patientDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Level    string `json:"level"`
	State    string `json:"state"`
	Location string `json:"location"`
	Room     int    `json:"room"` // 0 = sin habitación
	Doctor   string `json:"doctor"`
}

type roomDTO struct {
	Number    int      `json:"number"`
	Capacity  int      `json:"capacity"`
	State     string   `json:"state"`
	Occupants []string `json:"occupants"`
}

type episodeDTO struct {
	ID          string `json:"id"`
	Time        string `json:"time"`
	PatientID   string `json:"patient_id"`
	PatientName string `json:"patient_name"`
	Attender    string `json:"attender"`
	Location    string `json:"location"`
	Room        int    `json:"room"` // 0 = sigue en el pasillo
}

type doctorDTO struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Specialty string       `json:"specialty"`
	Episodes  []episodeDTO `json:"episodes"`
}

type severeDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Episodes int    `json:"episodes"`
}

type stateDTO struct {
	Hospital string       `json:"hospital"`
	Staff    []staffDTO   `json:"staff"`
	Patients []patientDTO `json:"patients"`
	Rooms    []roomDTO    `json:"rooms"`
	Hallway  []string     `json:"hallway"` // consulta 5.1
	Doctors  []doctorDTO  `json:"doctors"` // consulta 5.2
	Severe   []severeDTO  `json:"severe"`  // consulta 5.4
	History  []episodeDTO `json:"history"`
}

type actionResult struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func main() {
	s := newServer()

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/state", s.handleState)
	mux.HandleFunc("/api/sleep", s.handleSleep)
	mux.HandleFunc("/api/wake", s.handleWake)
	mux.HandleFunc("/api/assign", s.handleAssign)
	mux.HandleFunc("/api/reset", s.handleReset)

	fmt.Println("Hospital de los Costeños con Narcolepsia")
	fmt.Println("Abre http://localhost:8080 en tu navegador (Ctrl+C para cerrar)")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

func newServer() *server {
	s := &server{}
	s.reset()
	return s
}

// check se usa solo al armar los datos iniciales, que son fijos y correctos.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// reset deja el hospital como al inicio del escenario: personal,
// 3 habitaciones y 5 pacientes despiertos.
func (s *server) reset() {
	h := hospital.NewHospital("Hospital de los Costeños con Narcolepsia")

	karen := hospital.NewDoctor("D-001", "Karen Ospina", 38, "Neurología")
	luis := hospital.NewDoctor("D-002", "Luis Mejía", 45, "Medicina del sueño")
	jorge := hospital.NewOrderly("O-001", "Jorge Pérez", 29)
	check(h.HireDoctor(karen))
	check(h.HireDoctor(luis))
	check(h.HireOrderly(jorge))

	rooms := []*hospital.Room{
		hospital.NewRoom(101, 1),
		hospital.NewRoom(102, 1),
		hospital.NewRoom(103, 1),
	}
	for _, r := range rooms {
		check(h.AddRoom(r))
	}

	patients := []*hospital.Patient{
		hospital.NewPatient("P-001", "Edwin Bolaños", 34, hospital.Mild),
		hospital.NewPatient("P-002", "Yeimy Padilla", 27, hospital.Severe),
		hospital.NewPatient("P-003", "Aníbal Cantillo", 52, hospital.Moderate),
		hospital.NewPatient("P-004", "Wilfrido Berrio", 41, hospital.Severe),
		hospital.NewPatient("P-005", "Dayana Pertuz", 23, hospital.Severe),
	}
	for _, p := range patients {
		check(h.AdmitPatient(p))
	}

	s.h = h
	s.doctors = []*hospital.Doctor{karen, luis}
	s.staff = []staffDTO{
		{karen.ID(), karen.Name(), "Doctor"},
		{luis.ID(), luis.Name(), "Doctor"},
		{jorge.ID(), jorge.Name(), "Orderly"},
	}
	s.patients = patients
	s.rooms = rooms
}

func (s *server) findPatient(id string) *hospital.Patient {
	for _, p := range s.patients {
		if p.ID() == id {
			return p
		}
	}
	return nil
}

func episodeToDTO(rec hospital.EpisodeRecord) episodeDTO {
	roomNumber := 0
	if rec.Room() != nil {
		roomNumber = rec.Room().Number()
	}
	return episodeDTO{
		ID:          rec.ID(),
		Time:        rec.DateTime().Format("15:04:05"),
		PatientID:   rec.Patient().ID(),
		PatientName: rec.Patient().Name(),
		Attender:    rec.Attender().Name(),
		Location:    rec.Location(),
		Room:        roomNumber,
	}
}

// snapshot arma el estado completo para la página. Debe llamarse con el
// mutex tomado.
func (s *server) snapshot() stateDTO {
	st := stateDTO{
		Hospital: s.h.Name(),
		Staff:    s.staff,
		Patients: make([]patientDTO, 0, len(s.patients)),
		Rooms:    make([]roomDTO, 0, len(s.rooms)),
		Hallway:  make([]string, 0),
		Doctors:  make([]doctorDTO, 0, len(s.doctors)),
		Severe:   make([]severeDTO, 0),
		History:  make([]episodeDTO, 0),
	}

	for _, p := range s.patients {
		roomNumber := 0
		if p.CurrentRoom() != nil {
			roomNumber = p.CurrentRoom().Number()
		}
		doctor := ""
		if p.AssignedDoctor() != nil {
			doctor = p.AssignedDoctor().Name()
		}
		st.Patients = append(st.Patients, patientDTO{
			ID: p.ID(), Name: p.Name(), Age: p.Age(),
			Level: p.Level().String(), State: p.State().String(),
			Location: p.Location(), Room: roomNumber, Doctor: doctor,
		})
	}

	for _, r := range s.rooms {
		occupants := make([]string, 0)
		for _, p := range s.patients {
			if p.CurrentRoom() == r {
				occupants = append(occupants, p.ID())
			}
		}
		st.Rooms = append(st.Rooms, roomDTO{
			Number: r.Number(), Capacity: r.Capacity(),
			State: r.State().String(), Occupants: occupants,
		})
	}

	// Consulta 5.1
	for _, p := range s.h.PatientsInHallway() {
		st.Hallway = append(st.Hallway, p.ID())
	}

	// Consulta 5.2
	for _, d := range s.doctors {
		eps := make([]episodeDTO, 0)
		for _, rec := range d.MyEpisodes() {
			eps = append(eps, episodeToDTO(rec))
		}
		st.Doctors = append(st.Doctors, doctorDTO{
			ID: d.ID(), Name: d.Name(), Specialty: d.Specialty(), Episodes: eps,
		})
	}

	// Consulta 5.4 (el map no tiene orden, se ordena por id)
	for p, n := range s.h.SevereReport() {
		st.Severe = append(st.Severe, severeDTO{ID: p.ID(), Name: p.Name(), Episodes: n})
	}
	sort.Slice(st.Severe, func(i, j int) bool { return st.Severe[i].ID < st.Severe[j].ID })

	for _, rec := range s.h.History() {
		st.History = append(st.History, episodeToDTO(rec))
	}
	return st
}

// ---- Handlers ----

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("error escribiendo JSON:", err)
	}
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(indexHTML); err != nil {
		log.Println("error escribiendo la página:", err)
	}
}

func (s *server) handleState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, s.snapshot())
}

// request es el cuerpo JSON que mandan los botones.
type request struct {
	PatientID string `json:"patient_id"`
	Location  string `json:"location"`
}

// readPatient valida el método, lee el cuerpo y busca al paciente.
// Si algo falla ya responde el error y devuelve nil.
func (s *server) readPatient(w http.ResponseWriter, r *http.Request) (*hospital.Patient, request) {
	var req request
	if r.Method != http.MethodPost {
		http.Error(w, "usa POST", http.StatusMethodNotAllowed)
		return nil, req
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return nil, req
	}
	p := s.findPatient(req.PatientID)
	if p == nil {
		http.Error(w, "el paciente no existe", http.StatusNotFound)
		return nil, req
	}
	return p, req
}

func (s *server) handleSleep(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, req := s.readPatient(w, r)
	if p == nil {
		return
	}
	room, err := s.h.ReportSleepAttack(p, req.Location)
	if err != nil {
		writeJSON(w, actionResult{false, err.Error()})
		return
	}
	writeJSON(w, actionResult{true, fmt.Sprintf("%s se durmió en %s y quedó en la habitación %d",
		p.ID(), req.Location, room.Number())})
}

func (s *server) handleWake(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, _ := s.readPatient(w, r)
	if p == nil {
		return
	}
	if err := p.WakeUp(); err != nil {
		writeJSON(w, actionResult{false, fmt.Sprintf("%s: %v", p.ID(), err)})
		return
	}
	writeJSON(w, actionResult{true, fmt.Sprintf("%s despertó y liberó su habitación", p.ID())})
}

func (s *server) handleAssign(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, _ := s.readPatient(w, r)
	if p == nil {
		return
	}
	room, err := s.h.AssignRoom(p)
	if err != nil {
		writeJSON(w, actionResult{false, err.Error()})
		return
	}
	writeJSON(w, actionResult{true, fmt.Sprintf("%s fue asignado a la habitación %d", p.ID(), room.Number())})
}

func (s *server) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "usa POST", http.StatusMethodNotAllowed)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reset()
	writeJSON(w, actionResult{true, "Hospital reiniciado"})
}
