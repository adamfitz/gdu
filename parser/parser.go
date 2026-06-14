package parser

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
	"unicode/utf8"
)

// DirSize represents a directory and its calculated size.
type DirSize struct {
	Path string
	Size int64 // in bytes
}

// formatBytes converts bytes to a human-readable format (KB, MB, GB).
func formatBytes(bytes int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/gb)
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/mb)
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/kb)
	default:
		return fmt.Sprintf("%d Bytes", bytes)
	}
}

func dirSize(path string) (int64, error) {
	var totalSize int64
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	return totalSize, err
}

func WalkDirSize(path string) (int64, error) {
	var totalSize int64

	err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			totalSize += info.Size()
		}
		return nil
	})

	return totalSize, err
}

// Return Dir Size result
func DirResult(rootPath string, flag string) error {
	var dirSizes []DirSize

	// timer
	spin := defaultSpinner(defaultChars, 100*time.Millisecond)
	spin.Start()

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v", path, err)
			return nil
		}
		if d.IsDir() && path != rootPath {
			size, err := dirSize(path)
			if err != nil {
				log.Printf("Error calculating size for %s: %v", path, err)
				return nil
			}
			dirSizes = append(dirSizes, DirSize{Path: path, Size: size})
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	rootSize, err := dirSize(rootPath)
	if err != nil {
		log.Fatalf("Error calculating root size: %v", err)
	}

	switch flag {
	case "-a":
		sortBySizeAscending(dirSizes)
	case "-d":
		sortBySizeDescending(dirSizes)
	default:
		sortBySizeAscending(dirSizes)
	}

	spin.Stop()

	fmt.Printf("Disk usage for %s (sorted by size):\n", rootPath)
	fmt.Println("")

	for _, ds := range dirSizes {
		relPath, err := filepath.Rel(rootPath, ds.Path)
		if err != nil {
			relPath = ds.Path
		}
		firstRune, _ := utf8.DecodeRuneInString(relPath)
		if string(firstRune) != "." {
			relPath = " " + relPath
		}
		fmt.Printf("%-12s %s\n", formatBytes(ds.Size), relPath)
	}

	fmt.Printf("\n%-12s %s\n", formatBytes(rootSize), rootPath)

	return nil
}

// sortBySize sorts a slice of DirSize in ascending order of Size.
func sortBySizeAscending(dirSizes []DirSize) {
	sort.Slice(dirSizes, func(i, j int) bool {
		return dirSizes[i].Size < dirSizes[j].Size
	})
}

// sortBySize sorts a slice of DirSize in descending order of Size.
func sortBySizeDescending(dirSizes []DirSize) {
	sort.Slice(dirSizes, func(i, j int) bool {
		return dirSizes[i].Size > dirSizes[j].Size
	})
}
