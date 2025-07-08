package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"tinypdf/shared"
)

type Service interface {
	GeneratePdftocairoCommand(inputFile, outputFile string) *exec.Cmd
	Close()
	GetTempFileName() string
}

type service struct {
	pdftocairoFile *os.File // Temporary file for pdftocairo output
}

func New() Service {
	if !shared.IsBinaryAvailable("pdftocairo") {
		fmt.Println("Error: pdftocairo is not installed or not found in PATH.")
		os.Exit(1)
	}
	pdftocairoTmpFile, err := os.CreateTemp("", "tinypdf-pdftocairo-*.pdf")
	if err != nil {
		fmt.Println("Error creating temp file for pdftocairo output:", err)
		os.Exit(1)
	}
	return &service{
		pdftocairoFile: pdftocairoTmpFile,
	}
}

func (s *service) GeneratePdftocairoCommand(inputFile, outputFile string) *exec.Cmd {
	// Source: https://manpages.ubuntu.com/manpages/focal/man1/pdftocairo.1.html
	cmd := exec.Command(
		"pdftocairo",
		"-pdf",
		inputFile,
		outputFile,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println(cmd.String())
	return cmd
}

func (s *service) Close() {
	if s.pdftocairoFile != nil {
		s.pdftocairoFile.Close()
		os.Remove(s.pdftocairoFile.Name())
	}
}

func (s *service) GetTempFileName() string {
	if s.pdftocairoFile != nil {
		return s.pdftocairoFile.Name()
	}
	return ""
}
