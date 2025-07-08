package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"tinypdf/shared"
	"tinypdf/vendors/gs/entities"
)

type Service interface {
	GenerateGSCommand(inputFile, outputFile string, config *entities.Config) *exec.Cmd
	Close()
	GetTempFileName() string
}

type service struct {
	gsFile *os.File // Temporary file for Ghostscript output
}

func New() Service {
	if !shared.IsBinaryAvailable("gs") {
		fmt.Println("Error: Ghostscript (gs) is not installed or not found in PATH.")
		os.Exit(1)
	}
	gsTmpFile, err := os.CreateTemp("", "tinypdf-gs-*.pdf")
	if err != nil {
		fmt.Println("Error creating temp file for Ghostscript output:", err)
		os.Exit(1)
	}
	return &service{
		gsFile: gsTmpFile,
	}
}

func (s *service) GenerateGSCommand(inputFile, outputFile string, config *entities.Config) *exec.Cmd {
	args := []string{
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dBATCH",
		"-dNOPAUSE",
		"-dQUIET",
		"-dDownsampleColorImages=true",
		"-dColorImageDownsampleType=/Average",
		"-dDownsampleMonoImages=true",
		"-dMonoImageDownsampleType=/Subsample",
		"-dDownsampleGrayImages=true",
		"-dGrayImageDownsampleType=/Average",
		// "-dTextAlphaBits=50",
		// "-dGraphicsAlphaBits=50",
	}

	// Add config options
	if config != nil {
		if config.Preset != "" {
			// args = append(args, "-dPDFSETTINGS=/"+config.Preset)
		}
		if config.ColorImageResolution != 0 {
			args = append(args, "-dColorImageResolution="+strconv.Itoa(config.ColorImageResolution))
		}
		if config.MonoImageResolution != 0 {
			args = append(args, "-dMonoImageResolution="+strconv.Itoa(config.MonoImageResolution))
		}
		if config.GrayImageResolution != 0 {
			args = append(args, "-dGrayImageResolution="+strconv.Itoa(config.GrayImageResolution))
		}
	}

	// Don't quote filenames — exec.Command handles them properly
	args = append(args, "-sOutputFile="+outputFile, inputFile)

	cmd := exec.Command("gs", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println(cmd.String())

	return cmd
}

func (s *service) Close() {
	if s.gsFile != nil {
		err := s.gsFile.Close()
		if err != nil {
			log.Println("Error closing Ghostscript temp file:", err)
		}
		err = os.Remove(s.gsFile.Name())
		if err != nil {
			log.Println("Error removing Ghostscript temp file:", err)
		}
	}
}

func (s *service) GetTempFileName() string {
	if s.gsFile != nil {
		return s.gsFile.Name()
	}
	return ""
}
