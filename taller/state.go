package main

// Estado del taller controlado por la mutua (códigos 0..9 del enunciado).
type TallerState struct {
	Activo  bool // false si 0 (inactivo) o 9 (cerrado)
	Cerrado bool

	// Si SoloCategoria != "" solo se permite entrar a esa categoría ("A","B","C").
	SoloCategoria string

	// Si PrioridadCategoria != "" esa categoría tiene prioridad preferente.
	PrioridadCategoria string
}

// Estado inicial por defecto: activo, sin restricciones, sin prioridad especial.
func defaultState() TallerState {
	return TallerState{
		Activo:             true,
		Cerrado:            false,
		SoloCategoria:      "",
		PrioridadCategoria: "",
	}
}

// applyCode actualiza el estado según el código recibido.
// 7 y 8: mantener estado anterior.
func (s *TallerState) applyCode(code int) {
	switch code {
	case 0:
		// Taller inactivo: no hay atención.
		s.Activo = false
		s.Cerrado = false
		// Mantener restricciones y prioridad como estaban (decisión razonable).
	case 1:
		s.Activo = true
		s.Cerrado = false
		s.SoloCategoria = CatA
	case 2:
		s.Activo = true
		s.Cerrado = false
		s.SoloCategoria = CatB
	case 3:
		s.Activo = true
		s.Cerrado = false
		s.SoloCategoria = CatC
	case 4:
		s.Activo = true
		s.Cerrado = false
		s.PrioridadCategoria = CatA
	case 5:
		s.Activo = true
		s.Cerrado = false
		s.PrioridadCategoria = CatB
	case 6:
		s.Activo = true
		s.Cerrado = false
		s.PrioridadCategoria = CatC
	case 7, 8:
		// No definido: se mantiene estado anterior.
		return
	case 9:
		// Taller cerrado: no hay atención.
		s.Activo = false
		s.Cerrado = true
	default:
		// Fuera de rango: ignorar.
		return
	}
}
