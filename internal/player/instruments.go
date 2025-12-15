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
// Note: This calls other internal methods that lock, so we shouldn't lock here to avoid deadlock
// BUT, since we modified them to lock, we have a problem.
// We should NOT lock here if we call public methods that lock.
// Instead, we should rely on the public methods locking themselves.
func (m *Module) GetInstrumentList() []Instrument {
	numInst := m.GetNumInstruments()
	numSamp := m.GetNumSamples()

	instruments := []Instrument{}

	// Try instruments first (XM, IT use instruments)
	if numInst > 0 {
		for i := 1; i <= numInst && i <= 100; i++ { // Limit to 100
			name := m.GetInstrumentName(i)
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
			name := m.GetSampleName(i)
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

// Instrument represents a sample or instrument
type Instrument struct {
	ID   int
	Name string
	Type string // "Inst" or "Samp"
}
