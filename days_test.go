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
	"1-1":  "1018",
	"1-2":  "5815",
	"2-1":  "18952700150",
	"2-2":  "28858486244",
	"3-1":  "17324",
	"3-2":  "171846613143331",
	"4-1":  "1543",
	"4-2":  "9038",
	"5-1":  "623",
	"5-2":  "353507173555373",
	"6-1":  "6417439773370",
	"6-2":  "11044319475191",
	"7-1":  "1553",
	"7-2":  "15811946526915",
	"8-1":  "66640",
	"8-2":  "78894156",
	"9-1":  "4776100539",
	"9-2":  "1476550548",
	"10-1": "425",
	// "10-2": "",
	"11-1": "508",
	"11-2": "315116216513280",
	"12-1": "575",
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
