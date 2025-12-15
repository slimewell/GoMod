package player

/*
#cgo pkg-config: libopenmpt
#include <libopenmpt/libopenmpt.h>
#include <libopenmpt/libopenmpt_ext.h>
#include <stdlib.h>
#include <string.h>

// Force declaration if missing from pkg-config header path or visibility
float openmpt_module_get_current_channel_vu_mono( openmpt_module * mod, int32_t channel );
double openmpt_module_get_position_seconds( openmpt_module * mod );
int openmpt_module_set_render_param( openmpt_module * mod, int param, int32_t value );

// Wrapper helpers to call interface function pointers safe from CGo
// We don't need to store the interface struct in Go, we can just fetch it when needed.
// It's a lightweight operation.

int ext_set_channel_mute(openmpt_module_ext *mod_ext, int32_t channel, int mute) {
    if (!mod_ext) return 0;
    openmpt_module_ext_interface_interactive interactive;
    memset(&interactive, 0, sizeof(interactive));

    // Retrieve the interactive interface
    if (openmpt_module_ext_get_interface(mod_ext, "interactive", &interactive, sizeof(interactive)) != 0) {
        if (interactive.set_channel_mute_status) {
            return interactive.set_channel_mute_status(mod_ext, channel, mute);
        }
    }
    return 0; // Failed
}

int ext_get_channel_mute(openmpt_module_ext *mod_ext, int32_t channel) {
    if (!mod_ext) return 0;
    openmpt_module_ext_interface_interactive interactive;
    memset(&interactive, 0, sizeof(interactive));

    if (openmpt_module_ext_get_interface(mod_ext, "interactive", &interactive, sizeof(interactive)) != 0) {
        if (interactive.get_channel_mute_status) {
            return interactive.get_channel_mute_status(mod_ext, channel);
        }
    }
    return 0; // Default to unmuted if failed
}

*/
import "C"
import (
	"errors"
	"fmt"
	"os"
	"sync"
	"unsafe"
)

// Module wraps an openmpt_module handle
type Module struct {
	modExt         *C.openmpt_module_ext
	mod            *C.openmpt_module
	mu             sync.Mutex
	patternCache   map[int]*CachedPattern
	cachedMetadata *Metadata
	channelMuted   []bool // Track which channels are muted
}

// LoadModule loads a tracker module from a file path
func LoadModule(path string) (*Module, error) {
	// Read file into memory
	filedata, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create module from memory using extended API
	// Note: We use openmpt_module_ext_create_from_memory (not ...2) as it is the standard ext API
	modExt := C.openmpt_module_ext_create_from_memory(
		unsafe.Pointer(&filedata[0]),
		C.size_t(len(filedata)),
		nil, // log func
		nil, // log user
		nil, // error func
		nil, // error user
		nil, // error code
		nil, // error message
		nil, // initial ctls
	)

	if modExt == nil {
		return nil, errors.New("failed to load module: unsupported format or corrupted file")
	}

	mod := C.openmpt_module_ext_get_module(modExt)
	if mod == nil {
		C.openmpt_module_ext_destroy(modExt)
		return nil, errors.New("failed to get module interface from extended module")
	}

	m := &Module{
		modExt: modExt,
		mod:    mod,
	}

	// Initialize channel muted state
	numChannels := int(C.openmpt_module_get_num_channels(mod))
	m.channelMuted = make([]bool, numChannels)

	return m, nil
}

// SetStereoSeparation sets the stereo separation percentage (0-200)
// 0 = mono, 100 = default, 200 = full separation
func (m *Module) SetStereoSeparation(percent int) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	// OPENMPT_MODULE_RENDER_STEREOSEPARATION_PERCENT = 2
	const OPENMPT_MODULE_RENDER_STEREOSEPARATION_PERCENT = 2

	// We use the standard API for render params, which works on the underlying module
	result := C.openmpt_module_set_render_param(m.mod, C.int(OPENMPT_MODULE_RENDER_STEREOSEPARATION_PERCENT), C.int32_t(percent))
	if result != 1 {
		return fmt.Errorf("failed to set stereo separation to %d", percent)
	}

	return nil
}

// SetInterpolationFilter sets the interpolation quality
// 0 = default, 1 = none, 2 = linear, 4 = cubic, 8 = windowed sinc (best quality)
func (m *Module) SetInterpolationFilter(length int) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	const OPENMPT_MODULE_RENDER_INTERPOLATIONFILTER_LENGTH = 3

	result := C.openmpt_module_set_render_param(m.mod, C.int(OPENMPT_MODULE_RENDER_INTERPOLATIONFILTER_LENGTH), C.int32_t(length))
	if result != 1 {
		return fmt.Errorf("failed to set interpolation filter to %d", length)
	}

	return nil
}

// GetMetadata returns metadata about the module
type Metadata struct {
	Title    string
	Artist   string
	Type     string
	Duration float64
	Channels int
}

