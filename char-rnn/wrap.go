package charrnn

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Sampler struct {
	Dir           string
	ModelFilePath string
	TorchPath     string
}

func NewSampler(torchPath, dir, modelFilePath string) (Sampler, error) {
	_, err := os.Stat(dir)

	if err != nil {
		return Sampler{}, fmt.Errorf("error while stat for working directory %v", err)
	}

	_, err = os.Stat(modelFilePath)
	if err != nil {
		return Sampler{}, fmt.Errorf("error while stat for model file %v", err)
	}

	_, err = os.Stat(torchPath)
	if err != nil {
		return Sampler{}, fmt.Errorf("error while stat for torch path %v", err)
	}

	return Sampler{Dir: dir, ModelFilePath: modelFilePath, TorchPath: torchPath}, nil
}

func (s *Sampler) Run(length uint16, temperature float64, seed string) ([]string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	if temperature <= 0.0 || temperature > 1.0 {
		return []string{}, fmt.Errorf("temperature should be [0.0, 1.0]")
	}

	args := []string{
		"sample.lua",
		s.ModelFilePath,
		"-gpuid", "-1",
		"-temperature", strconv.FormatFloat(temperature, 'f', 2, 64),
		"-length", strconv.Itoa(int(length)),
	}

	if seed != "" {
		args = append(args, "-primetext", seed)
	}

	cmd := exec.Command(s.TorchPath, args...)
	cmd.Dir = s.Dir

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return []string{stderr.String()}, err
	}

	poem := strings.TrimSpace(out.String())

	verses := strings.Split(poem, "\n")

	var startIndex int
	var verse string

	for startIndex, verse = range verses {
		if strings.Contains(verse, "-----------") {
			break
		}
	}

	return verses[startIndex+1 : len(verses)-1], nil
}
