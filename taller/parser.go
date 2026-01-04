package main

import (
	"strconv"
	"strings"
)

// parseIncoming recibe trozos (pueden venir fragmentados) y reconstruye por '\n'.
// Extrae solo enteros 0..9; el resto (textos del servidor) se ignora.
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

	// Si se cerrase incoming en algún momento, cerramos out.
	close(out)
}
