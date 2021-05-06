package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CLI struct {
	reader *bufio.Reader
}

func Init() *CLI {
	r := bufio.NewReader(os.Stdin)
	return &CLI{reader: r}
}

func (c *CLI) Prompt(msg string) string {
	fmt.Println(msg)
	return c.GetInput()
}

func (c *CLI) GetInput() string {
	fmt.Print("> ")
	line, _ := c.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func (c *CLI) Print(msg string) {
	fmt.Println(msg)
}
