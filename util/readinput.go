package util

import (
	"bytes"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Iterator = func(s []byte)

func filename(dayPart string, isSample bool) string {
	folder := filepath.Join("inputs")
	if isSample {
		folder = filepath.Join("samples")
	}

	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join(folder, fmt.Sprintf("%s.txt", dayPart))
	}
	prefixPath := filepath.Dir(fileName)
	basePath := filepath.Join(prefixPath, "..", folder)

	return filepath.Join(basePath, fmt.Sprintf("%s.txt", dayPart))
}

func ReadBytes(dayPart string, isSample bool) ([]byte, error) {
	return os.ReadFile(filename(dayPart, isSample))
}

func ReadStrings(dayPart string, isSample bool, delim string) (iter.Seq[string], error) {
	partIter, err := read(filename(dayPart, isSample), delim)
	if err != nil {
		return nil, err
	}

	return func(yield func(string) bool) {
		for part := range partIter {
			partStr := string(part)
			if !yield(partStr) {
				return
			}
		}
	}, nil
}

func ReadIntLists(dayPart string, isSample bool, delim string) (iter.Seq[[]int], error) {
	partIter, err := read(filename(dayPart, isSample), delim)
	if err != nil {
		return nil, err
	}

	return func(yield func([]int) bool) {
		for part := range partIter {
			var list []int
			for number := range strings.SplitSeq(string(part), ",") {
				partInt, err := strconv.Atoi(number)
				if err != nil {
					return
				}
				list = append(list, partInt)
			}

			if !yield(list) {
				return
			}
		}
	}, nil
}

func ReadInts(dayPart string, isSample bool, delim string) (iter.Seq[int], error) {
	partIter, err := read(filename(dayPart, isSample), delim)
	if err != nil {
		return nil, err
	}

	return func(yield func(int) bool) {
		for part := range partIter {
			partInt, err := strconv.Atoi(string(part))
			if err != nil {
				return
			}
			if !yield(partInt) {
				return
			}
		}
	}, nil
}

func read(file string, delim string) (iter.Seq[[]byte], error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return func(yield func([]byte) bool) {
		for _, row := range bytes.Split(b, []byte(delim)) {
			if !yield(row) {
				return
			}
		}
	}, nil
}
