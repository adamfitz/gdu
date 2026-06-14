package parser

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
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
func DirResult(rootPath string) error {
	var dirSizes []DirSize

	fmt.Println("Calculating disk usage for: ", os.Args[1])

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v", path, err)
			return nil // Continue walking even if an error occurs for one entry
		}
		if d.IsDir() {
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

	// Sort directories by size in ascending order
	sort.Slice(dirSizes, func(i, j int) bool {
		return dirSizes[i].Size < dirSizes[j].Size
	})

	// Print in tree-like format
	fmt.Printf("Disk usage for %s (sorted by size):\n", rootPath)

	for idx, ds := range dirSizes {
		relPath, err := filepath.Rel(rootPath, ds.Path)
		if err != nil {
			relPath = ds.Path // Fallback if relative path calculation fails
		}
		// if a hidden file add a space to align
		firstRune, _ := utf8.DecodeRuneInString(relPath)
		if string(firstRune) != "." {
			relPath = " " + relPath
		}

		if relPath == "." { // Root directory itself
			relPath = rootPath
		}
		// if its the last element then add a carriage return before printing
		if idx == len(dirSizes)-1 {
			fmt.Printf("\n%-12s %s\n", formatBytes(ds.Size), relPath)
		} else {
			fmt.Printf("%-12s %s\n", formatBytes(ds.Size), relPath)
		}
	}

	return nil
}
