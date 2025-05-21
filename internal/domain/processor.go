package domain

import (
	"fmt"
	"os/exec"
)

// ImageProcessor defines the interface for image processing
type ImageProcessor interface {
	Process(filePath string) error
}

// ImageMagickProcessor implements ImageProcessor using ImageMagick
type ImageMagickProcessor struct {
	compressionLevel int
}

// NewImageMagickProcessor creates a new ImageMagick processor
func NewImageMagickProcessor(compressionLevel int) *ImageMagickProcessor {
	return &ImageMagickProcessor{
		compressionLevel: compressionLevel,
	}
}

// Process shrinks an image using ImageMagick
func (p *ImageMagickProcessor) Process(filePath string) error {
	// Create a temporary file for processing
	tempFile := filePath + ".tmp"

	// Use convert to resize the image with the desired quality
	quality := fmt.Sprintf("%d%%", p.compressionLevel)
	cmd := exec.Command("convert", filePath, "-quality", quality, tempFile)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}

	// Move the temp file to the original file
	cmd = exec.Command("mv", tempFile, filePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to replace original file: %w", err)
	}

	return nil
}
