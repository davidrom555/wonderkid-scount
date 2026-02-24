<script lang="ts">
  import {
    watchlistIds, predictions, showPredictions,
    getPhotoUrl, getFlagEmoji, getFlagUrl, getClubInfo, getLeagueInfo,
    getClubLogoUrl, getLeagueLogoUrl, getCardTier, calcROI
  } from '../stores/players'
  import { playNav, playWatchlistOn, playWatchlistOff } from '../lib/sounds'
  import type { Player } from '../stores/players'

  export let player: Player

  $: isInWatchlist = $watchlistIds.has(player.id)
  $: prediction = $predictions.get(player.id) ?? null
  $: isUndervalued = prediction?.isUndervalued ?? false
  $: tier = getCardTier(player.potential)
  $: photoUrl = getPhotoUrl(player.sofifaId)
  $: flag = getFlagEmoji(player.nationality)
  $: flagUrl = getFlagUrl(player.nationality)
  $: clubInfo = getClubInfo(player.club)
  $: leagueInfo = getLeagueInfo(player.club)
  $: clubLogoUrl = getClubLogoUrl(player.club)
  $: leagueLogoUrl = getLeagueLogoUrl(player.club)

  let imgLoaded = false
  let imgError = false
  let clubLogoError = false
  let leagueLogoError = false

  // ─── Referencias DOM para navegación por teclado ───
  let cardEl: HTMLDivElement
  let wlBtnEl: HTMLButtonElement

  function allCards(): HTMLElement[] {
    return Array.from(document.querySelectorAll('[data-card]'))
  }
  function myIndex(): number {
    return allCards().indexOf(cardEl)
  }
  function focusCardAt(i: number) {
    const cards = allCards()
    if (i >= 0 && i < cards.length) {
      playNav()
      cards[i].focus()
      cards[i].scrollIntoView({ block: 'nearest', behavior: 'smooth' })
    }
  }

  // Navega a la card más cercana en la fila de abajo (misma columna X)
  function focusCardBelow() {
    const cards = allCards()
    const myRect = cardEl.getBoundingClientRect()
    const myCX = myRect.left + myRect.width / 2
    const below = cards.filter(c => c.getBoundingClientRect().top > myRect.top + myRect.height * 0.5)
    if (!below.length) return
    const minTop = Math.min(...below.map(c => c.getBoundingClientRect().top))
    const row = below.filter(c => c.getBoundingClientRect().top < minTop + 10)
    row.sort((a, b) => {
      const ar = a.getBoundingClientRect(), br = b.getBoundingClientRect()
      return Math.abs(ar.left + ar.width / 2 - myCX) - Math.abs(br.left + br.width / 2 - myCX)
    })
    if (row[0]) { playNav(); row[0].focus(); row[0].scrollIntoView({ block: 'nearest', behavior: 'smooth' }) }
  }

  // Navega a la card más cercana en la fila de arriba (misma columna X)
  function focusCardAbove() {
    const cards = allCards()
    const myRect = cardEl.getBoundingClientRect()
    const myCX = myRect.left + myRect.width / 2
    const above = cards.filter(c => c.getBoundingClientRect().top < myRect.top - myRect.height * 0.5)
    if (!above.length) return
    const maxTop = Math.max(...above.map(c => c.getBoundingClientRect().top))
    const row = above.filter(c => c.getBoundingClientRect().top > maxTop - 10)
    row.sort((a, b) => {
      const ar = a.getBoundingClientRect(), br = b.getBoundingClientRect()
      return Math.abs(ar.left + ar.width / 2 - myCX) - Math.abs(br.left + br.width / 2 - myCX)
    })
    if (row[0]) { playNav(); row[0].focus(); row[0].scrollIntoView({ block: 'nearest', behavior: 'smooth' }) }
  }

  function handleCardKeydown(e: KeyboardEvent) {
    switch (e.key) {
      case 'Tab':
        if (!e.shiftKey) { e.preventDefault(); wlBtnEl?.focus() }
        break
      case 'ArrowRight':
        e.preventDefault(); focusCardAt(myIndex() + 1); break
      case 'ArrowLeft':
        e.preventDefault(); focusCardAt(myIndex() - 1); break
      case 'ArrowDown':
        e.preventDefault(); focusCardBelow(); break
      case 'ArrowUp':
        e.preventDefault(); focusCardAbove(); break
      case 'Enter':
      case ' ':
        e.preventDefault(); toggleWatchlist(); break
    }
  }

  function handleWlKeydown(e: KeyboardEvent) {
    switch (e.key) {
      case 'ArrowLeft':
      case 'Escape':
        e.preventDefault(); cardEl?.focus(); break
      case 'Tab':
        e.preventDefault()
        e.stopPropagation()
        if (e.shiftKey) { cardEl?.focus() }
        else { focusCardAt(myIndex() + 1) }
        break
    }
  }

  function handleImgLoad() { imgLoaded = true }
  function handleImgError() { imgError = true; imgLoaded = true }
  function hideBadge(e: Event) {
    const el = e.currentTarget as HTMLElement
    el.style.display = 'none'
  }
  function removeBadgeWrap(e: Event) {
    const el = e.currentTarget as HTMLElement
    const wrap = el.closest('.badge-icon') as HTMLElement | null
    if (wrap) wrap.style.display = 'none'
  }

  async function toggleWatchlist() {
    isInWatchlist ? playWatchlistOff() : playWatchlistOn()
    const method = isInWatchlist ? 'DELETE' : 'POST'
    await fetch(`/api/watchlist/${player.id}`, { method })
    watchlistIds.update(s => {
      const n = new Set(s)
      isInWatchlist ? n.delete(player.id) : n.add(player.id)
      return n
    })
  }

  async function loadPrediction() {
    if ($predictions.has(player.id)) return
    const res = await fetch(`/api/predict/${player.id}`)
    if (res.ok) {
      const data = await res.json()
      predictions.update(m => new Map(m).set(player.id, data))
    }
  }

  $: if ($showPredictions) loadPrediction()

  const stats = [
    { key: 'pac', label: 'PAC' },
    { key: 'sho', label: 'SHO' },
    { key: 'pas', label: 'PAS' },
    { key: 'dri', label: 'DRI' },
    { key: 'def', label: 'DEF' },
    { key: 'phy', label: 'PHY' },
  ] as const

  // Apellido o nombre completo si es corto
  $: displayName = (() => {
    const parts = player.name.trim().split(' ')
    return parts.length > 1 ? parts[parts.length - 1] : parts[0]
  })()

  function fmtVal(v: number): string {
    return v >= 1_000_000 ? (v / 1_000_000).toFixed(0) + 'M' : (v / 1_000).toFixed(0) + 'K'
  }
