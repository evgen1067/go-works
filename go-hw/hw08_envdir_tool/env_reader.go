package main

import (
	"bytes"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	envs := make(Environment)
	for _, f := range files {
		fPath := path.Join(dir, f.Name())
		in, err := os.ReadFile(fPath)
		if err != nil {
			return nil, err
		}
		// если файл полностью пустой (длина - 0 байт), то envdir удаляет переменную окружения с именем S
		// len(strings.TrimSpace(string(in))) == 0
		if len(in) == 0 {
			envs[f.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		in = bytes.Split(in, []byte("\n"))[0]                   // 1 row
		in = bytes.ReplaceAll(in, []byte("\x00"), []byte("\n")) // терминальные нули (0x00) заменяются на перевод строки (\n);

		envVal := strings.TrimRight(string(in), "\t\n") // пробелы и табуляция в конце T удаляются;

		if len(strings.TrimSpace(envVal)) == 0 {
			envVal = ""
		}

		envs[f.Name()] = EnvValue{
			Value:      envVal,
			NeedRemove: false,
		}
	}
	return envs, err
}
