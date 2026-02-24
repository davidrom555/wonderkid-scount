import { writable, derived } from 'svelte/store'

export interface Player {
  id: number
  name: string
  position: string
  age: number
  overall: number
  potential: number
  marketValue: number
  club: string
  nationality: string
  sofifaId: number
  pac: number
  sho: number
  pas: number
  dri: number
  def: number
  phy: number
}

export interface Prediction {
  growthRate: number
  projectRating: number
  isUndervalued: boolean
  roiScore: number
  badge: string
}

export const players = writable<Player[]>([])
export const watchlistIds = writable<Set<number>>(new Set())
export const predictions = writable<Map<number, Prediction>>(new Map())

export const searchTerm = writable('')
export const minPotential = writable(75)
export const selectedPosition = writable('')
export const showPredictions = writable(false)
export const activeTab = writable<'dashboard' | 'watchlist'>('dashboard')

export const filteredPlayers = derived(
  [players, searchTerm, minPotential, selectedPosition],
  ([$players, $search, $minPot, $position]) =>
    $players.filter(p => {
      if (p.potential < $minPot) return false
      if ($position && p.position !== $position) return false
      if ($search) {
        const s = $search.toLowerCase()
        if (!p.name.toLowerCase().includes(s) && !p.club.toLowerCase().includes(s)) return false
      }
      return true
    })
)

export function calcROI(player: Player): number {
  return player.potential / (player.marketValue / 10_000_000)
}

// ─── Foto del jugador (resolución 240px para mejor calidad)
export function getPhotoUrl(sofifaId: number): string {
  const id = sofifaId.toString()
  return `https://cdn.sofifa.net/players/${id.slice(0, 3)}/${id.slice(3)}/26_240.png`
}

// ─── Emoji de bandera (fallback)
export function getFlagEmoji(nationality: string): string {
  const flags: Record<string, string> = {
    'England': '🏴󠁧󠁢󠁥󠁮󠁧󠁿', 'Norway': '🇳🇴', 'Spain': '🇪🇸', 'Germany': '🇩🇪',
    'Netherlands': '🇳🇱', 'Switzerland': '🇨🇭', 'Portugal': '🇵🇹',
    'France': '🇫🇷', 'Brazil': '🇧🇷', 'Croatia': '🇭🇷', 'Argentina': '🇦🇷',
    'Morocco': '🇲🇦', 'Belgium': '🇧🇪', 'USA': '🇺🇸', 'Italy': '🇮🇹',
    'Turkey': '🇹🇷', 'Georgia': '🇬🇪', 'Ireland': '🇮🇪', 'Poland': '🇵🇱',
    'Japan': '🇯🇵', 'Colombia': '🇨🇴', 'Uruguay': '🇺🇾', 'Paraguay': '🇵🇾',
    'Nigeria': '🇳🇬', 'Egypt': '🇪🇬', 'Serbia': '🇷🇸',
  }
  return flags[nationality] ?? '🌍'
}

// ─── URL de bandera real (flagcdn.com — imagen rectangular ISO)
export function getFlagUrl(nationality: string): string {
  const codes: Record<string, string> = {
    'England': 'gb-eng', 'Norway': 'no', 'Spain': 'es', 'Germany': 'de',
    'Netherlands': 'nl', 'Switzerland': 'ch', 'Portugal': 'pt', 'France': 'fr',
    'Brazil': 'br', 'Croatia': 'hr', 'Argentina': 'ar', 'Morocco': 'ma',
    'Belgium': 'be', 'USA': 'us', 'Italy': 'it', 'Sweden': 'se',
    'Denmark': 'dk', 'Austria': 'at', 'Turkey': 'tr', 'Georgia': 'ge',
    'Ireland': 'ie', 'Poland': 'pl', 'Japan': 'jp', 'Colombia': 'co',
    'Uruguay': 'uy', 'Paraguay': 'py',
    'Nigeria': 'ng', 'Egypt': 'eg', 'Serbia': 'rs',
  }
  const code = codes[nationality]
  return code ? `https://flagcdn.com/w40/${code}.png` : ''
}

// ─── Alias de nombres cortos de clubs
function resolveClub(club: string): string {
  const aliases: Record<string, string> = {
    'Man City':    'Manchester City',
    'Man United':  'Manchester United',
    'Bayern':      'Bayern Munich',
    'Leverkusen':  'Bayer Leverkusen',
    'Dortmund':    'Borussia Dortmund',
    'Atletico':    'Atletico Madrid',
    'Inter':       'Inter Milan',
    'Milan':       'AC Milan',
    'Spurs':       'Tottenham',
  }
  return aliases[club] ?? club
}

