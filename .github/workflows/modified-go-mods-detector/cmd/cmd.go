package cmd

import (
	"encoding/json"
	"fmt"

	"modified-go-mods-detector/internal"
)

func Run(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("missing github refs")
	}

	d, err := internal.GetModifiedGoModsBetweenCommits(args[1], args[2])
	if err != nil {
		return err
	}

	j, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}
	fmt.Print(string(j))

	return nil
}
