from .dashboard.app import Dashboard
import asyncio

if __name__ == "__main__":
    asyncio.run(Dashboard().run())
