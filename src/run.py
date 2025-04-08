# run.py
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent / "src"))

from dashboard.app import Dashboard

if __name__ == "__main__":
    Dashboard().run()
