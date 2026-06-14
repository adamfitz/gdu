package help

import (
	"fmt"
)

func GduHelp() {
	fmt.Println("Go Disk Usage Help")
	fmt.Println("")
	fmt.Println("gdu <target_dir> [flags]")
	fmt.Println("")
	fmt.Println("-a\t\tsort ascending")
	fmt.Println("-d\t\tsort descending (default)")
	fmt.Println("-h\t\thelp")
	fmt.Println("-v\t\tgdu version")
}
