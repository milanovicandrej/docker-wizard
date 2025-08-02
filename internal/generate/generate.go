package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func DetectLanguage() string {
	if _, err := os.Stat("requirements.txt"); err == nil {
		return "python"
	}
	if _, err := os.Stat("package.json"); err == nil {
		return "nodejs"
	}
	if _, err := os.Stat("go.mod"); err == nil {
		return "golang"
	}
	return ""
}

func getPythonVersion() string {
	data, err := os.ReadFile("requirements.txt")
	if err != nil {
		return "3.9"
	}
	re := regexp.MustCompile(`python(?:[=><!]+)?([\d\.]+)?`)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		if len(m) > 1 {
			return m[1]
		}
	}
	return "3.9"
}

func getNodeVersion() string {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return "20"
	}
	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "20"
	}
	if engines, ok := pkg["engines"].(map[string]interface{}); ok {
		if node, ok := engines["node"].(string); ok {
			re := regexp.MustCompile(`(\d+)`)
			m := re.FindStringSubmatch(node)
			if len(m) > 1 {
				return m[1]
			}
		}
	}
	return "20"
}

func getGoVersion() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "1.21"
	}
	re := regexp.MustCompile(`go\s+(\d+\.\d+)`)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		if len(m) > 1 {
			return m[1]
		}
	}
	return "1.21"
}

func getBaseImage(language string) string {
	switch language {
	case "python":
		return fmt.Sprintf("python:%s-slim", getPythonVersion())
	case "nodejs":
		return fmt.Sprintf("node:%s-slim", getNodeVersion())
	case "golang":
		return fmt.Sprintf("golang:%s-alpine", getGoVersion())
	default:
		return "ubuntu:latest"
	}
}

func CreateDockerfileContent(language string) string {
var lines []string
isReact := false
isAngular := false
isVue := false
if language == "nodejs" {
	// Detect React, Angular, Vue by looking for dependencies in package.json
	data, err := os.ReadFile("package.json")
	if err == nil {
		var pkg map[string]interface{}
		if err := json.Unmarshal(data, &pkg); err == nil {
			deps := map[string]interface{}{}
			if d, ok := pkg["dependencies"].(map[string]interface{}); ok {
				deps = d
			}
			if d, ok := pkg["devDependencies"].(map[string]interface{}); ok {
				for k, v := range d {
					deps[k] = v
				}
			}
			for k := range deps {
				if strings.Contains(strings.ToLower(k), "react") {
					isReact = true
				}
				if strings.Contains(strings.ToLower(k), "@angular") {
					isAngular = true
				}
				if strings.Contains(strings.ToLower(k), "vue") {
					isVue = true
				}
			}
		}
	}
}

switch language {
case "python":
	lines = append(lines,
		"# Multistage build for Python",
		fmt.Sprintf("FROM %s AS builder", getBaseImage(language)),
		"WORKDIR /app",
		"COPY . .",
		"RUN pip install --upgrade pip",
		"RUN pip install -r requirements.txt",
		"\n# Final image",
		"FROM python:3.9-slim",
		"WORKDIR /app",
		"COPY --from=builder /app /app",
		"CMD [\"python\", \"main.py\"]",
	)
case "nodejs":
	buildStep := ""
	serveStep := "CMD [\"node\", \"main.js\"]"
	if isReact {
		buildStep = "RUN npm run build"
		serveStep = "# Serve React build\nCMD [\"npx\", \"serve\", \"-s\", \"build\"]"
	} else if isAngular {
		buildStep = "RUN npm run build"
		serveStep = "# Serve Angular build\nCMD [\"npx\", \"http-server\", \"dist\"]"
	} else if isVue {
		buildStep = "RUN npm run build"
		serveStep = "# Serve Vue build\nCMD [\"npx\", \"http-server\", \"dist\"]"
	}
	lines = append(lines,
		"# Multistage build for Node.js",
		fmt.Sprintf("FROM %s AS builder", getBaseImage(language)),
		"WORKDIR /app",
		"COPY . .",
		"RUN npm install",
		buildStep,
		"\n# Final image",
		"FROM node:20-slim",
		"WORKDIR /app",
		"COPY --from=builder /app /app",
		serveStep,
	)
case "golang":
	lines = append(lines,
		"# Multistage build for Go",
		fmt.Sprintf("FROM %s AS builder", getBaseImage(language)),
		"WORKDIR /src",
		"COPY . .",
		"RUN go build -o app .",
		"\n# Final image",
		"FROM alpine:latest",
		"WORKDIR /app",
		"COPY --from=builder /src/app /app/app",
		"CMD [\"/app/app\"]",
	)
default:
	lines = append(lines,
		"# Add custom instructions for other languages",
		fmt.Sprintf("FROM %s", getBaseImage(language)),
		"WORKDIR /app",
		"COPY . .",
	)
}
return strings.Join(lines, "\n\n")
}
