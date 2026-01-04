package main

import (
	"math/rand"
	"time"
)

// Categorías.
const (
	CatA = "A" // Mecánica
	CatB = "B" // Eléctrica
	CatC = "C" // Carrocería
)

// Fases.
const (
	FaseEsperaPlaza = iota // 0
	FaseMecanico           // 1
	FaseLimpieza           // 2
	FaseEntrega            // 3
)

func categoriaTipo(cat string) string {
	switch cat {
	case CatA:
		return "Mecanica"
	case CatB:
		return "Electrica"
	default:
		return "Carroceria"
	}
}

// Tiempo base por fase según categoría (segundos): A=5, B=3, C=1.
func categoriaBaseDur(cat string) time.Duration {
	switch cat {
	case CatA:
		return 5 * time.Second
	case CatB:
		return 3 * time.Second
	default:
		return 1 * time.Second
	}
}

// categoriaDurConVariacion aplica una variación simple al tiempo base.
// Usamos un jitter de +/-20% del tiempo base (acotado para que nunca sea <=0).
func categoriaDurConVariacion(cat string) time.Duration {
	base := categoriaBaseDur(cat)

	// jitter en milisegundos: +/-20% del base
	baseMs := int(base / time.Millisecond)
	if baseMs <= 0 {
		return base
	}

	maxJitter := baseMs / 5 // 20%
	j := rand.Intn(2*maxJitter+1) - maxJitter

	d := time.Duration(baseMs+j) * time.Millisecond
	if d < 50*time.Millisecond {
		d = 50 * time.Millisecond
	}
	return d
}

type Coche struct {
	ID        int
	Categoria string
}

type LogEvent struct {
	Elapsed    time.Duration
	CocheID    int
	Incidencia string
	Fase       int
	Estado     string // "Entra" o "Sale"
}
