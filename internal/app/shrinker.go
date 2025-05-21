package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gkwa/poppingpenguin/internal/domain"
	"github.com/gkwa/poppingpenguin/internal/logging"
)

// ImageShrinker is responsible for shrinking images and reporting results
type ImageShrinker struct {
	compressionLevel int
	concurrencyLevel int
	logger           logging.Logger
	processor        domain.ImageProcessor
	reporter         domain.ShrinkReporter
}

// NewImageShrinker creates a new image shrinker with the given compression level
func NewImageShrinker(compressionLevel, concurrencyLevel int, logger logging.Logger) *ImageShrinker {
	processor := domain.NewImageMagickProcessor(compressionLevel)
	reporter := domain.NewConsoleReporter()

	return &ImageShrinker{
		compressionLevel: compressionLevel,
		concurrencyLevel: concurrencyLevel,
		logger:           logger,
		processor:        processor,
		reporter:         reporter,
	}
}

// ShrinkImages processes all the provided file patterns
func (s *ImageShrinker) ShrinkImages(patterns []string) error {
	files, err := s.expandFilePatterns(patterns)
	if err != nil {
		return err
	}

	s.logger.Debug("Found %d files to process", len(files))

	if len(files) == 0 {
		s.logger.Warning("No files found matching the provided patterns")
		return nil
	}

	// Use a wait group to process files concurrently
	var wg sync.WaitGroup

	// Channel for collecting results
	resultsChan := make(chan domain.ShrinkResult, len(files))
	errorsChan := make(chan error, len(files))

	// Semaphore channel to limit concurrency
	semaphore := make(chan struct{}, s.concurrencyLevel)

	s.logger.Debug("Processing with concurrency level: %d", s.concurrencyLevel)

	// Process files concurrently
	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			s.logger.Debug("Processing file: %s", filePath)

			result, err := s.processFile(filePath)
			if err != nil {
				s.logger.Error("Failed to process %s: %v", filePath, err)
				errorsChan <- err
				return
			}

			resultsChan <- result
		}(file)
	}

	// Wait for all files to be processed
	wg.Wait()
	close(resultsChan)
	close(errorsChan)

	// Collect results and errors
	var results []domain.ShrinkResult
	var totalOriginalSize int64
	var totalNewSize int64

	for result := range resultsChan {
		results = append(results, result)
		totalOriginalSize += result.OriginalSize
		totalNewSize += result.NewSize
	}

	// Check if there were any errors
	errorCount := len(errorsChan)
	if errorCount > 0 {
		s.logger.Warning("Failed to process %d files", errorCount)
	}

	// Report results if we have any
	if len(results) > 0 {
		s.reporter.ReportResults(results, totalOriginalSize, totalNewSize)
	}

	return nil
}

// expandFilePatterns expands all file patterns into a list of actual files
func (s *ImageShrinker) expandFilePatterns(patterns []string) ([]string, error) {
	var allFiles []string

	for _, pattern := range patterns {
		s.logger.Debug("Expanding pattern: %s", pattern)

		files, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %s: %w", pattern, err)
		}

		if len(files) == 0 {
			s.logger.Warning("No files found matching pattern: %s", pattern)
			continue
		}

		allFiles = append(allFiles, files...)
	}

	return allFiles, nil
}

// processFile shrinks a single image file
func (s *ImageShrinker) processFile(filePath string) (domain.ShrinkResult, error) {
	// Get original file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return domain.ShrinkResult{}, err
	}

	originalSize := fileInfo.Size()

	// Process the image
	err = s.processor.Process(filePath)
	if err != nil {
		return domain.ShrinkResult{}, err
	}

	// Get new file info
	fileInfo, err = os.Stat(filePath)
	if err != nil {
		return domain.ShrinkResult{}, err
	}

	newSize := fileInfo.Size()

	return domain.ShrinkResult{
		FilePath:     filePath,
		OriginalSize: originalSize,
		NewSize:      newSize,
	}, nil
}
