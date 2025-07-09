package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"tinypdf/shared"
)

type Service interface {
	GenerateOcrmypdfCommand(inputFile, outputFile string) *exec.Cmd
	Close()
	GetTempFileName() string
}

type service struct {
	ocrmypdfFile *os.File // Temporary file for ocrmypdf output
}

func New() Service {
	if !shared.IsBinaryAvailable("ocrmypdf") {
		fmt.Println("Error: ocrmypdf is not installed or not found in PATH.")
		os.Exit(1)
	}
	ocrmypdfTmpFile, err := os.CreateTemp("", "tinypdf-ocrmypdf-*.pdf")
	if err != nil {
		fmt.Println("Error creating temp file for ocrmypdf output:", err)
		os.Exit(1)
	}
	return &service{
		ocrmypdfFile: ocrmypdfTmpFile,
	}
}

func (s *service) GenerateOcrmypdfCommand(inputFile, outputFile string) *exec.Cmd {
	// Source: https://ocrmypdf.readthedocs.io/en/latest/optimizer.html

	cmd := exec.Command(
		"ocrmypdf",
		"--quiet",
		"--skip-text",
		"--optimize 3",
		inputFile,
		outputFile,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println(cmd.String())
	return cmd
}

func (s *service) Close() {
	if s.ocrmypdfFile != nil {
		s.ocrmypdfFile.Close()
		os.Remove(s.ocrmypdfFile.Name())
	}
}

func (s *service) GetTempFileName() string {
	if s.ocrmypdfFile != nil {
		return s.ocrmypdfFile.Name()
	}
	return ""
}