// ─── Club: mapeo a abreviatura para badge CSS (más fiable que CDN externo)
export function getClubInfo(club: string): { abbr: string; color: string } {
  const map: Record<string, { abbr: string; color: string }> = {
    'Real Madrid':        { abbr: 'RM',  color: '#ffffff' },
    'Manchester City':    { abbr: 'MCI', color: '#6CABDD' },
    'Barcelona':          { abbr: 'FCB', color: '#A50044' },
    'Arsenal':            { abbr: 'ARS', color: '#EF0107' },
    'Bayern Munich':      { abbr: 'FCB', color: '#DC052D' },
    'Bayer Leverkusen':   { abbr: 'B04', color: '#E32221' },
    'PSG':                { abbr: 'PSG', color: '#003170' },
    'Paris Saint-Germain':{ abbr: 'PSG', color: '#003170' },
    'Liverpool':          { abbr: 'LFC', color: '#C8102E' },
    'Borussia Dortmund':  { abbr: 'BVB', color: '#FDE100' },
    'Atletico Madrid':    { abbr: 'ATM', color: '#CB3524' },
    'Chelsea':            { abbr: 'CHE', color: '#034694' },
    'Juventus':           { abbr: 'JUV', color: '#000000' },
    'Inter Milan':        { abbr: 'INT', color: '#0068A8' },
    'AC Milan':           { abbr: 'ACM', color: '#FB090B' },
    'Tottenham':          { abbr: 'TOT', color: '#132257' },
    'Manchester United':  { abbr: 'MUN', color: '#DA020E' },
    'Napoli':             { abbr: 'NAP', color: '#12A0C8' },
    'Benfica':            { abbr: 'SLB', color: '#EE1B24' },
    'Porto':              { abbr: 'FCP', color: '#003087' },
    'Ajax':               { abbr: 'AFC', color: '#C8102E' },
    'RB Leipzig':         { abbr: 'RBL', color: '#CC0A2F' },
  }
  return map[resolveClub(club)] ?? { abbr: club.slice(0, 3).toUpperCase(), color: '#888888' }
}

// ─── Liga: abreviatura y color para badge
export function getLeagueInfo(club: string): { abbr: string; color: string } {
  const c = resolveClub(club)
  const laLiga     = ['Real Madrid', 'Barcelona', 'Atletico Madrid', 'Villarreal', 'Sevilla', 'Athletic Bilbao']
  const premier    = ['Manchester City', 'Liverpool', 'Arsenal', 'Chelsea', 'Tottenham', 'Manchester United', 'Newcastle', 'Aston Villa', 'Brighton']
  const bundesliga = ['Bayern Munich', 'Bayer Leverkusen', 'Borussia Dortmund', 'RB Leipzig']
  const ligue1     = ['PSG', 'Paris Saint-Germain', 'Monaco', 'Lyon', 'Marseille']
  const serieA     = ['Juventus', 'Inter Milan', 'AC Milan', 'Napoli', 'Roma', 'Atalanta']
  const ligaPor    = ['Benfica', 'Porto', 'Sporting CP']
  const eredivisie = ['Ajax', 'PSV', 'Feyenoord']

  if (laLiga.includes(c))      return { abbr: 'LL', color: '#FF4B00' }
  if (premier.includes(c))     return { abbr: 'PL', color: '#3D195B' }
  if (bundesliga.includes(c))  return { abbr: 'BL', color: '#D3010C' }
  if (ligue1.includes(c))      return { abbr: 'L1', color: '#0033A0' }
  if (serieA.includes(c))      return { abbr: 'SA', color: '#024494' }
  if (ligaPor.includes(c))     return { abbr: 'LP', color: '#00A859' }
  if (eredivisie.includes(c))  return { abbr: 'ER', color: '#D40000' }
  return { abbr: '??', color: '#666666' }
}

// ─── Logo del club (ESPN CDN — acceso público)
export function getClubLogoUrl(club: string): string {
  const ids: Record<string, number> = {
    'Real Madrid':         86,  'Manchester City':    382,
    'Barcelona':           83,  'Arsenal':            359,
    'Bayern Munich':      132,  'Bayer Leverkusen':   126,
    'PSG':                160,  'Paris Saint-Germain':160,
    'Liverpool':          364,  'Borussia Dortmund':  124,
    'Atletico Madrid':   1068,  'Chelsea':            363,
    'Juventus':           111,  'Inter Milan':        110,
    'AC Milan':           103,  'Tottenham':          367,
    'Manchester United':  360,  'Napoli':             106,
    'Benfica':            210,  'Porto':              211,
    'Ajax':               169,  'RB Leipzig':       11109,
    'Atalanta':           122,  'Galatasaray':      383,
    'Aston Villa':        362,
  }
  const id = ids[resolveClub(club)]
  return id ? `https://a.espncdn.com/i/teamlogos/soccer/500/${id}.png` : ''
}

// ─── Logo de la liga (ESPN CDN — acceso público)
export function getLeagueLogoUrl(club: string): string {
  const leagueIds: Record<string, number> = {
    'Manchester City': 23,  'Liverpool': 23,       'Arsenal': 23,
    'Chelsea': 23,          'Tottenham': 23,        'Manchester United': 23,
    'Newcastle': 23,        'Aston Villa': 23,
    'Real Madrid': 15480,   'Barcelona': 15480,     'Atletico Madrid': 15480,
    'Bayern Munich': 10,    'Bayer Leverkusen': 10, 'Borussia Dortmund': 10, 'RB Leipzig': 10,
    'PSG': 9,               'Paris Saint-Germain': 9,
    'Juventus': 12,         'Inter Milan': 12,      'AC Milan': 12,
    'Napoli': 12,           'Roma': 12,        'Atalanta': 12,
    'Benfica': 197,         'Porto': 197,           'Sporting CP': 197,
    'Ajax': 11,             'PSV': 11,              'Feyenoord': 11,
  }
  const id = leagueIds[resolveClub(club)]
  return id ? `https://a.espncdn.com/i/leaguelogos/soccer/500/${id}.png` : ''
}

export function getCardTier(potential: number): 'diamond' | 'gold' | 'silver' {
  if (potential >= 95) return 'diamond'
  if (potential >= 88) return 'gold'
  return 'silver'
}
