package main

import "fmt"

type stateRequest struct {
	reply chan TallerState
}

// controller mantiene el TallerState actualizado y permite consultarlo.
// - codes: stream de 0..9 desde la mutua
// - queries: peticiones de “dame el estado actual”
func controller(codes <-chan int, queries <-chan stateRequest) {
	state := defaultState()

	for {
		select {
		case code, ok := <-codes:
			if !ok {
				return
			}
			state.applyCode(code)
			fmt.Println("ESTADO ACTUAL:", stateSummary(state)) // debug temporal
		case req := <-queries:
			req.reply <- state
		}
	}
}

func stateSummary(s TallerState) string {
	if s.Cerrado {
		return "CERRADO"
	}
	if !s.Activo {
		return "INACTIVO"
	}
	if s.SoloCategoria != "" {
		return "SOLO " + s.SoloCategoria
	}
	if s.PrioridadCategoria != "" {
		return "PRIORIDAD " + s.PrioridadCategoria
	}
	return "NORMAL"
}
