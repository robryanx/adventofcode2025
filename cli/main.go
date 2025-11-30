package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

//go:embed template.go.tmpl
var templateFS embed.FS

type TemplateData struct {
	Year int
	Day  int
	Part int
}

func main() {
	dayFlag := flag.String("day", "1", "Day number or inclusive range (e.g. 9 or 1-9)")
	partFlag := flag.String("part", "1-2", "Part number or inclusive range (e.g. 1 or 1-2)")
	submitFlag := flag.Bool("submit", false, "Run local solver and submit answers")
	yearFlag := flag.Int("year", 2025, "Event year (e.g. 2024)")
	flag.Parse()

	loginKey := os.Getenv("ADVENT_OF_CODE_SESSION")
	if loginKey == "" {
		panic("ADVENT_OF_CODE_SESSION environment variable is not set")
	}

	days, err := parseSelection(*dayFlag, 1, 25)
	if err != nil {
		panic(err)
	}

	parts, err := parseSelection(*partFlag, 1, 2)
	if err != nil {
		panic(err)
	}

	year := *yearFlag
	yearDir := strconv.Itoa(year)
	if year == 2025 {
		yearDir = "."
	}

	for _, day := range days {
		if err := processDay(year, yearDir, day, loginKey, *submitFlag, parts); err != nil {
			panic(err)
		}
	}
}

func processDay(year int, yearDir string, day int, loginKey string, submit bool, parts []int) error {
	// Ensure input exists
	if err := ensureInput(year, yearDir, day, loginKey); err != nil {
		return err
	}

	// Ensure samples exist
	if err := ensureSamples(year, yearDir, day, parts, loginKey); err != nil {
		fmt.Printf("Warning: failed to fetch samples: %v\n", err)
	}

	for _, part := range parts {
		if err := ensureDayTemplate(year, yearDir, day, part); err != nil {
			return err
		}

		if submit {
			answer, err := runSolution(yearDir, day, part)
			if err != nil {
				return err
			}
			fmt.Printf("Day %d part %d answer: %s\n", day, part, answer)

			err = submitAnswer(year, day, part, answer, loginKey)
			if err != nil {
				fmt.Printf("Day %d part %d submission failed: %v\n\n", day, part, err)
			}
		}
	}

	return nil
}

func ensureInput(year int, yearDir string, day int, loginKey string) error {
	inputDir := filepath.Join(yearDir, "inputs")
	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		return fmt.Errorf("ensure dir %s: %w", inputDir, err)
	}

	filename := fmt.Sprintf("%d.txt", day)
	path := filepath.Join(inputDir, filename)

	if _, err := os.Stat(path); err == nil {
		return nil // already exists
	}

	fmt.Printf("Fetching input for day %d...\n", day)
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	body, err := fetchWithCookie(url, loginKey)
	if err != nil {
		return fmt.Errorf("fetch input %d: %w", day, err)
	}

	return os.WriteFile(path, body, 0o644)
}

