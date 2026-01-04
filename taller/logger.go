package main

import "fmt"

// runLogger imprime logs en un único punto para evitar interleaving.
// Formato exigido: Tiempo {t} Coche {N} Incidencia {Tipo} Fase {Fase} Estado {Entra|Sale}
func runLogger(logs <-chan LogEvent) {
	for ev := range logs {
		fmt.Printf("Tiempo %v Coche %d Incidencia %s Fase %d Estado %s\n",
			ev.Elapsed, ev.CocheID, ev.Incidencia, ev.Fase, ev.Estado)
	}
}