func (m *Module) GetMetadata() Metadata {
	if m == nil {
		return Metadata{
			Title:    "Unknown",
			Artist:   "",
			Type:     "Unknown",
			Duration: 0,
			Channels: 0,
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return cached metadata if available
	if m.cachedMetadata != nil {
		return *m.cachedMetadata
	}

	// Safety check
	if m.mod == nil {
		return Metadata{
			Title:    "Unknown",
			Artist:   "",
			Type:     "Unknown",
			Duration: 0,
			Channels: 0,
		}
	}

	// Cache metadata on first call (reduces CGo overhead)
	metadata := Metadata{
		Title:    m.getMetadataString("title"),
		Artist:   m.getMetadataString("artist"),
		Type:     m.getMetadataString("type_long"),
		Duration: float64(C.openmpt_module_get_duration_seconds(m.mod)),
		Channels: int(C.openmpt_module_get_num_channels(m.mod)),
	}
	m.cachedMetadata = &metadata

	return metadata
}

// GetPositionSeconds returns the current playback position in seconds
func (m *Module) GetPositionSeconds() float64 {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mod == nil {
		return 0
	}
	return float64(C.openmpt_module_get_position_seconds(m.mod))
}

func (m *Module) getMetadataString(key string) string {
	// Mutex is expected to be held by caller (GetMetadata)

	if m == nil || m.mod == nil {
		return ""
	}

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cValue := C.openmpt_module_get_metadata(m.mod, cKey)
	if cValue == nil {
		return ""
	}
	defer C.openmpt_free_string(cValue)

	return C.GoString(cValue)
}

// GetCurrentRow returns the current row being rendered (note: ahead of audio playback)
func (m *Module) GetCurrentRow() int {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mod == nil {
		return 0
	}
	return int(C.openmpt_module_get_current_row(m.mod))
}

// GetCurrentPattern returns the current pattern being played
func (m *Module) GetCurrentPattern() int {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mod == nil {
		return 0
	}
	return int(C.openmpt_module_get_current_pattern(m.mod))
}

// GetNumPatterns returns the total number of patterns
func (m *Module) GetNumPatterns() int {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mod == nil {
		return 0
	}
	return int(C.openmpt_module_get_num_patterns(m.mod))
}

// GetPatternNumRows returns the number of rows in a pattern
func (m *Module) GetPatternNumRows(pattern int) int {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mod == nil {
		return 0
	}
	return int(C.openmpt_module_get_pattern_num_rows(m.mod, C.int(pattern)))
}

// GetPatternRowChannelCommand gets pattern data for a specific row/channel
func (m *Module) GetPatternRowChannelCommand(pattern, row, channel, command int) int {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mod == nil {
		return 0
	}
	return int(C.openmpt_module_get_pattern_row_channel_command(
		m.mod,
		C.int(pattern),
		C.int(row),
		C.int(channel),
		C.int(command),
	))
}

// Read renders audio into the provided buffer (interleaved stereo int16)
// Returns the number of frames read
func (m *Module) Read(buf []int16) int {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(buf) == 0 {
		return 0
	}

	frames := len(buf) / 2 // stereo
	count := C.openmpt_module_read_interleaved_stereo(
		m.mod,
		C.int(44100), // sample rate
		C.size_t(frames),
		(*C.int16_t)(unsafe.Pointer(&buf[0])),
	)

	return int(count)
}

// Close frees the module
func (m *Module) Close() {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// Destroy the extended module, which also destroys the underlying module
	if m.modExt != nil {
		C.openmpt_module_ext_destroy(m.modExt)
		m.modExt = nil
		m.mod = nil
	} else if m.mod != nil {
		// Fallback if we somehow only have mod
		C.openmpt_module_destroy(m.mod)
		m.mod = nil
	}
}

// ToggleChannelMute toggles the mute state of a channel
func (m *Module) ToggleChannelMute(channel int) bool {
	if m == nil || m.modExt == nil {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if channel < 0 || channel >= len(m.channelMuted) {
		return false
	}

	// Get current state from OpenMPT via helper
	// Note: We need to check if the helper returns valid data, but here we trust our tracking mostly,
	// except we want to sync with the engine.
	currentMute := C.ext_get_channel_mute(m.modExt, C.int32_t(channel))
	newMute := 1 - currentMute // Toggle

	C.ext_set_channel_mute(m.modExt, C.int32_t(channel), C.int(newMute))

	m.channelMuted[channel] = (newMute == 1)
	return m.channelMuted[channel]
}

// SoloChannel solos the specified channel (unmutes it, mutes all others)
func (m *Module) SoloChannel(channel int) {
	if m == nil || m.modExt == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if channel < 0 || channel >= len(m.channelMuted) {
		return
	}

	// Check if this channel is currently the ONLY one unmuted
	isSoloed := true
	for i, muted := range m.channelMuted {
		if i == channel {
			if muted {
				isSoloed = false // The target channel is muted
				break
			}
		} else {
			if !muted {
				isSoloed = false // Another channel is unmuted
				break
			}
		}
	}

	if isSoloed {
		// Unsolo all: Unmute everyone
		for i := range m.channelMuted {
			C.ext_set_channel_mute(m.modExt, C.int32_t(i), 0)
			m.channelMuted[i] = false
		}
	} else {
		// Solo this channel: Mute everyone else, unmute this one
		for i := range m.channelMuted {
			if i == channel {
				C.ext_set_channel_mute(m.modExt, C.int32_t(i), 0)
				m.channelMuted[i] = false
			} else {
				C.ext_set_channel_mute(m.modExt, C.int32_t(i), 1)
				m.channelMuted[i] = true
			}
		}
	}
}

// IsChannelMuted returns the mute state of a channel
func (m *Module) IsChannelMuted(channel int) bool {
	if m == nil || channel < 0 || channel >= len(m.channelMuted) {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.channelMuted[channel]
}

// Data structures for snapshot
type PatternSnapshot struct {
	CurrentRow     int
	CurrentPattern int
	NumChannels    int
	ChannelVolumes []float64
	Rows           []PatternRow
}

type PatternRow struct {
	RowNumber int
	Channels  []PatternCell
}

type PatternCell struct {
	Note       int
	Instrument int
	Volume     int
	Effect     int
}

// CachedPattern stores the full content of a pattern in Go memory
type CachedPattern struct {
	Rows []PatternRow
}

// GetPatternSnapshot efficiently retrieves a range of pattern data
func (m *Module) GetPatternSnapshot(visibleRows int) PatternSnapshot {
	if m == nil {
		return PatternSnapshot{}
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.mod == nil {
		return PatternSnapshot{}
	}

	currentRow := int(C.openmpt_module_get_current_row(m.mod))
	currentPattern := int(C.openmpt_module_get_current_pattern(m.mod))
	numChannels := int(C.openmpt_module_get_num_channels(m.mod))
	// No CGo call needed for numRows if we have it cached, but we need it to check cache or load
	numRows := int(C.openmpt_module_get_pattern_num_rows(m.mod, C.int(currentPattern)))

	snapshot := PatternSnapshot{
		CurrentRow:     currentRow,
		CurrentPattern: currentPattern,
		NumChannels:    numChannels,
		ChannelVolumes: make([]float64, numChannels),
	}

	// Fetch VU data (fast/essential, so we keep doing this every frame)
	for ch := 0; ch < numChannels; ch++ {
		snapshot.ChannelVolumes[ch] = float64(C.openmpt_module_get_current_channel_vu_mono(m.mod, C.int(ch)))
	}

	// --- PATTERN CACHING START ---

	// Ensure cache is initialized
	if m.patternCache == nil {
		m.patternCache = make(map[int]*CachedPattern)
	}

	// Check if this pattern is cached
	cached, exists := m.patternCache[currentPattern]
	if !exists {
		// Cache Miss: Load the ENTIRE pattern now
		// This is a heavy operation, but only happens once per pattern visit
		rows := make([]PatternRow, numRows)

		for r := 0; r < numRows; r++ {
			rowStr := PatternRow{
				RowNumber: r,
				Channels:  make([]PatternCell, numChannels),
			}
			for c := 0; c < numChannels; c++ {
				rowStr.Channels[c] = PatternCell{
					Note:       int(C.openmpt_module_get_pattern_row_channel_command(m.mod, C.int(currentPattern), C.int(r), C.int(c), 0)), // Note
					Instrument: int(C.openmpt_module_get_pattern_row_channel_command(m.mod, C.int(currentPattern), C.int(r), C.int(c), 1)), // Inst
					Volume:     int(C.openmpt_module_get_pattern_row_channel_command(m.mod, C.int(currentPattern), C.int(r), C.int(c), 2)), // Vol
					Effect:     int(C.openmpt_module_get_pattern_row_channel_command(m.mod, C.int(currentPattern), C.int(r), C.int(c), 3)), // Effect
				}
			}
			rows[r] = rowStr
		}

		cached = &CachedPattern{Rows: rows}
		m.patternCache[currentPattern] = cached
	}

	// --- PATTERN CACHING END ---

	// Calculate range for "Typewriter" style scrolling
	half := visibleRows / 2
	startRow := currentRow - half
	endRow := currentRow + half

	// Fill snapshot.Rows from Cache
	for r := startRow; r <= endRow; r++ {
		if r >= 0 && r < len(cached.Rows) {
			// Row exists in module/cache
			snapshot.Rows = append(snapshot.Rows, cached.Rows[r])
		} else {
			// Out of bounds (before start or after end of pattern)
			// Return empty row structure with correct row number for UI
			snapshot.Rows = append(snapshot.Rows, PatternRow{
				RowNumber: r,
				// Empty channels
			})
		}
	}

	return snapshot
}
