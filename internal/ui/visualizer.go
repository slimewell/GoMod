package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderVUMeters creates a 3-row tall vertical bar per channel (fills bottom-to-top)
// Muted channels are dimmed/grayed out
func RenderVUMeters(volumes []float64, mutedChannels []bool, width int, palette ColorPalette) string {
	numChannels := len(volumes)
	if numChannels == 0 {
		return ""
	}

	// Styles
	meterStyle := lipgloss.NewStyle().Foreground(palette.Note)
	highStyle := lipgloss.NewStyle().Foreground(palette.InfoValue).Bold(true)
	peakStyle := lipgloss.NewStyle().Foreground(palette.CurrentRow).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(palette.InfoLabel)
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Dimmed gray

	// Build 3 rows - each channel gets a 3-char tall vertical bar
	var rows [3]strings.Builder

	// Row labels
	labels := []string{"VU▲", "VU│", "VU▼"}

	for rowIdx := 0; rowIdx < 3; rowIdx++ {
		// Start with row label
		rows[rowIdx].WriteString(labelStyle.Render(fmt.Sprintf("%-4s", labels[rowIdx])))
		rows[rowIdx].WriteString(" │")

		// For each channel, render that row of its vertical bar
		for chIdx, vol := range volumes {
			// Check if this channel is muted
			isMuted := chIdx < len(mutedChannels) && mutedChannels[chIdx]

			// Apply perceptual scaling
			if vol < 0 {
				vol = 0
			}
			if vol > 1 {
				vol = 1
			}

			// Perceptual curve
			perceived := vol * vol
			if perceived < 0.1 {
				perceived = perceived * 1.5
			}

			// Map to 0-9 level range (3 rows × 3 levels = 9 total levels)
			level := int(perceived * 9.0)
			if level > 9 {
				level = 9
			}

			// Determine which rows should be lit based on level
			// 0-2: bottom row only
			// 3-5: bottom + middle
			// 6-9: all three rows

			var char string
			var style lipgloss.Style

			// If muted, override style with dimmed gray
			if isMuted {
				style = mutedStyle
			}

			switch rowIdx {
			case 0: // Top row
				if level >= 7 {
					// Peak levels
					char = "█"
					if !isMuted {
						if level >= 9 {
							style = peakStyle
						} else {
							style = highStyle
						}
					}
				} else if level == 6 {
					// Just touching top
					char = "▄"
					if !isMuted {
						style = highStyle
					}
				} else {
					// Below top threshold
					char = " "
					if !isMuted {
						style = meterStyle
					}
				}

			case 1: // Middle row
				if level >= 6 {
					// Full bar in middle
					char = "█"
					if !isMuted {
						style = highStyle
					}
				} else if level >= 4 {
					// Partial fill
					if level == 5 {
						char = "▇"
					} else {
						char = "▄"
					}
					if !isMuted {
						style = meterStyle
					}
				} else if level == 3 {
					// Just reached middle
					char = "▂"
					if !isMuted {
						style = meterStyle
					}
				} else {
					// Below middle threshold
					char = " "
					if !isMuted {
						style = meterStyle
					}
				}

			case 2: // Bottom row
				if level >= 3 {
					// Full bar at bottom (if we're in middle or above)
					char = "█"
					if !isMuted {
						style = meterStyle
					}
				} else if level > 0 {
					// Partial bottom fill
					switch level {
					case 1:
						char = "▂"
					case 2:
						char = "▄"
					}
					if !isMuted {
						style = meterStyle
					}
				} else {
					// Silent
					char = " "
					if !isMuted {
						style = meterStyle
					}
				}
			}

			// Column width: 14 chars to match pattern
			cellContent := fmt.Sprintf("%s%s%s",
				strings.Repeat(" ", 6),
				style.Render(char),
				strings.Repeat(" ", 7),
			)

			rows[rowIdx].WriteString(cellContent)
			rows[rowIdx].WriteString(" │")
		}
	}

	// Join all 3 rows with newlines
	return rows[0].String() + "\n" + rows[1].String() + "\n" + rows[2].String()
}
