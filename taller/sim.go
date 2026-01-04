package main

import (
	"math/rand"
	"time"
)

// Config agrupa todos los parámetros del taller para poder testear escenarios.
type Config struct {
	NumA int
	NumB int
	NumC int

	NumPlazas    int
	NumMecanicos int
	NumLimpieza  int
	NumEntrega   int

	CapQ1 int
	CapQ2 int
	CapQ3 int
}

// DefaultConfig para ejecución manual (go run ./taller).
func DefaultConfig() Config {
	return Config{
		NumA: 4, NumB: 4, NumC: 4,

		NumPlazas:    4,
		NumMecanicos: 2,
		NumLimpieza:  1,
		NumEntrega:   1,

		CapQ1: 10, CapQ2: 10, CapQ3: 10,
	}
}

// startSimulation genera coches y los hace pasar por 4 fases.
// Cada fase usa: cola con prioridad + recurso limitado (semáforo).
func startSimulation(start time.Time, logs chan<- LogEvent, cfg Config) {
	// Recursos físicos.
	plazas := make(chan struct{}, cfg.NumPlazas)
	mecanicos := make(chan struct{}, cfg.NumMecanicos)
	limpieza := make(chan struct{}, cfg.NumLimpieza)
	entrega := make(chan struct{}, cfg.NumEntrega)

	// Colas por fase con capacidad máxima.
	q1 := NewPhaseQueue(cfg.CapQ1)
	q2 := NewPhaseQueue(cfg.CapQ2)
	q3 := NewPhaseQueue(cfg.CapQ3)

	// Workers por fase.
	for i := 0; i < cfg.NumMecanicos; i++ {
		go phase1Worker(start, q1, q2, mecanicos, logs)
	}
	for i := 0; i < cfg.NumLimpieza; i++ {
		go phase2Worker(start, q2, q3, limpieza, logs)
	}
	for i := 0; i < cfg.NumEntrega; i++ {
		go phase3Worker(start, q3, entrega, logs)
	}

	// Generamos coches por categoría (A/B/C) y orden aleatorio.
	coches := genCoches(cfg.NumA, cfg.NumB, cfg.NumC)

	// Fase 0: un goroutine por coche.
	for _, c := range coches {
		coche := c
		go fase0Plaza(start, coche, plazas, q1, logs)
	}
}

// genCoches crea NumA, NumB, NumC y mezcla el orden de llegada.
func genCoches(numA, numB, numC int) []Coche {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	total := numA + numB + numC
	out := make([]Coche, 0, total)

	id := 1
	for i := 0; i < numA; i++ {
		out = append(out, Coche{ID: id, Categoria: CatA})
		id++
	}
	for i := 0; i < numB; i++ {
		out = append(out, Coche{ID: id, Categoria: CatB})
		id++
	}
	for i := 0; i < numC; i++ {
		out = append(out, Coche{ID: id, Categoria: CatC})
		id++
	}

	// Shuffle (orden aleatorio).
	for i := len(out) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		out[i], out[j] = out[j], out[i]
	}
	return out
}
