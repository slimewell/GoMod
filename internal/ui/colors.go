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
		Title:     lipgloss.Color("86"),  // Light cyan
		InfoLabel: lipgloss.Color("241"), // Gray
		InfoValue: lipgloss.Color("147"), // Light purple

		RowNumber:    lipgloss.Color("240"), // Dark gray
		CurrentRow:   lipgloss.Color("15"),  // White
		CurrentRowBg: lipgloss.Color("239"), // Medium gray background
		Note:         lipgloss.Color("51"),  // Cyan
		Instrument:   lipgloss.Color("226"), // Yellow
		Volume:       lipgloss.Color("141"), // Purple
		Effect:       lipgloss.Color("118"), // Green

		Border:   lipgloss.Color("240"), // Dark gray
		Controls: lipgloss.Color("240"), // Dark gray
	}
}

// CyberpunkPalette returns a vibrant cyberpunk color scheme
func CyberpunkPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("201"), // Hot pink
		InfoLabel: lipgloss.Color("93"),  // Purple
		InfoValue: lipgloss.Color("51"),  // Cyan

		RowNumber:    lipgloss.Color("240"), // Dark gray
		CurrentRow:   lipgloss.Color("15"),  // White
		CurrentRowBg: lipgloss.Color("53"),  // Dark purple background
		Note:         lipgloss.Color("51"),  // Bright cyan
		Instrument:   lipgloss.Color("201"), // Hot pink
		Volume:       lipgloss.Color("213"), // Light pink
		Effect:       lipgloss.Color("46"),  // Bright green

		Border:   lipgloss.Color("93"), // Purple
		Controls: lipgloss.Color("93"), // Purple
	}
}

// RetroAmberPalette returns a retro amber monochrome scheme
func RetroAmberPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("214"), // Orange
		InfoLabel: lipgloss.Color("58"),  // Dark orange
		InfoValue: lipgloss.Color("220"), // Light orange

		RowNumber:    lipgloss.Color("94"),  // Dark brown
		CurrentRow:   lipgloss.Color("220"), // Bright orange
		CurrentRowBg: lipgloss.Color("58"),  // Dark orange background
		Note:         lipgloss.Color("214"), // Orange
		Instrument:   lipgloss.Color("220"), // Light orange
		Volume:       lipgloss.Color("178"), // Medium orange
		Effect:       lipgloss.Color("172"), // Darker orange

		Border:   lipgloss.Color("58"), // Dark orange
		Controls: lipgloss.Color("58"), // Dark orange
	}
}

// GreenScreenPalette returns a classic green monochrome scheme
func GreenScreenPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("46"), // Bright green
		InfoLabel: lipgloss.Color("22"), // Dark green
		InfoValue: lipgloss.Color("83"), // Light green

		RowNumber:    lipgloss.Color("28"), // Medium dark green
		CurrentRow:   lipgloss.Color("46"), // Bright green
		CurrentRowBg: lipgloss.Color("22"), // Dark green background
		Note:         lipgloss.Color("46"), // Bright green
		Instrument:   lipgloss.Color("83"), // Light green
		Volume:       lipgloss.Color("77"), // Medium green
		Effect:       lipgloss.Color("71"), // Pale green

		Border:   lipgloss.Color("22"), // Dark green
		Controls: lipgloss.Color("22"), // Dark green
	}
}

// OceanPalette returns a calm blue/teal color scheme
func OceanPalette() ColorPalette {
	return ColorPalette{
		Title:     lipgloss.Color("81"),  // Bright cyan
		InfoLabel: lipgloss.Color("24"),  // Dark blue
		InfoValue: lipgloss.Color("117"), // Light blue

		RowNumber:    lipgloss.Color("239"), // Gray
		CurrentRow:   lipgloss.Color("15"),  // White
		CurrentRowBg: lipgloss.Color("24"),  // Dark blue background
		Note:         lipgloss.Color("81"),  // Bright cyan
		Instrument:   lipgloss.Color("117"), // Light blue
		Volume:       lipgloss.Color("75"),  // Teal
		Effect:       lipgloss.Color("43"),  // Aqua

		Border:   lipgloss.Color("24"), // Dark blue
		Controls: lipgloss.Color("24"), // Dark blue
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
		Title:     lipgloss.Color("#08CB00"), // Bright green
		InfoLabel: lipgloss.Color("#056600"), // Dark green
		InfoValue: lipgloss.Color("#0aff00"), // Neon green

		RowNumber:    lipgloss.Color("#056600"), // Dark green
		CurrentRow:   lipgloss.Color("#0aff00"), // Neon green
		CurrentRowBg: lipgloss.Color("#003300"), // Very dark green
		Note:         lipgloss.Color("#08CB00"), // Bright green
		Instrument:   lipgloss.Color("#0aff00"), // Neon green
		Volume:       lipgloss.Color("#05a600"), // Medium green
		Effect:       lipgloss.Color("#03d600"), // Lime green

		Border:   lipgloss.Color("#08CB00"),
		Controls: lipgloss.Color("#08CB00"),
	}
}

// GetPalette returns a palette by name
func GetPalette(name string) ColorPalette {
	switch name {
	case "cyberpunk":
		return CyberpunkPalette()
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
