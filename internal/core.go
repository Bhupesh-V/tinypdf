package internal

import (
	"os"
	"tinypdf/shared"
	gsEntities "tinypdf/vendors/gs/entities"
	gs "tinypdf/vendors/gs/service"
	poppler "tinypdf/vendors/poppler/service"
	qpdf "tinypdf/vendors/qpdf/service"
)

func Ghostscript(inputFilePath string, config shared.Config) string {
	gsService := gs.New()
	gsTmpFile := gsService.GetTempFileName()

	preset := config.GSConfig.(*gsEntities.Config).Preset
	colorRes := config.GSConfig.(*gsEntities.Config).ColorImageResolution
	monoRes := config.GSConfig.(*gsEntities.Config).MonoImageResolution
	grayRes := config.GSConfig.(*gsEntities.Config).GrayImageResolution

	err := gsService.GenerateGSCommand(inputFilePath, gsTmpFile, &gsEntities.Config{
		Preset:               preset,
		ColorImageResolution: colorRes,
		MonoImageResolution:  monoRes,
		GrayImageResolution:  grayRes,
	}).Run()
	if err != nil {
		os.Exit(1)
	}

	return gsTmpFile
}

func Poppler(inputFilePath string, _ shared.Config) string {
	popplerService := poppler.New()
	pdftocairoTmpFile := popplerService.GetTempFileName()

	err := popplerService.GeneratePdftocairoCommand(inputFilePath, pdftocairoTmpFile).Run()
	if err != nil {
		os.Exit(1)
	}

	return pdftocairoTmpFile
}

func QPDF(inputFilePath string, _ shared.Config) string {
	qpdfService := qpdf.New()
	qpdfTmpFile := qpdfService.GetTempFileName()

	err := qpdfService.GenerateQpdfCommand(inputFilePath, qpdfTmpFile).Run()
	if err != nil {
		os.Exit(1)
	}

	return qpdfTmpFile
}