</script>

<!-- EA FC 26 FUT Card -->
<div class="wrap" style="animation: fadeUp 0.4s ease-out both;">
  <div
    class="card"
    class:gold={tier === 'gold'}
    class:diamond={tier === 'diamond'}
    class:silver={tier === 'silver'}
    class:undervalued={isUndervalued}
    tabindex="0"
    role="article"
    aria-label="{player.name}, {player.overall} overall, {player.position}"
    data-card
    bind:this={cardEl}
    on:keydown={handleCardKeydown}
    on:click={() => {}}
  >
    <!-- ── Rayos radiales (over background, under everything) ── -->
    <div class="rays"></div>

    <!-- ── Botón watchlist (esquina superior derecha) ── -->
    <button
      class="wl-btn"
      class:active={isInWatchlist}
      on:click={toggleWatchlist}
      on:keydown={handleWlKeydown}
      tabindex="-1"
      bind:this={wlBtnEl}
      title={isInWatchlist ? 'Quitar de watchlist' : 'Añadir a watchlist'}
    >★</button>

    <!-- ════════════════════════════════
         SECCIÓN FOTO (parte superior)
         Rating + posición superpuestos
         ════════════════════════════════ -->
    <div class="photo-area">

      <!-- Foto del jugador -->
      {#if !imgError}
        <img
          class="photo" class:show={imgLoaded}
          src={photoUrl}
          alt={player.name}
          referrerpolicy="no-referrer"
          on:load={handleImgLoad}
          on:error={handleImgError}
        />
      {:else}
        <div class="avatar">
          {player.name.split(' ').map(w => w[0]).join('').slice(0, 2)}
        </div>
      {/if}

      {#if !imgLoaded && !imgError}
        <div class="skel"></div>
      {/if}

      <!-- Rating + posición superpuesto sobre la foto -->
      <div class="rating-box">
        <span class="ovr">{player.overall}</span>
        <span class="pos">{player.position}</span>
        {#if isUndervalued}<span class="fire">🔥</span>{/if}
      </div>

    </div>

    <!-- ════════════════════════════════
         SECCIÓN INFERIOR
         Nombre · Stats · Badges
         ════════════════════════════════ -->
    <div class="bottom">

      <!-- Nombre -->
      <h2 class="pname">{displayName}</h2>

      <div class="sep"></div>

      <!-- Stats: label arriba, valor abajo -->
      <div class="stats">
        {#each stats as s}
          <div class="sc">
            <span class="slbl">{s.label}</span>
            <span class="sval">{player[s.key]}</span>
          </div>
        {/each}
      </div>

      <div class="sep"></div>

      <!-- Badges: bandera · liga · club (estilo FC 26) -->
      <div class="badges">

        <!-- Bandera: imagen rectangular flagcdn.com -->
        <div class="badge-flag">
          {#if flagUrl}
            <img
              src={flagUrl}
              alt={player.nationality}
              class="flag-img"
              referrerpolicy="no-referrer"
              on:error={hideBadge}
            />
          {:else}
            <span class="flag-fallback">{flag}</span>
          {/if}
        </div>

        <!-- Liga: imagen real, si falla no mostrar nada -->
        {#if leagueLogoUrl && !leagueLogoError}
          <img
            src={leagueLogoUrl}
            alt={leagueInfo.abbr}
            class="badge-icon"
            referrerpolicy="no-referrer"
            on:error={() => leagueLogoError = true}
          />
        {/if}

        <!-- Club: imagen real, si falla no mostrar nada -->
        {#if clubLogoUrl && !clubLogoError}
          <img
            src={clubLogoUrl}
            alt={player.club}
            class="badge-icon"
            referrerpolicy="no-referrer"
            on:error={() => clubLogoError = true}
          />
        {/if}

      </div>

      <!-- Predicción (si está activa) -->
      {#if $showPredictions && prediction}
        <div class="pred">
          <span>Peak {prediction.projectRating}</span>
          <span>+{prediction.growthRate.toFixed(1)}/yr</span>
        </div>
      {/if}

    </div>
  </div>
</div>

<style>
  /* ════════ WRAPPER ════════ */
  .wrap { display: flex; justify-content: center; }

  /* ════════ CARD BASE ════════ */
  .card {
    position: relative;
    width: 205px;
    height: 285px;
    border-radius: 14px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    cursor: pointer;
    transition: transform 0.25s cubic-bezier(.34,1.56,.64,1), box-shadow 0.25s ease;
    user-select: none;
  }
  .card:hover { transform: translateY(-10px) scale(1.05); z-index: 10; }
  :global(body.kb-nav) .card:hover { transform: none; z-index: auto; }

  /* ── Focus ring dorado (navegación por teclado) ── */
  .card:focus { outline: none; }
  .card:focus-visible {
    outline: none;
    transform: translateY(-6px) scale(1.03);
    z-index: 10;
  }
  .card.gold:focus-visible {
    box-shadow:
      0 0 0 3px #f5c818,
      0 0 0 5px rgba(245,200,24,0.3),
      0 14px 50px rgba(0,0,0,0.75),
      inset 0 1.5px 0 rgba(255,240,120,0.6);
  }
  .card.diamond:focus-visible {
    box-shadow:
      0 0 0 3px #c084fc,
      0 0 0 5px rgba(192,132,252,0.3),
      0 0 30px rgba(124,58,237,0.45),
      0 14px 50px rgba(0,0,0,0.85);
  }
  .card.silver:focus-visible {
    box-shadow:
      0 0 0 3px #94aac0,
      0 0 0 5px rgba(148,170,192,0.3),
      0 14px 45px rgba(0,0,0,0.65);
  }

  /* ── Focus en el botón watchlist ── */
  .wl-btn:focus { outline: none; }
  .wl-btn:focus-visible { outline: none; }
  .gold .wl-btn:focus-visible        { color: #f5c818; transform: scale(1.35); }
  .diamond .wl-btn:focus-visible     { color: #c084fc; transform: scale(1.35); }
  .silver .wl-btn:focus-visible      { color: rgba(255,255,255,0.85); transform: scale(1.35); }

  /* ════════ GOLD ════════ */
  .card.gold {
    background: rgba(245, 200, 24, 0.07);
    backdrop-filter: blur(18px);
    -webkit-backdrop-filter: blur(18px);
    box-shadow:
      0 0 0 1px rgba(245,200,24,0.55),
      0 0 24px rgba(245,200,24,0.18),
      0 14px 50px rgba(0,0,0,0.65),
      inset 0 1px 0 rgba(255,240,120,0.2),
      inset 0 0 40px rgba(245,200,24,0.04);
  }

  /* ════════ DIAMOND ════════ */
  .card.diamond {
    background: rgba(124, 58, 237, 0.09);
    backdrop-filter: blur(18px);
    -webkit-backdrop-filter: blur(18px);
    box-shadow:
      0 0 0 1px rgba(167,107,254,0.65),
      0 0 30px rgba(124,58,237,0.35),
      0 14px 50px rgba(0,0,0,0.75),
      inset 0 1px 0 rgba(192,132,252,0.15),
      inset 0 0 40px rgba(124,58,237,0.05);
  }

  /* ════════ SILVER ════════ */
  .card.silver {
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(18px);
    -webkit-backdrop-filter: blur(18px);
    box-shadow:
      0 0 0 1px rgba(200,220,240,0.3),
      0 14px 45px rgba(0,0,0,0.55),
      inset 0 1px 0 rgba(255,255,255,0.12),
      inset 0 0 40px rgba(255,255,255,0.02);
  }

  /* ════════ RAYOS ════════ */
  .rays {
    position: absolute;
    inset: 0;
    border-radius: inherit;
    pointer-events: none;
    z-index: 0;
  }
  .gold .rays {
    background-image: repeating-conic-gradient(
      from 0deg at 50% 28%,
      rgba(255,245,100,0.14) 0deg 5deg,
      rgba(0,0,0,0.07) 5deg 10deg
    );
  }
  .diamond .rays {
    background-image: repeating-conic-gradient(
      from 0deg at 50% 28%,
      rgba(192,132,252,0.11) 0deg 5deg,
      rgba(0,0,0,0.08) 5deg 10deg
    );
  }
  .silver .rays {
    background-image: repeating-conic-gradient(
      from 0deg at 50% 28%,
      rgba(255,255,255,0.14) 0deg 5deg,
      rgba(0,0,0,0.05) 5deg 10deg
    );
  }

  /* ════════ BOTÓN WATCHLIST ════════ */
  .wl-btn {
    position: absolute;
    top: 9px; right: 10px;
    z-index: 20;
    background: none; border: none;
    cursor: pointer;
    font-size: 1.15rem;
    line-height: 1;
    padding: 2px;
    transition: color 0.18s, transform 0.18s;
  }
  .gold .wl-btn         { color: rgba(245,200,24,0.3); }
  .gold .wl-btn:hover   { color: rgba(245,200,24,0.9); transform: scale(1.35); }
  .gold .wl-btn.active  { color: #f5c818; }
  .diamond .wl-btn        { color: rgba(192,132,252,0.3); }
  .diamond .wl-btn:hover  { color: #c084fc; transform: scale(1.35); }
  .diamond .wl-btn.active { color: #c084fc; }
  .silver .wl-btn         { color: rgba(255,255,255,0.25); }
  .silver .wl-btn:hover   { color: rgba(255,255,255,0.85); transform: scale(1.35); }
  .silver .wl-btn.active  { color: #ffffff; }

  /* ════════ ÁREA FOTO ════════ */
  .photo-area {
    position: relative;
    height: 165px;
    flex-shrink: 0;
    overflow: hidden;
    z-index: 2;
  }

  .photo {
    position: absolute;
    bottom: 0; left: 0;
    width: 100%; height: 100%;
    object-fit: contain;
    object-position: center bottom;
    opacity: 0;
    transition: opacity 0.3s ease;
    filter: drop-shadow(0 5px 14px rgba(0,0,0,0.5));
    z-index: 2;
  }
  .photo.show { opacity: 1; }

  .avatar {
    position: absolute; inset: 0;
    display: flex; align-items: center; justify-content: center;
    font-size: 2.4rem; font-weight: 900;
    font-style: italic; text-transform: uppercase;
    z-index: 2;
  }
  .gold .avatar    { color: rgba(245,200,24,0.4); }
  .diamond .avatar { color: rgba(192,132,252,0.45); }
  .silver .avatar  { color: rgba(255,255,255,0.3); }

  .skel {
    position: absolute; inset: 0;
    background: rgba(0,0,0,0.09);
    animation: skel 1.4s ease-in-out infinite;
    z-index: 1;
  }

  /* ════════ RATING SUPERPUESTO ════════ */
  .rating-box {
    position: absolute;
    top: 9px; left: 10px;
    z-index: 10;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    line-height: 1;
  }

  .ovr {
    font-size: 3rem;
    font-weight: 900;
    font-style: italic;
    font-family: 'Rajdhani', sans-serif;
    line-height: 0.85;
    letter-spacing: -1px;
  }
  .gold .ovr    { color: #f5c818; text-shadow: 0 0 18px rgba(245,200,24,0.7), 0 2px 8px rgba(0,0,0,0.5); }
  .diamond .ovr { color: #e0c8ff; text-shadow: 0 0 18px rgba(192,132,252,0.9); }
  .silver .ovr  { color: #ffffff; text-shadow: 0 0 12px rgba(255,255,255,0.35), 0 2px 8px rgba(0,0,0,0.5); }

  .pos {
    font-size: 0.62rem;
    font-weight: 900;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    margin-top: 1px;
  }
  .gold .pos    { color: rgba(245,200,24,0.75); }
  .diamond .pos { color: #a78bfa; }
  .silver .pos  { color: rgba(255,255,255,0.55); }

  .fire { font-size: 0.78rem; margin-top: 3px; }

  /* ════════ SECCIÓN INFERIOR ════════ */
  .bottom {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 6px 10px 8px;
    position: relative;
    z-index: 2;
    gap: 0;
  }

  /* Sección inferior — glass con acento de color por tier */
  .gold .bottom {
    background: rgba(245,200,24,0.06);
    border-top: 1px solid rgba(245,200,24,0.2);
  }
  .diamond .bottom {
    background: rgba(124,58,237,0.08);
    border-top: 1px solid rgba(167,107,254,0.25);
  }
  .silver .bottom {
    background: rgba(255,255,255,0.04);
    border-top: 1px solid rgba(255,255,255,0.1);
  }

  /* ── Nombre ── */
  .pname {
    text-align: center;
    font-size: 0.9rem;
    font-weight: 900;
    text-transform: uppercase;
    letter-spacing: 0.07em;
    line-height: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-bottom: 3px;
  }
  /* Nombre — texto claro sobre fondo glass */
  .gold .pname    { color: #f5c818; }
  .diamond .pname { color: #e0c8ff; }
  .silver .pname  { color: rgba(255,255,255,0.88); }

  /* ── Separador ── */
  .sep {
    height: 1px;
    margin: 3px 0;
    flex-shrink: 0;
  }
  .gold .sep    { background: rgba(245,200,24,0.2); }
  .diamond .sep { background: rgba(167,107,254,0.25); }
  .silver .sep  { background: rgba(255,255,255,0.1); }

  /* ── Stats ── */
  .stats {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 1px;
    flex-shrink: 0;
  }
  .sc {
    display: flex; flex-direction: column;
    align-items: center; gap: 1px;
  }
  /* Label ARRIBA (pequeño) */
  .slbl {
    font-size: 0.38rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    line-height: 1;
  }
  .gold .slbl    { color: rgba(245,200,24,0.55); }
  .diamond .slbl { color: rgba(192,132,252,0.55); }
  .silver .slbl  { color: rgba(255,255,255,0.38); }

  /* Valor ABAJO (grande y negrita) */
  .sval {
    font-size: 0.73rem;
    font-weight: 900;
    line-height: 1;
  }
  .gold .sval    { color: #fff8dc; }
  .diamond .sval { color: #ead5ff; }
  .silver .sval  { color: rgba(255,255,255,0.9); }

  /* ── Badges ── */
  .badges {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 5px;
    flex-shrink: 0;
  }

  /* ── Badges: todos 22px de alto para tamaño uniforme ── */
  .badge-flag {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 22px;
    border-radius: 2px;
    overflow: hidden;
    box-shadow: 0 1px 4px rgba(0,0,0,0.55);
    flex-shrink: 0;
  }
  .flag-img {
    height: 20px;
    width: auto;
    display: block;
  }
  .flag-fallback {
    font-size: 1.1rem;
    line-height: 1;
  }

  /* Liga y club: mismo alto que la bandera */
  .badge-icon {
    height: 22px;
    width: 22px;
    object-fit: contain;
    filter: drop-shadow(0 1px 4px rgba(0,0,0,0.6));
    flex-shrink: 0;
  }

  /* (fallback CSS eliminado — si no carga la imagen, no se muestra nada) */

  /* ── Predicción ── */
  .pred {
    display: flex; justify-content: space-between;
    font-size: 0.42rem; font-weight: 900;
    text-transform: uppercase; letter-spacing: 0.03em;
    padding: 2px 5px; border-radius: 3px; margin-top: 2px;
  }
  .gold .pred    { background: rgba(245,200,24,0.1);   color: rgba(245,200,24,0.85); }
  .diamond .pred { background: rgba(192,132,252,0.1);  color: rgba(192,132,252,0.85); }
  .silver .pred  { background: rgba(255,255,255,0.06); color: rgba(255,255,255,0.55); }

  /* ════════ UNDERVALUED GLOW ════════ */
  .card.gold.undervalued {
    animation: glowGold 2.2s ease-in-out infinite;
  }
  .card.diamond.undervalued {
    animation: glowDiamond 2.2s ease-in-out infinite;
  }
  .card.silver.undervalued {
    animation: glowSilver 2.2s ease-in-out infinite;
  }

  @keyframes glowGold {
    0%,100% { box-shadow: 0 0 0 2px #b88a00, 0 14px 50px rgba(0,0,0,0.75), 0 0 18px rgba(249,115,22,0.2); }
    50%      { box-shadow: 0 0 0 2px #f97316, 0 14px 50px rgba(0,0,0,0.75), 0 0 44px rgba(249,115,22,0.65); }
  }
  @keyframes glowDiamond {
    0%,100% { box-shadow: 0 0 0 2px #7c3aed, 0 0 30px rgba(124,58,237,0.45), 0 14px 50px rgba(0,0,0,0.85); }
    50%      { box-shadow: 0 0 0 2px #f97316, 0 0 44px rgba(249,115,22,0.55), 0 14px 50px rgba(0,0,0,0.85); }
  }
  @keyframes glowSilver {
    0%,100% { box-shadow: 0 0 0 2px #8aaac0, 0 14px 45px rgba(0,0,0,0.65), 0 0 14px rgba(249,115,22,0.15); }
    50%      { box-shadow: 0 0 0 2px #f97316, 0 14px 45px rgba(0,0,0,0.65), 0 0 38px rgba(249,115,22,0.58); }
  }

  /* ════════ ANIMACIONES ════════ */
  @keyframes fadeUp {
    from { opacity: 0; transform: translateY(16px); }
    to   { opacity: 1; transform: translateY(0); }
  }
  @keyframes skel {
    0%,100% { opacity: 0.35; }
    50%      { opacity: 0.7; }
  }
</style>
