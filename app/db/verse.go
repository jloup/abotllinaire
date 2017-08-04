package db

import (
	"crypto/md5"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Verse struct {
	Id   bson.ObjectId `bson:"_id"`
	S    string
	Hash string
}

func (v Verse) MD5Hash() [md5.Size]byte {
	return md5.Sum([]byte(v.S))
}

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

func GetPoems(query interface{}, poems *[]Poem) error {
	op := NewPoemOp()

	op.Find(query).All(poems)

	return op.Err()
}
