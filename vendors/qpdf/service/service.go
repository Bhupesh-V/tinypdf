package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"tinypdf/shared"
)

type Service interface {
	GenerateQpdfCommand(inputFile, outputFile string) *exec.Cmd
}

type service struct {
}

func New() Service {
	if !shared.IsBinaryAvailable("qpdf") {
		fmt.Println("Error: QPDF is not installed or not found in PATH.")
		fmt.Println("Please install QPDF to use this tool.")
		os.Exit(1)
	}
	return &service{}
}

func (s *service) GenerateQpdfCommand(inputFile, outputFile string) *exec.Cmd {
	// Source: https://qpdf.readthedocs.io/en/stable/cli.html#optimizing-file-size
	cmd := exec.Command(
		"qpdf",
		"--object-streams=generate",
		"--stream-data=compress",
		// "--optimize-images",
		// "--jpeg-quality=85", // default is 75

		// TODO: Enable zopfli support in qpdf
		// "--compression-level=9",
		// "--decode-level=generalized",
		inputFile,
		outputFile,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println(cmd.String())
	return cmd
}
