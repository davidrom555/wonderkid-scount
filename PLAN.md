# 🗺️ PLAN DE DESARROLLO — The Wonderkid Scout
## FIFA 26 Dashboard · Arquitectura Fullstack

> **Estado general:** ✅ Stack completo funcional — en expansión
> **Stack:** Svelte + Go (Fiber v2) + MariaDB

---

## 📊 ESTADO POR FASE

| Fase | Descripción | Estado |
|------|-------------|--------|
| 0 | Setup de estructura y archivos de persistencia | ✅ Completado |
| 1 | Backend Go: modelos GORM + endpoints Fiber | ✅ Completado |
| 2 | Frontend Svelte: App.svelte + PlayerCard.svelte | ✅ Completado |
| 3 | Integración Frontend ↔ Backend | ✅ Completado |
| 4 | Features avanzadas (ROI, Watchlist, Fotos, Audio) | ✅ Completado |
| 5 | Polish UI + expansión de datos | 🔄 En progreso |

---

## 🏗️ ARQUITECTURA DEL SISTEMA

```
FIFA2/
├── PLAN.md                     ← Este archivo
├── PROMPTS.md                  ← Historial de prompts del usuario
├── README.md                   ← Guía de inicio rápido
├── START.bat                   ← Lanza backend (Go :8080) + frontend (Vite :5173)
├── BUILD.bat                   ← Recompila fifa-scout.exe
│
├── backend/                    ← Go 1.21 + Fiber v2.52.5 + GORM + MariaDB
│   ├── main.go                 ← Entry point, rutas, CORS, static /sound
│   ├── fifa-scout.exe          ← Binario compilado (go build)
│   ├── go.mod                  ← fiber v2.52.5, gorm v1.25.12, mysql driver v1.5.7
│   ├── go.sum
│   ├── models/
│   │   ├── player.go           ← Player (19 campos) + PredictionResponse
│   │   └── watchlist.go        ← WatchlistItem (FK a Player, uniqueIndex)
│   ├── handlers/
│   │   ├── players.go          ← GET /api/players, GET /api/players/filter
│   │   ├── predict.go          ← GET /api/predict/:id
│   │   └── watchlist.go        ← GET/POST/DELETE /api/watchlist/:id
│   ├── database/
│   │   └── db.go               ← Conexión MariaDB + fallback mock (88 jugadores)
│   └── sound/
│       └── fifa.mp3            ← Música de fondo (servida como static en /sound)
│
└── frontend/                   ← Svelte 4.2.20 + Vite 5.4.10 + TypeScript 5.6.3
    ├── src/
    │   ├── main.ts             ← Entry point
    │   ├── App.svelte          ← Layout completo: header + sidebar + grid + tabs
    │   │                          (filtros, stats overview, top prospects, watchlist)
    │   ├── components/
    │   │   └── PlayerCard.svelte   ← Tarjeta EA FC 26 FUT (205×285px)
    │   ├── stores/
    │   │   └── players.ts          ← 8 writable + 1 derived + helpers
    │   └── lib/
    │       └── sounds.ts           ← Web Audio API: UI sounds + música de fondo
    ├── public/
    │   └── sound/
    │       └── fifa.mp3
    ├── dist/                   ← Build de producción
    ├── index.html
    ├── package.json            ← version 2.0.0
    ├── vite.config.ts          ← vitePreprocess() + proxy /api → :8080
    ├── tsconfig.json
    ├── tsconfig.node.json
    ├── tailwind.config.js
    └── postcss.config.js
```

---

## 🔌 API CONTRACT

### Endpoints Backend (Puerto 8080)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/players` | Todos los jugadores |
| GET | `/api/players/filter?minPotential=80&position=ST&search=bell` | Filtrado por nombre, club, posición y potencial |
| GET | `/api/predict/:id` | Predicción de crecimiento |
| GET | `/api/watchlist` | Shortlist → `{"players": [...], "count": N}` |
| POST | `/api/watchlist/:id` | Añadir jugador (idempotente) |
| DELETE | `/api/watchlist/:id` | Quitar jugador |
| GET | `/sound/fifa.mp3` | Música de fondo (static) |

### Modelo de Datos — Player (19 campos)

```json
{
  "id": 1,
  "name": "Jude Bellingham",
  "sofifaId": 234396,
  "position": "CAM",
  "age": 20,
  "overall": 88,
  "potential": 94,
  "marketValue": 180000000,
  "club": "Real Madrid",
  "nationality": "England",
  "pac": 78,
  "sho": 82,
  "pas": 85,
  "dri": 88,
  "def": 65,
  "phy": 80,
  "createdAt": "2026-02-19T00:00:00Z"
}
```

