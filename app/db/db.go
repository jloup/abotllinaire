package db

import (
	"os"

	"github.com/go-redis/redis"
	"github.com/jloup/utils"
	"gopkg.in/mgo.v2"
)

const (
	SERVER_URL = "127.0.0.1:27017"
	VERSE_C    = "verse"
	POEM_C     = "poems"
)

var (
	SESSION  *mgo.Session
	DB       *mgo.Database
	DbError  = utils.NewErrorFlag("DbError")
	notFound = utils.NewErrorFlag("NotFound")
	isDup    = utils.NewErrorFlag("Duplicate")
	NotFound = utils.Join("NotFound", DbError, notFound)
	IsDup    = utils.Join("IsDup", DbError, isDup)
	log      = utils.StandardL().WithField("module", "db")

	Redis *redis.Client
)

func ensureIndexes() {
	/*if err := DB.C(DOCUMENT_C).EnsureIndex(mgo.Index{Key: []string{"scheme", "host", "phash"}, Unique: true}); err != nil {
		log.Fatalln(err)
	}*/
}

func InitDB() {
	SESSION, err := mgo.Dial(SERVER_URL)
	if err != nil {
		log.Errorf("SESSION MGO ERROR %s\n", err)
		os.Exit(-1)
	}
	SESSION.SetMode(mgo.Monotonic, true)
	SESSION.SetSafe(&mgo.Safe{})

	DB = SESSION.DB("abotllinaire")
	ensureIndexes()

	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = Redis.Ping().Result()

	if err != nil {
		log.Errorf("redis init error %s", err)
		os.Exit(-1)
	}
}

func NewPoemOp() DbOp {
	return DbOp{c: DB.C(POEM_C)}
}
