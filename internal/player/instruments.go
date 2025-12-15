package player

/*
#cgo pkg-config: libopenmpt
#include <libopenmpt/libopenmpt.h>
#include <stdlib.h>
*/
import "C"
import "fmt"

// GetNumInstruments returns the number of instruments
func (m *Module) GetNumInstruments() int {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mod == nil {
		return 0
	}
	return int(C.openmpt_module_get_num_instruments(m.mod))
}

// GetInstrumentName returns the name of an instrument by index (1-based)
func (m *Module) GetInstrumentName(index int) string {
	if m == nil {
		return ""
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.getInstrumentNameLocked(index)
}

func (m *Module) getInstrumentNameLocked(index int) string {
	if m.mod == nil {
		return ""
	}

	cIndex := C.int(index - 1) // openmpt uses 0-based indexing
	cName := C.openmpt_module_get_instrument_name(m.mod, cIndex)
	if cName == nil {
		return ""
	}
	defer C.openmpt_free_string(cName)

	name := C.GoString(cName)
	if name == "" {
		return fmt.Sprintf("Instrument %02X", index)
	}
	return name
}

// GetNumSamples returns the number of samples
func (m *Module) GetNumSamples() int {
	if m == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.mod == nil {
		return 0
	}
	return int(C.openmpt_module_get_num_samples(m.mod))
}

// GetSampleName returns the name of a sample by index (1-based)
func (m *Module) GetSampleName(index int) string {
	if m == nil {
		return ""
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.getSampleNameLocked(index)
}

func (m *Module) getSampleNameLocked(index int) string {
	if m.mod == nil {
		return ""
	}

	cIndex := C.int(index - 1) // openmpt uses 0-based indexing
	cName := C.openmpt_module_get_sample_name(m.mod, cIndex)
	if cName == nil {
		return ""
	}
	defer C.openmpt_free_string(cName)

	name := C.GoString(cName)
	if name == "" {
		return fmt.Sprintf("Sample %02X", index)
	}
	return name
}

// GetInstrumentList returns all instruments with their names
// Optimized to hold a single lock for the entire duration
func (m *Module) GetInstrumentList() []Instrument {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.mod == nil {
		return []Instrument{}
	}

	numInst := int(C.openmpt_module_get_num_instruments(m.mod))
	numSamp := int(C.openmpt_module_get_num_samples(m.mod))

	instruments := []Instrument{}

	// Try instruments first (XM, IT use instruments)
	if numInst > 0 {
		for i := 1; i <= numInst && i <= 100; i++ { // Limit to 100
			name := m.getInstrumentNameLocked(i)
			if name != "" {
				instruments = append(instruments, Instrument{
					ID:   i,
					Name: name,
					Type: "Inst",
				})
			}
		}
	} else if numSamp > 0 {
		// Fallback to samples (MOD, S3M use samples)
		for i := 1; i <= numSamp && i <= 100; i++ { // Limit to 100
			name := m.getSampleNameLocked(i)
			if name != "" {
				instruments = append(instruments, Instrument{
					ID:   i,
					Name: name,
					Type: "Samp",
				})
			}
		}
	}

	return instruments
}

// GetRowInstruments returns a map of active instruments for the given row
// The map keys are instrument IDs, and values are initial brightness (100)
// This logic was moved from the UI to the player to separate concerns
func (m *Module) GetRowInstruments(row int) map[int]int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.mod == nil {
		return nil
	}

	activeMap := make(map[int]int)
	currentPattern := int(C.openmpt_module_get_current_pattern(m.mod))
	numChannels := int(C.openmpt_module_get_num_channels(m.mod))

	// Iterate all channels to see what instruments are triggered
	for ch := 0; ch < numChannels; ch++ {
		// Command 1 = Instrument (OPENMPT_MODULE_COMMAND_INSTRUMENT)
		inst := int(C.openmpt_module_get_pattern_row_channel_command(m.mod, C.int(currentPattern), C.int(row), C.int(ch), 1))
		if inst > 0 {
			activeMap[inst] = 100
		}
	}
	return activeMap
}

// Instrument represents a sample or instrument
type Instrument struct {
	ID   int
	Name string
	Type string // "Inst" or "Samp"
}