package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func gitDiff(ctx context.Context, gitRefBefore string, gitRefAfter string) ([]string, error) {
	cmdArgs := []string{"diff", "--name-only", gitRefBefore, gitRefAfter}

	cmd := exec.CommandContext(ctx, "git", cmdArgs...)

	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		errs := errOut.String()
		err2 := errors.New(errs)
		return nil, fmt.Errorf("failed to get git diff: %w: %w", err, err2)
	}

	errs := errOut.String()
	if errs != "" {
		return nil, fmt.Errorf("failed to get git diff: %v", errs)
	}

	s := out.String()
	return strings.Split(s, "\n"), nil
}
