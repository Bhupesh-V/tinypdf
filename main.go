package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"tinypdf/shared"
	gsEntities "tinypdf/vendors/gs/entities"
	gs "tinypdf/vendors/gs/service"
	qpdf "tinypdf/vendors/qpdf/service"

	"github.com/dustin/go-humanize"
)

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

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func mapToAllRanges(percent float64, ranges ...[2]float64) []int {
	result := make([]int, len(ranges))
	for i, r := range ranges {
		diff := r[1] - r[0]
		val := r[0] + (percent/100.0)*diff
		result[i] = int(val)
	}
	return result
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

	*quality = clamp(*quality, 10, 90)

	resolutions := mapToAllRanges(*quality, [2]float64{0, 300}, [2]float64{0, 1200}, [2]float64{0, 300})
	colorRes, monoRes, grayRes := resolutions[0], resolutions[1], resolutions[2]

	// Create temp files for intermediary steps
	gsTmpFile, err := os.CreateTemp("", "tinypdf-gs-*.pdf")
	if err != nil {
		fmt.Println("Error creating temp file for Ghostscript output:", err)
		os.Exit(1)
	}
	defer os.Remove(gsTmpFile.Name())
	gsTmpFile.Close()

	qpdfTmpFile, err := os.CreateTemp("", "tinypdf-qpdf-*.pdf")
	if err != nil {
		fmt.Println("Error creating temp file for QPDF output:", err)
		os.Exit(1)
	}
	defer os.Remove(qpdfTmpFile.Name())
	qpdfTmpFile.Close()

	err = gsService.GenerateGSCommand(*inputPath, gsTmpFile.Name(), &gsEntities.Config{
		Preset:               *preset,
		ColorImageResolution: colorRes,
		MonoImageResolution:  monoRes,
		GrayImageResolution:  grayRes,
	}).Run()
	if err != nil {
		os.Exit(1)
	}

	err = qpdfService.GenerateQpdfCommand(gsTmpFile.Name(), qpdfTmpFile.Name()).Run()
	if err != nil {
		os.Exit(1)
	}

	// Move final file to user's directory with correct name
	finalOutputFile := fmt.Sprintf("tinypdf-%s", *inputPath)
	err = os.Rename(qpdfTmpFile.Name(), finalOutputFile)
	if err != nil {
		// If os.Rename fails (e.g., cross-device), fallback to copy using standard library
		src, openErr := os.Open(qpdfTmpFile.Name())
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
