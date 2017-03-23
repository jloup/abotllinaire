package api

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	s, err := NewWorkerPool(3, "/home/jam/lab/char-rnn", "/home/jam/lab/abotllinaire/lm_lstm_epoch24.91_1.1464.t7")
	if err != nil {
		t.Fatal(err)
	}

	s.Run()

	go func() {
		res, err := s.Request(0.65, 1000, "Fin de soirée, bord de mer, casino métallique")
		fmt.Println("ERR", err)
		for _, verse := range res {
			fmt.Println(verse)

		}
		fmt.Println("")
	}()
	go func() {
		res, err := s.Request(0.66, 200, "L'amour au chinois")
		fmt.Println("ERR", err)
		for _, verse := range res {
			fmt.Println(verse)

		}
		fmt.Println("")
	}()
	go func() {
		res, err := s.Request(0.68, 200, "L'amour au chinois")
		fmt.Println("ERR", err)
		for _, verse := range res {
			fmt.Println(verse)

		}
		fmt.Println("")
	}()

	time.Sleep(30 * time.Second)
}
