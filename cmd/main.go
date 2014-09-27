package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/xconstruct/know"
)

func main() {
	flag.Parse()

	in := bufio.NewReader(os.Stdin)
	for {
		line, _, err := in.ReadLine()
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
		}
		answers, _ := know.Ask(string(line))
		go func() {
			for ans := range answers {
				fmt.Printf("[%s] %s: %s\n", ans.Provider, ans.Question, ans.Answer)
			}
		}()
	}
}
