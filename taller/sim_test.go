package main

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type testCase struct {
	name             string
	numA, numB, numC int
}

type cfgCase struct {
	name         string
	numPlazas    int
	numMecanicos int
}

// Escala de tiempo para tests: 20 => duerme 20 veces menos.
// (Si quieres aún más rápido, sube a 50).
const timeScale = 20

func scaledSleep(d time.Duration) {
	if d <= 0 {
		return
	}
	time.Sleep(d / timeScale)
}

// runScenario ejecuta una simulación SIN TCP ni servidor.
// Devuelve duración total y throughput.
func runScenario(t *testing.T, cfg Config) (time.Duration, float64) {
	t.Helper()

	// Estado fijo NORMAL durante el test.
	stateProvider = func() TallerState {
		return TallerState{Activo: true, Cerrado: false}
	}

	// Sleep acelerado (para que no peten los timeouts).
	sleepFn = scaledSleep

	logCh := make(chan LogEvent, 8192)

	totalCoches := int32(cfg.NumA + cfg.NumB + cfg.NumC)
	var finished int32

	start := time.Now()

	done := make(chan struct{})
	go func() {
		for ev := range logCh {
			if ev.Fase == FaseEntrega && ev.Estado == "Sale" {
				if atomic.AddInt32(&finished, 1) == totalCoches {
					close(done)
					return
				}
			}
		}
	}()

	startSimulation(start, logCh, cfg)

	// Timeout del escenario (ya no debería saltar con timeScale).
	select {
	case <-done:
		// NO cerramos logCh para evitar que goroutines antiguas de otros subtests hagan send a canal cerrado.
		// (Los workers se quedan vivos, pero en tests no nos afecta y evitamos panics.)
	case <-time.After(2 * time.Minute):
		t.Fatalf("timeout: finalizaron %d/%d coches", finished, totalCoches)
	}

	dur := time.Since(start)
	throughput := float64(totalCoches) / dur.Seconds()
	return dur, throughput
}

func TestComparativas_6Casos(t *testing.T) {
	tests := []testCase{
		{name: "T1_10_10_10", numA: 10, numB: 10, numC: 10},
		{name: "T2_20_5_5", numA: 20, numB: 5, numC: 5},
		{name: "T3_5_5_20", numA: 5, numB: 5, numC: 20},
	}

	cfgs := []cfgCase{
		{name: "CFG1_P6_M3", numPlazas: 6, numMecanicos: 3},
		{name: "CFG2_P4_M4", numPlazas: 4, numMecanicos: 4},
	}

	for _, tc := range tests {
		for _, cc := range cfgs {
			name := fmt.Sprintf("%s_%s", tc.name, cc.name)

			t.Run(name, func(t *testing.T) {
				cfg := DefaultConfig()

				cfg.NumA, cfg.NumB, cfg.NumC = tc.numA, tc.numB, tc.numC
				cfg.NumPlazas = cc.numPlazas
				cfg.NumMecanicos = cc.numMecanicos

				// Mantenemos limpieza/entrega constantes para aislar el impacto de plazas/mecánicos.
				cfg.NumLimpieza = 1
				cfg.NumEntrega = 1

				// Colas grandes para que la capacidad de cola no afecte a la comparativa.
				cfg.CapQ1, cfg.CapQ2, cfg.CapQ3 = 1000, 1000, 1000

				dur, th := runScenario(t, cfg)
				t.Logf("%s (timeScale=%dx) -> dur=%v | throughput=%.2f coches/s", name, timeScale, dur, th)
			})
		}
	}
}
