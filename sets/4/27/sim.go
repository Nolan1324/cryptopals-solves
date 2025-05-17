package main

import "cryptopals/internal/randx"

type Simulation struct {
	inputChan  chan []byte
	outputChan chan []byte
	errorChan  chan error
	key        []byte
}

type AttackerChans struct {
	InputChan  <-chan []byte
	OutputChan chan<- []byte
	ErrorChan  <-chan error
}

func MakeSimulation() Simulation {
	return Simulation{
		inputChan:  make(chan []byte),
		outputChan: make(chan []byte),
		errorChan:  make(chan error),
		key:        randx.RandBytes(16),
	}
}

func (s Simulation) GetAttackerChans() AttackerChans {
	return AttackerChans{
		InputChan:  s.inputChan,
		OutputChan: s.outputChan,
		ErrorChan:  s.errorChan,
	}
}

func (s Simulation) GetSenderChan() chan<- []byte {
	return s.inputChan
}

func (s Simulation) GetReceiverChans() (<-chan []byte, chan<- error) {
	return s.outputChan, s.errorChan
}
