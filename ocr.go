package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func ocrDocker(image string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(image))
	if err != nil {
		return "", fmt.Errorf("failed to determine absolute path to %s: %v", image, err)
	}

	base := filepath.Base(image)
	args := []string{
		"run",
		"--rm",
		"-it",
		"-v", fmt.Sprintf("%s:/work", dir),
		"jitesoft/tesseract-ocr",
		fmt.Sprintf("/work/%s", base),
		"stdout",
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run Tesseract: %v", err)
	}

	return string(output), nil
}

func ocr(image string) (string, error) {
	cmd := exec.Command("tesseract", image, "stdout", "-l", "deu")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run Tesseract: %v", err)
	}

	return string(output), nil
}
