package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrCantReadDir         = errors.New("can't read directory")
	ErrCantGetFileStat     = errors.New("can't get file stat")
	ErrCantOpenFile        = errors.New("can't open file")
	ErrHasForbiddenSymbols = errors.New("filename has forbidden symbols")
)

const (
	forbiddenSymbol = "="
	zero            = 0x00
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
	// Place your code here
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrCantReadDir
	}

	env := make(Environment)
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		// имя S не должно содержать =
		if strings.Contains(fileInfo.Name(), forbiddenSymbol) {
			return nil, ErrHasForbiddenSymbols
		}

		stat, err := fileInfo.Info()
		if err != nil {
			return nil, ErrCantGetFileStat
		}

		// если файл полностью пустой (длина - 0 байт)
		if stat.Size() == 0 {
			env[fileInfo.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		file, err := os.Open(filepath.Join(dir, fileInfo.Name()))
		if err != nil {
			return nil, ErrCantOpenFile
		}
		defer file.Close()
		reader := bufio.NewReader(file)

		firstLine, _, err := reader.ReadLine()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("can't read first line: %s", err.Error())
		}

		// терминальные нули (0x00) заменяются на перевод строки (\n)
		envValue := bytes.ReplaceAll(firstLine, []byte{zero}, []byte{'\n'})
		// пробелы и табуляция в конце T удаляются
		envValue = bytes.TrimRight(envValue, "\t ")

		env[fileInfo.Name()] = EnvValue{Value: string(envValue)}
	}

	return env, nil
}
