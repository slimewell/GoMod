package player

import (
	"context"
	"sync"
	"time"

	"github.com/hajimehoshi/oto/v2"
)

const (
	sampleRate   = 44100
	channelCount = 2 // stereo
	// Buffer size tuned for low latency (~23ms)
	// We rely on hardware sync to keep the UI tight even with small buffers
	bufferSize = 1024
)

// Player manages audio playback
type Player struct {
	module     *Module
	otoContext *oto.Context
	otoPlayer  oto.Player
	mu         sync.RWMutex
	playing    bool

	// Sync mechanism
	stateQueue     []SyncState
	queueMu        sync.Mutex
	startSample    int64
	samplesWritten int64
}

// SyncState represents the state of the engine at a specific sample time
type SyncState struct {
	SampleCount    int64
	Row            int
	Pattern        int
	ChannelVolumes []float64
}

// NewPlayer creates a new player for the given module
func NewPlayer(module *Module) (*Player, error) {
	// Initialize Oto with small buffer for low latency
	// We MUST use NewContextWithOptions to control the device buffer size
	// Default is often too large (200ms+), causing mute lag
	options := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channelCount,
		Format:       2,                     // FormatSignedInt16LE
		BufferSize:   60 * time.Millisecond, // ~60ms total device buffer
	}

	otoContext, ready, err := oto.NewContextWithOptions(options)
	if err != nil {
		return nil, err
	}
	<-ready

	p := &Player{
		module:     module,
		otoContext: otoContext,
		playing:    false,
		stateQueue: make([]SyncState, 0, 100),
	}

	return p, nil
}

// GetSyncedState returns the module state corresponding to the CURRENT playback time
// precise to the audio buffer latency using hardware feedback
func (p *Player) GetSyncedState() (int, int, []float64) {
	p.mu.RLock()
	if p.otoPlayer == nil || !p.playing {
		p.mu.RUnlock()
		return 0, 0, nil
	}
	p.mu.RUnlock() // Release generic lock before acquiring queue lock

	p.queueMu.Lock()
	defer p.queueMu.Unlock()

	// Calculate what sample the hardware is currently playing
	// samplesWritten = total samples sent to oto
	// UnplayedBufferSize = bytes buffered in driver (convert to samples)
	// currentSample = samplesWritten - samplesBuffered
	if len(p.stateQueue) == 0 {
		return 0, 0, nil
	}

	unplayedBytes := p.otoPlayer.UnplayedBufferSize()
	unplayedSamples := int64(unplayedBytes) / 4 // 4 bytes per stereo sample (16bit * 2)
	currentSample := p.samplesWritten - unplayedSamples

	if currentSample < 0 {
		currentSample = 0
	}

	// Find the state that matches the current hardware sample count
	// We want the LAST state where SampleCount <= currentSample
	var bestState SyncState
	bestState = p.stateQueue[0] // Default to first available

	idx := -1
	for i, state := range p.stateQueue {
		if state.SampleCount > currentSample {
			break
		}
		bestState = state
		idx = i
	}

	// Cleanup old states to prevent memory leak
	// Keep a few previous states just in case
	if idx > 10 {
		// Keep from idx-5 onwards
		keepFrom := idx - 5
		if keepFrom < 0 {
			keepFrom = 0
		}
		// Create new slice to allow GC of underlying array eventually
		newQueue := make([]SyncState, len(p.stateQueue)-keepFrom)
		copy(newQueue, p.stateQueue[keepFrom:])
		p.stateQueue = newQueue
	}

	return bestState.Pattern, bestState.Row, bestState.ChannelVolumes
}

// GetSyncedTime returns the current playback time in seconds, sync'd to hardware
func (p *Player) GetSyncedTime() float64 {
	p.mu.RLock()
	if p.otoPlayer == nil || !p.playing {
		p.mu.RUnlock()
		return 0
	}
	p.mu.RUnlock()

	p.queueMu.Lock()
	defer p.queueMu.Unlock()

	unplayedBytes := p.otoPlayer.UnplayedBufferSize()
	unplayedSamples := int64(unplayedBytes) / 4
	currentSample := p.samplesWritten - unplayedSamples

	if currentSample < 0 {
		return 0
	}

	return float64(currentSample) / float64(sampleRate)
}

// Play starts playback
func (p *Player) Play(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.playing {
		return nil
	}

	p.otoPlayer = p.otoContext.NewPlayer(&audioReader{
		module: p.module,
		ctx:    ctx,
		player: p,
	})

	p.otoPlayer.Play()
	p.playing = true
	// Reset sync
	p.queueMu.Lock()
	p.stateQueue = p.stateQueue[:0]
	p.samplesWritten = 0
	p.queueMu.Unlock()

	return nil
}

// IsPlaying returns whether audio is currently playing
func (p *Player) IsPlaying() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.playing
}

// TogglePause toggles playback state
func (p *Player) TogglePause() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.otoPlayer == nil {
		return false
	}

	if p.playing {
		p.otoPlayer.Pause()
	} else {
		p.otoPlayer.Play()
	}
	p.playing = !p.playing
	return p.playing
}

// Close cleans up resources
func (p *Player) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.otoPlayer != nil {
		if err := p.otoPlayer.Close(); err != nil {
			return err
		}
	}

	// Context doesn't need closing in Oto v2

	return nil
}

