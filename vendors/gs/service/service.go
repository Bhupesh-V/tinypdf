package service

import (
	"os"
	"os/exec"
	"strconv"
	"tinypdf/vendors/gs/entities"
)

type Service interface {
	GenerateGSCommand(inputFile, outputFile string, config *entities.Config) (*exec.Cmd, error)
}

type service struct{}

func New() Service {
	return &service{}
}

func (s *service) GenerateGSCommand(inputFile, outputFile string, config *entities.Config) (*exec.Cmd, error) {
	args := []string{
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dBATCH",
		"-dNOPAUSE",
		"-dQUIET",
		"-dDownsampleColorImages=true",
		"-dColorImageDownsampleType=/Bicubic",
		"-dDownsampleMonoImages=true",
		"-dMonoImageDownsampleType=/Subsample",
		"-dDownsampleGrayImages=true",
		"-dGrayImageDownsampleType=/Bicubic",
	}

	// Add config options
	if config != nil {
		if config.Preset != "" {
			args = append(args, "-dPDFSETTINGS=/"+config.Preset)
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

	return cmd, nil
}
