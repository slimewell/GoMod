package ui

import "github.com/charmbracelet/lipgloss"

// ColorPalette defines the color scheme for the TUI
type ColorPalette struct {
	// Header colors
	Title     lipgloss.Color
	InfoLabel lipgloss.Color
	InfoValue lipgloss.Color

	// Pattern view colors
	RowNumber    lipgloss.Color
	CurrentRow   lipgloss.Color
	CurrentRowBg lipgloss.Color
	Note         lipgloss.Color
	Instrument   lipgloss.Color
	Volume       lipgloss.Color
	Effect       lipgloss.Color

	// UI elements
	Border   lipgloss.Color
	Controls lipgloss.Color
}

// DefaultPalette returns the default color scheme
func DefaultPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("#EAE0CF"), // Cream
		InfoLabel: lipgloss.Color("#94A3B8"), // Slate gray
		InfoValue: lipgloss.Color("#E2E8F0"), // Light slate

		RowNumber:    lipgloss.Color("#64748B"), // Medium slate
		CurrentRow:   lipgloss.Color("#EAE0CF"), // Cream
		CurrentRowBg: lipgloss.Color("#334155"), // Dark slate
		Note:         lipgloss.Color("#E2E8F0"), // Light slate
		Instrument:   lipgloss.Color("#CBD5E1"), // Slate 300
		Volume:       lipgloss.Color("#94A3B8"), // Slate 400
		Effect:       lipgloss.Color("#64748B"), // Slate 500

		Border:   lipgloss.Color("#475569"), // Slate 600
		Controls: lipgloss.Color("#94A3B8"), // Slate 400
	}
}

// RetroAmberPalette returns a retro amber monochrome scheme
func RetroAmberPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("#FFB000"), // Warm amber
		InfoLabel: lipgloss.Color("#D99058"), // Dark orange
		InfoValue: lipgloss.Color("#FFC947"), // Light amber

		RowNumber:    lipgloss.Color("#D99058"), // Clean dark orange
		CurrentRow:   lipgloss.Color("#FFFFFF"), // White
		CurrentRowBg: lipgloss.Color("#E65100"), // Deep rich orange
		Note:         lipgloss.Color("#FFB000"), // Warm amber
		Instrument:   lipgloss.Color("#FFC947"), // Light amber
		Volume:       lipgloss.Color("#FF9500"), // Deep amber
		Effect:       lipgloss.Color("#FFCC80"), // Pale amber

		Border:   lipgloss.Color("#D99058"),
		Controls: lipgloss.Color("#FFB000"),
	}
}

// GreenScreenPalette returns a classic green monochrome scheme
func GreenScreenPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("#33FF33"), // Classic phosphor green
		InfoLabel: lipgloss.Color("#2E8B57"), // SeaGreen
		InfoValue: lipgloss.Color("#90EE90"), // LightGreen

		RowNumber:    lipgloss.Color("#228B22"), // ForestGreen
		CurrentRow:   lipgloss.Color("#FFFFFF"), // White
		CurrentRowBg: lipgloss.Color("#006400"), // DarkGreen
		Note:         lipgloss.Color("#33FF33"), // Classic phosphor green
		Instrument:   lipgloss.Color("#98FB98"), // PaleGreen
		Volume:       lipgloss.Color("#32CD32"), // LimeGreen
		Effect:       lipgloss.Color("#00FF7F"), // SpringGreen

		Border:   lipgloss.Color("#2E8B57"), // SeaGreen
		Controls: lipgloss.Color("#3CB371"), // MediumSeaGreen
	}
}

// OceanPalette returns a calm blue/teal color scheme
func OceanPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("#06B6D4"), // Vibrant cyan
		InfoLabel: lipgloss.Color("#155E75"), // Deep teal
		InfoValue: lipgloss.Color("#67E8F9"), // Light cyan

		RowNumber:    lipgloss.Color("#0E7490"), // Dark cyan
		CurrentRow:   lipgloss.Color("#E0F2FE"), // Very light blue
		CurrentRowBg: lipgloss.Color("#164E63"), // Deep ocean blue
		Note:         lipgloss.Color("#22D3EE"), // Bright cyan
		Instrument:   lipgloss.Color("#67E8F9"), // Light cyan
		Volume:       lipgloss.Color("#06B6D4"), // Vibrant cyan
		Effect:       lipgloss.Color("#5EEAD4"), // Teal

		Border:   lipgloss.Color("#155E75"), // Deep teal
		Controls: lipgloss.Color("#0E7490"), // Dark cyan
	}
}

