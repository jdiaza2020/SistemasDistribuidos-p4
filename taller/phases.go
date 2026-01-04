package main

import "time"

// stateProvider permite sustituir el origen del estado en tests.
// Por defecto apunta a getState (controlador real).
var stateProvider = func() TallerState { return getState() }

// sleepFn permite acelerar tests si hiciera falta en el futuro.
var sleepFn = time.Sleep

// fase0Plaza: respeta estado (inactivo/cerrado/solo categoría), usa plazas y al salir ENCOLA en fase 1.
func fase0Plaza(start time.Time, c Coche, plazas chan struct{}, q1 *PhaseQueue, logs chan<- LogEvent) {
	for {
		st := stateProvider()

		// Si cerrado o inactivo: espera y reintenta.
		if st.Cerrado || !st.Activo {
			sleepFn(200 * time.Millisecond)
			continue
		}

		// Restricción “solo categoría X”.
		if st.SoloCategoria != "" && st.SoloCategoria != c.Categoria {
			sleepFn(200 * time.Millisecond)
			continue
		}

		// Coge plaza (bloquea si no hay).
		plazas <- struct{}{}

		// Re-chequeo por si cambió justo después.
		st2 := stateProvider()
		if st2.Cerrado || !st2.Activo || (st2.SoloCategoria != "" && st2.SoloCategoria != c.Categoria) {
			<-plazas
			sleepFn(200 * time.Millisecond)
			continue
		}

		inc := categoriaTipo(c.Categoria)
		logs <- LogEvent{Elapsed: time.Since(start), CocheID: c.ID, Incidencia: inc, Fase: FaseEsperaPlaza, Estado: "Entra"}

		sleepFn(categoriaDurConVariacion(c.Categoria))

		logs <- LogEvent{Elapsed: time.Since(start), CocheID: c.ID, Incidencia: inc, Fase: FaseEsperaPlaza, Estado: "Sale"}

		<-plazas

		// Encolamos en la Fase 1 con prioridad.
		q1.Enqueue(c)
		return
	}
}

// phase1Worker: worker de mecánico.
// - Saca coches de la cola q1 aplicando PRIORIDAD del estado.
// - Respeta inactivo/cerrado/solo categoría antes de empezar un trabajo.
// - Usa "mecanicos" como semáforo de recurso físico.
// - Al terminar, encola en Fase 2.
func phase1Worker(start time.Time, q1 *PhaseQueue, q2 *PhaseQueue, mecanicos chan struct{}, logs chan<- LogEvent) {
	for {
		st := stateProvider()
		car := q1.Dequeue(st)

		// Espera a que el estado permita atender.
		for {
			st2 := stateProvider()
			if st2.Cerrado || !st2.Activo {
				sleepFn(200 * time.Millisecond)
				continue
			}
			if st2.SoloCategoria != "" && st2.SoloCategoria != car.Categoria {
				sleepFn(200 * time.Millisecond)
				continue
			}
			break
		}

		// Espera mecánico libre.
		mecanicos <- struct{}{}

		// Re-chequeo antes del trabajo real.
		st3 := stateProvider()
		if st3.Cerrado || !st3.Activo || (st3.SoloCategoria != "" && st3.SoloCategoria != car.Categoria) {
			<-mecanicos
			// Devolvemos el coche a la cola para no perderlo.
			q1.Enqueue(car)
			sleepFn(200 * time.Millisecond)
			continue
		}

		inc := categoriaTipo(car.Categoria)
		logs <- LogEvent{Elapsed: time.Since(start), CocheID: car.ID, Incidencia: inc, Fase: FaseMecanico, Estado: "Entra"}

		sleepFn(categoriaDurConVariacion(car.Categoria))

		logs <- LogEvent{Elapsed: time.Since(start), CocheID: car.ID, Incidencia: inc, Fase: FaseMecanico, Estado: "Sale"}

		<-mecanicos

		// Pasa a Fase 2
		q2.Enqueue(car)
	}
}

// phase2Worker: limpieza.
// Mismo patrón que fase 1, con su propio recurso y cola.
func phase2Worker(start time.Time, q2 *PhaseQueue, q3 *PhaseQueue, limpieza chan struct{}, logs chan<- LogEvent) {
	for {
		st := stateProvider()
		car := q2.Dequeue(st)

		for {
			st2 := stateProvider()
			if st2.Cerrado || !st2.Activo {
				sleepFn(200 * time.Millisecond)
				continue
			}
			if st2.SoloCategoria != "" && st2.SoloCategoria != car.Categoria {
				sleepFn(200 * time.Millisecond)
				continue
			}
			break
		}

		limpieza <- struct{}{}

		st3 := stateProvider()
		if st3.Cerrado || !st3.Activo || (st3.SoloCategoria != "" && st3.SoloCategoria != car.Categoria) {
			<-limpieza
			q2.Enqueue(car)
			sleepFn(200 * time.Millisecond)
			continue
		}

		inc := categoriaTipo(car.Categoria)
		logs <- LogEvent{Elapsed: time.Since(start), CocheID: car.ID, Incidencia: inc, Fase: FaseLimpieza, Estado: "Entra"}

		sleepFn(categoriaDurConVariacion(car.Categoria))

		logs <- LogEvent{Elapsed: time.Since(start), CocheID: car.ID, Incidencia: inc, Fase: FaseLimpieza, Estado: "Sale"}

		<-limpieza

		// Pasa a Fase 3
		q3.Enqueue(car)
	}
}

// phase3Worker: entrega.
// Última fase del pipeline.
func phase3Worker(start time.Time, q3 *PhaseQueue, entrega chan struct{}, logs chan<- LogEvent) {
	for {
		st := stateProvider()
		car := q3.Dequeue(st)

		for {
			st2 := stateProvider()
			if st2.Cerrado || !st2.Activo {
				sleepFn(200 * time.Millisecond)
				continue
			}
			if st2.SoloCategoria != "" && st2.SoloCategoria != car.Categoria {
				sleepFn(200 * time.Millisecond)
				continue
			}
			break
		}

		entrega <- struct{}{}

		st3 := stateProvider()
		if st3.Cerrado || !st3.Activo || (st3.SoloCategoria != "" && st3.SoloCategoria != car.Categoria) {
			<-entrega
			q3.Enqueue(car)
			sleepFn(200 * time.Millisecond)
			continue
		}

		inc := categoriaTipo(car.Categoria)
		logs <- LogEvent{Elapsed: time.Since(start), CocheID: car.ID, Incidencia: inc, Fase: FaseEntrega, Estado: "Entra"}

		sleepFn(categoriaDurConVariacion(car.Categoria))

		logs <- LogEvent{Elapsed: time.Since(start), CocheID: car.ID, Incidencia: inc, Fase: FaseEntrega, Estado: "Sale"}

		<-entrega
	}
}
