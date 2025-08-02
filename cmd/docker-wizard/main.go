
package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/milanovicandrej/docker-wizard/internal/generate"
)

const Version = "0.4.0"

func main() {
	app := &cli.App{
		Name:  "docker-wizard",
		Usage: fmt.Sprintf("Generate a Dockerfile for Python, Node.js, or Go projects (v%s)", Version),
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "Dockerfile",
				Usage:   "Specify the output filename for the Dockerfile",
			},
		},
		Action: func(c *cli.Context) error {
			color.Magenta(fmt.Sprintf("🚀 Docker Wizard v%s: Generating Dockerfile...", Version))
			language := generate.DetectLanguage()
			if language == "" {
				color.Red("✗ Could not detect project type. Please ensure you have requirements.txt, package.json, or go.mod in your project.")
				return nil
			}
			content := generate.CreateDockerfileContent(language)
			output := c.String("output")
			if err := os.WriteFile(output, []byte(content), 0644); err != nil {
				color.Red("✗ Failed to write Dockerfile: %v", err)
				return err
			}
			color.Green("✔ Dockerfile generated!")
			color.Cyan("Project type detected: %s", language)
			color.Yellow("Output file: %s", output)
			color.Magenta("----------------------------------------")
			color.Blue(content)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error:", err)
	}
}
