import os
from pathlib import Path

REPO_HUB = Path(os.path.expanduser(os.getenv("REPO_HUB", "~/Neonware")))

print(f"Using repository hub: {REPO_HUB}")

try:
    if not REPO_HUB.exists():
        raise FileNotFoundError(f"Repository hub '{REPO_HUB}' does not exist.")
    if not REPO_HUB.is_dir():
        raise NotADirectoryError(f"'{REPO_HUB}' is not a directory.")
    REPOS = sorted([p for p in REPO_HUB.iterdir() if p.is_dir()], key=lambda p: p.name)

except Exception as e:
    REPOS = []
    print(f"Warning: {e}")