### Modelo de Datos — WatchlistItem

```json
{
  "id": 1,
  "playerId": 5,
  "player": { ... }
}
```

### Respuesta de Predicción

```json
{
  "growthRate": 2.14,
  "projectRating": 94,
  "isUndervalued": true,
  "roiScore": 1.85,
  "badge": "🔥 UNDERVALUED TARGET"
}
```

---

## 🧮 ALGORITMOS CLAVE

### Predicción de Crecimiento
```
ageFactor     = max((age - 15) * 0.8, 0.8)
growthRate    = (potential - overall) / ageFactor
projectRating = min(overall + growthRate * ageFactor, 99)
```

### ROI Score
```
roiScore      = potential / (marketValue / 10_000_000)
isUndervalued = roiScore > 1.2
```

### Sistema de Badges
| ROI Score | Badge | Comportamiento |
|-----------|-------|----------------|
| > 1.2 | 🔥 UNDERVALUED TARGET | Glow naranja pulsante |
| 0.8 – 1.2 | 📊 FAIR VALUE | Sin glow |
| < 0.8 | 💎 PREMIUM ASSET | Sin glow |

### Card Tiers (por Potential)
| Potential | Tier | Gradiente |
|-----------|------|-----------|
| ≥ 95 | Diamond | Gradiente púrpura |
| 88 – 94 | Gold | Gradiente dorado |
| < 88 | Silver | Gradiente plateado |

---

## 🎨 DESIGN SYSTEM

### Layout
- **Header:** 64px fijo — "WONDERKID SCOUT" con branding EA FC 26, tabs, botón audio
- **Sidebar:** 265px fijo a la izquierda — filtros, stats, top prospects, undervalued
- **Main grid:** Responsivo — cards de 205×285px

### Paleta de Colores
```css
--bg-primary:    #020617   /* Azul casi negro */
--bg-secondary:  #1e1b4b   /* Índigo oscuro */
--glass-bg:      rgba(255,255,255,0.05)
--glass-border:  rgba(255,255,255,0.10)
--accent-green:  #00f080   /* "SCOUT" en el header */
--neon-purple:   #a855f7
--neon-orange:   #f97316
--neon-blue:     #3b82f6
--neon-green:    #22c55e
```

### Componente PlayerCard
```
Dimensiones: 205px × 285px (estilo FUT)
Foto:        cdn.sofifa.net/players/.../26_240.png
Fallback:    Avatar con iniciales si la foto falla
Tier glow:   glowDiamond / glowGold / glowSilver (2.2s ease-in-out infinite)
```

```css
/* Estructura por secciones */
[Foto 165px]
  → Rating OVR (top-left, 3rem bold) + POS + 🔥 si undervalued
  → Botón ★ watchlist (top-right)
[Sección info]
  → Nombre (apellido o completo si corto)
  → Stats grid 6 cols: PAC SHO PAS DRI DEF PHY
  → Badges: flag (flagcdn.com) + league logo (ESPN) + club logo (ESPN)
  → Predicción (si showPredictions): "Peak {rating} +{growthRate}/yr"
```

### Tipografía
- **Títulos:** `Rajdhani, sans-serif; font-weight: 900; font-style: italic; text-transform: uppercase;`
- **Stats:** `Inter, sans-serif; font-weight: 800;`
- **Labels:** `font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase;`

### CDNs de imágenes
| Recurso | CDN |
|---------|-----|
| Fotos jugadores | `cdn.sofifa.net/players/{id1}/{id2}/26_240.png` |
| Banderas | `flagcdn.com/w40/{code}.png` |
| Logos de clubs | ESPN CDN |
| Logos de ligas | ESPN CDN |

---

## 🔊 SISTEMA DE AUDIO

Implementado en `frontend/src/lib/sounds.ts` — **0 dependencias externas**, solo Web Audio API nativa del browser.

### Arquitectura

```
sounds.ts
├── AudioContext singleton (_ctx)  ← se crea lazy, se resume si está suspended
├── tone()                         ← utilidad base: oscilador + gain + ramp
├── Sonidos UI (Web Audio API)
│   ├── playNav()          → navegación con teclado entre cards
│   ├── playSelect()       → 3 tonos simultáneos (acorde C-E-G)
│   ├── playWatchlistOn()  → 3 notas ascendentes al añadir a watchlist
│   └── playWatchlistOff() → 2 frecuencias descendentes al quitar
├── Jingle EA Sports
│   └── playEAIntro()      → sawtooth + delay/feedback + 9 notas + fade-out
└── Música de fondo (HTMLAudioElement)
    ├── playBGMusic()  → carga /sound/fifa.mp3, loop, fade-in a 12% en 80ms steps
    ├── stopBGMusic()  → fade-out en 60ms steps hasta 0, luego pause
    └── isBGMusicPlaying() → estado booleano
```

