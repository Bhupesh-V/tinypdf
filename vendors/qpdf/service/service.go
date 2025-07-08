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
	Close()
	GetTempFileName() string
}

type service struct {
	qpdfFile *os.File // Temporary file for QPDF output
}

func New() Service {
	if !shared.IsBinaryAvailable("qpdf") {
		fmt.Println("Error: QPDF is not installed or not found in PATH.")
		os.Exit(1)
	}
	qpdfTmpFile, err := os.CreateTemp("", "tinypdf-qpdf-*.pdf")
	if err != nil {
		fmt.Println("Error creating temp file for QPDF output:", err)
		os.Exit(1)
	}
	return &service{
		qpdfFile: qpdfTmpFile,
	}
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

func (s *service) Close() {
	if s.qpdfFile != nil {
		s.qpdfFile.Close()
		os.Remove(s.qpdfFile.Name())
	}
}

func (s *service) GetTempFileName() string {
	if s.qpdfFile != nil {
		return s.qpdfFile.Name()
	}
	return ""
}
