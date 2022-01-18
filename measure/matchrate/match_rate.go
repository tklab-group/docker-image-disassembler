package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()

	err := run(args, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	if len(args) != 2 {
		return fmt.Errorf("2 args are required")
	}

	fileA := args[0]
	fileB := args[1]

	a, err := fileToMap(fileA)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", fileA, err)
	}

	b, err := fileToMap(fileB)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", fileB, err)
	}

	comparedResult := compare(a, b)

	err = outResult(comparedResult, fileA, fileB, out)
	if err != nil {
		return fmt.Errorf("failed to output: %w", err)
	}

	return nil
}

type ComparedInfo struct {
	Path    string
	SizeInA int
	SizeInB int
}

type ComparedResult struct {
	Both  []*ComparedInfo
	OnlyA []*ComparedInfo
	OnlyB []*ComparedInfo
}

func compare(a map[string]int, b map[string]int) *ComparedResult {
	result := &ComparedResult{
		Both:  make([]*ComparedInfo, 0),
		OnlyA: make([]*ComparedInfo, 0),
		OnlyB: make([]*ComparedInfo, 0),
	}

	for path, sizeInA := range a {
		sizeInB, ok := b[path]
		if ok {
			comparedInfo := &ComparedInfo{
				Path:    path,
				SizeInA: sizeInA,
				SizeInB: sizeInB,
			}
			result.Both = append(result.Both, comparedInfo)
			delete(b, path)
		} else {
			comparedInfo := &ComparedInfo{
				Path:    path,
				SizeInA: sizeInA,
				SizeInB: -1,
			}
			result.OnlyA = append(result.OnlyA, comparedInfo)
		}
	}

	for path, sizeInB := range b {
		comparedInfo := &ComparedInfo{
			Path:    path,
			SizeInA: -1,
			SizeInB: sizeInB,
		}
		result.OnlyB = append(result.OnlyB, comparedInfo)
	}

	return result
}

func fileToMap(path string) (map[string]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	m := map[string]int{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		split := strings.Split(line, " ")
		if len(split) != 2 {
			return nil, fmt.Errorf("failed to parse: `%s`", line)
		}

		path := split[0]
		dataSize, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse as datasize: line: `%s` err: %w", line, err)
		}

		m[path] = dataSize
	}

	return m, nil
}

func outResult(result *ComparedResult, fileA string, fileB string, out io.Writer) error {
	var perfectMatchCount int
	for _, comparedInfo := range result.Both {
		if comparedInfo.SizeInA == comparedInfo.SizeInB {
			perfectMatchCount++
		}
	}

	templateText := `file A: {{.FileA}}
file B: {{.FileB}}

それぞれのパスの数
	A: {{.TotalA}}
	B: {{.TotalB}}

一致しているパスの数: {{.BothTotalCount}}
データサイズも一致しているパスの数: {{.PerfectMatchCount}}
`

	v := struct {
		FileA             string
		FileB             string
		TotalA            int
		TotalB            int
		BothTotalCount    int
		PerfectMatchCount int
	}{
		FileA:             fileA,
		FileB:             fileB,
		TotalA:            len(result.Both) + len(result.OnlyA),
		TotalB:            len(result.Both) + len(result.OnlyB),
		BothTotalCount:    len(result.Both),
		PerfectMatchCount: perfectMatchCount,
	}

	tpl, err := template.New("").Parse(templateText)
	if err != nil {
		return err
	}

	err = tpl.Execute(out, v)
	if err != nil {
		return err
	}

	return nil
}
