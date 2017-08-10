package api

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

// We have two version of the same file - one is the lower case of the fisrt one - it will speed up the grep.
func SearchVerse(word string, nbOfVerses int, filepath string, filepath_lower string) (string, error) {
	cmd := exec.Command("grep", "-E", "-b", fmt.Sprintf("(^|\\W)%s($|\\W)", word), filepath_lower)

	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// select a random line
	source := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(source)
	offset := randGen.Intn(len(out) - 2)

	// in case offset point to end of line
	for offset < len(out)-1 && out[offset] == '\n' {
		offset += 1
	}

	var start, end int
	for end = offset; end < len(out) && out[end] != '\n'; end += 1 {
	}

	for start = offset; start > 0 && out[start] != '\n'; start -= 1 {
	}

	if out[start] == '\n' {
		start += 1
	}

	line := out[start:end]
	reg := regexp.MustCompile("^([0-9]+):")

	m := reg.FindSubmatch(line)
	if m == nil {
		return "", fmt.Errorf("could not find byte offset in '%s'", string(line))
	}

	readOffset, err := strconv.ParseInt(string(m[1]), 10, 64)
	if err != nil {
		return "", fmt.Errorf("cannot convert %v to int (line %s)", string(m[1]), string(line))
	}

	_, err = f.Seek(readOffset, 0)
	if err != nil {
		return "", err
	}

	verses := make([]byte, 0)

	rd := bufio.NewReader(f)

	for i := 0; i < nbOfVerses; i += 1 {
		line, err := rd.ReadBytes('\n')
		if err != nil {
			return string(verses), err
		}
		verses = append(verses, line...)
	}

	verses = bytes.TrimSpace(verses)

	// replace , : ; byt a '.'
	regPunct := regexp.MustCompile("((?: :)|(?: ;)|[,'\"])$")
	verses = regPunct.ReplaceAll(verses, []byte("."))

	// put a final '.' if needed
	regWord := regexp.MustCompile("\\w$")
	verses = regWord.ReplaceAll(verses, []byte{verses[len(verses)-1], '.'})

	return string(verses), nil
}