### Tabla de funciones

| Función | Tecnología | Descripción |
|---------|-----------|-------------|
| `playNav()` | Web Audio API | Tone 680→480Hz, 38ms, sine wave |
| `playSelect()` | Web Audio API | 3 tonos simultáneos en secuencia (C5-E5-G5) |
| `playWatchlistOn()` | Web Audio API | 3 notas subiendo (C5→E5→G5), 0.09s entre notas |
| `playWatchlistOff()` | Web Audio API | 2 frecuencias descendentes (600→380Hz, 400→280Hz) |
| `playEAIntro()` | Web Audio API | Jingle 9 notas tipo EA Sports, delay + feedback loop |
| `playBGMusic()` | HTMLAudioElement | `fifa.mp3` en loop, fade-in gradual hasta vol 0.12 |
| `stopBGMusic()` | HTMLAudioElement | Fade-out progresivo hasta vol 0, luego `pause()` |

### Integración en UI
- **Header:** botón con animación ecualizador (4 barras CSS) cuando la música está activa
- **PlayerCard:** `playNav()` en navegación de teclado, `playWatchlistOn/Off()` al marcar/desmarcar
- **Todos los sonidos** envueltos en `try/catch` → nunca rompen la UI si el browser bloquea audio

---

## 🗃️ STATE MANAGEMENT (Svelte Stores)

Definido en `frontend/src/stores/players.ts`:

```ts
// Writable stores
players          = writable<Player[]>([])
watchlistIds     = writable<Set<number>>(new Set())
predictions      = writable<Map<number, Prediction>>(new Map())
searchTerm       = writable('')
minPotential     = writable(75)
selectedPosition = writable('')
showPredictions  = writable(false)
activeTab        = writable<'dashboard' | 'watchlist'>('dashboard')

// Derived store
filteredPlayers  = derived([players, searchTerm, minPotential, selectedPosition], ...)
```

### Funciones helper
- `calcROI(player)` → potential / (marketValue / 10M)
- `getPhotoUrl(sofifaId)` → URL CDN SoFIFA
- `getFlagEmoji(nationality)` → emoji de bandera
- `getFlagUrl(nationality)` → flagcdn.com URL
- `getClubInfo(club)` → `{ abbr, color }`
- `getLeagueInfo(club)` → `{ abbr, color }`
- `getClubLogoUrl(club)` → ESPN CDN
- `getLeagueLogoUrl(club)` → ESPN CDN
- `getCardTier(potential)` → `'diamond' | 'gold' | 'silver'`

### Mappings incluidos
- 35+ clubs con logos ESPN
- 8 ligas (La Liga, Premier, Bundesliga, Ligue 1, Serie A, Liga Por, Eredivisie)
- 24+ países con emojis y códigos flagcdn

---

## ⌨️ NAVEGACIÓN POR TECLADO

| Tecla | Acción |
|-------|--------|
| ←→↑↓ Arrow keys | Navega entre tarjetas del grid |
| Tab | Enfoca botón watchlist de la tarjeta activa |
| Enter / Space | Toggle watchlist del jugador enfocado |
| Escape | Limpia foco |

Detección automática de teclado vs mouse para activar/desactivar hover effects.

---

## 📋 CHECKLIST DE IMPLEMENTACIÓN

### FASE 1 — Backend Go
- ✅ `go.mod` con dependencias (fiber v2.52.5, gorm v1.25.12, mysql driver v1.5.7)
- ✅ `go.sum` generado
- ✅ Modelo `Player` (19 campos: id, name, position, age, overall, potential, marketValue, club, nationality, sofifaId, pac, sho, pas, dri, def, phy, createdAt)
- ✅ Modelo `WatchlistItem` con uniqueIndex en PlayerID
- ✅ Conexión a MariaDB (`localhost:3306/fifa_scouting`) con fallback mock automático
- ✅ Seed con 88 jugadores wonderkids reales
- ✅ Handler `GET /api/players`
- ✅ Handler `GET /api/players/filter` (minPotential + position + search en nombre y club)
- ✅ Handler `GET /api/predict/:id` (ageFactor mínimo 0.8, badge en respuesta)
- ✅ Handler `GET /api/watchlist` (retorna `{players, count}`)
- ✅ Handler `POST /api/watchlist/:id` (idempotente)
- ✅ Handler `DELETE /api/watchlist/:id`
- ✅ CORS middleware
- ✅ Static `/sound` (sirve fifa.mp3)
- ✅ Server en :8080
- ✅ `backend/fifa-scout.exe` compilado

