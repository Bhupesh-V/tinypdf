package service

import (
	"os"
	"os/exec"
)

type Service interface {
	GenerateQpdfCommand(inputFile, outputFile string) (*exec.Cmd, error)
}

type service struct {
}

func New() Service {
	return &service{}
}

func (s *service) GenerateQpdfCommand(inputFile, outputFile string) (*exec.Cmd, error) {
	cmd := exec.Command("qpdf", "--object-streams=generate", "--stream-data=compress", inputFile, outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}
