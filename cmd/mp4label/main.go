package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xd/mp4label/pkg/config"
	"github.com/xd/mp4label/pkg/server"
)

//go:embed web
var webFS embed.FS

// Version info (injected at build time, default v0.2.9)
var version = "v0.2.9"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "web":
		runWebServer()
	case "view":
		runViewServer()
	case "version", "--version", "-v":
		printVersion()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("mp4Label %s\n", version)
	fmt.Println("Video annotation tool for tutorial videos")
}

func printUsage() {
	fmt.Println("mp4Label - Video Annotation Tool")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  mp4label web [options]     Start web server (edit mode)")
	fmt.Println("  mp4label view [options]    Start web server (read-only mode)")
	fmt.Println("  mp4label version           Show version info")
	fmt.Println("  mp4label help              Show this help message")
	fmt.Println()
	fmt.Println("Web server options (web):")
	fmt.Println("  -port string               Server port (default: 8080)")
	fmt.Println("  -task-dir string           Task directory with multiple .txt task files (optional)")
	fmt.Println()
	fmt.Println("View server options (view):")
	fmt.Println("  -port string               Server port (default: 8080)")
	fmt.Println("  -video-dir string          Video directory (required)")
	fmt.Println("  -output-dir string         Output/annotation directory (required)")
	fmt.Println("  -pre-annotation-dir string Pre-annotation directory (optional)")
	fmt.Println("  -task-file string          Single task file path (optional)")
	fmt.Println("  -task-dir string           Task directory with multiple .txt task files (optional)")
	fmt.Println("  -model-annotation-dir string Model annotation directory (optional)")
	fmt.Println()
	fmt.Println("  Note: If both -task-dir and -task-file are specified, -task-dir takes priority.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mp4label web                           # Edit mode on port 8080")
	fmt.Println("  mp4label web -task-dir /path/to/tasks  # Edit mode with task groups")
	fmt.Println("  mp4label view -video-dir /videos -output-dir /output")
	fmt.Println("  mp4label view -video-dir /videos -output-dir /output -task-dir /tasks")
}

// runWebServer starts the server in edit mode (loads config from file, with optional CLI overrides)
func runWebServer() {
	webCmd := flag.NewFlagSet("web", flag.ExitOnError)
	port := webCmd.String("port", "8080", "Server port")
	taskDir := webCmd.String("task-dir", "", "Task directory with multiple .txt task files (optional)")

	webCmd.Parse(os.Args[2:])

	// Load config from file
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// CLI -task-dir overrides config file
	if *taskDir != "" {
		if cfg.TaskFile != "" {
			fmt.Println("WARNING: CLI -task-dir specified, overriding task_file from config. Using -task-dir.")
			cfg.TaskFile = ""
		}
		cfg.TaskDir = *taskDir
	}

	// If config has both TaskDir and TaskFile set, TaskDir takes priority
	if cfg.TaskDir != "" && cfg.TaskFile != "" {
		fmt.Println("WARNING: Both task_dir and task_file configured. Using task_dir (ignoring task_file).")
		cfg.TaskFile = ""
	}

	srv, err := server.NewServerWithConfig(webFS, cfg, false)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Printf("mp4Label %s\n", version)
	fmt.Printf("Server started at http://localhost:%s\n", *port)
	fmt.Println("Mode: Edit (read-write)")
	if cfg.TaskDir != "" {
		fmt.Printf("Task directory: %s\n", cfg.TaskDir)
	}
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	if err := srv.Start(*port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// runViewServer starts the server in read-only mode (config from CLI flags)
func runViewServer() {
	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	port := viewCmd.String("port", "8080", "Server port")
	videoDir := viewCmd.String("video-dir", "", "Video directory (required)")
	outputDir := viewCmd.String("output-dir", "", "Output/annotation directory (required)")
	preAnnotationDir := viewCmd.String("pre-annotation-dir", "", "Pre-annotation directory (optional)")
	taskFile := viewCmd.String("task-file", "", "Single task file path (optional)")
	taskDir := viewCmd.String("task-dir", "", "Task directory with multiple .txt task files (optional)")
	modelAnnotationDir := viewCmd.String("model-annotation-dir", "", "Model annotation directory (optional)")

	viewCmd.Parse(os.Args[2:])

	// Validate required flags
	if *videoDir == "" {
		fmt.Println("Error: -video-dir is required for view mode")
		fmt.Println()
		viewCmd.PrintDefaults()
		os.Exit(1)
	}
	if *outputDir == "" {
		fmt.Println("Error: -output-dir is required for view mode")
		fmt.Println()
		viewCmd.PrintDefaults()
		os.Exit(1)
	}

	// Handle -task-dir and -task-file conflict: -task-dir wins
	effectiveTaskFile := *taskFile
	effectiveTaskDir := *taskDir
	if effectiveTaskDir != "" && effectiveTaskFile != "" {
		fmt.Println("WARNING: Both -task-dir and -task-file specified. Using -task-dir (ignoring -task-file).")
		effectiveTaskFile = ""
	}

	// Build config from CLI flags
	cfg := &config.Config{
		VideoDir:           *videoDir,
		OutputDir:          *outputDir,
		PreAnnotationDir:   *preAnnotationDir,
		TaskFile:           effectiveTaskFile,
		TaskDir:            effectiveTaskDir,
		ModelAnnotationDir: *modelAnnotationDir,
	}

	// Validate config (checks directory existence, etc.)
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	// Create server in read-only mode
	srv, err := server.NewServerWithConfig(webFS, cfg, true)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Printf("mp4Label %s\n", version)
	fmt.Printf("Server started at http://localhost:%s\n", *port)
	fmt.Println("Mode: View (read-only)")
	fmt.Printf("Video directory: %s\n", *videoDir)
	fmt.Printf("Output directory: %s\n", *outputDir)
	if *preAnnotationDir != "" {
		fmt.Printf("Pre-annotation directory: %s\n", *preAnnotationDir)
	}
	if effectiveTaskFile != "" {
		fmt.Printf("Task file: %s\n", effectiveTaskFile)
	}
	if effectiveTaskDir != "" {
		fmt.Printf("Task directory: %s\n", effectiveTaskDir)
	}
	if *modelAnnotationDir != "" {
		fmt.Printf("Model annotation directory: %s\n", *modelAnnotationDir)
	}
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	if err := srv.Start(*port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