// PeachyPalette returns a soft peachy-pink gradient palette
func PeachyPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("#FD7979"), // Coral pink
		InfoLabel: lipgloss.Color("#FFCDC9"), // Light peach
		InfoValue: lipgloss.Color("#FEEAC9"), // Cream peach

		RowNumber:    lipgloss.Color("#FFCDC9"), // Light peach
		CurrentRow:   lipgloss.Color("#FD7979"), // Coral pink
		CurrentRowBg: lipgloss.Color("#3d2929"), // Dark pink-brown
		Note:         lipgloss.Color("#FD7979"), // Coral pink
		Instrument:   lipgloss.Color("#FEEAC9"), // Cream peach
		Volume:       lipgloss.Color("#FFCDC9"), // Light peach
		Effect:       lipgloss.Color("#f5a5a0"), // Medium coral

		Border:   lipgloss.Color("#FFCDC9"),
		Controls: lipgloss.Color("#FFCDC9"),
	}
}

// PurpleDreamPalette returns a vibrant purple-to-pink gradient
func PurpleDreamPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("#9B5DE0"), // Purple
		InfoLabel: lipgloss.Color("#4E56C0"), // Deep blue-purple
		InfoValue: lipgloss.Color("#FDCFFA"), // Light pink

		RowNumber:    lipgloss.Color("#7961b8"), // Medium purple
		CurrentRow:   lipgloss.Color("#FDCFFA"), // Light pink
		CurrentRowBg: lipgloss.Color("#4E56C0"), // Deep blue-purple
		Note:         lipgloss.Color("#FDCFFA"), // Light pink
		Instrument:   lipgloss.Color("#9B5DE0"), // Purple
		Volume:       lipgloss.Color("#c199eb"), // Light purple
		Effect:       lipgloss.Color("#4E56C0"), // Deep blue-purple

		Border:   lipgloss.Color("#9B5DE0"),
		Controls: lipgloss.Color("#9B5DE0"),
	}
}

// PastelPalette returns a soft pastel multi-color palette
func PastelPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("#F39EB6"), // Pink
		InfoLabel: lipgloss.Color("#B8DB80"), // Green
		InfoValue: lipgloss.Color("#FFE4EF"), // Light pink

		RowNumber:    lipgloss.Color("#d4d2b3"), // Muted beige
		CurrentRow:   lipgloss.Color("#F39EB6"), // Pink
		CurrentRowBg: lipgloss.Color("#5a6d4a"), // Dark green
		Note:         lipgloss.Color("#F39EB6"), // Pink
		Instrument:   lipgloss.Color("#B8DB80"), // Green
		Volume:       lipgloss.Color("#FFE4EF"), // Light pink
		Effect:       lipgloss.Color("#F7F6D3"), // Cream

		Border:   lipgloss.Color("#B8DB80"),
		Controls: lipgloss.Color("#B8DB80"),
	}
}

// MatrixPalette returns a bright matrix-style green palette
func MatrixPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("#00FF41"), // Matrix green
		InfoLabel: lipgloss.Color("#003B00"), // Very dark green
		InfoValue: lipgloss.Color("#8AFF8A"), // Bright lime

		RowNumber:    lipgloss.Color("#00802B"), // Forest green
		CurrentRow:   lipgloss.Color("#FFFFFF"), // White
		CurrentRowBg: lipgloss.Color("#004400"), // Brighter dark green
		Note:         lipgloss.Color("#00FF41"), // Matrix green
		Instrument:   lipgloss.Color("#39FF14"), // Neon green
		Volume:       lipgloss.Color("#00D936"), // Vivid green
		Effect:       lipgloss.Color("#7FFF00"), // Chartreuse

		Border:   lipgloss.Color("#00802B"), // Forest green
		Controls: lipgloss.Color("#00B33C"), // Medium green
	}
}

// GetPalette returns a palette by name
func GetPalette(name string) ColorPalette {
	switch name {
	case "amber":
		return RetroAmberPalette()
	case "green":
		return GreenScreenPalette()
	case "ocean":
		return OceanPalette()
	case "peachy", "peach":
		return PeachyPalette()
	case "purple", "dream":
		return PurpleDreamPalette()
	case "pastel":
		return PastelPalette()
	case "matrix":
		return MatrixPalette()
	default:
		return DefaultPalette()
	}
}
