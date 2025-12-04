package test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectations = map[string]string{
	"1-1": "1018",
	"1-2": "5815",
	"2-1": "18952700150",
	"2-2": "28858486244",
	"3-1": "17324",
	"3-2": "171846613143331",
	"4-1": "1543",
	"4-2": "9038",
}

func TestDays(t *testing.T) {
	for day, expect := range expectations {
		t.Run(day, func(t *testing.T) {
			t.Parallel()
			runCmd := exec.Command("go", "run", ".")
			runCmd.Dir = filepath.Join("days", day)
			output, err := runCmd.CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
			}

			assert.NoError(t, err)
			assert.Equal(t, expect, strings.TrimRight(string(output), "\n"), fmt.Sprintf("Day %s", day))
		})
	}
}
