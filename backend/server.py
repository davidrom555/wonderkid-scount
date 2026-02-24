from http.server import HTTPServer, BaseHTTPRequestHandler
import json
import math
from urllib.parse import urlparse, parse_qs

PLAYERS = [
    {"id": 1,  "name": "Jude Bellingham", "position": "CAM", "age": 20, "overall": 88, "potential": 94, "marketValue": 180000000, "club": "Real Madrid",  "nationality": "England",     "sofifaId": 252371, "pac": 82, "sho": 82, "pas": 87, "dri": 90, "def": 68, "phy": 80},
    {"id": 2,  "name": "E. Haaland",       "position": "ST",  "age": 23, "overall": 91, "potential": 95, "marketValue": 200000000, "club": "Man City",      "nationality": "Norway",      "sofifaId": 239085, "pac": 89, "sho": 94, "pas": 66, "dri": 82, "def": 44, "phy": 88},
    {"id": 3,  "name": "Gavi",             "position": "CM",  "age": 19, "overall": 85, "potential": 92, "marketValue":  90000000, "club": "Barcelona",     "nationality": "Spain",       "sofifaId": 264240, "pac": 71, "sho": 74, "pas": 87, "dri": 88, "def": 76, "phy": 72},
    {"id": 4,  "name": "Pedri",            "position": "CAM", "age": 21, "overall": 86, "potential": 93, "marketValue": 120000000, "club": "Barcelona",     "nationality": "Spain",       "sofifaId": 251854, "pac": 73, "sho": 76, "pas": 88, "dri": 91, "def": 71, "phy": 71},
    {"id": 5,  "name": "J. Musiala",       "position": "CAM", "age": 20, "overall": 86, "potential": 93, "marketValue": 130000000, "club": "Bayern",        "nationality": "Germany",     "sofifaId": 256790, "pac": 83, "sho": 78, "pas": 82, "dri": 90, "def": 62, "phy": 68},
    {"id": 6,  "name": "B. Saka",          "position": "RW",  "age": 22, "overall": 87, "potential": 91, "marketValue": 150000000, "club": "Arsenal",       "nationality": "England",     "sofifaId": 246669, "pac": 90, "sho": 83, "pas": 84, "dri": 88, "def": 59, "phy": 71},
    {"id": 7,  "name": "V. Koopmeiners",   "position": "CM",  "age": 25, "overall": 84, "potential": 88, "marketValue":  70000000, "club": "Juventus",      "nationality": "Netherlands", "sofifaId": 240679, "pac": 74, "sho": 80, "pas": 84, "dri": 82, "def": 74, "phy": 78},
    {"id": 8,  "name": "M. Akanji",        "position": "CB",  "age": 28, "overall": 87, "potential": 89, "marketValue":  65000000, "club": "Man City",      "nationality": "Switzerland", "sofifaId": 229237, "pac": 78, "sho": 49, "pas": 71, "dri": 68, "def": 87, "phy": 84},
    {"id": 9,  "name": "R. Dias",          "position": "CB",  "age": 26, "overall": 88, "potential": 91, "marketValue":  85000000, "club": "Man City",      "nationality": "Portugal",    "sofifaId": 239818, "pac": 77, "sho": 46, "pas": 72, "dri": 70, "def": 89, "phy": 85},
    {"id": 10, "name": "L. Yamal",         "position": "RW",  "age": 17, "overall": 82, "potential": 96, "marketValue": 180000000, "club": "Barcelona",     "nationality": "Spain",       "sofifaId": 277643, "pac": 92, "sho": 77, "pas": 80, "dri": 89, "def": 33, "phy": 58},
    {"id": 11, "name": "W. Foden",         "position": "CAM", "age": 23, "overall": 87, "potential": 92, "marketValue": 140000000, "club": "Man City",      "nationality": "England",     "sofifaId": 209658, "pac": 81, "sho": 83, "pas": 87, "dri": 89, "def": 58, "phy": 69},
    {"id": 12, "name": "K. Mbappé",        "position": "ST",  "age": 25, "overall": 92, "potential": 94, "marketValue": 220000000, "club": "Real Madrid",   "nationality": "France",      "sofifaId": 231747, "pac": 97, "sho": 91, "pas": 82, "dri": 93, "def": 45, "phy": 78},
    {"id": 13, "name": "V. Junior",        "position": "LW",  "age": 23, "overall": 89, "potential": 92, "marketValue": 200000000, "club": "Real Madrid",   "nationality": "Brazil",      "sofifaId": 238794, "pac": 95, "sho": 84, "pas": 79, "dri": 93, "def": 32, "phy": 70},
    {"id": 14, "name": "J. Sancho",        "position": "LW",  "age": 24, "overall": 84, "potential": 87, "marketValue":  40000000, "club": "Man United",    "nationality": "England",     "sofifaId": 233049, "pac": 83, "sho": 78, "pas": 80, "dri": 86, "def": 44, "phy": 63},
    {"id": 15, "name": "C. Gakpo",         "position": "LW",  "age": 24, "overall": 86, "potential": 91, "marketValue":  75000000, "club": "Liverpool",     "nationality": "Netherlands", "sofifaId": 242516, "pac": 84, "sho": 82, "pas": 79, "dri": 84, "def": 54, "phy": 72},
    {"id": 16, "name": "F. De Jong",       "position": "CDM", "age": 26, "overall": 86, "potential": 90, "marketValue":  80000000, "club": "Barcelona",     "nationality": "Netherlands", "sofifaId": 228702, "pac": 74, "sho": 71, "pas": 87, "dri": 86, "def": 78, "phy": 74},
]

