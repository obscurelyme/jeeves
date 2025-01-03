package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
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

// Executes a basic select prompt for quick end-user input
func SelectPrompt[T interface{}](label string, options []T, template *promptui.SelectTemplates) (T, error) {
	selectPrompt := promptui.Select{
		Label:     label,
		Items:     options,
		Templates: template,
	}

	index, _, err := selectPrompt.Run()

	return options[index], err
}
