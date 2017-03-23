package db

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type FacebookMessageSeqStatus string

const SeqReceived FacebookMessageSeqStatus = "rec"
const SeqResponded FacebookMessageSeqStatus = "resp"

func (f FacebookMessageSeqStatus) MarshalBinary() ([]byte, error) {
	return []byte(f), nil
}

const SEQ_TTL = time.Hour * 1

func BuildSeqKey(senderId string, seq int) string {
	return fmt.Sprintf("%s:%v", senderId, seq)
}

func FacebookSeqExists(senderId string, seq int) (bool, error) {
	_, err := Redis.Get(BuildSeqKey(senderId, seq)).Result()

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("redis error %s", err)
	}

	return true, nil
}

func SetFacebookSeq(senderId string, seq int, status FacebookMessageSeqStatus) error {
	return Redis.Set(BuildSeqKey(senderId, seq), status, SEQ_TTL).Err()
}
