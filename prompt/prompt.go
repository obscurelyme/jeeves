package prompt

import (
	"errors"

	"github.com/manifoldco/promptui"
)

// Executes a basic prompt for quick end-user input
func QuickPrompt(label string) (string, error) {
	p := promptui.Prompt{
		Label: label,
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("input is required")
			}
			return nil
		},
	}

	return p.Run()
}
