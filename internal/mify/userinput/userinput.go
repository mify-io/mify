package userinput

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type UserInput struct {
}

func (ui UserInput) AskInput(format string, a ...interface{}) (string, error) {
	println(fmt.Sprintf(format, a...))
	print("> ")

	reader := bufio.NewReader(os.Stdin)
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(data), nil
}
