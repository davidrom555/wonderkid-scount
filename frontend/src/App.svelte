<script lang="ts">
  import { onMount } from 'svelte'
  import PlayerCard from './components/PlayerCard.svelte'
  import { playBGMusic, stopBGMusic, isBGMusicPlaying } from './lib/sounds'
  import {
    players, watchlistIds, filteredPlayers, calcROI,
    searchTerm, minPotential, selectedPosition, showPredictions, activeTab,
    getPhotoUrl, type Player
  } from './stores/players'

  let isLoading = true
  let apiError = ''
  let musicOn = false
  let showFilters = false
  let menuOpen = false

  function toggleMusic() {
    if (musicOn) {
      stopBGMusic()
      musicOn = false
    } else {
      playBGMusic()
      musicOn = true
    }
  }
  const positions = ['ST', 'CAM', 'CB', 'CM', 'RW', 'LW', 'GK', 'LB', 'CDM', 'RB']

  onMount(async () => {
    // Detectar navegación por teclado vs ratón para desactivar hover en cards
    const onKey   = () => document.body.classList.add('kb-nav')
    const onMouse = () => document.body.classList.remove('kb-nav')
    window.addEventListener('keydown', onKey)
    window.addEventListener('mousemove', onMouse)

    await Promise.all([fetchPlayers(), fetchWatchlist()])
    isLoading = false
  })

  async function fetchPlayers() {
    try {
      const res = await fetch('/api/players')
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const data: Player[] = await res.json()
      players.set(data)
    } catch (e) {
      apiError = 'Cannot connect to backend on port 8080. Start the Go server first.'
      isLoading = false
    }
  }

  async function fetchWatchlist() {
    try {
      const res = await fetch('/api/watchlist')
      if (!res.ok) return
      const data = await res.json()
      watchlistIds.set(new Set((data.players as Player[]).map(p => p.id)))
    } catch {}
  }

  $: displayPlayers = $activeTab === 'watchlist'
    ? $filteredPlayers.filter(p => $watchlistIds.has(p.id))
    : $filteredPlayers

  $: topProspects = [...$filteredPlayers].sort((a, b) => b.potential - a.potential).slice(0, 5)
  $: undervalued = [...$filteredPlayers].filter(p => calcROI(p) > 1.2).sort((a, b) => calcROI(b) - calcROI(a)).slice(0, 5)
  $: avgPotential = $players.length ? ($players.reduce((s, p) => s + p.potential, 0) / $players.length).toFixed(1) : '—'
</script>

