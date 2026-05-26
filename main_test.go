package main

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed templates/*
var testTemplates embed.FS

func TestGenerateDomain(t *testing.T) {
	t.Run("Successfully generates domain files inside sandbox", func(t *testing.T) {
		oldTmpFiles := tmpFiles
		tmpFiles = testTemplates
		defer func() { tmpFiles = oldTmpFiles }()

		tempWorkspace := t.TempDir()

		originalWD, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current working directory: %v", err)
		}

		err = os.Chdir(tempWorkspace)
		if err != nil {
			t.Fatalf("Failed to switch to sandbox directory: %v", err)
		}
		defer os.Chdir(originalWD)

		domainName := "Billing"
		expectedTargetFolder := filepath.Join("internal", "billing")

		err = GenerateDomain(domainName)
		if err != nil {
			t.Fatalf("GenerateDomain failed unexpectedly: %v", err)
		}

		expectedFile := filepath.Join(expectedTargetFolder, "service.go")
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Expected file to be generated at %s, but it was missing", expectedFile)
		}

		contentBytes, err := os.ReadFile(expectedFile)
		if err != nil {
			t.Fatalf("Failed to read generated output file: %v", err)
		}

		fileContent := string(contentBytes)
		if !strings.Contains(fileContent, "package billing") {
			t.Errorf("Template parsing failed! Expected file content to contain 'package billing', got:\n%s", fileContent)
		}
	})

	t.Run("Fails validation or execution if directory tree is locked", func(t *testing.T) {
		err := GenerateDomain("invalid/domain/name")
		if err == nil {
			t.Error("Expected an error when passing nested directory slashes, but got nil")
		}
	})
}
