# ⚽ Wonderkid Scout — FIFA 26 Dashboard
**Fractalia Technical Exam · Fullstack (Svelte + Go + MariaDB)**

---

## Inicio Rápido

```
Doble clic → START.bat
Abre → http://localhost:5173
```

---

## Stack
- **Frontend:** Svelte
- **Backend:** Go + Fiber v2
- **DB:** MariaDB

## API Endpoints
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/players` | Todos los jugadores |
| GET | `/api/players/filter` | Filtrar por potencial/posición |
| GET | `/api/predict/:id` | Predicción de crecimiento |
| GET | `/api/watchlist` | Lista de seguimiento |
| POST | `/api/watchlist/:id` | Añadir a watchlist |
| DELETE | `/api/watchlist/:id` | Quitar de watchlist |

## Archivos de Contexto
- `PLAN.md` — Arquitectura, roadmap y decisiones técnicas
- `PROMPTS.md` — Historial de prompts para persistencia entre sesiones