watchlist: set[int] = set()


def calc_prediction(player: dict) -> dict:
    age_factor = (player["age"] - 15) * 0.8
    if age_factor <= 0:
        age_factor = 0.8
    growth_rate = (player["potential"] - player["overall"]) / age_factor
    project_rating = min(int(player["overall"] + growth_rate * age_factor), 99)
    roi_score = player["potential"] / (player["marketValue"] / 10_000_000)
    is_undervalued = roi_score > 1.2
    badge = "🔥 UNDERVALUED TARGET" if roi_score > 1.2 else ("📊 FAIR VALUE" if roi_score > 0.8 else "💎 PREMIUM ASSET")
    return {
        "growthRate": round(growth_rate, 2),
        "projectRating": project_rating,
        "isUndervalued": is_undervalued,
        "roiScore": round(roi_score, 2),
        "badge": badge,
    }


class Handler(BaseHTTPRequestHandler):
    def log_message(self, format, *args):
        print(f"[{self.command}] {self.path}")

    def send_json(self, data, status=200):
        body = json.dumps(data, ensure_ascii=False).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type")
        self.end_headers()
        self.wfile.write(body)

    def do_OPTIONS(self):
        self.send_response(204)
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type")
        self.end_headers()

    def do_GET(self):
        parsed = urlparse(self.path)
        path = parsed.path.rstrip("/")
        qs = parse_qs(parsed.query)

        if path == "/api/players":
            return self.send_json(PLAYERS)

        if path == "/api/players/filter":
            min_pot = int(qs.get("minPotential", ["0"])[0])
            position = qs.get("position", [""])[0]
            search = qs.get("search", [""])[0].lower()
            result = [
                p for p in PLAYERS
                if p["potential"] >= min_pot
                and (not position or p["position"] == position)
                and (not search or search in p["name"].lower() or search in p["club"].lower())
            ]
            return self.send_json(result)

        if path.startswith("/api/predict/"):
            pid = int(path.split("/")[-1])
            player = next((p for p in PLAYERS if p["id"] == pid), None)
            if not player:
                return self.send_json({"error": "Player not found"}, 404)
            return self.send_json(calc_prediction(player))

        if path == "/api/watchlist":
            wl_players = [p for p in PLAYERS if p["id"] in watchlist]
            return self.send_json({"players": wl_players, "count": len(wl_players)})

        self.send_json({"error": "Not found"}, 404)

    def do_POST(self):
        path = urlparse(self.path).path.rstrip("/")
        if path.startswith("/api/watchlist/"):
            pid = int(path.split("/")[-1])
            if any(p["id"] == pid for p in PLAYERS):
                watchlist.add(pid)
                return self.send_json({"success": True})
            return self.send_json({"error": "Player not found"}, 404)
        self.send_json({"error": "Not found"}, 404)

    def do_DELETE(self):
        path = urlparse(self.path).path.rstrip("/")
        if path.startswith("/api/watchlist/"):
            pid = int(path.split("/")[-1])
            watchlist.discard(pid)
            return self.send_json({"success": True})
        self.send_json({"error": "Not found"}, 404)


if __name__ == "__main__":
    server = HTTPServer(("0.0.0.0", 8080), Handler)
    print("=" * 45)
    print("  FIFA WONDERKID SCOUT — Backend v2")
    print(f"  Players: {len(PLAYERS)} wonderkids loaded")
    print("  API:     http://localhost:8080/api")
    print("=" * 45)
    server.serve_forever()
