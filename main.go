package main

import (
	"fmt"
	"sort"

	"example.com/costenos-narcolepsia/hospital"
)

func main() {
	// 1. Crear el hospital, personal, habitaciones y pacientes
	h := hospital.NewHospital("Hospital de los Costeños con Narcolepsia")

	karen := hospital.NewDoctor("D-001", "Karen Ospina", 38, "Neurología")
	luis := hospital.NewDoctor("D-002", "Luis Mejía", 45, "Medicina del sueño")
	jorge := hospital.NewOrderly("O-001", "Jorge Pérez", 29)

	revisar(h.HireDoctor(karen))
	revisar(h.HireDoctor(luis))
	revisar(h.HireOrderly(jorge))

	for _, numero := range []int{101, 102, 103} {
		revisar(h.AddRoom(hospital.NewRoom(numero, 1)))
	}

	p1 := hospital.NewPatient("P-001", "Edwin Bolaños", 34, hospital.Mild)
	p2 := hospital.NewPatient("P-002", "Yeimy Padilla", 27, hospital.Severe)
	p3 := hospital.NewPatient("P-003", "Aníbal Cantillo", 52, hospital.Moderate)
	p4 := hospital.NewPatient("P-004", "Wilfrido Berrio", 41, hospital.Severe)
	p5 := hospital.NewPatient("P-005", "Dayana Pertuz", 23, hospital.Severe)
	for _, p := range []*hospital.Patient{p1, p2, p3, p4, p5} {
		revisar(h.AdmitPatient(p))
	}
	fmt.Printf("== %s ==\n", h.Name())
	fmt.Println("Hospital listo: 2 doctores, 1 camillero, 3 habitaciones y 5 pacientes")

	// 2. Cuatro pacientes se duermen; el cuarto no consigue habitación
	fmt.Println("\n== Ataques de sueño ==")
	ataques := []struct {
		paciente *hospital.Patient
		lugar    string
	}{
		{p1, "cafeteria"},
		{p2, "radiology queue"},
		{p3, "pharmacy line"},
		{p4, "hallway 2, second floor"},
	}
	var asignaciones []string // resultados de AssignRoom, se imprimen en la consulta 5.3
	for _, a := range ataques {
		room, err := h.ReportSleepAttack(a.paciente, a.lugar)
		if err != nil {
			fmt.Printf("  %s (%s): %v\n", a.paciente.ID(), a.lugar, err)
			asignaciones = append(asignaciones, fmt.Sprintf("%s: %v", a.paciente.ID(), err))
			continue
		}
		fmt.Printf("  %s (%s) -> habitación %d\n", a.paciente.ID(), a.lugar, room.Number())
		asignaciones = append(asignaciones, fmt.Sprintf("%s -> habitación %d", a.paciente.ID(), room.Number()))
	}
	fmt.Printf("  Pacientes en pasillo ahora: %d\n", len(h.PatientsInHallway()))

	// 3. Un paciente despierta, libera la habitación y el del pasillo la toma
	fmt.Println("\n== Despierta P-001 y se reasigna su habitación ==")
	revisar(p1.WakeUp())
	room, err := h.AssignRoom(p4)
	if err != nil {
		fmt.Println(" ", err)
		asignaciones = append(asignaciones, fmt.Sprintf("%s: %v", p4.ID(), err))
	} else {
		fmt.Printf("  %s -> habitación %d\n", p4.ID(), room.Number())
		asignaciones = append(asignaciones, fmt.Sprintf("%s -> habitación %d (tras la liberación)", p4.ID(), room.Number()))
	}

	// 4. Resultado de las cuatro consultas
	fmt.Println("\n== Consulta 5.1: pacientes dormidos en pasillo ==")
	enPasillo := h.PatientsInHallway()
	if len(enPasillo) == 0 {
		fmt.Println("  (ninguno)")
	}
	for _, p := range enPasillo {
		fmt.Printf("  %s  %-18s (%s)\n", p.ID(), p.Name(), p.Location())
	}

	fmt.Println("\n== Consulta 5.2: episodios por doctor ==")
	for _, d := range []*hospital.Doctor{karen, luis} {
		fmt.Printf("  Dr. %s\n", d.Name())
		for _, ep := range d.MyEpisodes() {
			fmt.Printf("    %s\n", ep.Summary())
		}
	}

	fmt.Println("\n== Consulta 5.3: asignación de habitación ==")
	for _, linea := range asignaciones {
		fmt.Printf("  %s\n", linea)
	}

	fmt.Println("\n== Consulta 5.4: reporte de severidad ==")
	reporte := h.SevereReport()
	severos := make([]*hospital.Patient, 0, len(reporte))
	for p := range reporte {
		severos = append(severos, p)
	}
	// El map no tiene orden, así que se ordena por id para imprimir
	sort.Slice(severos, func(i, j int) bool { return severos[i].ID() < severos[j].ID() })
	for _, p := range severos {
		fmt.Printf("  %s  %-18s %d episodes today\n", p.ID(), p.Name(), reporte[p])
	}
}

// revisar imprime el error de la preparación del escenario, si lo hay.
func revisar(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
