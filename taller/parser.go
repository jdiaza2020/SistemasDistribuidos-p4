package main

import (
	"strconv"
	"strings"
)

// parseIncoming reconstruye líneas con '\n' y extrae solo enteros 0..9.
// Ignora textos del servidor (p.ej. "Taller localizado en ...").
func parseIncoming(incoming <-chan string, out chan<- int) {
	pending := ""

	for chunk := range incoming {
		pending += chunk

		for {
			idx := strings.IndexByte(pending, '\n')
			if idx == -1 {
				break
			}

			line := strings.TrimSpace(pending[:idx])
			pending = pending[idx+1:]

			if line == "" {
				continue
			}

			// Recomendación del profe: strconv.Atoi.
			n, err := strconv.Atoi(line)
			if err != nil {
				continue
			}
			if n < 0 || n > 9 {
				continue
			}

			out <- n
		}
	}

	// Si cerramos incoming en el futuro, cerramos out.
	close(out)
}