<div class="app">

  <!-- ══════════════════════════════════════════
       HEADER — EA FC 26 Gradient
       ══════════════════════════════════════════ -->
  <header class="fc-header">
    <div class="header-inner">

      <!-- ── Logo izquierda ── -->
      <div class="logo-block">
        <h1 class="logo-title">
          WONDERKID<span class="logo-green"> SCOUT</span>
        </h1>
        <p class="logo-sub">EA FC 26 &nbsp;·&nbsp; SCOUTING PLATFORM</p>
      </div>

      <!-- ── Centro: EA FC 26 brand ── -->
      <div class="brand-center">
        <div class="brand-mark">
          <span class="bm-ea">EA</span>
          <span class="bm-div"></span>
          <span class="bm-fc">FC<em>26</em></span>
        </div>
      </div>

      <!-- ── Hamburguesa (solo mobile) ── -->
      <button
        class="hamburger"
        class:ham-open={menuOpen}
        on:click={() => menuOpen = !menuOpen}
        aria-label="Menú"
      >
        <span></span><span></span><span></span>
      </button>

      <!-- ── Nav derecha ── -->
      <div class="nav-right">
        <nav class="nav-tabs">
          <button
            class="nav-btn"
            class:nav-active={$activeTab === 'dashboard'}
            on:click={() => activeTab.set('dashboard')}
          >
            <svg class="tab-ico" viewBox="0 0 14 14" fill="none">
              <rect x="1" y="1" width="5" height="5" rx="1" fill="currentColor"/>
              <rect x="8" y="1" width="5" height="5" rx="1" fill="currentColor"/>
              <rect x="1" y="8" width="5" height="5" rx="1" fill="currentColor"/>
              <rect x="8" y="8" width="5" height="5" rx="1" fill="currentColor"/>
            </svg>
            Dashboard
          </button>
          <button
            class="nav-btn"
            class:nav-active={$activeTab === 'watchlist'}
            class:nav-gold={$activeTab === 'watchlist'}
            on:click={() => activeTab.set('watchlist')}
          >
            <svg class="tab-ico" viewBox="0 0 14 14" fill="none">
              <path d="M7 1l1.8 3.6L13 5.3l-3 2.9.7 4.1L7 10.4l-3.7 1.9.7-4.1-3-2.9 4.2-.7z" fill="currentColor"/>
            </svg>
            Watchlist
            {#if $watchlistIds.size > 0}
              <span class="nav-badge">{$watchlistIds.size}</span>
            {/if}
          </button>
        </nav>

        <div class="nav-sep"></div>

        <!-- Sound pill -->
        <button
          class="sound-btn"
          class:sound-on={musicOn}
          on:click={toggleMusic}
          title={musicOn ? 'Desactivar música' : 'Activar música'}
        >
          <svg class="sound-ico" viewBox="0 0 16 16" fill="none">
            {#if musicOn}
              <path d="M2 5.5h2.5L8 2v12l-3.5-3.5H2V5.5z" fill="currentColor"/>
              <path d="M10 5a4 4 0 0 1 0 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              <path d="M12 3a7 7 0 0 1 0 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            {:else}
              <path d="M2 5.5h2.5L8 2v12l-3.5-3.5H2V5.5z" fill="currentColor" opacity="0.4"/>
              <path d="M11 6l3 4M14 6l-3 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            {/if}
          </svg>
          {#if musicOn}
            <span class="eq">
              <span class="bar b1"></span>
              <span class="bar b2"></span>
              <span class="bar b3"></span>
              <span class="bar b4"></span>
            </span>
          {/if}
        </button>
      </div>

    </div>
  </header>

  <!-- ══════════════════════════════════════════
       MENÚ MÓVIL
       ══════════════════════════════════════════ -->
  {#if menuOpen}
    <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
    <div class="mm-backdrop" on:click={() => menuOpen = false}></div>
    <nav class="mobile-menu">
      <button
        class="mm-item"
        class:mm-active={$activeTab === 'dashboard'}
        on:click={() => { activeTab.set('dashboard'); menuOpen = false }}
      >
        <svg viewBox="0 0 14 14" fill="none" width="16" height="16">
          <rect x="1" y="1" width="5" height="5" rx="1" fill="currentColor"/>
          <rect x="8" y="1" width="5" height="5" rx="1" fill="currentColor"/>
          <rect x="1" y="8" width="5" height="5" rx="1" fill="currentColor"/>
          <rect x="8" y="8" width="5" height="5" rx="1" fill="currentColor"/>
        </svg>
        Dashboard
      </button>
      <button
        class="mm-item"
        class:mm-active={$activeTab === 'watchlist'}
        on:click={() => { activeTab.set('watchlist'); menuOpen = false }}
      >
        <svg viewBox="0 0 14 14" fill="none" width="16" height="16">
          <path d="M7 1l1.8 3.6L13 5.3l-3 2.9.7 4.1L7 10.4l-3.7 1.9.7-4.1-3-2.9 4.2-.7z" fill="currentColor"/>
        </svg>
        Watchlist
        {#if $watchlistIds.size > 0}
          <span class="mm-badge">{$watchlistIds.size}</span>
        {/if}
      </button>
      <div class="mm-sep"></div>
      <button class="mm-item mm-sound" on:click={() => { toggleMusic(); menuOpen = false }}>
        {#if musicOn}
          <svg viewBox="0 0 16 16" fill="none" width="16" height="16">
            <path d="M2 5.5h2.5L8 2v12l-3.5-3.5H2V5.5z" fill="currentColor"/>
            <path d="M10 5a4 4 0 0 1 0 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
          Música ON
        {:else}
          <svg viewBox="0 0 16 16" fill="none" width="16" height="16">
            <path d="M2 5.5h2.5L8 2v12l-3.5-3.5H2V5.5z" fill="currentColor" opacity="0.4"/>
            <path d="M11 6l3 4M14 6l-3 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
          Música OFF
        {/if}
      </button>
    </nav>
  {/if}

  <!-- ══════════════════════════════════════════
       LOADING
       ══════════════════════════════════════════ -->
  {#if isLoading}
    <div class="full-center">
      <div class="spinner"></div>
      <p class="state-lbl">Scouting Players...</p>
    </div>

  <!-- ══════════════════════════════════════════
       ERROR
       ══════════════════════════════════════════ -->
  {:else if apiError}
    <div class="full-center">
      <div class="error-box">
        <div class="error-ico">⚠️</div>
        <h3 class="error-title">Connection Error</h3>
        <p class="error-msg">{apiError}</p>
        <button on:click={fetchPlayers} class="retry-btn">Retry</button>
      </div>
    </div>

  <!-- ══════════════════════════════════════════
       MAIN — SIDEBAR + CONTENT (FC26 Layout)
       ══════════════════════════════════════════ -->
  {:else}
    <div class="page-layout">

      <!-- ────────── SIDEBAR ────────── -->
      <aside class="sidebar" class:filters-open={showFilters}>

        <!-- Quick stats 4-col -->
        <div class="stats-row">
          <div class="qstat">
            <span class="qval">{$players.length}</span>
            <span class="qlbl">Players</span>
          </div>
          <div class="qstat">
            <span class="qval gold-t">{avgPotential}</span>
            <span class="qlbl">Avg Pot</span>
          </div>
          <div class="qstat">
            <span class="qval orange-t">{$watchlistIds.size}</span>
            <span class="qlbl">Watch</span>
          </div>
          <div class="qstat">
            <span class="qval green-t">{undervalued.length}</span>
            <span class="qlbl">ROI</span>
          </div>
        </div>

        <div class="sep-line"></div>

        <!-- ── Filters heading ── -->
        <div class="section-hd">
          <span class="hd-bar gold-bar"></span>
          Scout Filters
          <span class="count-chip">{displayPlayers.length}</span>
        </div>

        <!-- Search -->
        <div class="filter-g">
          <label class="g-lbl" for="s-inp">Search</label>
          <div class="inp-wrap">
            <span class="inp-ico">⚲</span>
            <input
              id="s-inp"
              type="text"
              bind:value={$searchTerm}
              placeholder="Bellingham, Real Madrid..."
              class="fc-inp"
            />
          </div>
        </div>

        <!-- Min Potential -->
        <div class="filter-g">
          <label class="g-lbl" for="pot-rng">
            Min Potential
            <span class="g-val">{$minPotential}</span>
          </label>
          <input
            id="pot-rng"
            type="range"
            bind:value={$minPotential}
            min={50} max={99}
            class="fc-range"
          />
          <div class="range-ends"><span>50</span><span>99</span></div>
        </div>

        <!-- Position picker -->
        <div class="filter-g">
          <label class="g-lbl">Position</label>
          <div class="pos-grid">
            <button
              class="pos-btn"
              class:pos-on={$selectedPosition === ''}
              on:click={() => selectedPosition.set('')}
            >ALL</button>
            {#each positions as pos}
              <button
                class="pos-btn"
                class:pos-on={$selectedPosition === pos}
                on:click={() => selectedPosition.set(pos)}
              >{pos}</button>
            {/each}
          </div>
        </div>

        <!-- AI Predictions -->
        <button
          class="pred-btn"
          class:pred-on={$showPredictions}
          on:click={() => showPredictions.update(v => !v)}
        >
          {$showPredictions ? '🔮 AI Predictions ON' : '📊 Show AI Predictions'}
        </button>

        <div class="sep-line"></div>

        <!-- Top Prospects -->
        <div class="section-hd sm-hd">
          <span class="hd-bar gold-bar"></span>
          ★ Top Prospects
        </div>
        <div class="mini-list">
          {#each topProspects as p, i}
            <div class="mini-row">
              <span class="mini-rank">{i + 1}</span>
              <div class="mini-info">
                <div class="mini-name">{p.name.split(' ').pop() ?? p.name}</div>
                <div class="mini-club">{p.club}</div>
              </div>
              <span class="mini-val gold-t">{p.potential}</span>
            </div>
          {/each}
        </div>

        <!-- Undervalued -->
        <div class="section-hd sm-hd" style="margin-top:8px">
          <span class="hd-bar orange-bar"></span>
          🔥 Undervalued
        </div>
        <div class="mini-list">
          {#each undervalued.slice(0, 4) as p}
            <div class="mini-row">
              <div class="mini-info">
                <div class="mini-name">{p.name.split(' ').pop() ?? p.name}</div>
                <div class="mini-club">{p.club}</div>
              </div>
              <span class="mini-val orange-t">{calcROI(p).toFixed(1)}x</span>
            </div>
          {:else}
            <div class="mini-empty">No undervalued targets</div>
          {/each}
        </div>

      </aside>

      <!-- ────────── MAIN CONTENT ────────── -->
      <main class="main-area">

        <div class="main-top">
          <button class="mobile-filter-btn" on:click={() => showFilters = !showFilters}>
            <svg viewBox="0 0 16 16" fill="none" width="14" height="14">
              <path d="M2 4h12M4 8h8M6 12h4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
            {showFilters ? 'Cerrar' : 'Filtros'}
          </button>
          <span class="result-lbl">{displayPlayers.length} Players</span>
        </div>

        {#if displayPlayers.length > 0}
          <div class="card-grid">
            {#each displayPlayers as player (player.id)}
              <PlayerCard {player} />
            {/each}
          </div>
        {:else}
          <div class="empty-state">
            <div class="empty-ico">⚽</div>
            <h3 class="empty-title">No Players Found</h3>
            <p class="empty-sub">
              {$activeTab === 'watchlist' ? 'Your watchlist is empty' : 'Try adjusting your filters'}
            </p>
          </div>
        {/if}

      </main>
    </div>
  {/if}
</div>

<style>
  /* Reserva siempre el espacio del scrollbar para evitar layout shift */
  :global(html) { scrollbar-gutter: stable; }

  /* ════════════════════════════════════════════
     HEADER — EA FC 26 Dark Glass
     ════════════════════════════════════════════ */
  .fc-header {
    position: fixed;
    top: 0; left: 0; right: 0;
    z-index: 50;
    height: 64px;
    background: rgba(4, 6, 18, 0.97);
    backdrop-filter: blur(28px);
    -webkit-backdrop-filter: blur(28px);
    border-bottom: 1px solid rgba(255,255,255,0.07);
    box-shadow: 0 8px 32px rgba(0,0,0,0.7);
  }
  /* Línea degradada en la parte inferior — sello EA FC */
  .fc-header::after {
    content: '';
    position: absolute;
    bottom: 0; left: 0; right: 0;
    height: 2px;
    background: linear-gradient(90deg,
      #00f080 0%, #00c8d8 22%, #4820e0 48%, #9018c8 72%, #e030c8 100%
    );
  }

  .header-inner {
    max-width: 1800px;
    margin: 0 auto;
    padding: 0 28px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 24px;
    position: relative;
  }

  /* ── Logo ── */
  .logo-block { flex-shrink: 0; }
  .logo-title {
    font-family: 'Rajdhani', sans-serif;
    font-size: 1.5rem; font-weight: 900;
    font-style: italic; text-transform: uppercase;
    letter-spacing: 0.04em; line-height: 1;
    color: #ffffff;
  }
  .logo-green { color: #00f080; }
  .logo-sub {
    font-size: 0.44rem; font-weight: 700;
    text-transform: uppercase; letter-spacing: 0.26em;
    color: rgba(255,255,255,0.22); margin-top: 3px;
  }

  /* ── Brand center EA FC 26 — siempre centrado respecto al header ── */
  .brand-center {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
    pointer-events: none;
  }
  .brand-mark {
    display: flex; align-items: stretch;
    background: rgba(255,255,255,0.05);
    border: 1px solid rgba(255,255,255,0.1);
    border-radius: 9px; overflow: hidden;
    box-shadow: inset 0 1px 0 rgba(255,255,255,0.08);
  }
  .bm-ea {
    padding: 5px 11px;
    font-size: 0.7rem; font-weight: 900; letter-spacing: 0.1em;
    color: rgba(255,255,255,0.4);
    display: flex; align-items: center;
    background: rgba(255,255,255,0.03);
  }
  .bm-div { width: 1px; background: rgba(255,255,255,0.09); flex-shrink: 0; }
  .bm-fc {
    padding: 3px 14px;
    font-family: 'Rajdhani', sans-serif;
    font-size: 1.6rem; font-weight: 900;
    color: #ffffff; letter-spacing: -0.02em; line-height: 1;
    display: flex; align-items: center;
  }
  .bm-fc em { color: #00f080; font-style: italic; margin-left: 1px; }

  /* ── Nav right ── */
  .nav-right {
    display: flex; align-items: center; gap: 8px; flex-shrink: 0;
  }
  .nav-sep {
    width: 1px; height: 24px;
    background: rgba(255,255,255,0.1);
  }

  /* Tabs — línea indicadora inferior */
  .nav-tabs { display: flex; align-items: stretch; height: 64px; gap: 0; }
  .nav-btn {
    position: relative;
    display: flex; align-items: center; gap: 7px;
    padding: 0 20px;
    font-family: 'Rajdhani', sans-serif;
    font-weight: 900; font-style: italic;
    font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.1em;
    cursor: pointer;
    border: none; background: transparent;
    color: rgba(255,255,255,0.38);
    transition: color 0.18s, background 0.18s;
    white-space: nowrap;
  }
  .nav-btn::after {
    content: '';
    position: absolute;
    bottom: 0; left: 14px; right: 14px;
    height: 2px; border-radius: 2px 2px 0 0;
    background: #00f080;
    opacity: 0; transition: opacity 0.18s;
  }
  .nav-btn:hover { color: rgba(255,255,255,0.72); background: rgba(255,255,255,0.04); }
  .nav-btn:hover::after { opacity: 0.35; }

  .nav-active { color: #ffffff !important; background: rgba(255,255,255,0.06) !important; }
  .nav-active::after { opacity: 1 !important; }
  .nav-gold { color: #f5c818 !important; }
  .nav-gold::after { background: #f5c818 !important; }

  .tab-ico { width: 13px; height: 13px; flex-shrink: 0; }

  .nav-badge {
    background: rgba(0,240,128,0.18);
    border: 1px solid rgba(0,240,128,0.35);
    color: #00f080;
    padding: 1px 7px; border-radius: 20px;
    font-size: 0.56rem; font-style: normal; font-weight: 900;
  }
  .nav-gold .nav-badge {
    background: rgba(245,200,24,0.15);
    border-color: rgba(245,200,24,0.35);
    color: #f5c818;
  }

  /* ── Sound button ── */
  .sound-btn {
    display: flex; align-items: center; gap: 5px;
    padding: 0 12px; height: 34px;
    border-radius: 17px;
    cursor: pointer;
    border: 1px solid rgba(255,255,255,0.12);
    background: rgba(255,255,255,0.05);
    color: rgba(255,255,255,0.35);
    transition: all 0.2s;
  }
  .sound-btn:hover {
    background: rgba(255,255,255,0.1);
    color: rgba(255,255,255,0.7);
    border-color: rgba(255,255,255,0.22);
  }
  .sound-on {
    background: rgba(0,240,128,0.1) !important;
    border-color: rgba(0,240,128,0.4) !important;
    color: #00f080 !important;
    box-shadow: 0 0 12px rgba(0,240,128,0.2);
  }
  .sound-ico { width: 15px; height: 15px; flex-shrink: 0; }

  /* Barras ecualizador */
  .eq { display: flex; align-items: flex-end; gap: 2px; height: 14px; flex-shrink: 0; }
  .bar { width: 2.5px; border-radius: 1px; background: #00f080; }
  .b1 { animation: eq 0.7s ease-in-out infinite alternate; }
  .b2 { animation: eq 0.7s ease-in-out 0.15s infinite alternate; }
  .b3 { animation: eq 0.7s ease-in-out 0.3s  infinite alternate; }
  .b4 { animation: eq 0.7s ease-in-out 0.45s infinite alternate; }
  @keyframes eq {
    from { height: 3px; }
    to   { height: 13px; }
  }

  /* ════════════════════════════════════════════
     LOADING / ERROR
     ════════════════════════════════════════════ */
  .full-center {
    display: flex; flex-direction: column;
    align-items: center; justify-content: center;
    min-height: 100vh; gap: 16px; padding-top: 64px;
  }
  .spinner {
    width: 48px; height: 48px;
    border: 4px solid rgba(245,200,24,0.2);
    border-top-color: #f5c818;
    border-radius: 50%;
    animation: spin 0.9s linear infinite;
  }
  .state-lbl {
    font-size: 0.62rem; font-weight: 700;
    text-transform: uppercase; letter-spacing: 0.28em;
    color: rgba(255,255,255,0.3);
  }
  .error-box {
    background: rgba(239,68,68,0.06);
    border: 1px solid rgba(239,68,68,0.2);
    border-radius: 14px; padding: 36px;
    max-width: 400px; text-align: center;
  }
  .error-ico { font-size: 3rem; margin-bottom: 12px; }
  .error-title {
    font-family: 'Rajdhani', sans-serif;
    font-size: 1.2rem; font-weight: 900; font-style: italic;
    text-transform: uppercase; color: #f87171; margin-bottom: 8px;
  }
  .error-msg { color: rgba(255,255,255,0.45); font-size: 0.82rem; margin-bottom: 16px; }
  .retry-btn {
    padding: 8px 20px;
    background: rgba(239,68,68,0.12);
    border: 1px solid rgba(239,68,68,0.28);
    border-radius: 7px; color: #f87171;
    font-weight: 900; font-size: 0.72rem;
    text-transform: uppercase; letter-spacing: 0.1em;
    cursor: pointer; transition: background 0.2s;
  }
  .retry-btn:hover { background: rgba(239,68,68,0.22); }

  /* ════════════════════════════════════════════
     PAGE LAYOUT — SIDEBAR + MAIN
     ════════════════════════════════════════════ */
  .page-layout {
    display: flex;
    padding-top: 64px;
    min-height: 100vh;
  }

  /* ════════════════════════════════════════════
     SIDEBAR
     ════════════════════════════════════════════ */
  .sidebar {
    width: 265px;
    flex-shrink: 0;
    position: sticky;
    top: 64px;
    height: calc(100vh - 64px);
    overflow-y: auto;
    background: rgba(255,255,255,0.022);
    border-right: 1px solid rgba(255,255,255,0.06);
    padding: 14px 13px;
    display: flex;
    flex-direction: column;
    gap: 11px;
  }
  .sidebar::-webkit-scrollbar { width: 3px; }
  .sidebar::-webkit-scrollbar-track { background: transparent; }
  .sidebar::-webkit-scrollbar-thumb { background: rgba(245,200,24,0.18); border-radius: 2px; }

  /* Quick stats */
  .stats-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 5px;
  }
  .qstat {
    background: rgba(255,255,255,0.04);
    border: 1px solid rgba(255,255,255,0.06);
    border-radius: 7px;
    padding: 7px 3px;
    display: flex; flex-direction: column;
    align-items: center; gap: 2px;
  }
  .qval {
    font-family: 'Rajdhani', sans-serif;
    font-size: 1rem; font-weight: 900;
    font-style: italic; color: white; line-height: 1;
  }
  .qlbl {
    font-size: 0.37rem; font-weight: 700;
    text-transform: uppercase; letter-spacing: 0.1em;
    color: rgba(255,255,255,0.28);
  }

  .sep-line { height: 1px; background: rgba(255,255,255,0.06); }

  /* Section heading */
  .section-hd {
    display: flex; align-items: center; gap: 6px;
    font-family: 'Rajdhani', sans-serif;
    font-size: 0.6rem; font-weight: 900;
    font-style: italic; text-transform: uppercase;
    letter-spacing: 0.16em; color: rgba(255,255,255,0.44);
  }
  .sm-hd { font-size: 0.54rem; letter-spacing: 0.12em; }
  .hd-bar { width: 3px; height: 11px; border-radius: 1px; flex-shrink: 0; }
  .gold-bar   { background: #f5c818; }
  .orange-bar { background: #fb923c; }
  .count-chip {
    margin-left: auto;
    font-size: 0.5rem; color: #f5c818;
    background: rgba(245,200,24,0.1);
    border: 1px solid rgba(245,200,24,0.22);
    padding: 1px 7px; border-radius: 20px; font-style: normal;
  }

  /* Filter group */
  .filter-g { display: flex; flex-direction: column; gap: 5px; }
  .g-lbl {
    display: flex; align-items: center; justify-content: space-between;
    font-size: 0.5rem; font-weight: 700;
    text-transform: uppercase; letter-spacing: 0.15em;
    color: rgba(255,255,255,0.3);
  }
  .g-val {
    font-family: 'Rajdhani', sans-serif;
    font-size: 0.74rem; font-weight: 900;
    font-style: italic; color: #f5c818;
  }

  /* Input */
  .inp-wrap { position: relative; display: flex; align-items: center; }
  .inp-ico {
    position: absolute; left: 9px;
    font-size: 0.85rem; opacity: 0.28;
    pointer-events: none; color: white; line-height: 1;
  }
  .fc-inp {
    width: 100%;
    padding: 8px 10px 8px 28px;
    background: rgba(255,255,255,0.05);
    border: 1px solid rgba(255,255,255,0.09);
    border-radius: 7px; color: white;
    font-size: 0.72rem; font-weight: 600;
    font-family: 'Inter', sans-serif;
    outline: none; transition: border-color 0.2s, background 0.2s;
  }
  .fc-inp::placeholder { color: rgba(255,255,255,0.18); font-weight: 400; }
  .fc-inp:focus {
    border-color: rgba(245,200,24,0.42);
    background: rgba(245,200,24,0.04);
    box-shadow: 0 0 0 2px rgba(245,200,24,0.07);
  }

  /* Range slider */
  .fc-range {
    width: 100%; height: 3px;
    appearance: none; -webkit-appearance: none;
    background: rgba(255,255,255,0.1);
    border-radius: 2px; outline: none;
    cursor: pointer; accent-color: #f5c818;
    margin-top: 5px;
  }
  .fc-range::-webkit-slider-thumb {
    appearance: none; -webkit-appearance: none;
    width: 15px; height: 15px; border-radius: 50%;
    background: #f5c818;
    box-shadow: 0 0 7px rgba(245,200,24,0.55);
    cursor: pointer; border: 2px solid #1a1000;
  }
  .fc-range::-moz-range-thumb {
    width: 13px; height: 13px; border-radius: 50%;
    background: #f5c818; border: 2px solid #1a1000; cursor: pointer;
  }
  .range-ends {
    display: flex; justify-content: space-between;
    font-size: 0.44rem; font-weight: 700;
    color: rgba(255,255,255,0.2); margin-top: 2px;
  }

  /* Position picker */
  .pos-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 4px; }
  .pos-btn {
    padding: 5px 2px;
    background: rgba(255,255,255,0.04);
    border: 1px solid rgba(255,255,255,0.07);
    border-radius: 5px; color: rgba(255,255,255,0.36);
    font-size: 0.44rem; font-weight: 900;
    text-transform: uppercase; letter-spacing: 0.02em;
    cursor: pointer; transition: all 0.14s;
    font-family: 'Rajdhani', sans-serif; font-style: italic;
    line-height: 1.3;
  }
  .pos-btn:hover {
    background: rgba(255,255,255,0.09);
    color: rgba(255,255,255,0.78);
    border-color: rgba(255,255,255,0.15);
  }
  .pos-on {
    background: rgba(245,200,24,0.12) !important;
    border-color: rgba(245,200,24,0.32) !important;
    color: #f5c818 !important;
  }

  /* AI Predictions toggle */
  .pred-btn {
    width: 100%; padding: 8px 12px; border-radius: 7px;
    font-weight: 700; font-size: 0.56rem;
    text-transform: uppercase; letter-spacing: 0.08em;
    cursor: pointer;
    border: 1px solid rgba(255,255,255,0.09);
    background: rgba(255,255,255,0.04);
    color: rgba(255,255,255,0.36);
    transition: all 0.2s;
    display: flex; align-items: center; justify-content: center; gap: 6px;
    font-family: 'Inter', sans-serif;
  }
  .pred-btn:hover { border-color: rgba(255,255,255,0.18); color: rgba(255,255,255,0.65); }
  .pred-on {
    background: rgba(245,200,24,0.09) !important;
    border-color: rgba(245,200,24,0.34) !important;
    color: #f5c818 !important;
  }

  /* Mini lists (prospects / undervalued) */
  .mini-list { display: flex; flex-direction: column; gap: 3px; }
  .mini-row {
    display: flex; align-items: center; gap: 5px;
    padding: 5px 6px; border-radius: 5px;
    border: 1px solid transparent; transition: all 0.14s;
  }
  .mini-row:hover {
    background: rgba(255,255,255,0.04);
    border-color: rgba(255,255,255,0.06);
  }
  .mini-rank {
    font-size: 0.5rem; font-weight: 900;
    color: rgba(255,255,255,0.14);
    width: 10px; flex-shrink: 0; text-align: right;
  }
  .mini-info { flex: 1; min-width: 0; }
  .mini-name {
    font-size: 0.65rem; font-weight: 700;
    color: rgba(255,255,255,0.78);
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  }
  .mini-club {
    font-size: 0.42rem; font-weight: 700;
    text-transform: uppercase; letter-spacing: 0.08em;
    color: rgba(255,255,255,0.25);
  }
  .mini-val {
    font-family: 'Rajdhani', sans-serif;
    font-size: 0.72rem; font-weight: 900;
    font-style: italic; flex-shrink: 0;
  }
  .mini-empty {
    font-size: 0.56rem; color: rgba(255,255,255,0.18);
    text-align: center; padding: 6px 0;
  }

  /* ════════════════════════════════════════════
     MAIN AREA
     ════════════════════════════════════════════ */
  .main-area { flex: 1; min-width: 0; padding: 20px; }

  .main-top {
    display: flex; align-items: center; justify-content: flex-end;
    margin-bottom: 18px;
  }
  .result-lbl {
    font-size: 0.54rem; font-weight: 700;
    text-transform: uppercase; letter-spacing: 0.18em;
    color: rgba(255,255,255,0.26);
  }

  /* Card grid */
  .card-grid { display: flex; flex-wrap: wrap; gap: 16px; }

  /* Empty state */
  .empty-state {
    display: flex; flex-direction: column;
    align-items: center; gap: 12px; padding: 80px 0;
  }
  .empty-ico { font-size: 3.5rem; }
  .empty-title {
    font-family: 'Rajdhani', sans-serif;
    font-size: 1.4rem; font-weight: 900;
    font-style: italic; text-transform: uppercase;
    letter-spacing: 0.1em; color: rgba(255,255,255,0.38);
  }
  .empty-sub {
    font-size: 0.64rem; font-weight: 700;
    text-transform: uppercase; letter-spacing: 0.2em;
    color: rgba(255,255,255,0.2);
  }

  /* ── Color utilities ── */
  .gold-t   { color: #f5c818; }
  .orange-t { color: #fb923c; }
  .green-t  { color: #4ade80; }

  /* ── Animations ── */
  @keyframes spin { to { transform: rotate(360deg); } }

  /* ════════════════════════════════════════════
     HAMBURGUESA (oculto en desktop)
     ════════════════════════════════════════════ */
  .hamburger {
    display: none;
    flex-direction: column; justify-content: center; align-items: center;
    gap: 5px;
    width: 38px; height: 38px;
    padding: 0;
    background: rgba(255,255,255,0.05);
    border: 1px solid rgba(255,255,255,0.1);
    border-radius: 8px;
    cursor: pointer;
    flex-shrink: 0;
    transition: background 0.18s, border-color 0.18s;
  }
  .hamburger:hover { background: rgba(255,255,255,0.1); border-color: rgba(255,255,255,0.2); }
  .hamburger span {
    display: block;
    width: 18px; height: 1.5px;
    background: rgba(255,255,255,0.7);
    border-radius: 2px;
    transition: transform 0.22s ease, opacity 0.22s ease, width 0.22s ease;
    transform-origin: center;
  }
  /* Estado abierto → X */
  .hamburger.ham-open { background: rgba(0,240,128,0.08); border-color: rgba(0,240,128,0.3); }
  .hamburger.ham-open span:nth-child(1) { transform: translateY(6.5px) rotate(45deg); }
  .hamburger.ham-open span:nth-child(2) { opacity: 0; width: 0; }
  .hamburger.ham-open span:nth-child(3) { transform: translateY(-6.5px) rotate(-45deg); }

  /* ════════════════════════════════════════════
     MENÚ MÓVIL
     ════════════════════════════════════════════ */
  .mm-backdrop {
    position: fixed; inset: 0; top: 64px;
    z-index: 48;
    background: rgba(0,0,0,0.45);
    backdrop-filter: blur(4px);
  }
  .mobile-menu {
    position: fixed;
    top: 64px; left: 0; right: 0;
    z-index: 49;
    background: rgba(4,6,18,0.98);
    backdrop-filter: blur(24px);
    border-bottom: 1px solid rgba(255,255,255,0.08);
    padding: 10px 14px 14px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    animation: mmSlide 0.18s ease-out both;
  }
  @keyframes mmSlide {
    from { opacity: 0; transform: translateY(-8px); }
    to   { opacity: 1; transform: translateY(0); }
  }
  .mm-item {
    display: flex; align-items: center; gap: 12px;
    width: 100%; padding: 13px 16px;
    background: rgba(255,255,255,0.03);
    border: 1px solid rgba(255,255,255,0.06);
    border-radius: 10px;
    color: rgba(255,255,255,0.55);
    font-family: 'Rajdhani', sans-serif;
    font-size: 1rem; font-weight: 900; font-style: italic;
    text-transform: uppercase; letter-spacing: 0.1em;
    cursor: pointer; transition: all 0.15s;
  }
  .mm-item:hover { background: rgba(255,255,255,0.07); color: rgba(255,255,255,0.88); }
  .mm-active {
    background: rgba(0,240,128,0.07) !important;
    border-color: rgba(0,240,128,0.25) !important;
    color: #00f080 !important;
  }
  .mm-sound { color: rgba(255,255,255,0.4); }
  .mm-sep { height: 1px; background: rgba(255,255,255,0.06); margin: 4px 0; }
  .mm-badge {
    margin-left: auto;
    background: rgba(245,200,24,0.15);
    border: 1px solid rgba(245,200,24,0.35);
    color: #f5c818;
    padding: 2px 10px; border-radius: 20px;
    font-size: 0.7rem; font-style: normal;
  }

  /* ════════════════════════════════════════════
     MOBILE FILTER BUTTON (oculto en desktop)
     ════════════════════════════════════════════ */
  .mobile-filter-btn {
    display: none;
    align-items: center; gap: 6px;
    padding: 7px 14px; border-radius: 20px;
    border: 1px solid rgba(255,255,255,0.12);
    background: rgba(255,255,255,0.05);
    color: rgba(255,255,255,0.55);
    font-size: 0.62rem; font-weight: 700;
    text-transform: uppercase; letter-spacing: 0.1em;
    cursor: pointer; transition: all 0.18s;
    font-family: 'Rajdhani', sans-serif; font-style: italic;
  }
  .mobile-filter-btn:hover {
    background: rgba(255,255,255,0.1);
    color: rgba(255,255,255,0.85);
  }

  /* ════════════════════════════════════════════
     RESPONSIVE — TABLET/MOBILE (≤900px)
     ════════════════════════════════════════════ */
  @media (max-width: 900px) {
    /* Header */
    .brand-center { display: none; }
    .logo-sub     { display: none; }
    .nav-right    { display: none; }
    .hamburger    { display: flex; }
    .header-inner { padding: 0 14px; gap: 10px; }

    /* Layout vertical */
    .page-layout { flex-direction: column; }

    /* Sidebar colapsable */
    .sidebar {
      width: 100%;
      flex-shrink: 0;
      position: static;
      height: auto;
      max-height: 0;
      overflow: hidden;
      border-right: none;
      border-bottom: 1px solid rgba(255,255,255,0.07);
      transition: max-height 0.3s ease, padding 0.3s ease;
      padding-top: 0; padding-bottom: 0;
    }
    .sidebar.filters-open {
      max-height: 75vh;
      overflow-y: auto;
      padding: 14px 13px;
    }

    /* Barra superior del área principal */
    .mobile-filter-btn { display: flex; }
    .main-top { justify-content: space-between; margin-bottom: 14px; }
    .main-area { padding: 14px 12px; }

    /* 2 tarjetas por fila */
    .card-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 10px;
    }
    /* Escala las cards para que quepan en la columna */
    :global(.wrap) {
      zoom: 0.88;
      justify-content: center;
    }
  }

  /* ════════════════════════════════════════════
     RESPONSIVE — MÓVIL PEQUEÑO (≤430px)
     ════════════════════════════════════════════ */
  @media (max-width: 430px) {
    .logo-title   { font-size: 1.1rem; }
    .header-inner { padding: 0 10px; gap: 8px; }

    .card-grid { gap: 6px; }
    :global(.wrap) { zoom: 0.80; }
  }

  /* ════════════════════════════════════════════
     RESPONSIVE — MÓVIL MUY PEQUEÑO (≤360px)
     ════════════════════════════════════════════ */
  @media (max-width: 360px) {
    :global(.wrap) { zoom: 0.73; }
  }
</style>
