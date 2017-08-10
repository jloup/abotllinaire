package charrnn

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
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

func (s *Sampler) PipedRun(length uint32, temperature float64, primetext string, seed int, out io.Writer) error {
	var stderr bytes.Buffer

	if temperature <= 0.0 || temperature > 1.0 {
		return fmt.Errorf("temperature should be [0.0, 1.0]")
	}

	args := []string{
		"sample.lua",
		s.ModelFilePath,
		"-gpuid", "-1",
		"-seed", strconv.Itoa(seed),
		"-verbose", "0",
		"-temperature", strconv.FormatFloat(temperature, 'f', 2, 64),
		"-length", strconv.Itoa(int(length)),
	}

	primetext = strings.TrimSpace(primetext)
	if primetext != "" {
		args = append(args, "-primetext", primetext)
	}

	cmd := exec.Command(s.TorchPath, args...)
	cmd.Dir = s.Dir

	cmd.Stdout = out
	cmd.Stderr = &stderr

	return cmd.Run()
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
		"-seed", strconv.Itoa(rand.Int()),
		"-gpuid", "-1",
		"-verbose", "0",
		"-temperature", strconv.FormatFloat(temperature, 'f', 2, 64),
		"-length", strconv.Itoa(int(length)),
	}

	startIndex := 0

	seed = strings.TrimSpace(seed)
	if seed != "" {
		startIndex = len(strings.Split(seed, "\n"))
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

	if len(verses) <= startIndex {
		return []string{}, nil
	}

	return verses[startIndex : len(verses)-1], nil
}
