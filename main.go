package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

type Project struct {
	Cmd  string `yaml:"cmd"`
	Path string `yaml:"path"`
}

// isUnixShell detects if we're running in a Unix-like shell on Windows (Git Bash, MSYS2, Cygwin)
func isUnixShell() bool {
	if runtime.GOOS != "windows" {
		return true
	}
	// MSYSTEM is set by Git Bash/MSYS2 (e.g., MINGW64, MINGW32, MSYS)
	if os.Getenv("MSYSTEM") != "" {
		return true
	}
	// CYGWIN/TERM check for Cygwin environments
	if strings.HasPrefix(os.Getenv("TERM"), "cygwin") {
		return true
	}
	return false
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		path = filepath.Join(home, path[2:])
	}
	// Convert forward slashes to OS-specific separator
	return filepath.FromSlash(path)
}

func loadProjects(configPath string) (map[string]Project, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var projects map[string]Project
	if err := yaml.Unmarshal(data, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func findConfig() string {
	// Look for apps.yaml in the same directory as the executable
	exe, err := os.Executable()
	if err == nil {
		configPath := filepath.Join(filepath.Dir(exe), "apps.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Fallback to current directory
	return "apps.yaml"
}

func runProject(project Project) error {
	path := expandPath(project.Path)

	// Verify path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Change to project directory
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	fmt.Printf("Changed to: %s\n", path)
	fmt.Printf("Running: %s\n\n", project.Cmd)

	// Prepare command based on shell type
	var cmd *exec.Cmd
	if isUnixShell() {
		cmd = exec.Command("sh", "-c", project.Cmd)
	} else {
		cmd = exec.Command("cmd", "/c", project.Cmd)
	}

	// Connect stdin/stdout/stderr for interactive use
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	projects, err := loadProjects(findConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: pj <project>")
		fmt.Println("\nAvailable projects:")
		for name, p := range projects {
			fmt.Printf("  %-12s -> %s\n", name, p.Path)
		}
		os.Exit(0)
	}

	projectName := os.Args[1]
	project, ok := projects[projectName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown project: %s\n", projectName)
		fmt.Println("\nAvailable projects:")
		for name := range projects {
			fmt.Printf("  %s\n", name)
		}
		os.Exit(1)
	}

	if err := runProject(project); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
