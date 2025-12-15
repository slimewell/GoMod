package main

import (
	"flag"
	"fmt"
	"github.com/slimewell/GoMod/internal/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load saved config
	cfg, _ := ui.LoadConfig() // Ignore errors, use defaults

	// Parse command-line flags
	stereoSep := flag.Int("separation", cfg.StereoSep, "Stereo separation percentage (0-100)")
	theme := flag.String("theme", cfg.Theme, "Color theme")
	flag.StringVar(theme, "t", cfg.Theme, "Color theme (shorthand)")
	flag.Parse()

	// Get file path argument
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <module-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSupported formats: .mod, .xm, .it, .s3m, and more\n")
		fmt.Fprintf(os.Stderr, "Available themes: default, cyberpunk, peachy, purple, pastel, matrix, amber, green, ocean\n")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	// Validate file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File not found: %s\n", filename)
		os.Exit(1)
	}

	// Validate stereo separation range
	if *stereoSep < 0 || *stereoSep > 100 {
		fmt.Fprintf(os.Stderr, "Error: Stereo separation must be between 0 and 100\n")
		os.Exit(1)
	}

	// Save config for next time
	cfg.Theme = *theme
	cfg.StereoSep = *stereoSep
	cfg.LastUsed = filename
	_ = ui.SaveConfig(cfg) // Ignore save errors

	// Create and run the TUI
	model := ui.NewModel(filename, *stereoSep, *theme)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
