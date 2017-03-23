package charrnn

import "testing"

func TestMain(t *testing.T) {
	s, err := NewSampler("/home/jam/torch/install/bin/th", "/home/jam/lab/char-rnn", "/home/jam/lab/abotllinaire/lm_lstm_epoch24.91_1.1464.t7")
	if err != nil {
		t.Fatal(err)
	}

	seed := `Des poèmes toujours des poèmes,
avec de la prose, des vers long,
ignorant les alexandrins 
et le conformisme du rythme monotone de l'octosyllabe
`

	ss, err := s.Run(100, 0.9, seed)
	t.Log(err)

	for _, sss := range ss {
		t.Log(sss)
	}

	t.Log("OK")
}
