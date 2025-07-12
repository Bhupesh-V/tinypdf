package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"tinypdf/entities"
	tools "tinypdf/internal"
	"tinypdf/shared"

	gsEntities "tinypdf/vendors/gs/entities"

	"github.com/dustin/go-humanize"
)

// ASCII is "miniwi"
func usage() {
	fmt.Println(`

▗ ▘       ▌▐▘
▜▘▌▛▌▌▌▛▌▛▌▜▘
▐▖▌▌▌▙▌▙▌▙▌▐
     ▄▌▌


Usage:
	tinypdf -i input.pdf [-preset screen|ebook|printer|prepress] [-quality 50]

Example:
	docker run --rm -v $(pwd):/app bhupeshimself/tinypdf -i input.pdf -preset printer -quality 40

Options:
	-preset         One of: screen, ebook, printer, prepress (default: screen)
	-quality        Quality percentage between 10 and 90 (default: 50)

Bugs:
	Please report any issues on the GitHub repository.
		https://github.com/Bhupesh-V/tinypdf/issues`)
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
	inputPath := flag.String("i", "", "Path to input PDF")
	preset := flag.String("preset", "screen", "Compression preset: screen, ebook, printer, prepress")
	quality := flag.Float64("quality", 50, "Quality percent (10-90)")

	flag.Usage = usage
	flag.Parse()

	if len(os.Args) == 1 {
		// just show usage if no args are provided
		flag.Usage()
		os.Exit(0)
	}

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

	pipeline := shared.Pipe(
		tools.Poppler,
		tools.Ghostscript,
		tools.QPDF,
		// collect the final output file
		func(inputFilePath string, c shared.Config) string {
			finalOutputFile := fmt.Sprintf("tinypdf-%s", c.OriginalFileName)

			// If os.Rename fails (e.g., cross-device), fallback to copy using standard library
			src, openErr := os.Open(inputFilePath)
			if openErr != nil {
				fmt.Println("Error opening temp file for copying:", openErr)
				os.Exit(1)
			}
			defer src.Close()

			originalFileDir := filepath.Dir(c.OriginalFilePath)
			finalOutputFile = filepath.Join(originalFileDir, finalOutputFile)

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

			return finalOutputFile
		},
	)

	var isDebug bool
	if os.Getenv("TINYPDF_DEBUG") != "" {
		isDebug = true
	}

	baseFileName := filepath.Base(*inputPath)

	finalOutputFile := pipeline(*inputPath, shared.Config{
		IsDebug:          isDebug,
		OriginalFilePath: *inputPath,
		OriginalFileName: baseFileName,
		GSConfig: &gsEntities.Config{
			Preset:               *preset,
			ColorImageResolution: colorRes,
			MonoImageResolution:  monoRes,
			GrayImageResolution:  grayRes,
		},
	})

	outputFileSizeBytes := shared.FileSizeBytes(finalOutputFile)
	printFileSizeReport(inputFileSizeBytes, outputFileSizeBytes)
}
