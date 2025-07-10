package main

import (
	"flag"
	"fmt"
	"os"
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

	pipeline := shared.Pipe(
		tools.Poppler,
		tools.Ghostscript,
		tools.QPDF,
		// collect the final output file
		func(inputFilePath string, c shared.Config) string {
			finalOutputFile := fmt.Sprintf("tinypdf-%s", c.OriginalFilePath)

			err := os.Rename(inputFilePath, finalOutputFile)
			if err != nil {
				os.Exit(1)
			}
			return finalOutputFile
		},
	)

	finalOutputFile := pipeline(*inputPath, shared.Config{
		OriginalFilePath: *inputPath,
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