// audioReader implements io.Reader for Oto
type audioReader struct {
	module *Module
	ctx    context.Context
	buf    [bufferSize]int16
	player *Player // Reference back to player to push states
}

func (r *audioReader) Read(p []byte) (int, error) {
	// Check context cancellation
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
	}

	// BEFORE rendering, capture the state that corresponds to the START of this buffer
	row := r.module.GetCurrentRow()
	pat := r.module.GetCurrentPattern()

	// Get VUs efficiently (module already has a helper or we can add one to Module if needed,
	// but GetPatternSnapshot does it. Let's assume we can get them.)
	// Actually we need to add GetChannelVolumes to Module to do this efficiently without building full snapshot.
	// For now let's reuse GetPatternSnapshot logic or adding a small helper in Module.
	// WAIT: We can use the C function directly if we were in the same package (which we are).
	// But Module struct hides C types.
	// Better approach: Module.GetSyncState()

	// Temporarily: We will use a small helper we add to Module in next step, or just use GetPatternSnapshot(0) which is cheap?
	// GetPatternSnapshot(0) builds a snapshot.
	// Let's assume we added GetCurrentVolumes to module.
	// Actually, let's just grab them via GetPatternSnapshot(0) for now, it's safe.
	snap := r.module.GetPatternSnapshot(0)

	// Push state to queue
	r.player.queueMu.Lock()
	r.player.stateQueue = append(r.player.stateQueue, SyncState{
		SampleCount:    r.player.samplesWritten,
		Row:            row,
		Pattern:        pat,
		ChannelVolumes: snap.ChannelVolumes,
	})
	r.player.queueMu.Unlock()

	// Render audio from openmpt
	frames := r.module.Read(r.buf[:])
	if frames == 0 {
		return 0, nil // End of module
	}

	// Convert int16 samples to bytes
	samples := frames * 2 // stereo
	bytesWritten := samples * 2

	r.player.queueMu.Lock()
	r.player.samplesWritten += int64(frames)
	r.player.queueMu.Unlock()

	for i := 0; i < samples; i++ {
		p[i*2] = byte(r.buf[i] & 0xff)
		p[i*2+1] = byte((r.buf[i] >> 8) & 0xff)
	}

	return bytesWritten, nil
}

// InstantMute toggles mute on a channel and performs a Flush & Seek to make it audible immediately
func (p *Player) InstantMute(channel int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.otoPlayer == nil || p.module == nil {
		return
	}

	// 1. Get current positions
	renderPos := p.module.GetPositionSeconds()

	// 2. Calculate latency and seek target
	unplayedBytes := p.otoPlayer.UnplayedBufferSize()
	bytesPerSec := float64(sampleRate * 4) // stereo, 16-bit = 4 bytes/sample
	bufferedSecs := float64(unplayedBytes) / bytesPerSec

	seekTarget := renderPos - bufferedSecs
	if seekTarget < 0 {
		seekTarget = 0
	}

	// 3. Log for debugging (Removed)
	// fmt.Printf("[MUTE] Render: %.3f, Buf: %.3f (%dB), Target: %.3f\n", renderPos, bufferedSecs, unplayedBytes, seekTarget)

	// 4. Toggle Mute
	p.module.ToggleChannelMute(channel)

	// 5. Seek Module to the "heard" position
	p.module.SetPositionSeconds(seekTarget)

	// 6. Flush Oto buffer
	// Reset clears the underlying buffer and pauses
	p.otoPlayer.Reset()

	// 7. Reset sync state to match the seek
	p.queueMu.Lock()
	p.stateQueue = p.stateQueue[:0]
	// Reset samplesWritten so GetSyncedTime() remains accurate to the new position
	// seekTarget is in seconds, samplesWritten is in frames (samples per channel)
	p.samplesWritten = int64(seekTarget * float64(sampleRate))
	p.queueMu.Unlock()

	// 8. Resume if we were playing
	if p.playing {
		p.otoPlayer.Play()
	}
}

// InstantSolo solos a channel (unmutes it, mutes others) with Flush & Seek
func (p *Player) InstantSolo(channel int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.otoPlayer == nil || p.module == nil {
		return
	}

	// 1. Get current positions
	renderPos := p.module.GetPositionSeconds()

	// 2. Calculate latency and seek target
	unplayedBytes := p.otoPlayer.UnplayedBufferSize()
	bytesPerSec := float64(sampleRate * 4)
	bufferedSecs := float64(unplayedBytes) / bytesPerSec

	seekTarget := renderPos - bufferedSecs
	if seekTarget < 0 {
		seekTarget = 0
	}

	// 3. Log for debugging (Removed)
	// fmt.Printf("[SOLO] Render: %.3f, Buf: %.3f, Target: %.3f\n", renderPos, bufferedSecs, seekTarget)

	// 4. Solo Channel
	p.module.SoloChannel(channel)

	// 5. Seek Module
	p.module.SetPositionSeconds(seekTarget)

	// 6. Flush Oto
	p.otoPlayer.Reset()

	// 7. Reset sync
	p.queueMu.Lock()
	p.stateQueue = p.stateQueue[:0]
	p.samplesWritten = int64(seekTarget * float64(sampleRate))
	p.queueMu.Unlock()

	// 8. Resume
	if p.playing {
		p.otoPlayer.Play()
	}
}
