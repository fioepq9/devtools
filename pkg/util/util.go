package util

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

func AllGoFiles(ctx context.Context, packages string) ([]string, error) {
	op := "/"
	if strings.HasPrefix(runtime.GOOS, "windows") {
		op = `\`
	}
	var buf bytes.Buffer
	golist := exec.CommandContext(ctx, "go",
		"list",
		"-f",
		`{{ $d := .Dir }}{{ range $f := .GoFiles }}{{ $d }}`+op+`{{ $f }}
{{ end }}`,
		packages,
	)
	golist.Stderr = &buf
	golist.Stdout = &buf
	err := golist.Run()
	if err != nil {
		return nil, errors.New(buf.String())
	}
	gofiles := strings.Split(strings.TrimSpace(buf.String()), "\n")
	return gofiles, nil
}
