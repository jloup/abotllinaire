package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Poem struct {
	Id          bson.ObjectId `bson:"_id"`
	Seed        string
	Temperature float32
	Content     string

	FacebookUserId string
}

func InsertPoem(poem *Poem) error {
	op := NewPoemOp()

	poem.Id = bson.NewObjectId()
	op.Insert(poem)
	if mgo.IsDup(op.err) {
		op.err = nil
	}

	op.Find(bson.M{"_id": poem.Id}).One(poem)

	return op.Err()
}
