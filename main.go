package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"tinypdf/entities"
	"tinypdf/shared"
	gsEntities "tinypdf/vendors/gs/entities"
	gs "tinypdf/vendors/gs/service"
	pdftocairo "tinypdf/vendors/pdftocairo/service"
	qpdf "tinypdf/vendors/qpdf/service"

	"github.com/dustin/go-humanize"
)

// ASCII is "miniwi"
func usage() {
	fmt.Println(`

▗ ▘       ▌▐▘
▜▘▌▛▌▌▌▛▌▛▌▜▘
▐▖▌▌▌▙▌▙▌▙▌▐
     ▄▌▌


Usage: tinypdf -i input.pdf [-preset screen|ebook|printer|prepress] [-quality 50]

Example: tinypdf -i input.pdf -preset ebook -quality 60

Options:
  -preset         One of: screen, ebook, printer, prepress (default: screen)
  -quality        Quality percentage between 10 and 90 (default: 50)`)
}

func printFileSizeReport(originalBytes, outputBytes int64) {
	fmt.Println()
	fmt.Printf("Original:   %s\n", humanize.Bytes(uint64(originalBytes)))
	fmt.Printf("Compressed: %s\n", humanize.Bytes(uint64(outputBytes)))

	if originalBytes > 0 {
		reduction := originalBytes - outputBytes
		reductionPercent := float64(reduction) / float64(originalBytes) * 100

		if reduction > 0 {
			fmt.Printf("Reduced by: %s (%.1f%%)\n",
				humanize.Bytes(uint64(reduction)),
				reductionPercent)
		} else if reduction < 0 {
			// Bruh
			fmt.Printf("Increased by: %s (%.1f%%)\n",
				humanize.Bytes(uint64(-reduction)),
				-reductionPercent)
		} else {
			fmt.Println("No size change")
		}
	}
}

func main() {
	gsService := gs.New()
	qpdfService := qpdf.New()
	pdftocairoService := pdftocairo.New()

	defer gsService.Close()
	defer qpdfService.Close()
	defer pdftocairoService.Close()

	inputPath := flag.String("i", "", "Path to input PDF")
	preset := flag.String("preset", "screen", "Compression preset: screen, ebook, printer, prepress")
	quality := flag.Float64("quality", 50, "Quality percent (10-90)")

	flag.Usage = usage
	flag.Parse()

	if *inputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if !shared.FileExists(*inputPath) {
		fmt.Printf("Error: Input file '%s' does not exist.\n", *inputPath)
		os.Exit(1)
	}
	inputFileSizeBytes := shared.FileSizeBytes(*inputPath)

	*quality = shared.Clamp(*quality, 10, 90)

	res, ok := entities.DefaultResolutionLimits[*preset]
	if !ok {
		os.Exit(1)
	}

	resolutions := shared.MapToRanges(
		*quality,
		[2]float64{res.ColorImageMin, res.ColorImageMax},
		[2]float64{res.MonoImageMin, res.MonoImageMax},
		[2]float64{res.GrayImageMin, res.GrayImageMax},
	)
	colorRes, monoRes, grayRes := resolutions[0], resolutions[1], resolutions[2]

	// Create temp files for intermediary steps
	gsTmpFile := gsService.GetTempFileName()
	qpdfTmpFile := qpdfService.GetTempFileName()
	pdftocairoTmpFile := pdftocairoService.GetTempFileName()

	err := pdftocairoService.GeneratePdftocairoCommand(*inputPath, pdftocairoTmpFile).Run()
	if err != nil {
		os.Exit(1)
	}

	err = gsService.GenerateGSCommand(pdftocairoTmpFile, gsTmpFile, &gsEntities.Config{
		Preset:               *preset,
		ColorImageResolution: colorRes,
		MonoImageResolution:  monoRes,
		GrayImageResolution:  grayRes,
	}).Run()
	if err != nil {
		os.Exit(1)
	}

	err = qpdfService.GenerateQpdfCommand(gsTmpFile, qpdfTmpFile).Run()
	if err != nil {
		os.Exit(1)
	}

	// Move final file to user's directory with correct name
	finalOutputFile := fmt.Sprintf("tinypdf-%s", *inputPath)
	err = os.Rename(qpdfTmpFile, finalOutputFile)
	if err != nil {
		// If os.Rename fails (e.g., cross-device), fallback to copy using standard library
		src, openErr := os.Open(qpdfTmpFile)
		if openErr != nil {
			fmt.Println("Error opening temp file for copying:", openErr)
			os.Exit(1)
		}
		defer src.Close()
		dst, createErr := os.Create(finalOutputFile)
		if createErr != nil {
			fmt.Println("Error creating final output file:", createErr)
			os.Exit(1)
		}
		defer dst.Close()
		_, copyErr := io.Copy(dst, src)
		if copyErr != nil {
			fmt.Println("Error copying final file to output location:", copyErr)
			os.Exit(1)
		}
	}

	outputFileSizeBytes := shared.FileSizeBytes(finalOutputFile)

	printFileSizeReport(inputFileSizeBytes, outputFileSizeBytes)
}
