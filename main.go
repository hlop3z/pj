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

// useUnixShell returns true if we should use bash (VS Code terminal on Windows, or Unix systems)
// Standalone Git Bash (MinTTY) has PTY issues with Windows console programs, so we use cmd there
func useUnixShell() bool {
	if runtime.GOOS != "windows" {
		return true
	}
	// VS Code terminal with Git Bash works well with bash
	if os.Getenv("TERM_PROGRAM") == "vscode" && os.Getenv("MSYSTEM") != "" {
		return true
	}
	// Windows Terminal with Git Bash also works
	if os.Getenv("WT_SESSION") != "" && os.Getenv("MSYSTEM") != "" {
		return true
	}
	return false
}

// toUnixPath converts Windows path to Unix-style for Git Bash (C:\Users -> /c/Users)
func toUnixPath(winPath string) string {
	path := strings.ReplaceAll(winPath, "\\", "/")
	if len(path) >= 2 && path[1] == ':' {
		drive := strings.ToLower(string(path[0]))
		path = "/" + drive + path[2:]
	}
	return path
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

	fmt.Printf("Changed to: %s\n", path)
	fmt.Printf("Running: %s\n\n", project.Cmd)

	var cmd *exec.Cmd
	if useUnixShell() {
		// For Unix systems or VS Code/Windows Terminal with Git Bash
		unixPath := path
		if runtime.GOOS == "windows" {
			unixPath = toUnixPath(path)
		}
		fullCmd := fmt.Sprintf("cd '%s' && %s", unixPath, project.Cmd)
		cmd = exec.Command("bash", "-c", fullCmd)
	} else {
		// For Windows CMD/PowerShell/standalone Git Bash
		cmd = exec.Command("cmd", "/c", project.Cmd)
		cmd.Dir = path
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
