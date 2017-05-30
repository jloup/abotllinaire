package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/jloup/abotllinaire/char-rnn"
)

func main() {
	sampler, err := charrnn.NewSampler("/home/jam/torch/install/bin/th", "/home/jam/lab/char-rnn", "/home/jam/lab/abotllinaire/lm_lstm_epoch24.91_1.1464.t7")
	if err != nil {
		return
	}

	f, _ := os.Open("seed.txt")
	b, _ := ioutil.ReadAll(f)

	temp, _ := strconv.ParseFloat(os.Args[2], 64)
	len, _ := strconv.ParseUint(os.Args[1], 10, 64)

	verses, err := sampler.Run(uint16(len), temp, string(b))
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, verse := range verses {
		fmt.Println(verse)
	}

}
