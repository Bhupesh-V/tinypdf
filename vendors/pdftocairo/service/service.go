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
}

type service struct {
}

func New() Service {
	if !shared.IsBinaryAvailable("pdftocairo") {
		fmt.Println("Error: pdftocairo is not installed or not found in PATH.")
		os.Exit(1)
	}
	return &service{}
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
