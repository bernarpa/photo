package main

import (
	"fmt"
	"os"

	"github.com/bernarpa/photo/operations"
)

func showHelp() {
	fmt.Println()
	fmt.Println("Usage: photo <OPERATION>")
	fmt.Println()
	fmt.Println("   OPERATION     available options: help, fix, filter, info, ignore, stats, update")
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}
	switch op := os.Args[1]; op {
	case "help":
		showHelp()
	case "localupdate":
		operations.RunCommandFunction(operations.LocalUpdate, operations.ShowHelpUpdate, true)
	case "update":
		operations.RunCommandFunction(operations.Update, operations.ShowHelpUpdate, true)
	case "stats":
		operations.RunCommandFunction(operations.Stats, operations.ShowHelpStats, true)
	case "filter":
		operations.RunCommandFunction(operations.Filter, operations.ShowHelpFilter, true)
	case "fix":
		operations.RunCommandFunction(operations.Fix, operations.ShowHelpFix, false)
	case "info":
		operations.RunCommandFunction(operations.Info, operations.ShowHelpInfo, false)
	case "ignore":
		operations.RunCommandFunction(operations.Ignore, operations.ShowHelpIgnore, false)
	default:
		fmt.Printf("Invalid operation: %s\n", op)
		showHelp()
		os.Exit(1)
	}
}
