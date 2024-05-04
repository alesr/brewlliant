package rebrew

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	brewListFilename = "brew_list.txt"
	brewfileFilename = "Brewfile"
)

// CheckBrew checks if brew is installed.
func CheckBrew() error {
	cmd := exec.Command("brew", "--version")

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("brew is not installed: %w", err)
	}
	return nil
}

// BrewList generates a list of installed brew packages.
func BrewList() error {
	cmd := exec.Command("brew", "list", "-1")

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running brew list: %w", err)
	}

	file, err := os.Create(brewListFilename)
	if err != nil {
		return fmt.Errorf("error creating brew_list.txt: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(output); err != nil {
		return fmt.Errorf("error writing to %s: %w", brewListFilename, err)
	}
	return nil
}

// InstallFromBrewList installs brew packages from a list.
func InstallFromBrewList() error {
	fmt.Printf("What would you like to do?\n\t1. Install all packages.\n\t2. Get description for each package.\n\nEnter your choice:")

	var choice int
	if _, err := fmt.Scan(&choice); err != nil {
		return fmt.Errorf("error reading choice: %w", err)
	}

	brewList, err := openBrewListFile()
	if err != nil {
		return fmt.Errorf("error opening %s: %w", brewListFilename, err)
	}
	defer brewList.Close()

	switch choice {
	case 1:
		return installAllPackages(brewList)
	case 2:
		return getPackageDescriptions(brewList)
	default:
		return fmt.Errorf("invalid choice: %d", choice)
	}
}

func openBrewListFile() (*os.File, error) {
	file, err := os.Open(brewListFilename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", brewListFilename, err)
	}
	return file, nil
}

func installAllPackages(brewList *os.File) error {
	brewfile, err := createBrewfile()
	if err != nil {
		return fmt.Errorf("error creating %s: %w", brewfileFilename, err)
	}
	defer brewfile.Close()

	scanner := bufio.NewScanner(brewList)
	for scanner.Scan() {
		if err := writePackageToBrewfile(brewfile, scanner.Text()); err != nil {
			return fmt.Errorf("error writing to %s: %w", brewfileFilename, err)
		}
	}

	cmd := exec.Command("brew", "bundle", "install", "--file", brewfileFilename)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error installing packages: %w", err)
	}

	if outb.Len() > 0 {
		fmt.Print(outb.String())
	}

	if errb.Len() > 0 {
		fmt.Print(errb.String())
	}
	return nil
}

func createBrewfile() (*os.File, error) {
	brewfile, err := os.Create(brewfileFilename)
	if err != nil {
		return nil, fmt.Errorf("error creating %s: %w", brewfileFilename, err)
	}
	return brewfile, nil
}

func writePackageToBrewfile(brewfile *os.File, pkg string) error {
	_, err := fmt.Fprintf(brewfile, "brew \"%s\"\n", pkg)
	if err != nil {
		return fmt.Errorf("error writing to %s: %w", brewfileFilename, err)
	}
	return nil
}

func getPackageDescriptions(brewList *os.File) error {
	scanner := bufio.NewScanner(brewList)
	for scanner.Scan() {
		pkg := scanner.Text()
		desc, err := getPackageDescription(pkg)
		if err != nil {
			return fmt.Errorf("error getting description for %s: %w", pkg, err)
		}
		fmt.Printf("Description for %s:\n%s\n\n", pkg, desc)
	}
	return nil
}

func getPackageDescription(pkg string) (string, error) {
	cmd := exec.Command("brew", "info", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running 'brew info %s': %s", pkg, err)
	}

	lines := strings.Split(string(output), "\n")

	return fmt.Sprintf("%s\n%s\n%s", lines[0], lines[1], lines[2]), nil
}
