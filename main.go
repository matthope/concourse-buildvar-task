package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const OutputFilename = "buildvars.yaml"

type Output struct {
	Time        TimeFormats       `yaml:"time"`
	Environment map[string]string `yaml:"environment"`
}

type TimeFormats struct {
	RFC3339 string `yaml:"rfc3339"`
	Date    string `yaml:"date"`
}

func main() {
	if err := run(".", envMap(os.Environ()), time.Now().UTC()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run(pwd string, enviro map[string]string, now time.Time) error {
	files, err := os.ReadDir(pwd)
	if err != nil {
		return fmt.Errorf("dir %q: %w", pwd, err)
	}

	out := io.MultiWriter(os.Stdout)

	var errs error

	for _, i := range files {
		if i.IsDir() {
			file, err := os.Create(filepath.Join(pwd, i.Name(), OutputFilename))
			if err != nil {
				errs = errors.Join(errs, err)
			}

			defer file.Close()

			out = io.MultiWriter(out, file)
		}
	}

	if errs != nil {
		return errs
	}

	if err = output(now, enviro, out); err != nil {
		return err
	}

	return nil
}

func envMap(e []string) map[string]string {
	enviro := make(map[string]string)

	const kvParts = 2

	for _, i := range e {
		a := strings.SplitN(i, "=", kvParts)

		key := strings.ToLower(a[0])
		key = strings.ReplaceAll(key, "-", "_")

		enviro[key] = a[1]
	}

	return enviro
}

func output(now time.Time, enviro map[string]string, out io.Writer) error {
	e := yaml.NewEncoder(out)

	err := e.Encode(
		Output{
			Time: TimeFormats{
				RFC3339: now.Format(time.RFC3339Nano),
				Date:    now.Format(time.DateOnly),
			},
			Environment: enviro,
		},
	)
	if err != nil {
		return fmt.Errorf("yaml: %w", err)
	}

	return nil
}