### FASE 2 — Frontend Svelte
- ✅ Vite + Svelte 4 + TypeScript scaffold
- ✅ `tsconfig.json` y `tsconfig.node.json`
- ✅ Tailwind CSS + PostCSS
- ✅ `vitePreprocess()` en vite.config.ts
- ✅ Fuentes Rajdhani + Inter
- ✅ `App.svelte`: header fijo + sidebar (filtros, stats, top prospects) + grid + tabs
- ✅ `PlayerCard.svelte`: 205×285px, foto CDN, rating, 6 stats, badges (flag/league/club), tiers, glow
- ✅ Filtrado por: texto (nombre + club), potencial mínimo (slider 50–99), posición (grid 4×3)
- ✅ Toggle "AI Predictions" en sidebar
- ✅ Tab Watchlist con contador
- ✅ `stores/players.ts`: 8 writable + 1 derived + funciones helper
- ✅ `lib/sounds.ts`: Web Audio API + música de fondo fifa.mp3

### FASE 3 — Integración
- ✅ Proxy Vite `/api` → `localhost:8080`
- ✅ `onMount` fetch de jugadores y watchlist al cargar
- ✅ Filtrado reactivo con derived store
- ✅ Toggle watchlist (POST/DELETE) con optimistic UI
- ✅ Predicciones cacheadas en `Map<number, Prediction>`
- ✅ Búsqueda en tiempo real (nombre + club)
- ✅ Error handling básico

### FASE 4 — Features Pro
- ✅ ROI scoring y badge system (undervalued / fair / premium)
- ✅ Glow pulsante por tier (glowDiamond / glowGold / glowSilver)
- ✅ Fotos reales via CDN SoFIFA (26_240.png)
- ✅ Fallback avatar con iniciales
- ✅ Stats FIFA 26: PAC / SHO / PAS / DRI / DEF / PHY en la tarjeta
- ✅ Banderas via flagcdn.com
- ✅ Logos de clubs y ligas via ESPN CDN
- ✅ Top Prospects en sidebar (5 mejores por potential)
- ✅ Undervalued panel en sidebar (4 mejores por ROI)
- ✅ Sistema de audio completo (UI sounds + música de fondo)
- ✅ Navegación por teclado completa
- ✅ Predicción "Peak {rating} +{growthRate}/yr" en la tarjeta

### FASE 5 — Polish & Expansión (en progreso)
- ✅ `START.bat` y `BUILD.bat`
- ✅ README con instrucciones
- 🔄 Responsive completo (mobile 1 col / tablet 2 col / desktop 4 col)
- 🔄 Loading spinner durante fetch
- 🔄 Estado vacío cuando no hay resultados

---

## 👥 BASE DE DATOS DE JUGADORES

- **Total en mock/DB:** 88 jugadores reales
- **Posiciones disponibles:** ST, CAM, CB, CM, RW, LW, GK, LB, CDM, RB
- **Ligas cubiertas:** La Liga, Premier League, Bundesliga, Ligue 1, Serie A, Liga Portuguesa, Eredivisie
- **Clubs incluidos:** Real Madrid, Man City, Barcelona, Arsenal, Liverpool, Bayern, PSG, Juventus, Chelsea, Man United, Leverkusen, Dortmund, Atlético, Inter, y más

---

## 🔄 DECISIONES TÉCNICAS

| Decisión | Opción elegida | Razón |
|----------|---------------|-------|
| Componentes UI | `App.svelte` + `PlayerCard.svelte` | Control total del layout, sin overhead |
| DB connection | MariaDB + fallback mock automático | Funciona con o sin MariaDB instalada |
| Filtering | Frontend (derived store Svelte) | Latencia cero, UX fluida |
| Predictions cache | `Map<number, Prediction>` en store | Evita re-llamadas innecesarias a la API |
| Image loading | Preload + filtrado si fallan | No muestra frames rotos |
| Audio | Web Audio API nativa | Sin dependencias externas |
| Player photos | CDN SoFIFA 240px (`26_240.png`) | Resolución óptima para card FUT |
| Flag images | flagcdn.com 40px | Banderas rectangulares reales |
| Club/league logos | ESPN CDN | Logos oficiales sin hosting propio |
| Keyboard nav | `tabindex` + keydown events | Accesibilidad completa |
| Seed strategy | Upsert (ON DUPLICATE KEY) | Añadir jugadores sin borrar watchlist |
| Go version | 1.21 | Compatibilidad y estabilidad |
| Frontend port | 5173 (Vite default) | Proxy a :8080 |
