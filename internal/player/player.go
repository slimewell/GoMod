package player

import (
	"context"
	"sync"

	"github.com/hajimehoshi/oto/v2"
)

const (
	sampleRate   = 44100
	channelCount = 2 // stereo
	bufferSize   = 2048
)

// Player manages audio playback
type Player struct {
	module     *Module
	otoContext *oto.Context
	otoPlayer  oto.Player
	mu         sync.RWMutex
	playing    bool
}

// NewPlayer creates a new player for the given module
func NewPlayer(module *Module) (*Player, error) {
	// Initialize Oto with small buffer for low latency
	otoContext, ready, err := oto.NewContext(sampleRate, channelCount, 2)
	if err != nil {
		return nil, err
	}
	<-ready

	p := &Player{
		module:     module,
		otoContext: otoContext,
		playing:    false,
	}

	return p, nil
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
	})

	p.otoPlayer.Play()
	p.playing = true

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
}

func (r *audioReader) Read(p []byte) (int, error) {
	// Check context cancellation
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
	}

	// Render audio from openmpt
	frames := r.module.Read(r.buf[:])
	if frames == 0 {
		return 0, nil // End of module
	}

	// Convert int16 samples to bytes
	samples := frames * 2 // stereo
	bytesWritten := samples * 2

	for i := 0; i < samples; i++ {
		p[i*2] = byte(r.buf[i] & 0xff)
		p[i*2+1] = byte((r.buf[i] >> 8) & 0xff)
	}

	return bytesWritten, nil
}