func ensureDayTemplate(year int, yearDir string, day, part int) error {
	dayDir := filepath.Join(yearDir, "days", fmt.Sprintf("%d-%d", day, part))
	mainGoPath := filepath.Join(dayDir, "main.go")

	if _, err := os.Stat(mainGoPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(dayDir, 0o755); err != nil {
		return fmt.Errorf("create day directory %s: %w", dayDir, err)
	}

	templateContent, err := generateMainGoTemplate(year, day, part)
	if err != nil {
		return err
	}

	if err := os.WriteFile(mainGoPath, []byte(templateContent), 0o644); err != nil {
		return fmt.Errorf("write main.go template: %w", err)
	}

	fmt.Printf("Created template: %s\n", mainGoPath)
	return nil
}

func generateMainGoTemplate(year, day, part int) (string, error) {
	tmplContent, err := templateFS.ReadFile("template.go.tmpl")
	if err != nil {
		return "", fmt.Errorf("read template file: %w", err)
	}

	tmpl, err := template.New("main.go").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	data := TemplateData{
		Year: year,
		Day:  day,
		Part: part,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

func ensureSamples(year int, yearDir string, day int, parts []int, loginKey string) error {
	samplesDir := filepath.Join(yearDir, "samples")
	if err := os.MkdirAll(samplesDir, 0o755); err != nil {
		return fmt.Errorf("ensure dir %s: %w", samplesDir, err)
	}

	needFetch := false
	for _, part := range parts {
		filename := fmt.Sprintf("%d-%d.txt", day, part)
		path := filepath.Join(samplesDir, filename)
		if _, err := os.Stat(path); err != nil {
			needFetch = true
			break
		}
	}

	if !needFetch {
		return nil
	}

	fmt.Printf("Fetching description for day %d to extract samples...\n", day)
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d", year, day)
	body, err := fetchWithCookie(url, loginKey)
	if err != nil {
		return fmt.Errorf("fetch description %d: %w", day, err)
	}

	htmlContent := string(body)
	part1Content := htmlContent
	part2Content := ""

	split := strings.SplitN(htmlContent, "--- Part Two ---", 2)
	if len(split) == 2 {
		part1Content = split[0]
		part2Content = split[1]
	}

	if s1 := extractSample(part1Content); s1 != "" {
		if err := os.WriteFile(filepath.Join(samplesDir, fmt.Sprintf("%d-1.txt", day)), []byte(s1), 0o644); err != nil {
			return err
		}
	}

	if part2Content != "" {
		if s2 := extractSample(part2Content); s2 != "" {
			if err := os.WriteFile(filepath.Join(samplesDir, fmt.Sprintf("%d-2.txt", day)), []byte(s2), 0o644); err != nil {
				return err
			}
		}
	}

	return nil
}

var sampleRE = regexp.MustCompile(`(?is)<pre><code>(.*?)</code></pre>`)

func extractSample(htmlContent string) string {
	matches := sampleRE.FindStringSubmatch(htmlContent)
	if len(matches) < 2 {
		return ""
	}
	sample := matches[1]
	sample = html.UnescapeString(sample)
	tagRE := regexp.MustCompile(`<[^>]*>`)
	sample = tagRE.ReplaceAllString(sample, "")
	return strings.TrimSpace(sample)
}

func parseSelection(value string, min, max int) ([]int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("selection is empty")
	}

	if strings.Contains(value, "-") {
		parts := strings.SplitN(value, "-", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range %q", value)
		}
		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", parts[0], err)
		}
		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", parts[1], err)
		}
		if start < min || end < min || end > max || end < start {
			return nil, fmt.Errorf("range %q outside supported bounds %d-%d", value, min, max)
		}
		days := make([]int, 0, end-start+1)
		for d := start; d <= end; d++ {
			days = append(days, d)
		}
		return days, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid number %q: %w", value, err)
	}
	if n < min || n > max {
		return nil, fmt.Errorf("value %d outside supported bounds %d-%d", n, min, max)
	}
	return []int{n}, nil
}

func runSolution(yearDir string, day, part int) (string, error) {
	target := fmt.Sprintf("./%s/days/%d-%d", yearDir, day, part)
	cmd := exec.Command("go", "run", target)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("run solver for day %d part %d: %w\n%s", day, part, err, string(output))
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			return line, nil
		}
	}
	return "", fmt.Errorf("solver for day %d part %d produced no output", day, part)
}

func fetchWithCookie(urlStr, key string) ([]byte, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "github.com/robryanx/adventofcode/cli")
	req.Header.Set("Cookie", "session="+key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func submitAnswer(year, day, part int, answer, loginKey string) error {
	fmt.Printf("Submitting day %d part %d answer: %s\n", day, part, answer)
	urlStr := fmt.Sprintf("https://adventofcode.com/%d/day/%d/answer", year, day)

	data := url.Values{}
	data.Set("level", strconv.Itoa(part))
	data.Set("answer", answer)

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("build answer request: %w", err)
	}
	req.Header.Set("User-Agent", "github.com/robryanx/adventofcode/cli")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "session="+loginKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("submit answer: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read answer response: %w", err)
	}

	content := string(body)
	if strings.Contains(content, "That's the right answer") {
		fmt.Printf("Day %d part %d: CORRECT!\n", day, part)
		return nil
	}
	if strings.Contains(content, "That's not the right answer") {
		return fmt.Errorf("incorrect answer")
	}
	if strings.Contains(content, "You gave an answer too recently") {
		return fmt.Errorf("rate limited (too recent)")
	}
	if strings.Contains(content, "You don't seem to be solving the right level") {
		return fmt.Errorf("already solved or wrong level")
	}

	return fmt.Errorf("unknown response: %s...", content[:min(len(content), 300)])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
