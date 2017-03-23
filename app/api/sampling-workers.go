package api

import (
	"fmt"

	"github.com/jloup/abotllinaire/char-rnn"
)

type SamplerWorkerResponse struct {
	WorkerId string
	Verses   []string
	Error    error
}

type SamplerCommand struct {
	Temperature float64
	Length      uint16
	Seed        string

	ChanOut chan SamplerWorkerResponse
}

type worker struct {
	Id       string
	Sampler  charrnn.Sampler
	ChanIn   chan SamplerCommand
	ChanStop chan struct{}
}

func (w *worker) Run() {
	for {

		select {
		case cmd := <-w.ChanIn:
			verses, err := w.Sampler.Run(cmd.Length, cmd.Temperature, cmd.Seed)
			cmd.ChanOut <- SamplerWorkerResponse{w.Id, verses, err}

		case <-w.ChanStop:
			break
		}
	}

}

type WorkerPool struct {
	workers  []*worker
	ChanStop chan struct{}
	ChanIn   chan SamplerCommand
}

func NewWorkerPool(n int, torchPath, workingDir, modelFilePath string) (*WorkerPool, error) {
	pool := WorkerPool{ChanStop: make(chan struct{}), ChanIn: make(chan SamplerCommand)}

	pool.workers = make([]*worker, n)

	for i := 0; i < n; i++ {
		sampler, err := charrnn.NewSampler(torchPath, workingDir, modelFilePath)
		if err != nil {
			return nil, fmt.Errorf("cannot create sampler %v", err)
		}

		pool.workers[i] = &worker{Id: fmt.Sprintf("worker#%v", i), Sampler: sampler, ChanIn: pool.ChanIn, ChanStop: pool.ChanStop}
	}

	return &pool, nil
}

func (w *WorkerPool) Run() {
	for i, _ := range w.workers {

		go w.workers[i].Run()
	}
}

func (w *WorkerPool) Stop() {
	close(w.ChanStop)
}

func (w *WorkerPool) Request(temperature float64, length uint16, seed string) ([]string, error) {
	cmd := SamplerCommand{Temperature: temperature, Length: length, Seed: seed, ChanOut: make(chan SamplerWorkerResponse, 1)}

	w.ChanIn <- cmd

	response := <-cmd.ChanOut

	close(cmd.ChanOut)

	return response.Verses, response.Error
}
