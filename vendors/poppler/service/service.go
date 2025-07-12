package service

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"tinypdf/shared"
	"tinypdf/vendors/poppler/entities"

	"github.com/dustin/go-humanize"
)

type Service interface {
	GeneratePdftocairoCommand(inputFile, outputFile string) *exec.Cmd
	Close()
	GetTempFileName() string
	ListImages(filePath string) ([]entities.Image, error)
}

type service struct {
	pdftocairoFile *os.File // Temporary file for pdftocairo output
	isDebug        bool
}

func New(isDebug bool) Service {
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
		isDebug:        isDebug,
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
	if s.isDebug {
		cmd.Stderr = os.Stdout
	}

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

func (s *service) ListImages(filePath string) ([]entities.Image, error) {
	cmd := exec.Command("pdfimages", "-list", filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var images []entities.Image
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[2:] {
		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}
		size, _ := humanize.ParseBigBytes(parts[14])
		width, _ := strconv.Atoi(parts[3])
		height, _ := strconv.Atoi(parts[4])
		images = append(images, entities.Image{
			Size:   size.Int64(),
			Color:  parts[5],
			Width:  width,
			Height: height,
		})
	}
	return images, nil
}
