package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "no input file or directory")
		os.Exit(1)
	}

	inputPath := args[1]

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if fileInfo.IsDir() {
		duplicateDirectory(inputPath)
	} else {
		duplicateFile(inputPath)
	}
}

func duplicateDirectory(inputPath string) {
	name := filepath.Base(inputPath)
	outputPath := ""

	for i := 1; ; i++ {
		outputPath = filepath.Join(inputPath, "..", fmt.Sprintf("%s_%d", name, i))
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			break
		}
	}

	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(inputPath, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(outputPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func duplicateFile(inputPath string) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer inputFile.Close()

	ext := filepath.Ext(inputPath)
	name := strings.TrimSuffix(filepath.Base(inputPath), ext)
	outputPath := ""

	for i := 1; ; i++ {
		outputPath = filepath.Join(filepath.Dir(inputPath), fmt.Sprintf("%s_%d%s", name, i, ext))
		if _, err = os.Stat(outputPath); os.IsNotExist(err) {
			break
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer outputFile.Close()

	if _, err := io.Copy(outputFile, inputFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, sourceInfo.Mode())
}
