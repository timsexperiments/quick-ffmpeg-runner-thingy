//go:build mage
// +build mage

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	output := "build/m3u8-combiner"
	if err := sh.Run("go", "build", "-o", output, "./cmd/m3u8-combiner"); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Println("Build completed successfully!")
	return nil
}

func Run() error {
	if err := Build(); err != nil {
		return err
	}

	cmd := exec.Command("./build/m3u8-combiner", os.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run the program: %w", err)
	}

	return nil
}

func Install() error {
	mg.Deps(Build)

	buildOutput := "./build/m3u8-combiner"
	targetDir := "/usr/local/bin"

	cmd := exec.Command("sudo", "mv", buildOutput, filepath.Join(targetDir, "m3u8-combiner"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install binary globally: %w", err)
	}

	fmt.Println("Installation completed successfully!")
	return nil
}

func Clean() error {
	output := "build"
	if err := os.RemoveAll(filepath.Join(".", output)); err != nil {
		if os.IsNotExist(err) {
			log.Println("No binary found to clean")
			return nil
		}
		return fmt.Errorf("failed to clean up: %w", err)
	}
	return nil
}
