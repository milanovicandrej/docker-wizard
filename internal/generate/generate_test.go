package generate

import (
	"os"
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	// Setup: create dummy files
	os.WriteFile("requirements.txt", []byte(""), 0644)
	if lang := DetectLanguage(); lang != "python" {
		t.Errorf("Expected python, got %s", lang)
	}
	os.Remove("requirements.txt")

	os.WriteFile("package.json", []byte("{}"), 0644)
	if lang := DetectLanguage(); lang != "nodejs" {
		t.Errorf("Expected nodejs, got %s", lang)
	}
	os.Remove("package.json")

	os.WriteFile("go.mod", []byte("go 1.21"), 0644)
	if lang := DetectLanguage(); lang != "golang" {
		t.Errorf("Expected golang, got %s", lang)
	}
	os.Remove("go.mod")
}

func TestCreateDockerfileContent_Python(t *testing.T) {
	os.WriteFile("requirements.txt", []byte("python==3.9"), 0644)
	content := CreateDockerfileContent("python")
	if !contains(content, "FROM python:3.9-slim AS builder") {
		t.Error("Python Dockerfile missing correct base image")
	}
	os.Remove("requirements.txt")
}

func TestCreateDockerfileContent_Nodejs(t *testing.T) {
	os.WriteFile("package.json", []byte(`{"dependencies":{"react":"^18.0.0"}}`), 0644)
	content := CreateDockerfileContent("nodejs")
	if !contains(content, "RUN npm run build") {
		t.Error("Nodejs Dockerfile missing build step for React")
	}
	os.Remove("package.json")
}

func TestCreateDockerfileContent_Golang(t *testing.T) {
	os.WriteFile("go.mod", []byte("go 1.21"), 0644)
	content := CreateDockerfileContent("golang")
	if !contains(content, "FROM golang:1.21-alpine AS builder") {
		t.Error("Golang Dockerfile missing correct base image")
	}
	os.Remove("go.mod")
}

func contains(s, substr string) bool {
	return len(s) > 0 && (s == substr || len(s) > len(substr) && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}
