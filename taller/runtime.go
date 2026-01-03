package main

import "sync"

var (
	startOnce sync.Once

	incomingMsgCh chan string
	stateCodeCh   chan int

	// Canal para pedir el estado actual desde las fases.
	stateQueryCh chan stateRequest
)

func initRuntime() {
	incomingMsgCh = make(chan string, 32)
	stateCodeCh = make(chan int, 16)
	stateQueryCh = make(chan stateRequest)

	go parseIncoming(incomingMsgCh, stateCodeCh)
	go controller(stateCodeCh, stateQueryCh)
}

func dispatch(msg string) {
	startOnce.Do(initRuntime)
	incomingMsgCh <- msg
}

// getState devuelve una copia del estado actual (consulta al controlador).
func getState() TallerState {
	reply := make(chan TallerState, 1)
	stateQueryCh <- stateRequest{reply: reply}
	return <-reply
}
