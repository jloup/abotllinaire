package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/jloup/abotllinaire/char-rnn"
)

func main() {
	sampler, err := charrnn.NewSampler("/home/jam/torch/install/bin/th", "/home/jam/lab/char-rnn", "/home/jam/lab/abotllinaire/lm_lstm_epoch24.91_1.1464.t7")
	if err != nil {
		return
	}

	outFile, err := os.Create("writer_out.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	stopChan := make(chan struct{})

	go func() {
		for i := 0; i < 500; i += 1 {
			seed := rand.Intn(1000)
			temperature := rand.Float64()*(0.65-0.50) + 0.50

			fmt.Printf("batch #%v with seed %v temp %v\n", i, seed, temperature)

			err = sampler.PipedRun(100000, temperature, "", seed, outFile)
			if err != nil {
				fmt.Println(err)
			}
		}
		close(stopChan)
	}()

	for {
		select {
		case <-stopChan:
			return
		}
	}
}
