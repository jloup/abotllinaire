package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jloup/abotllinaire/app/db"

	"gopkg.in/mgo.v2/bson"
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

type ApiGetFacebookVerses struct {
	In struct {
		FromId bson.ObjectId
		Count  int
	}

	Out []struct {
		Id          string  `json:"id"`
		Seed        string  `json:"seed"`
		Verses      string  `json:"verses"`
		Temperature float32 `json:"temperature"`
	}

	Error
}

func (a *ApiGetFacebookVerses) GetRawOut() interface{} {
	if a.err != nil {
		return nil
	}

	return a.Out
}

func (a *ApiGetFacebookVerses) GetJSONOut() json.RawMessage {
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

func (a *ApiGetFacebookVerses) SetFromIdParam(param string) {
	if a.err != nil {
		return
	}

	if param == "" || !bson.IsObjectIdHex(param) {
		a.Error = Error{INVALID_REQUEST, fmt.Errorf("invalid id")}
	}

	a.In.FromId = bson.ObjectIdHex(param)
}

func (a *ApiGetFacebookVerses) SetCountParam(param string) {
	if a.err != nil {
		return
	}

	count, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		a.Error = Error{INVALID_REQUEST, fmt.Errorf("param '%v' not a int %v", param, err)}
		return
	}

	a.In.Count = int(count)
}

func (a *ApiGetFacebookVerses) Run() {
	if a.err != nil {
		return
	}

	op := db.NewPoemOp()

	var poems []db.Poem

	query := bson.M{}
	if a.In.FromId.Hex() != "" {
		query["_id"] = bson.M{"$lt": a.In.FromId}
	}

	op.Find(query).Sort("-_id").Limit(a.In.Count).All(&poems)

	if op.Err() != nil {
		a.Error = Error{INTERNAL_ERROR, fmt.Errorf("cannot fetch verses %v", op.Err())}
		return
	}

	for _, poem := range poems {
		a.Out = append(a.Out, struct {
			Id          string  `json:"id"`
			Seed        string  `json:"seed"`
			Verses      string  `json:"verses"`
			Temperature float32 `json:"temperature"`
		}{
			poem.Id.Hex(),
			poem.Seed,
			poem.Content,
			poem.Temperature,
		})
	}
}
