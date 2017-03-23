package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var workerPool *WorkerPool

func InitWorkerPool(n int, torchPath, workingDir, modelFilePath string) error {
	var err error

	workerPool, err = NewWorkerPool(n, torchPath, workingDir, modelFilePath)

	if err != nil {
		return err
	}

	workerPool.Run()

	return nil
}

type ApiMakeVerses struct {
	In struct {
		Temperature float64
		Length      uint16
		Seed        string
	}

	Out struct {
		Verses []string `json:"verses"`
	}

	Error
}

func (a *ApiMakeVerses) GetRawOut() interface{} {
	if a.err != nil {
		return nil
	}

	return a.Out
}

func (a *ApiMakeVerses) GetJSONOut() json.RawMessage {
	if a.err != nil {
		return nil
	}

	var out json.RawMessage
	var err error
	out, err = json.Marshal(a.Out)

	if err != nil {
		a.Error = Error{INTERNAL_ERROR, err}
	}

	return out
}

func (a *ApiMakeVerses) SetTemperatureParam(param string) {
	if a.err != nil {
		return
	}

	var err error
	a.In.Temperature, err = strconv.ParseFloat(param, 64)

	if err != nil {
		a.Error = Error{INVALID_REQUEST, fmt.Errorf("not a float %v", err)}
		return
	}
}

func (a *ApiMakeVerses) SetLengthParam(param string) {
	if a.err != nil {
		return
	}

	var err error
	var l uint64
	l, err = strconv.ParseUint(param, 10, 32)

	a.In.Length = uint16(l)

	if err != nil {
		a.Error = Error{INVALID_REQUEST, fmt.Errorf("not a uint %v", err)}
		return
	}
}

func (a *ApiMakeVerses) SetSeedParam(param string) {
	if a.err != nil {
		return
	}

	a.In.Seed = param
}

func (a *ApiMakeVerses) Run() {
	if a.err != nil {
		return
	}

	var err error

	a.Out.Verses, err = workerPool.Request(a.In.Temperature, a.In.Length, a.In.Seed)

	if err != nil {
		a.Error = Error{INTERNAL_ERROR, fmt.Errorf("cannot generate verses %v", err)}
		return
	}
}
