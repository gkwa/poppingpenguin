package domain

import (
	"fmt"
	"sort"
)

// ShrinkResult contains the results of a shrink operation
type ShrinkResult struct {
	FilePath     string
	OriginalSize int64
	NewSize      int64
}

// ShrinkPercentage calculates the reduction percentage
func (r ShrinkResult) ShrinkPercentage() float64 {
	if r.OriginalSize == 0 {
		return 0
	}
	return (1 - float64(r.NewSize)/float64(r.OriginalSize)) * 100
}

// ShrinkReporter defines the interface for reporting shrink results
type ShrinkReporter interface {
	ReportResults(results []ShrinkResult, totalOriginalSize, totalNewSize int64)
}

// ConsoleReporter implements ShrinkReporter for console output
type ConsoleReporter struct{}

// NewConsoleReporter creates a new console reporter
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

// ReportResults reports the shrink results to the console
func (r *ConsoleReporter) ReportResults(results []ShrinkResult, totalOriginalSize, totalNewSize int64) {
	// Sort results by filepath for consistent output
	sort.Slice(results, func(i, j int) bool {
		return results[i].FilePath < results[j].FilePath
	})

	for _, result := range results {
		fmt.Printf("%s: %.2f MB → %.2f MB (%.2f%% smaller)\n",
			result.FilePath,
			float64(result.OriginalSize)/1024/1024,
			float64(result.NewSize)/1024/1024,
			result.ShrinkPercentage())
	}

	totalShrinkPercentage := 0.0
	if totalOriginalSize > 0 {
		totalShrinkPercentage = (1 - float64(totalNewSize)/float64(totalOriginalSize)) * 100
	}

	fmt.Printf("\nSummary: %.2f MB → %.2f MB (%.2f%% smaller)\n",
		float64(totalOriginalSize)/1024/1024,
		float64(totalNewSize)/1024/1024,
		totalShrinkPercentage)
}
