package hospital

import (
	"errors"
	"testing"
)

// nuevoHospital arma un hospital con n habitaciones de capacidad 1,
// un doctor y un camillero.
func nuevoHospital(t *testing.T, n int) *Hospital {
	t.Helper()
	h := NewHospital("Test")
	if err := h.HireDoctor(NewDoctor("D-001", "Karen", 38, "Neurología")); err != nil {
		t.Fatal(err)
	}
	if err := h.HireOrderly(NewOrderly("O-001", "Jorge", 29)); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < n; i++ {
		if err := h.AddRoom(NewRoom(101+i, 1)); err != nil {
			t.Fatal(err)
		}
	}
	return h
}

func admitir(t *testing.T, h *Hospital, id string, nivel NarcolepsyLevel) *Patient {
	t.Helper()
	p := NewPatient(id, "Paciente "+id, 30, nivel)
	if err := h.AdmitPatient(p); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestAssignRoomConHabitacionLibre(t *testing.T) {
	h := nuevoHospital(t, 1)
	p := admitir(t, h, "P-001", Mild)
	p.SufferSleepAttack("cafeteria")

	room, err := h.AssignRoom(p)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if room.State() != Occupied {
		t.Errorf("la habitación debería quedar Occupied, quedó %v", room.State())
	}
	if p.State() != AsleepInBed {
		t.Errorf("el paciente debería quedar AsleepInBed, quedó %v", p.State())
	}
}

func TestAssignRoomSinHabitacion(t *testing.T) {
	h := nuevoHospital(t, 0)
	p := admitir(t, h, "P-001", Severe)
	p.SufferSleepAttack("radiology queue")

	room, err := h.AssignRoom(p)
	if !errors.Is(err, ErrNoRoom) {
		t.Fatalf("se esperaba ErrNoRoom, llegó %v", err)
	}
	if room != nil {
		t.Error("no se debe inventar una habitación")
	}
	if p.State() != AsleepInHallway {
		t.Errorf("el paciente debe seguir AsleepInHallway, está %v", p.State())
	}
}

func TestPatientsInHallway(t *testing.T) {
	h := nuevoHospital(t, 1)
	p1 := admitir(t, h, "P-001", Mild)
	p2 := admitir(t, h, "P-002", Mild)
	admitir(t, h, "P-003", Mild) // este nunca se duerme

	h.ReportSleepAttack(p1, "cafeteria") // toma la única habitación
	h.ReportSleepAttack(p2, "pharmacy")  // se queda en el pasillo

	enPasillo := h.PatientsInHallway()
	if len(enPasillo) != 1 || enPasillo[0] != p2 {
		t.Fatalf("debería haber solo a P-002 en el pasillo, llegó %v", enPasillo)
	}
}

func TestWakeUpLiberaHabitacion(t *testing.T) {
	h := nuevoHospital(t, 1)
	p1 := admitir(t, h, "P-001", Mild)
	p2 := admitir(t, h, "P-002", Mild)
	h.ReportSleepAttack(p1, "cafeteria")
	h.ReportSleepAttack(p2, "pharmacy")

	if err := p1.WakeUp(); err != nil {
		t.Fatal(err)
	}
	room, err := h.AssignRoom(p2)
	if err != nil || room.Number() != 101 {
		t.Fatalf("P-002 debería tomar la habitación 101, llegó %v, %v", room, err)
	}
}

func TestDespachoPolimorfico(t *testing.T) {
	h := nuevoHospital(t, 3)
	p1 := admitir(t, h, "P-001", Mild)
	p2 := admitir(t, h, "P-002", Mild)
	h.ReportSleepAttack(p1, "a")
	h.ReportSleepAttack(p2, "b")

	hist := h.History()
	if len(hist) != 2 {
		t.Fatalf("se esperaban 2 episodios, hay %d", len(hist))
	}
	_, esDoctor := hist[0].Attender().(*Doctor)
	_, esCamillero := hist[1].Attender().(*Orderly)
	if !esDoctor || !esCamillero {
		t.Error("el primero debería atenderlo el Doctor y el segundo el Orderly")
	}
}

func TestSevereReport(t *testing.T) {
	h := nuevoHospital(t, 2)
	grave := admitir(t, h, "P-001", Severe)
	leve := admitir(t, h, "P-002", Mild)
	h.ReportSleepAttack(grave, "a")
	h.ReportSleepAttack(leve, "b")

	reporte := h.SevereReport()
	if len(reporte) != 1 {
		t.Fatalf("solo debe haber pacientes Severe, hay %d", len(reporte))
	}
	if reporte[grave] != 1 {
		t.Errorf("P-001 debería tener 1 episodio, tiene %d", reporte[grave])
	}
}
