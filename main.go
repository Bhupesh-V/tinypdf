package main

import (
	"flag"
	"fmt"
	"os"
	"tinypdf/shared"
	gsEntities "tinypdf/vendors/gs/entities"
	gs "tinypdf/vendors/gs/service"
	qpdf "tinypdf/vendors/qpdf/service"

	"github.com/dustin/go-humanize"
)

func usage() {
	fmt.Println(`Usage: pdftrim -in input.pdf -out output.pdf [-preset screen|ebook|printer|prepress] [-quality 50]

Example: pdftrim -in input.pdf -out output.pdf -preset ebook -quality 75

Trim Options:
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
			fmt.Printf("Increased by: %s (%.1f%%)\n",
				humanize.Bytes(uint64(-reduction)),
				-reductionPercent)
		} else {
			fmt.Println("No size change")
		}
	}
}

func main() {
	inputPath := flag.String("i", "", "Path to input PDF")
	outputPath := flag.String("o", "", "Path to output PDF")
	preset := flag.String("preset", "screen", "Compression preset: screen, ebook, printer, prepress")
	quality := flag.Float64("quality", 50, "Quality percent (10-90)")

	flag.Usage = usage
	flag.Parse()

	if !shared.IsBinaryAvailable("gs") {
		fmt.Println("Error: Ghostscript (gs) is not installed or not found in PATH.")
		fmt.Println("Please install Ghostscript to use this tool.")
		os.Exit(1)
	}

	if !shared.IsBinaryAvailable("qpdf") {
		fmt.Println("Error: QPDF is not installed or not found in PATH.")
		fmt.Println("Please install QPDF to use this tool.")
		os.Exit(1)
	}

	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Error: -in and -out are required.")
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

	gsService := gs.New()
	gsCmd, _ := gsService.GenerateGSCommand(*inputPath, *outputPath, &gsEntities.Config{
		Preset:               *preset,
		ColorImageResolution: colorRes,
		MonoImageResolution:  monoRes,
		GrayImageResolution:  grayRes,
	})
	err := gsCmd.Run()
	if err != nil {
		fmt.Println("Error executing Ghostscript:", err)
		os.Exit(1)
	}

	qpdfService := qpdf.New()
	qpdfCmd, _ := qpdfService.GenerateQpdfCommand(*outputPath, fmt.Sprintf("tinypdf-%s", *outputPath))
	err = qpdfCmd.Run()
	if err != nil {
		fmt.Println("Error executing QPDF:", err)
		os.Exit(1)
	}

	outputFilePath := fmt.Sprintf("tinypdf-%s", *outputPath)
	outputFileSizeBytes := shared.FileSizeBytes(outputFilePath)

	printFileSizeReport(inputFileSizeBytes, outputFileSizeBytes)
}
