package promp

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Promp struct{}

func NewPromp() *Promp {
	return &Promp{}
}

func (p *Promp) AgrFlag() (string, int, error) {
	command := flag.String("c", "", "Command to run/execute")
	flag.StringVar(command, "command", "", "Command to run/execute")
	bufferSize := flag.Int("b", 64, "Size of the command history buffer")
	flag.IntVar(bufferSize, "bufferSize", 64, "Size of the command history buffer")
	flag.Parse()

	return *command, *bufferSize, nil
}

func (p *Promp) CommandPrompt(history []string, files []*os.File, historyIndex *int, bufferSize int) (string, []string) {
	if len(files) > 0 {
		fmt.Printf("\nLoaded Files:\n")
		for i := 0; i < len(files); i++ {
			fmt.Printf("\t[%d] %s\n", i, files[i].Name())
		}
	}

	fmt.Printf("\n> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	args := strings.Fields(input)

	if input != "" && input != "history" {
		if !(*historyIndex > bufferSize-1) {
			history[*historyIndex] = input
			*historyIndex++
		} else {
			history = append(history[:0], history[1:]...)
			history[bufferSize-1] = input
		}
	}

	if len(args) == 0 {
		return "", nil
	}
	return args[0], args[1:]
}
