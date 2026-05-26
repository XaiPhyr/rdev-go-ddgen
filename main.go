package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var tmpFiles embed.FS

type GeneratorData struct {
	Package string
	Domain  string
}

var DomainErr = errors.New("Domain Already Exists")

func main() {
	domainFlag := flag.String("d", "", "Domain name")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of ddgen:\n")
		fmt.Fprintf(os.Stderr, "  ddgen init\n")
		fmt.Fprintf(os.Stderr, "  ddgen -d <domain_name>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if os.Args[1] == "init" {
		err := GenerateFolderStructure()
		if err != nil {
			fmt.Println(fmt.Errorf("Initializing error %v", err))
		}
		return
	}

	if *domainFlag == "" {
		fmt.Println("Error: domain name flag (-d) is required")
		flag.Usage()
		os.Exit(1)
	}

	err := GenerateDomain(*domainFlag)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v\n", err))
		if errors.Is(err, DomainErr) {
			os.Exit(1)
		}

		flag.Usage()
		os.Exit(1)
	}
}

func GenerateDomain(domainName string) error {
	domain := strings.ToLower(domainName)

	templates := map[string]string{
		"templates/test.tmpl":       fmt.Sprintf("%s_test.go", domain),
		"templates/handler.tmpl":    "handler.go",
		"templates/service.tmpl":    "service.go",
		"templates/repository.tmpl": "repository.go",
		"templates/types.tmpl":      "types.go",
	}

	info, err := os.Stat(filepath.Join("internal", domain))
	if err == nil {
		if info.IsDir() {
			return DomainErr
		}
	}

	if err := os.MkdirAll(filepath.Join("internal", domain), 0755); err != nil {
		return fmt.Errorf("Cannot proceed creating domain folder %v", err)
	}

	cap := strings.ToUpper(domain[0:1])
	capitalizedDomain := fmt.Sprintf("%s%s", cap, domain[1:])
	data := GeneratorData{Package: domain, Domain: capitalizedDomain}

	for tmpPath, outputName := range templates {
		GenerateAndParse(domain, "internal", outputName, tmpPath, &data)
	}

	return nil
}

func GenerateFolderStructure() error {
	fmt.Println("Initializing folder structure...")

	foldersToGenerate := map[string]string{
		"cmd":        "cmd",
		"scripts":    "scripts",
		"config":     "internal/config",
		"migration":  "internal/db/migrations",
		"middleware": "internal/middleware",
		"templates":  "internal/templates",
		"server":     "internal/server",
		"dto":        "internal/shared/dto",
		"helpers":    "internal/shared/helpers",
		"models":     "internal/shared/models",
	}

	for _, v := range foldersToGenerate {
		if _, err := os.Stat(v); err == nil {
			fmt.Printf(" -> Folder already exists: %s\n", v)
			continue
		}

		fmt.Printf(" -> Created: %s\n", v)

		if err := os.MkdirAll(v, 0755); err != nil {
			return fmt.Errorf("Cannot proceed creating domain folder %v", err)
		}
	}

	_, err := os.Stat("internal/config/config.go")
	if err != nil {
		err = GenerateAndParse("", "internal/config", "config.go", "templates/config.tmpl", nil)
		if err != nil {
			fmt.Println(fmt.Errorf("Config file not created %v", err))
		}
	}

	_, err = os.Stat("internal/db/migrations.go")
	if err != nil {
		err = GenerateAndParse("", "internal/db", "migrations.go", "templates/migrations.tmpl", nil)
		if err != nil {
			fmt.Println(fmt.Errorf("Routes file not created %v", err))
		}
	}

	_, err = os.Stat("internal/server/routes.go")
	if err != nil {
		err = GenerateAndParse("", "internal/server", "routes.go", "templates/routes.tmpl", nil)
		if err != nil {
			fmt.Println(fmt.Errorf("Routes file not created %v", err))
		}
	}

	return nil
}

func GenerateAndParse(domain, folder, outputName, tmpPath string, data *GeneratorData) error {
	targetFilePath := filepath.Join(domain, outputName)

	tmplBytes, err := tmpFiles.ReadFile(tmpPath)

	tmpl, err := template.New(outputName).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", outputName, err)
	}

	// issue: os.Create(targetFilepath)
	// what?: not creating targetFilePath error failed to create specific path
	// how?: because targetFilePath only generates the domain and the file,
	// since we are generating it for internal, we need to make sure to join it using filepath.Join
	file, err := os.Create(filepath.Join(folder, targetFilePath))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", targetFilePath, err)
	}

	err = tmpl.Execute(file, data)
	file.Close()

	if err != nil {
		return fmt.Errorf("failed to write template %s: %w", targetFilePath, err)
	}

	fmt.Printf(" -> Created: %s\n", targetFilePath)

	return nil
}
