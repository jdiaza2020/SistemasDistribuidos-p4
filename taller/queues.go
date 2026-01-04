package main

// PhaseQueue es una cola con prioridad y capacidad máxima.
// Implementación estilo "actor": una goroutine es la dueña de los datos.
// No usamos mutex, solo canales.
type PhaseQueue struct {
	capacity int

	enq chan enqReq
	deq chan deqReq
}

type enqReq struct {
	car   Coche
	reply chan struct{} // se cierra cuando el coche queda encolado
}

type deqReq struct {
	state TallerState
	reply chan Coche
}

// NewPhaseQueue crea una cola con capacidad máxima y arranca su goroutine interna.
func NewPhaseQueue(capacity int) *PhaseQueue {
	q := &PhaseQueue{
		capacity: capacity,
		enq:      make(chan enqReq),
		deq:      make(chan deqReq),
	}
	go q.loop()
	return q
}

// Enqueue bloquea hasta que el coche se encola (si está llena, espera hueco).
func (q *PhaseQueue) Enqueue(c Coche) {
	done := make(chan struct{})
	q.enq <- enqReq{car: c, reply: done}
	<-done
}

// Dequeue bloquea hasta que haya un coche disponible y lo devuelve.
// La elección respeta el estado actual (prioridad/solo categoría).
func (q *PhaseQueue) Dequeue(state TallerState) Coche {
	reply := make(chan Coche)
	q.deq <- deqReq{state: state, reply: reply}
	return <-reply
}

// loop mantiene las colas internas y resuelve encolados y desencolados.
func (q *PhaseQueue) loop() {
	// Tres colas simples (A/B/C).
	var a, b, c []Coche

	// Encolados pendientes cuando la cola está llena.
	var pending []enqReq

	// Peticiones de dequeue pendientes cuando no hay coches.
	var waiting []deqReq

	totalLen := func() int { return len(a) + len(b) + len(c) }
	hasAny := func() bool { return totalLen() > 0 }

	// Selecciona el siguiente coche en función del estado.
	pick := func(st TallerState) (Coche, bool) {
		// Si hay "solo categoría", solo sacamos de esa.
		if st.SoloCategoria != "" {
			switch st.SoloCategoria {
			case CatA:
				if len(a) > 0 {
					x := a[0]
					a = a[1:]
					return x, true
				}
			case CatB:
				if len(b) > 0 {
					x := b[0]
					b = b[1:]
					return x, true
				}
			case CatC:
				if len(c) > 0 {
					x := c[0]
					c = c[1:]
					return x, true
				}
			}
			return Coche{}, false
		}

		// Si hay prioridad forzada, esa categoría va primero.
		if st.PrioridadCategoria != "" {
			switch st.PrioridadCategoria {
			case CatA:
				if len(a) > 0 {
					x := a[0]
					a = a[1:]
					return x, true
				}
			case CatB:
				if len(b) > 0 {
					x := b[0]
					b = b[1:]
					return x, true
				}
			case CatC:
				if len(c) > 0 {
					x := c[0]
					c = c[1:]
					return x, true
				}
			}
			// Si la prioritaria está vacía, seguimos orden normal.
		}

		// Orden normal A -> B -> C (A es la más prioritaria).
		if len(a) > 0 {
			x := a[0]
			a = a[1:]
			return x, true
		}
		if len(b) > 0 {
			x := b[0]
			b = b[1:]
			return x, true
		}
		if len(c) > 0 {
			x := c[0]
			c = c[1:]
			return x, true
		}
		return Coche{}, false
	}

	// Encola el coche en su cola por categoría.
	push := func(car Coche) {
		switch car.Categoria {
		case CatA:
			a = append(a, car)
		case CatB:
			b = append(b, car)
		default:
			c = append(c, car)
		}
	}

	// Intenta resolver dequeues en espera mientras haya coches.
	flushWaiting := func() {
		for len(waiting) > 0 {
			req := waiting[0]
			car, ok := pick(req.state)
			if !ok {
				return
			}
			waiting = waiting[1:]
			req.reply <- car
		}
	}

	// Intenta meter encolados pendientes si hay hueco.
	flushPending := func() {
		for len(pending) > 0 && totalLen() < q.capacity {
			r := pending[0]
			pending = pending[1:]
			push(r.car)
			close(r.reply)
		}
	}

	for {
		select {
		case r := <-q.enq:
			// Si hay hueco, encolamos; si no, guardamos como pendiente.
			if totalLen() < q.capacity {
				push(r.car)
				close(r.reply)
				// Si alguien estaba esperando, intentamos servirle.
				if hasAny() && len(waiting) > 0 {
					flushWaiting()
				}
			} else {
				pending = append(pending, r)
			}

		case r := <-q.deq:
			// Si hay coches, sacamos según el estado. Si no hay, se queda esperando.
			if hasAny() {
				car, ok := pick(r.state)
				if ok {
					r.reply <- car
					// Tras sacar, puede haber hueco: metemos pendientes.
					flushPending()
				} else {
					// No hay coche elegible para ese estado (p.ej. SOLO B y no hay B).
					waiting = append(waiting, r)
				}
			} else {
				waiting = append(waiting, r)
			}
		}
	}
}
