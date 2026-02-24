// ─── Web Audio API — sonidos UI + música de fondo estilo FIFA ────────────────

let _ctx: AudioContext | null = null

function ctx(): AudioContext {
  if (!_ctx || _ctx.state === 'closed') {
    _ctx = new (window.AudioContext || (window as any).webkitAudioContext)()
  }
  if (_ctx.state === 'suspended') _ctx.resume()
  return _ctx
}

// ── Utilidad: tono simple ─────────────────────────────────────────────────────
function tone(
  freq: number, endFreq: number, duration: number, volume: number,
  type: OscillatorType = 'sine', startOffset = 0, dest?: AudioNode
) {
  const ac = ctx()
  const t0 = ac.currentTime + startOffset
  const osc  = ac.createOscillator()
  const gain = ac.createGain()
  osc.connect(gain)
  gain.connect(dest ?? ac.destination)
  osc.type = type
  osc.frequency.setValueAtTime(freq, t0)
  if (endFreq !== freq) osc.frequency.exponentialRampToValueAtTime(endFreq, t0 + duration)
  gain.gain.setValueAtTime(volume, t0)
  gain.gain.exponentialRampToValueAtTime(0.0001, t0 + duration)
  osc.start(t0)
  osc.stop(t0 + duration + 0.02)
}

// ══════════════════════════════════════════════════════════════════════════════
//  SONIDOS DE UI
// ══════════════════════════════════════════════════════════════════════════════

export function playNav()         { try { tone(680, 480, 0.038, 0.13) } catch {} }
export function playSelect()      { try { tone(523,523,0.08,0.20,'sine',0); tone(659,659,0.07,0.16,'sine',0.02); tone(784,784,0.06,0.12,'sine',0.04) } catch {} }
export function playWatchlistOn() { try { tone(523,523,0.11,0.26,'sine',0); tone(659,659,0.11,0.26,'sine',0.09); tone(784,784,0.15,0.28,'sine',0.18) } catch {} }
export function playWatchlistOff(){ try { tone(600,380,0.13,0.20,'sine',0); tone(400,280,0.10,0.14,'sine',0.07) } catch {} }

// ══════════════════════════════════════════════════════════════════════════════
//  JINGLE EA SPORTS (intro)
// ══════════════════════════════════════════════════════════════════════════════

export function playEAIntro() {
  try {
    const ac  = ctx()
    const t0  = ac.currentTime + 0.2
    const dly = ac.createDelay(0.6)
    const fb  = ac.createGain()
    const wet = ac.createGain()
    const mst = ac.createGain()
    dly.delayTime.value = 0.22; fb.gain.value = 0.32; wet.gain.value = 0.28; mst.gain.value = 0.9
    dly.connect(fb); fb.connect(dly); dly.connect(wet); wet.connect(mst); mst.connect(ac.destination)

    const note = (freq: number, start: number, dur: number, vol: number) => {
      const osc = ac.createOscillator(); const g = ac.createGain()
      osc.connect(g); g.connect(mst); g.connect(dly)
      osc.type = 'sawtooth'; osc.frequency.value = freq
      g.gain.setValueAtTime(vol, t0 + start); g.gain.exponentialRampToValueAtTime(0.0001, t0 + start + dur)
      osc.start(t0 + start); osc.stop(t0 + start + dur + 0.05)
    }
    note(293.66, 0.00, 0.11, 0.20); note(369.99, 0.11, 0.11, 0.20)
    note(440.00, 0.22, 0.11, 0.22); note(587.33, 0.33, 0.13, 0.24)
    note(293.66, 0.60, 0.60, 0.26); note(369.99, 0.60, 0.60, 0.26)
    note(440.00, 0.60, 0.60, 0.26); note(587.33, 0.60, 0.60, 0.28); note(880.00, 0.60, 0.55, 0.14)
    mst.gain.setValueAtTime(0.9, ac.currentTime + t0 + 0.6)
    mst.gain.linearRampToValueAtTime(0, ac.currentTime + t0 + 1.9)
  } catch {}
}

// ══════════════════════════════════════════════════════════════════════════════
//  MÚSICA DE FONDO — fifa.mp3 (archivo real del backend)
// ══════════════════════════════════════════════════════════════════════════════

let bgRunning = false
let _audio: HTMLAudioElement | null = null
let _fadeTimer: ReturnType<typeof setInterval> | null = null

export function playBGMusic() {
  if (bgRunning) return
  if (!_audio) {
    _audio = new Audio('/sound/fifa.mp3')
    _audio.loop = true
    _audio.volume = 0
  }
  const audio = _audio
  audio.currentTime = 0
  audio.volume = 0

  audio.play().then(() => {
    bgRunning = true
    if (_fadeTimer) clearInterval(_fadeTimer)
    let vol = 0
    _fadeTimer = setInterval(() => {
      vol = Math.min(vol + 0.008, 0.12)
      audio.volume = vol
      if (vol >= 0.12) { clearInterval(_fadeTimer!); _fadeTimer = null }
    }, 80)
  }).catch(() => {
    bgRunning = false
  })
}

export function stopBGMusic() {
  if (!bgRunning) return
  bgRunning = false
  if (_fadeTimer) { clearInterval(_fadeTimer); _fadeTimer = null }
  if (!_audio) return

  const audio = _audio
  let vol = audio.volume
  _fadeTimer = setInterval(() => {
    vol = Math.max(vol - 0.03, 0)
    audio.volume = vol
    if (vol <= 0) {
      audio.pause()
      clearInterval(_fadeTimer!); _fadeTimer = null
    }
  }, 60)
}

export function isBGMusicPlaying() { return bgRunning }
