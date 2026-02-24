# 📋 PROMPTS HISTORY — The Wonderkid Scout
> Este archivo se actualiza automáticamente con cada prompt del usuario para garantizar la persistencia del contexto entre sesiones.

---

## [P001] — PROMPT INICIAL

```
# ROLE: Senior Fullstack Engineer & Game UI Architect (EA Sports Specialist)
# PROJECT: "The Wonderkid Scout" - FIFA 26 Dashboard

Actúa como un experto en Svelte, Go (Fiber) y MariaDB. Vamos a desarrollar un panel de scouting de alto rendimiento con estética "Broadcasting Glass" de EA Sports.

## Tech Stack:
- Frontend: Svelte (Vite) + Tailwind CSS
- Backend: Go con el framework Fiber
- DB: MariaDB
- Estilo: Glassmorphism extremo, desenfoques, bordes neón, tipografías gruesas (Rajdhani/Inter Black)

## Core Features:
1. Dashboard de Scouting: Grid de tarjetas con backdrop-filter: blur(12px)
2. Algoritmo de Predicción (Go): /api/predict — Crecimiento = (Potencial - Actual) / ((Edad - 15) * 0.8)
3. Filtrado Reactivo: Svelte — filtra por potencial y posición en tiempo real
4. Watchlist System: Persistencia de shortlist

## Pro Feature:
- Market ROI Heatmap: badge "🔥 UNDERVALUED TARGET" con glow animado si potencial alto + valor bajo

## UI/UX:
- Background: #020617 a #1e1b4b
- Cards: bg-white/5, border-white/10, clip-path esquinas FIFA
- Tipografía: Mayúsculas, Itálica, Peso 900

## Instrucciones especiales:
- Desarrollar buen plan ANTES de codificar
- Guardar cada prompt automáticamente en MD
- Crear dos MD: uno para el PLAN, otro para los PROMPTS
- Persistencia de historia entre sesiones
```

**Respuesta:** Creación de PLAN.md y PROMPTS.md en FIFA_26/, estructura completa del proyecto con backend Go + frontend Svelte.

---

## [P002] — GUARDAR EN FIFA2

```
todo guardalo en C:\Users\DavidRom\Desktop\proyectos\FIFA_26
```

**Respuesta:** Toda la estructura del proyecto (PLAN.md, PROMPTS.md, backend Go, frontend Svelte) creada en FIFA2/.

---

**Verificación:**
- ✅ Frontend: Svelte + Vite + TypeScript + Tailwind CSS
- ✅ Backend: Go + Fiber v2
- ✅ BBDD: MariaDB (`fifa_scouting`)
- ✅ Estilo Broadcasting Glass: glassmorphism, blur, bordes neón, Rajdhani/Inter Black
- ✅ Filtrado por potencial y posición en tiempo real
- ✅ Watchlist con persistencia en DB
- ✅ Algoritmo de predicción de crecimiento implementado

---

## [P004] — FOTOS + REDISEÑO FIFA 26

```
# FASE 4: Player Cards con identidad visual FIFA FUT

Las tarjetas deben mostrar la foto real del jugador exactamente como aparecen
en FIFA 26 Ultimate Team. Para lograrlo:

1. Investigar el CDN oficial de SoFIFA para obtener las fotos reales de cada
   jugador (formato: id → path dinámico, resolución 26_240.png).

2. Mapear los SofifaIDs de todos los jugadores del dataset actual contra la
   base de datos pública de EA Sports / SoFIFA.

3. Rediseñar la card siguiendo el estándar visual de FIFA FUT:
   - Dimensiones canónicas: 205×285px
   - Foto del jugador como elemento central dominante
   - Rating overall en esquina superior izquierda (tipografía bold/italic)
   - Stats de atributos (PAC, SHO, PAS, DRI, DEF, PHY) en franja inferior
   - Tier system con glow: Diamond ≥95 / Gold 88-94 / Silver <88
   - Bandera de nacionalidad + logo de club / liga

4. Incluir fallback con iniciales si la imagen CDN no resuelve.

5. Validar contra la referencia visual del examen Fractalia para asegurar
   coherencia con el brief original.
```

