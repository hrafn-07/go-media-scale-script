package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/joho/godotenv"
)

const DefaultENV = ".env"

func main() {
	log.Println("########################################")
	log.Println("Starting image processing application...")

	// Command-line flags
	envFlag := flag.String("env", DefaultENV, "Path to the .env file")
	watermarkFlag := flag.Bool("w", false, "Add watermark")
	allSizesFlag := flag.Bool("a", false, "Process all sizes")
	smallFlag := flag.Bool("s", false, "Process small size")
	mediumFlag := flag.Bool("m", false, "Process medium size")
	largeFlag := flag.Bool("l", false, "Process large size")
	xlargeFlag := flag.Bool("xl", false, "Process extra-large size")
	flag.Parse()

	// Validate input arguments
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("[ERROR] No input file provided. Usage: %s [options] <file>", os.Args[0])
	}
	file := args[0]
	log.Printf("[INFO] Processing file: %s", file)

	// Define size processing flags
	sizes := map[string]bool{
		"s":  *smallFlag || *allSizesFlag,
		"m":  *mediumFlag || *allSizesFlag,
		"l":  *largeFlag || *allSizesFlag,
		"xl": *xlargeFlag || *allSizesFlag,
	}

	// Load environment variables
	envPath := *envFlag
	log.Printf("[INFO] Loading environment variables from %s", envPath)
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("[ERROR] Failed to load .env file: %v", err)
	}

	// Read required environment variables
	outputBaseDir := getEnvOrFail("OUTPUT_BASE_DIR")
	ownerUser := getEnvOrFail("OWNER_USER")
	watermarkFile := os.Getenv("WATERMARK_FILE")

	dimensions := map[string]string{
		"s":  os.Getenv("DIMENSION_S"),
		"m":  os.Getenv("DIMENSION_M"),
		"l":  os.Getenv("DIMENSION_L"),
		"xl": os.Getenv("DIMENSION_XL"),
	}

	// Validate input file type
	if !isImage(file) {
		log.Fatalf("[ERROR] File %s is not a valid image", file)
	}

	// Process each enabled size
	for size, enabled := range sizes {
		if !enabled {
			continue
		}

		dimension := dimensions[size]
		if dimension == "" {
			log.Printf("[WARNING] No dimension found for size %s. Skipping.", size)
			continue
		}

		outputDir := filepath.Join(outputBaseDir, size)
		outputFile := filepath.Join(outputDir, filepath.Base(file))
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("[ERROR] Failed to create directory %s: %v", outputDir, err)
		}

		startTime := time.Now()
		log.Printf("[INFO] Processing %s as %s (%s pixels)", file, size, dimension)

		err := processImage(file, watermarkFile, outputFile, dimension, size, *watermarkFlag)
		if err != nil {
			log.Printf("[ERROR] Failed to process %s as %s: %v", file, size, err)
			continue
		}

		duration := time.Since(startTime)
		log.Printf("[INFO] Successfully processed %s as %s in %v", file, size, duration)

		if err := changeOwnership(outputFile, ownerUser); err != nil {
			log.Printf("[ERROR] Failed to change ownership for %s: %v", outputFile, err)
		}
	}
}

func processImage(inputFile, watermarkFile, outputFile, dimension, size string, addWatermark bool) error {
	srcImage, err := imaging.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input image: %w", err)
	}

	dim, err := strconv.Atoi(dimension)
	if err != nil {
		return fmt.Errorf("invalid dimension: %w", err)
	}

	dstImage := imaging.Resize(srcImage, dim, 0, imaging.Lanczos)

	if addWatermark && (size == "xl" || size == "l" || size == "m") {
		watermark, err := imaging.Open(watermarkFile)
		if err != nil {
			return fmt.Errorf("failed to open watermark image: %w", err)
		}

		scaleFactor := getWatermarkScaleFactor(size)
		resizedWatermark := imaging.Resize(watermark, watermark.Bounds().Dx()*scaleFactor/100, 0, imaging.Lanczos)
		dstImage = imaging.OverlayCenter(dstImage, resizedWatermark, 1.0)
	}

	if err := imaging.Save(dstImage, outputFile); err != nil {
		return fmt.Errorf("failed to save output image: %w", err)
	}

	log.Printf("[INFO] Image saved: %s", outputFile)
	return nil
}

func isImage(file string) bool {
	cmd := exec.Command("file", "--mime-type", "-b", file)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("[ERROR] Failed to determine file type: %v", err)
	}
	return strings.Contains(string(output), "image")
}

func getWatermarkScaleFactor(size string) int {
	switch size {
	case "l":
		return 66
	case "m":
		return 33
	default:
		return 100
	}
}

func changeOwnership(file, ownerUser string) error {
	cmd := exec.Command("chown", fmt.Sprintf("%s:%s", ownerUser, ownerUser), file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to change ownership: %w, output: %s", err, string(output))
	}

	log.Printf("[INFO] Ownership changed for %s to %s", file, ownerUser)
	return nil
}

func getEnvOrFail(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("[ERROR] Environment variable %s is not set. Exiting.", key)
	}
	log.Printf("[INFO] Loaded environment variable: %s=%s", key, value)
	return value
}
