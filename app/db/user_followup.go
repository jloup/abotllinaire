package db

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type UserAction struct {
	Type UserActionType
	Meta string
}

type UserActionType string

const NoneAction UserActionType = "N"
const GreetingsAction UserActionType = "g"
const FreestyleAction UserActionType = "f"
const SubjectAction UserActionType = "s"

func (f UserAction) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func BuildUserActionKey(senderId string) string {
	return fmt.Sprintf("%s:action", senderId)
}

func GetLastUserAction(senderId string) (UserAction, error) {
	s, err := Redis.Get(BuildUserActionKey(senderId)).Bytes()

	if err == redis.Nil {
		return UserAction{NoneAction, ""}, nil
	} else if err != nil {
		return UserAction{NoneAction, ""}, fmt.Errorf("redis error %s", err)
	}

	var action UserAction
	err = json.Unmarshal(s, &action)
	return action, err
}

func SetLastUserAction(senderId string, action UserAction) error {
	return Redis.Set(BuildUserActionKey(senderId), action, 0).Err()
}
