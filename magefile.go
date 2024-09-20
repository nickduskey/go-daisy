//go:build mage
// +build mage

package main

import (
	"fmt"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
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
	return sh.RunV("templ", "generate", "--watch")
}

type Tailwind mg.Namespace

func (Tailwind) Build() error {
	return sh.RunV("tailwindcss", "-i", "./static/css/input.css", "-o", "./static/css/output.css")
}

func (Tailwind) Watch() error {
	return sh.RunV("tailwindcss", "-i", "./static/css/input.css", "-o", "./static/css/output.css", "--watch")
}