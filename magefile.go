//go:build mage
// +build mage

package main

import (
	"fmt"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/nickduskey/go-daisy/internal/server"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Dev

func executableName() string {
	executableName := "go-daisy"
	if runtime.GOOS == "windows" {
		executableName += ".exe"
	}
	return executableName
}

func InstallDeps() error {
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}
	return sh.RunV("pnpm", "install")
}

// Build builds production binaries for the application
func Build() error {
	fmt.Println("Executing `Build` task...")

	executableName := executableName()
	fmt.Printf("Outputting executable: %s\n", executableName)

	fmt.Println("Building tailwind css...")
	var tailwind Tailwind
	if err := tailwind.Build(); err != nil {
		return err
	}

	fmt.Println("Generating templ templates...")
	var templ Templ
	if err := templ.Gen(); err != nil {
		return err
	}

	fmt.Println("Building server executable...")
	return sh.RunV("go", "build", "-o", executableName, ".")
}

// Dev serves the application. Meant for a local development environment
func Dev() error {
	var tailwind Tailwind
	go tailwind.Watch()
	var templ Templ
	go templ.Watch()

	return sh.RunV("air")
}

type Templ mg.Namespace

func (Templ) Gen() error {
	return sh.RunV("templ", "generate")
}

func (Templ) Watch() error {
	_ = godotenv.Load()
	address := server.BuildServerAddress()
	proxyFlag := fmt.Sprintf("--proxy=http://%s", address)
	return sh.RunV("templ", "generate", "--watch", proxyFlag, "--open-browser=false")
}

type Tailwind mg.Namespace

func (Tailwind) Build() error {
	return sh.RunV("tailwindcss", "-i", "./static/css/input.css", "-o", "./static/css/output.css")
}

func (Tailwind) Watch() error {
	return sh.RunV("tailwindcss", "-i", "./static/css/input.css", "-o", "./static/css/output.css", "--watch")
}
