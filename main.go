package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"gdu/parser"
	"gdu/help"
)

// DirSize represents a directory and its calculated size.
type DirSize struct {
	Path string
	Size int64 // in bytes
}

/*
// calculateDirSize recursively calculates the size of a directory and its contents.
func calculateDirSize(path string) (int64, error) {
	var totalSize int64
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
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

// formatBytes converts bytes to a human-readable format.
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
*/

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\tgdu <directory_path>")
		fmt.Println("Help:\tgdu -h")
		return
	}
	if os.Args[1] != "-h" {
		rootPath := os.Args[1]

	var dirSizes []DirSize

	fmt.Println("Calculating disk usage for: ", os.Args[1])

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v", path, err)
			return nil // Continue walking even if an error occurs for one entry
		}
		if d.IsDir() {
			size, err := parser.DirSize(path)
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
		if relPath == "." { // Root directory itself
			relPath = rootPath
		}
		// if its the last element then add a carriage return before printing
		if idx == len(dirSizes) -1 {
			fmt.Printf("\n%-12s %s\n", parser.FormatBytes(ds.Size), relPath)
		} else {
			fmt.Printf("%-12s %s\n", parser.FormatBytes(ds.Size), relPath)
		}
	}
	} else {
		help.GduHelp()
	}

	
}
