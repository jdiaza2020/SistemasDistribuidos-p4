package main

const (
	CatA = "A" // Mecánica
	CatB = "B" // Eléctrica
	CatC = "C" // Carrocería
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