**Implementado:**
- ✅ Fotos reales via CDN SoFIFA: `cdn.sofifa.net/players/{id1}/{id2}/26_240.png`
- ✅ SofifaID añadido al modelo Player
- ✅ Stats FIFA 26: PAC, SHO, PAS, DRI, DEF, PHY (valores reales)
- ✅ Card rediseñada estilo FIFA FUT (205×285px): foto grande, rating top-left, stats abajo
- ✅ Tier glow: Diamond (púrpura ≥95) / Gold (dorado 88-94) / Silver (plateado <88)
- ✅ Fallback avatar con iniciales si la foto no carga
- ✅ Badge de potencial sobre la foto
- ✅ Banderas via flagcdn.com + logos clubs/ligas via ESPN CDN
- ✅ Header: "Wonderkid Scouting Dashboard"
- ✅ Build exitoso: 0 errores

---

## [P005] — INSTALACIÓN GO + MARIADB

```
ya instale los programas necesarios
```

**Stack activado:**
- ✅ Go 1.21 detectado — `go mod tidy` ejecutado (Fiber v2, MySQL driver)
- ✅ MariaDB detectada — base de datos `fifa_scouting` creada
- ✅ `go build` → `backend/fifa-scout.exe` compilado sin errores
- ✅ API Go funcionando: jugadores desde MariaDB, predicción, watchlist
- ✅ `START.bat` → lanza `fifa-scout.exe` + frontend Vite
- ✅ `BUILD.bat` → recompila el binario Go

**Stack 100% real:**
- Frontend: Svelte ✅
- Backend: Go + Fiber ✅
- BBDD: MariaDB ✅

---

## [P006] — AÑADIR MÁS JUGADORES

```
AGREGA MAS JUDADORES
```

**Jugadores añadidos:**
- ✅ F. Wirtz — Leverkusen — Germany — CAM — OVR:87 POT:93
- ✅ X. Simons — RB Leipzig — Netherlands — CAM — OVR:83 POT:90
- ✅ Endrick — Real Madrid — Brazil — ST — OVR:78 POT:91
- ✅ W. Zaïre-Emery — PSG — France — CM — OVR:80 POT:91
- ✅ A. Garnacho — Man United — Argentina — LW — OVR:82 POT:88
- ✅ K. Mainoo — Man United — England — CM — OVR:79 POT:87
- ✅ M. Tel — Bayern — France — ST — OVR:78 POT:88
- ✅ J. Neves — PSG — Portugal — CDM — OVR:82 POT:90
- ✅ Savinho — Man City — Brazil — RW — OVR:81 POT:88
- ✅ A. Güler — Real Madrid — Turkey — CAM — OVR:79 POT:88
- ✅ SeedDB con upsert — la watchlist existente no se borra
- ✅ Backend recompilado: 0 errores

```

**Resultado del análisis:**
- ✅ Frontend: Svelte
- ✅ Backend: Go + Fiber v2
- ✅ BBDD: MariaDB
- ✅ Estilo Broadcasting Glass: header/sidebar glassmorphism + cards FIFA FUT (EA Sports)
- ✅ Panel de scouting wonderkids con fotos reales, stats y tiers
- ✅ Filtrado por potencial (slider 50–99), posición (11 botones) y búsqueda texto
- ✅ Watchlist con persistencia en MariaDB (POST/DELETE idempotente)
- ✅ Predicción de crecimiento: growthRate = (potential − overall) / ageFactor
- ✅ Feature adicional: Market ROI Heatmap (UNDERVALUED / FAIR VALUE / PREMIUM ASSET)


<!-- PRÓXIMOS PROMPTS SE AÑADEN AQUÍ -->
