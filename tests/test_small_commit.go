// Test script for small commit - run with: go run test_small_commit.go
//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("🧪 Testing Diny with small diff...")
	fmt.Println("================================")

	// Build diny first
	fmt.Println("🔨 Building diny...")
	buildCmd := exec.Command("go", "build", "-o", "diny")
	err := buildCmd.Run()
	if err != nil {
		fmt.Printf("❌ Failed to build diny: %v\n", err)
		os.Exit(1)
	}

	// First, make sure we're clean
	cleanCmd := exec.Command("git", "reset", "HEAD", "test-diff/small.md")
	cleanCmd.Run() // Ignore errors if file not staged

	// Stage the small test file
	fmt.Println("📁 Staging test-diff/small.md...")
	stageCmd := exec.Command("git", "add", "test-diff/small.md")
	err = stageCmd.Run()
	if err != nil {
		fmt.Printf("❌ Failed to stage file: %v\n", err)
		os.Exit(1)
	}

	// Show what's staged
	fmt.Println("📋 Staged changes:")
	diffCmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, _ := diffCmd.Output()
	fmt.Printf("%s\n", output)

	// Show diff size
	diffSizeCmd := exec.Command("git", "diff", "--cached")
	diffOutput, _ := diffSizeCmd.Output()
	fmt.Printf("📏 Diff size: %d characters\n\n", len(diffOutput))

	// Run diny commit (this will be interactive)
	fmt.Println("🦖 Running diny commit...")
	fmt.Println()

	commitCmd := exec.Command("./diny", "commit")
	commitCmd.Stdin = os.Stdin
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr

	err = commitCmd.Run()
	if err != nil {
		fmt.Printf("❌ Diny commit failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Test completed!")
}
