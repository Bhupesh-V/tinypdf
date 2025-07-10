package internal

import (
	"log"
	"os"
	"tinypdf/shared"
	gsEntities "tinypdf/vendors/gs/entities"
	gs "tinypdf/vendors/gs/service"
	poppler "tinypdf/vendors/poppler/service"
	qpdf "tinypdf/vendors/qpdf/service"
)

func Ghostscript(inputFilePath string, c shared.Config) string {
	gsService := gs.New()
	gsTmpFile := gsService.GetTempFileName()

	preset := c.GSConfig.(*gsEntities.Config).Preset
	colorRes := c.GSConfig.(*gsEntities.Config).ColorImageResolution
	monoRes := c.GSConfig.(*gsEntities.Config).MonoImageResolution
	grayRes := c.GSConfig.(*gsEntities.Config).GrayImageResolution

	cmd := gsService.GenerateGSCommand(inputFilePath, gsTmpFile, &gsEntities.Config{
		Preset:               preset,
		ColorImageResolution: colorRes,
		MonoImageResolution:  monoRes,
		GrayImageResolution:  grayRes,
	})

	if c.IsDebug {
		log.Println(cmd.String())
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Error running Ghostscript command: %v\n", err)
		os.Exit(1)
	}

	return gsTmpFile
}

func Poppler(inputFilePath string, c shared.Config) string {
	popplerService := poppler.New()
	pdftocairoTmpFile := popplerService.GetTempFileName()

	cmd := popplerService.GeneratePdftocairoCommand(inputFilePath, pdftocairoTmpFile)

	if c.IsDebug {
		log.Println(cmd.String())
	}

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	return pdftocairoTmpFile
}

func QPDF(inputFilePath string, c shared.Config) string {
	qpdfService := qpdf.New()
	qpdfTmpFile := qpdfService.GetTempFileName()

	cmd := qpdfService.GenerateQpdfCommand(inputFilePath, qpdfTmpFile)

	if c.IsDebug {
		log.Println(cmd.String())
	}

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	return qpdfTmpFile
}
