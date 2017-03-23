package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hoisie/mustache"
)

func main() {

	in := flag.String("i", "poem.txt", "input file")
	out := flag.String("o", "poem.html", "output file")
	template := flag.String("template", "index.mustache", "template file")
	center := flag.Bool("center", false, "center text")

	flag.Parse()

	outFile, err := os.Create(*out)
	if err != nil {
		fmt.Printf("cannot create '%s' file\n", *out)
		return
	}
	defer outFile.Close()

	inFile, err := os.Open(*in)
	if err != nil {
		fmt.Printf("cannot open '%s' file\n", *in)
		return
	}
	defer inFile.Close()

	b, err := ioutil.ReadAll(inFile)
	if err != nil {
		fmt.Printf("cannot read '%s' content\n", *in)
		return
	}

	poem := strings.Join(strings.Split(string(b), "\n"), "<br>\n")

	var templateData = map[string]string{
		"poem": poem,
	}

	if *center {
		templateData["text-center"] = "text-center"
	}

	html := mustache.RenderFile(*template, templateData)

	outFile.WriteString(html)
}
