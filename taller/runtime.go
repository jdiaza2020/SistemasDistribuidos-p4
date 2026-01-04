package main

import (
	"sync"
	"time"
)

var (
	startOnce sync.Once

	incomingMsgCh chan string
	stateCodeCh   chan int
	stateQueryCh  chan stateRequest

	logCh     chan LogEvent
	startTime time.Time
)

func initRuntime() {
	incomingMsgCh = make(chan string, 32)
	stateCodeCh = make(chan int, 16)
	stateQueryCh = make(chan stateRequest)
	logCh = make(chan LogEvent, 1024)

	startTime = time.Now()

	go parseIncoming(incomingMsgCh, stateCodeCh)
	go controller(stateCodeCh, stateQueryCh)
	go runLogger(logCh)

	// Config de desarrollo (luego en tests se pasará otro).
	cfg := DefaultConfig()
	go startSimulation(startTime, logCh, cfg)
}

func dispatch(msg string) {
	startOnce.Do(initRuntime)
	incomingMsgCh <- msg
}

func getState() TallerState {
	reply := make(chan TallerState, 1)
	stateQueryCh <- stateRequest{reply: reply}
	return <-reply
}
