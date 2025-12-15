Audio Engine Architectural Handoff
Project Status
We are building GoMod, a tracker music player (MOD/XM/S3M) using libopenmpt (CGo) and oto/v2 (Audio). The core playback is working, but we have a critical UX issue with Channel Mute/Solo Latency.

The Problem: "Mute Delay"
When the user mutes a channel (Keys 1-9), the change happens in libopenmpt's internal state immediately, but the user doesn't hear it for ~200-300ms.

Cause: libopenmpt renders audio into a buffer. oto buffers this before sending to the OS. The OS buffers it again. The "Unmuted" audio is already in the pipeline when the user presses Mute.
What We Tried & Why It Failed
Attempt 1: "Flush & Seek" (The theoretically correct fix)
Logic:

Calculate UnplayedBufferSize (how much audio is queued).
Seek the module BACK by that amount.
Flush the buffer.
Resume rendering (re-rendering the lost audio with Mute applied).
Result:

Audio Skipping: The user described "skipping every other half of the song".
Diagnosis: The seek calculation currentRenderPos - bufferedSeconds was likely incorrect or UnplayedBufferSize is unreliable on macOS, causing us to seek to the wrong spot (too far forward), effectively deleting time.
Artifacts: Re-triggering envelopes (pads restarting) is a risk with this method even if timing is perfect.
Attempt 2: "Low Latency Buffer" (The current state)
Logic:

Reduce oto buffer to almost nothing (tried 1024 samples / 23ms).
Use NewContextWithOptions to force OS buffer to ~60ms.
Result:

Still Delayed: User reports "unmuting right before a pattern i start hearing it after the new one starts". The total rendering pipeline latency is still too high for "musical" timing.
Recommended Path for New Agent
You need to implement Instant Mute without Skipping.

Option A: Fix the "Flush & Seek" Math
If you can get UnplayedBufferSize to be accurate, this is the cleanest path.

Debug: You need to verify what oto reports vs real time.
Math: SeekPos = RenderPos - (BytesInOto / BytesPerSec) - (InternalRenderBuffer / BytesPerSec).
Note: This will ALWAYs cause "Pad Retrigger" artifacts because libopenmpt resets envelopes on seek. If the user hates that, Option A is dead.
Option B: The "Ghost" Buffer (Recommended)
Don't use libopenmpt's mute. Implement Post-Render Muting.

Bitmask: Keep a [128]bool mute array in Go.
Interceptor: In the Audio Read loop (audioReader.Read), after getting samples from CGo, manually zero out the muted channels?
Problem: libopenmpt returns mixed stereo. We can't un-mix Channel 1 to mute it.
The only way Option B works: You must configure libopenmpt to render Quad/Surround/Multi-channel if possible, or intercept the channel volume before mixing? (Not possible in standard API).
Option C: Parallel Rendering (Hard but Perfect)
If Option A fails (artifacts) and B is impossible (pre-mixed):

You might need to accept the delay OR implement a "Lookahead" UI.
Actually: The user wants the SOUND to stop.
Hypothesis: The delay might be in the libopenmpt internal mixing buffer. Check openmpt_module_set_render_param for RENDER_PARAM_MIXING_INTERPOLATION or buffer config.
Current Code State

internal/player/player.go
: Contains the NewContextWithOptions logic and the "Hardware Sync" (locking UI time to unplayed buffer).

internal/player/openmpt.go
: Contains the CGo bindings for set_channel_mute_status.
Your Goal
Make the Mute button feel INSTANT (~20ms response), without skipping beats or glitching the song position.