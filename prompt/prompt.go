package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Executes a basic prompt for quick end-user input
func QuickPrompt(label string) (string, error) {
	var s string
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	return strings.TrimSpace(s), nil
}
