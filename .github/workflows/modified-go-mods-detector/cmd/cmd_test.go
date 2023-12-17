package cmd

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"sort"
	"testing"

	"modified-go-mods-detector/internal"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	err := os.Chdir("../../../../")
	require.NoError(t, err)

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			name:           "Push on main event detected",
			args:           []string{"main", "4675207", "89d0267"},
			expectedOutput: `[]`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stdoutReader, originalStdout, stdoutWriter := setupStdoutReader(t)

			err := Run(tt.args)
			require.NoError(t, err)

			output := readStdOut(t, stdoutReader, stdoutWriter)

			var outputJSON []internal.GoMod
			err = json.Unmarshal([]byte(output), &outputJSON)
			require.NoError(t, err)
			sort.Slice(outputJSON, func(i, j int) bool {
				return outputJSON[i].GoModName < outputJSON[j].GoModName
			})

			var expectedJSON []internal.GoMod
			err = json.Unmarshal([]byte(tt.expectedOutput), &expectedJSON)
			require.NoError(t, err)
			sort.Slice(expectedJSON, func(i, j int) bool {
				return expectedJSON[i].GoModName < expectedJSON[j].GoModName
			})

			eq := reflect.DeepEqual(outputJSON, expectedJSON)
			require.True(t, eq)
			os.Stdout = originalStdout
		})
	}
}

func setupStdoutReader(t *testing.T) (*os.File, *os.File, *os.File) {
	originalStdout := os.Stdout

	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w

	return r, originalStdout, w
}

func readStdOut(t *testing.T, r *os.File, w *os.File) string {
	w.Close()

	output, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(output)
}
