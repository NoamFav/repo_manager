from pathlib import Path
import subprocess
import time
import os
from typing import Dict, List, Optional, Any, Tuple
import asyncio

from rich.console import Console
from rich.layout import Layout
from rich.panel import Panel
from rich.table import Table
from rich.text import Text
from rich.box import ROUNDED, HEAVY, DOUBLE, SIMPLE
from rich.prompt import Prompt, Confirm
from rich.live import Live
from rich import print as rprint
from rich.columns import Columns
from rich.spinner import Spinner
from rich.progress import (
    Progress,
    SpinnerColumn,
    TextColumn,
    BarColumn,
    TaskID,
    TimeElapsedColumn,
)
from rich.align import Align
from rich.style import Style
from rich.rule import Rule
from rich.syntax import Syntax
from rich.markdown import Markdown
from rich.tree import Tree
from rich.emoji import Emoji

# Import your REPOS list
from .. import REPOS, REPO_HUB


class Theme:
    """Theme class to maintain consistent colors and styles throughout the app."""

    # Base colors
    PRIMARY = "deep_sky_blue1"
    SECONDARY = "cyan"
    ACCENT = "bright_magenta"
    SUCCESS = "green3"
    WARNING = "gold1"
    ERROR = "red"
    INFO = "steel_blue1"

    # Background colors
    BG_DARK = "grey15"
    BG_MEDIUM = "grey23"
    BG_LIGHT = "grey30"

    # Text colors
    TEXT_BRIGHT = "white"
    TEXT_NORMAL = "grey85"
    TEXT_MUTED = "grey66"
    TEXT_DIM = "grey50"

    # Styles for specific elements
    HEADER = Style(color=TEXT_BRIGHT, bgcolor=PRIMARY)
    SIDEBAR = Style(color=TEXT_NORMAL, bgcolor=BG_MEDIUM)
    SIDEBAR_SELECTED = Style(color=TEXT_BRIGHT, bgcolor=SECONDARY)

    # Box styles
    BOX_HEADER = ROUNDED
    BOX_CONTENT = SIMPLE
    BOX_FOCUSED = DOUBLE

    # Panel styles
    @classmethod
    def panel_header(cls) -> Dict[str, Any]:
        return {"box": cls.BOX_HEADER, "border_style": cls.PRIMARY, "padding": (1, 2)}

    @classmethod
    def panel_sidebar(cls) -> Dict[str, Any]:
        return {
            "box": cls.BOX_CONTENT,
            "border_style": cls.SECONDARY,
            "padding": (1, 1),
        }

    @classmethod
    def panel_content(cls) -> Dict[str, Any]:
        return {"box": cls.BOX_CONTENT, "border_style": cls.INFO, "padding": (1, 2)}

    @classmethod
    def panel_action(cls) -> Dict[str, Any]:
        return {"box": cls.BOX_CONTENT, "border_style": cls.ACCENT, "padding": (1, 2)}


class RepoCard:
    """Class to format and display repository information."""

    def __init__(self, repo_path: Path, is_selected: bool = False):
        self.repo_path = repo_path
        self.is_selected = is_selected
        self.name = repo_path.name
        self.info = {}
        self.stats = {}
        self._fetch_info()

    def _fetch_info(self) -> None:
        """Fetch repository information."""
        try:
            self.stats = self._get_repo_stats()
            self.info = self._get_repo_info()
        except Exception as e:
            self.error = str(e)

    def _get_repo_info(self) -> Dict[str, Any]:
        """Get the information of the repository using onefetch."""
        try:
            # Run onefetch with --no-art flag to avoid ASCII art issues in Rich layouts
            onefetch_info = subprocess.run(
                ["onefetch", "--no-art", str(self.repo_path)],
                capture_output=True,
                text=True,
                timeout=5,
            )

            # Optional: Get only the ASCII art for separate display
            onefetch_art = subprocess.run(
                ["onefetch", "--disabled-fields=all", str(self.repo_path)],
                capture_output=True,
                text=True,
                timeout=5,
            )

            return {
                "name": self.name,
                "path": str(self.repo_path),
                "info": onefetch_info.stdout,
                "art": onefetch_art.stdout if onefetch_art.returncode == 0 else "",
                "error": None,
            }
        except (subprocess.SubprocessError, subprocess.TimeoutExpired) as e:
            return {
                "name": self.name,
                "path": str(self.repo_path),
                "info": "",
                "art": "",
                "error": str(e),
            }

    def _get_repo_stats(self) -> Dict[str, Any]:
        """Get additional repository statistics."""
        if not self.repo_path.exists():
            return {"error": "Repository path does not exist"}

        try:
            # Get commit count
            commit_count = subprocess.run(
                ["git", "-C", str(self.repo_path), "rev-list", "--count", "HEAD"],
                capture_output=True,
                text=True,
                timeout=2,
            ).stdout.strip()

            # Get last commit date
            last_commit = subprocess.run(
                [
                    "git",
                    "-C",
                    str(self.repo_path),
                    "log",
                    "-1",
                    "--format=%cd",
                    "--date=relative",
                ],
                capture_output=True,
                text=True,
                timeout=2,
            ).stdout.strip()

            # Get branch count
            branch_count = (
                subprocess.run(
                    ["git", "-C", str(self.repo_path), "branch", "--list"],
                    capture_output=True,
                    text=True,
                    timeout=2,
                ).stdout.count("\n")
                + 1
            )

            # Get file count - using a more optimized approach with .gitignore respected
            try:
                file_count = int(
                    subprocess.run(
                        [
                            "git",
                            "-C",
                            str(self.repo_path),
                            "ls-files",
                            "--exclude-standard",
                            "|",
                            "wc",
                            "-l",
                        ],
                        capture_output=True,
                        text=True,
                        shell=True,
                        timeout=3,
                    ).stdout.strip()
                )
            except:
                # Fallback method if git ls-files fails
                file_count = sum(1 for _ in self.repo_path.glob("**/*"))

            # Get languages breakdown
            languages = {}
            try:
                lang_output = subprocess.run(
                    ["tokei", str(self.repo_path), "--output", "json"],
                    capture_output=True,
                    text=True,
                    timeout=3,
                )
                if lang_output.returncode == 0:
                    import json

                    langs_data = json.loads(lang_output.stdout)
                    total_lines = 0
                    for lang, info in langs_data.items():
                        if lang != "Total" and isinstance(info, dict):
                            languages[lang] = info.get("code", 0)
                            total_lines += info.get("code", 0)

                    # Convert to percentages
                    if total_lines > 0:
                        languages = {
                            lang: round(lines / total_lines * 100, 1)
                            for lang, lines in languages.items()
                        }
            except:
                # Fallback if tokei is not installed
                languages = {"unknown": 100}

            return {
                "commit_count": commit_count,
                "last_commit": last_commit,
                "branch_count": branch_count,
                "file_count": file_count,
                "languages": languages,
                "error": None,
            }
        except Exception as e:
            return {"error": str(e)}

    def to_panel(self) -> Panel:
        """Convert the repository info to a rich Panel."""
        # Create a table for the stats
        stats_table = Table(show_header=False, box=None, padding=(0, 1), expand=True)
        stats_table.add_column("Stat", style=Theme.TEXT_MUTED)
        stats_table.add_column("Value", style=Theme.TEXT_BRIGHT)

        # Add repository stats to the table
        if self.stats.get("error") is None:
            stats_table.add_row("Commits", str(self.stats.get("commit_count", "N/A")))
            stats_table.add_row(
                "Last activity", str(self.stats.get("last_commit", "N/A"))
            )
            stats_table.add_row("Branches", str(self.stats.get("branch_count", "N/A")))
            stats_table.add_row("Files", str(self.stats.get("file_count", "N/A")))

            # Language breakdown
            languages = self.stats.get("languages", {})
            if languages and not isinstance(languages, str):
                lang_text = ""
                for lang, percentage in sorted(
                    languages.items(), key=lambda x: x[1], reverse=True
                )[:3]:
                    lang_text += f"{lang}: {percentage}% "
                if len(languages) > 3:
                    lang_text += "..."
                stats_table.add_row("Languages", lang_text)
        else:
            stats_table.add_row("Error", self.stats.get("error", "Unknown error"))

        # Create the main content
        content = Columns(
            [
                Align(
                    Text(
                        f"{Emoji('file_folder')} {self.name}",
                        style=Theme.TEXT_BRIGHT,
                        justify="center",
                    ),
                    align="center",
                    vertical="middle",
                ),
                stats_table,
            ]
        )

        # Select the appropriate border style based on selection status
        border_style = Theme.SECONDARY if self.is_selected else Theme.TEXT_MUTED

        return Panel(
            content,
            title=self.name,
            border_style=border_style,
            box=ROUNDED,
            padding=(1, 2),
            title_align="left",
            highlight=True,
        )

    def to_detailed_view(self) -> Panel:
        """Create a detailed view panel for the repository."""
        content = []

        # Header with repository name
        content.append(
            Align(
                Text(
                    f"{Emoji('file_folder')} {self.name}",
                    style=f"bold {Theme.PRIMARY}",
                    justify="center",
                ),
                align="center",
            )
        )
        content.append(Rule(style=Theme.PRIMARY))

        # Repository path
        path_text = Text()
        path_text.append("Path: ", style=f"bold {Theme.TEXT_MUTED}")
        path_text.append(str(self.repo_path), style=Theme.TEXT_NORMAL)
        content.append(path_text)

        # Repository stats in a table
        if self.stats.get("error") is None:
            stats_table = Table(
                show_header=True,
                box=SIMPLE,
                title="Statistics",
                title_style=f"bold {Theme.INFO}",
                expand=True,
            )
            stats_table.add_column("Metric", style=Theme.TEXT_MUTED)
            stats_table.add_column("Value", style=Theme.TEXT_BRIGHT)

            stats_table.add_row("Commits", str(self.stats.get("commit_count", "N/A")))
            stats_table.add_row(
                "Last activity", str(self.stats.get("last_commit", "N/A"))
            )
            stats_table.add_row("Branches", str(self.stats.get("branch_count", "N/A")))
            stats_table.add_row("Files", str(self.stats.get("file_count", "N/A")))

            content.append(stats_table)

            # Language breakdown
            languages = self.stats.get("languages", {})
            if languages and not isinstance(languages, str):
                lang_table = Table(
                    show_header=True,
                    box=SIMPLE,
                    title="Languages",
                    title_style=f"bold {Theme.SECONDARY}",
                    expand=True,
                )
                lang_table.add_column("Language", style=Theme.TEXT_NORMAL)
                lang_table.add_column("Percentage", style=Theme.TEXT_BRIGHT)

                for lang, percentage in sorted(
                    languages.items(), key=lambda x: x[1], reverse=True
                )[:5]:
                    lang_table.add_row(lang, f"{percentage}%")

                content.append(lang_table)

        # Add ASCII art if available
        if self.info.get("art"):
            art_text = Text(self.info.get("art"), style=Theme.TEXT_BRIGHT)
            content.append(Align(art_text, align="center"))

        # Onefetch info
        if self.info.get("info"):
            onefetch_info = Text(self.info.get("info"), style=Theme.TEXT_NORMAL)
            content.append(onefetch_info)

        # Wrap everything in a panel
        return Panel(
            Columns([Align(item, align="center") for item in content], equal=True),
            title=f"Repository Details: {self.name}",
            **Theme.panel_content(),
        )


class Dashboard:
    """Enhanced repository dashboard with a modern UI."""

    def __init__(self):
        self.console = Console()
        self.selected_repo_index = 0
        self.main_menu_index = 0
        self.main_menu = [
            ("Repos", "Manage Git repositories"),
            ("Actions", "Perform repository actions"),
            ("Tasks", "Schedule automated tasks"),
            ("Settings", "Configure dashboard"),
            ("Exit", "Close the dashboard"),
        ]

        self.action_menu = [
            ("Create Repo", "Create a new Git repository"),
            ("Clean Repo", "Clean up unnecessary files"),
            ("Modify with GH", "Modify using GitHub CLI"),
            ("AI Commit", "Generate commit messages with AI"),
            ("Back", "Return to main menu"),
        ]

        self.task_menu = [
            ("Schedule Backup", "Schedule regular backups"),
            ("Schedule Sync", "Sync repositories automatically"),
            ("Run Linters", "Run code quality tools"),
            ("Back", "Return to main menu"),
        ]

        self.settings_menu = [
            ("Appearance", "Customize dashboard appearance"),
            ("Notifications", "Configure notification settings"),
            ("Git Settings", "Configure global Git settings"),
            ("Back", "Return to main menu"),
        ]

        self.current_view = "main"
        self.layout = self._make_layout()
        self.running = True
        self.repo_cards = []
        self._update_repo_cards()

        # Status message display
        self.status_message = Text("Ready", style=Theme.TEXT_MUTED)
        self.status_timestamp = time.time()

    def _update_repo_cards(self):
        """Update repository cards."""
        self.repo_cards = []
        for i, repo in enumerate(REPOS):
            self.repo_cards.append(RepoCard(repo, i == self.selected_repo_index))

    def _make_layout(self) -> Layout:
        """Create the layout for the dashboard."""
        layout = Layout(name="root")

        # Split the main layout into header, body, and footer
        layout.split(
            Layout(name="header", size=3),
            Layout(name="body", ratio=1),
            Layout(name="footer", size=3),
        )

        # Split the body into left sidebar and right content
        layout["body"].split_row(
            Layout(name="sidebar", ratio=1),
            Layout(name="content", ratio=3),
        )

        # Split the content into info and actions
        layout["content"].split(
            Layout(name="info", ratio=2),
            Layout(name="actions", ratio=1),
        )

        return layout

    def _render_header(self) -> Panel:
        """Render the header with the main menu."""
        header_text = Text()

        # Add title
        header_text.append("🚀 ", style=Theme.TEXT_BRIGHT)
        header_text.append("Repository Dashboard", style=f"bold {Theme.TEXT_BRIGHT}")
        header_text.append(" | ", style=Theme.TEXT_DIM)

        # Add menu items
        for i, (item, desc) in enumerate(self.main_menu):
            if i == self.main_menu_index and self.current_view == "main":
                header_text.append(
                    f"{item}", style=f"bold {Theme.PRIMARY} on {Theme.BG_LIGHT}"
                )
            else:
                header_text.append(f"{item}", style=Theme.TEXT_NORMAL)

            if i < len(self.main_menu) - 1:
                header_text.append(" • ", style=Theme.TEXT_DIM)

        return Panel(
            Align(header_text, align="center", vertical="middle"),
            **Theme.panel_header(),
            title="[bold]Git Repository Manager[/bold]",
            subtitle=f"Repositories: {len(REPOS)}",
        )

    def _render_sidebar(self) -> Panel:
        """Render the sidebar with repository list."""
        if not REPOS:
            return Panel(
                Align(
                    Text("No repositories found", style=Theme.TEXT_DIM),
                    align="center",
                    vertical="middle",
                ),
                title="Repositories",
                **Theme.panel_sidebar(),
            )

        # Create a list of repository cards
        repo_panels = []
        for i, repo_card in enumerate(self.repo_cards):
            repo_card.is_selected = i == self.selected_repo_index
            repo_panels.append(repo_card.to_panel())

        # Group them in columns for display
        repo_list = Columns(repo_panels, equal=True)

        return Panel(
            repo_list, title=f"Repositories ({len(REPOS)})", **Theme.panel_sidebar()
        )

    def _render_info(self) -> Panel:
        """Render repository information."""
        if not REPOS:
            return Panel(
                Align(
                    Text("No repositories available", style=Theme.TEXT_DIM),
                    align="center",
                    vertical="middle",
                ),
                title="Repository Information",
                **Theme.panel_content(),
            )

        selected_repo = REPOS[self.selected_repo_index]
        repo_card = self.repo_cards[self.selected_repo_index]

        return repo_card.to_detailed_view()

    def _render_actions(self) -> Panel:
        """Render the actions panel."""
        if self.current_view == "main":
            action_text = Table.grid(padding=(0, 2))
            action_text.add_column("Icon", style=Theme.ACCENT)
            action_text.add_column("Action", style=Theme.TEXT_NORMAL)
            action_text.add_column("Description", style=Theme.TEXT_MUTED)

            action_text.add_row("→", "Select menu item", "Navigate using arrow keys")
            action_text.add_row("↵", "Confirm selection", "Execute selected action")
            action_text.add_row("1-5", "Quick menu access", "Jump to menu section")
            action_text.add_row("q", "Quit application", "Exit the dashboard")

            title = "Available Actions"

        elif self.current_view == "actions":
            action_text = Table.grid(padding=(0, 2))
            action_text.add_column("", style=Theme.TEXT_NORMAL)
            action_text.add_column("Action", style=Theme.TEXT_NORMAL)

            for i, (action, desc) in enumerate(self.action_menu):
                icon = "▶ " if i == self.main_menu_index else "  "
                style = (
                    f"bold {Theme.ACCENT}"
                    if i == self.main_menu_index
                    else Theme.TEXT_NORMAL
                )
                action_text.add_row(icon, Text(action, style=style))

            title = "Repository Actions"

        elif self.current_view == "tasks":
            action_text = Table.grid(padding=(0, 2))
            action_text.add_column("", style=Theme.TEXT_NORMAL)
            action_text.add_column("Task", style=Theme.TEXT_NORMAL)

            for i, (task, desc) in enumerate(self.task_menu):
                icon = "▶ " if i == self.main_menu_index else "  "
                style = (
                    f"bold {Theme.ACCENT}"
                    if i == self.main_menu_index
                    else Theme.TEXT_NORMAL
                )
                action_text.add_row(icon, Text(task, style=style))

            title = "Scheduled Tasks"

        elif self.current_view == "settings":
            action_text = Table.grid(padding=(0, 2))
            action_text.add_column("", style=Theme.TEXT_NORMAL)
            action_text.add_column("Setting", style=Theme.TEXT_NORMAL)

            for i, (setting, desc) in enumerate(self.settings_menu):
                icon = "▶ " if i == self.main_menu_index else "  "
                style = (
                    f"bold {Theme.ACCENT}"
                    if i == self.main_menu_index
                    else Theme.TEXT_NORMAL
                )
                action_text.add_row(icon, Text(setting, style=style))

            title = "Dashboard Settings"

        else:
            action_text = Text(
                "Use arrow keys to navigate, Enter to select", style=Theme.TEXT_MUTED
            )
            title = "Help"

        return Panel(
            Align(action_text, align="left", vertical="top"),
            title=title,
            **Theme.panel_action(),
        )

    def _render_footer(self) -> Panel:
        """Render the footer with status and key bindings."""
        # Create the key binding help text
        key_help = Table.grid(padding=(0, 1))
        key_help.add_column("Key", style=f"bold {Theme.PRIMARY}")
        key_help.add_column("Action", style=Theme.TEXT_MUTED)

        # Vim-style navigation keys
        key_help.add_row("h j k l", "Navigate (Vim-style)")
        key_help.add_row("Enter/Space", "Select")
        key_help.add_row("1-5", "Menu")
        key_help.add_row("q", "Quit")

        # Status message display
        status_age = time.time() - self.status_timestamp
        status_style = Theme.TEXT_MUTED if status_age > 5 else Theme.SUCCESS

        footer_content = Columns(
            [
                Align(
                    self.status_message,
                    align="left",
                    vertical="middle",
                    style=status_style,
                ),
                Align(key_help, align="right", vertical="middle"),
            ]
        )

        return Panel(
            footer_content, border_style=Theme.PRIMARY, box=ROUNDED, padding=(0, 1)
        )

    def set_status(self, message: str, style: str = Theme.SUCCESS):
        """Set a status message to display in the footer."""
        self.status_message = Text(message, style=style)
        self.status_timestamp = time.time()

    def _update_screen(self, live: Live):
        """Update the screen layout."""
        self.layout["header"].update(self._render_header())
        self.layout["sidebar"].update(self._render_sidebar())
        self.layout["info"].update(self._render_info())
        self.layout["actions"].update(self._render_actions())
        self.layout["footer"].update(self._render_footer())
        live.update(self.layout)

    def handle_input(self, key: str):
        """Handle key input with Vim-style hjkl navigation."""
        if key.lower() == "q":
            self.running = False
            return

        # Vim-style keys: k = up, j = down, h = left, l = right
        if key == "k" or key == "up":  # Up navigation
            if self.current_view == "main":
                self.selected_repo_index = max(0, self.selected_repo_index - 1)
                self._update_repo_cards()
            else:
                self.main_menu_index = max(0, self.main_menu_index - 1)
        elif key == "j" or key == "down":  # Down navigation
            if self.current_view == "main":
                self.selected_repo_index = min(
                    len(REPOS) - 1, self.selected_repo_index + 1
                )
                self._update_repo_cards()
            elif self.current_view == "actions":
                self.main_menu_index = min(
                    len(self.action_menu) - 1, self.main_menu_index
                )
            elif self.current_view == "tasks":
                self.main_menu_index = min(
                    len(self.task_menu) - 1, self.main_menu_index
                )
            elif self.current_view == "settings":
                self.main_menu_index = min(
                    len(self.settings_menu) - 1, self.main_menu_index
                )

        # Handle h/l keys for left/right navigation
        elif key == "h" or key == "left":  # Left navigation (go back)
            if self.current_view in ["actions", "tasks", "settings"]:
                # Return to main view
                self.current_view = "main"

                # Set focus to the appropriate main menu item
                if self.current_view == "actions":
                    self.main_menu_index = 1  # Actions menu
                elif self.current_view == "tasks":
                    self.main_menu_index = 2  # Tasks menu
                elif self.current_view == "settings":
                    self.main_menu_index = 3  # Settings menu

        elif key == "l" or key == "right":  # Right navigation (go deeper)
            if self.current_view == "main":
                # Enter the currently selected menu
                menu_item = self.main_menu[self.main_menu_index][0]
                if menu_item == "Actions":
                    self.current_view = "actions"
                    self.main_menu_index = 0
                elif menu_item == "Tasks":
                    self.current_view = "tasks"
                    self.main_menu_index = 0
                elif menu_item == "Settings":
                    self.current_view = "settings"
                    self.main_menu_index = 0

        # Number keys for main menu
        if key in "12345" and self.current_view == "main":
            idx = int(key) - 1
            if 0 <= idx < len(self.main_menu):
                self.main_menu_index = idx
                menu_item = self.main_menu[idx][0]

                if menu_item == "Actions":
                    self.current_view = "actions"
                    self.main_menu_index = 0
                elif menu_item == "Tasks":
                    self.current_view = "tasks"
                    self.main_menu_index = 0
                elif menu_item == "Settings":
                    self.current_view = "settings"
                    self.main_menu_index = 0
                elif menu_item == "Exit":
                    self.running = False

        # Enter key or Space bar handling for selection
        if key == "enter" or key == "space":
            if self.current_view == "main":
                menu_item = self.main_menu[self.main_menu_index][0]

                if menu_item == "Actions":
                    self.current_view = "actions"
                    self.main_menu_index = 0
                elif menu_item == "Tasks":
                    self.current_view = "tasks"
                    self.main_menu_index = 0
                elif menu_item == "Settings":
                    self.current_view = "settings"
                    self.main_menu_index = 0
                elif menu_item == "Exit":
                    self.running = False

            elif self.current_view == "actions":
                action = self.action_menu[self.main_menu_index][0]
                if action == "Back":
                    self.current_view = "main"
                    self.main_menu_index = 1  # Actions menu
                else:
                    # Here you would implement the action functionality
                    self.execute_action(action)

            elif self.current_view == "tasks":
                task = self.task_menu[self.main_menu_index][0]
                if task == "Back":
                    self.current_view = "main"
                    self.main_menu_index = 2  # Tasks menu
                else:
                    # Here you would implement the task scheduling functionality
                    self.schedule_task(task)

            elif self.current_view == "settings":
                setting = self.settings_menu[self.main_menu_index][0]
                if setting == "Back":
                    self.current_view = "main"
                    self.main_menu_index = 3  # Settings menu
                else:
                    # Here you would implement the settings functionality
                    self.modify_setting(setting)

    def execute_action(self, action: str):
        """Execute a repository action."""
        if not REPOS:
            self.set_status("No repositories to perform actions on", Theme.WARNING)
            return

        selected_repo = REPOS[self.selected_repo_index]

        with self.console.screen():
            with Progress(
                SpinnerColumn(spinner_name="dots"),
                TextColumn("[bold {task.fields[style]}]{task.description}"),
                BarColumn(bar_width=40),
                TextColumn("[bold {task.fields[style]}]{task.percentage:>3.0f}%"),
                TimeElapsedColumn(),
                expand=True,
            ) as progress:
                task_id = progress.add_task(
                    f"[{action}] {selected_repo.name}...",
                    total=100,
                    style=Theme.SUCCESS,
                )

                # Simulate work with a animated progress
                for i in range(101):
                    progress.update(task_id, completed=i)
                    time.sleep(0.02)

                    # Update status message as work progresses
                    if i == 25:
                        progress.update(
                            task_id, description=f"[{action}] Analyzing repository..."
                        )
                    elif i == 50:
                        progress.update(
                            task_id, description=f"[{action}] Processing files..."
                        )
                    elif i == 75:
                        progress.update(
                            task_id, description=f"[{action}] Finalizing changes..."
                        )

            # Show results
            if action == "Create Repo":
                result_text = f"""
                ## Repository Creation Complete
                
                **Repository:** {selected_repo.name}
                **Location:** {selected_repo}
                **Template:** Standard Git repository
                
                The repository has been initialized with the following:
                - README.md file
                - .gitignore for common file types
                - MIT License file
                - Initial commit
                
                You can now start adding your project files.
                """
            elif action == "Clean Repo":
                result_text = f"""
                ## Repository Cleaning Complete
                
                **Repository:** {selected_repo.name}
                **Location:** {selected_repo}
                
                The following has been cleaned:
                - Removed 15 temporary files
                - Cleared 3 empty directories
                - Optimized Git storage (saved 7.2MB)
                - Updated .gitignore rules
                
                Your repository is now optimized and ready for use.
                """
            elif action == "Modify with GH":
                result_text = f"""
                ## GitHub CLI Operation Complete
                
                **Repository:** {selected_repo.name}
                **GitHub URL:** https://github.com/user/{selected_repo.name}
                
                The following changes were made:
                - Updated repository description
                - Added 2 collaborators
                - Created development branch protection rule
                - Enabled issue templates
                
                GitHub repository settings have been updated successfully.
                """
            elif action == "AI Commit":
                result_text = f"""
                ## AI Commit Message Generation
                
                **Repository:** {selected_repo.name}
                **Changes analyzed:** 7 files modified
                
                **Generated commit message:**
                ```
                feat(core): implement user authentication system
                
                - Add JWT token-based authentication
                - Create user registration API endpoint
                - Implement password hashing with bcrypt
                - Add unit tests for authentication routes
                ```
                
                Commit message has been applied to your staging area.
                """

            # Display result with rich markdown
            result_md = Markdown(result_text)
            self.console.print(result_md)

            # Ask for user confirmation
            self.console.print()
            self.console.print(
                "[bold green]Operation completed successfully![/bold green]"
            )
            self.console.print("\nPress any key to continue...")
            self.console.input()

        # Update status message
        self.set_status(f"{action} completed for {selected_repo.name}", Theme.SUCCESS)

    def schedule_task(self, task: str):
        """Schedule a repository task."""
        if not REPOS:
            self.set_status("No repositories to schedule tasks for", Theme.WARNING)
            return

        selected_repo = REPOS[self.selected_repo_index]

        # Create a form-like interface for task scheduling
        with self.console.screen():
            self.console.print(
                f"[bold {Theme.PRIMARY}]Schedule Task: {task}[/bold {Theme.PRIMARY}]\n"
            )

            task_config = {}

            if task == "Schedule Backup":
                # Create a form for backup configuration
                task_config["repo"] = selected_repo.name
                task_config["frequency"] = Prompt.ask(
                    "Backup frequency",
                    choices=["daily", "weekly", "monthly"],
                    default="daily",
                )
                task_config["time"] = Prompt.ask("Time of day (HH:MM)", default="03:00")
                task_config["location"] = Prompt.ask(
                    "Backup location", default=f"{REPO_HUB.parent}/backups"
                )
                task_config["compress"] = Confirm.ask("Compress backup?", default=True)
                task_config["retain"] = Prompt.ask(
                    "Backups to retain",
                    default="5",
                    choices=["1", "3", "5", "10", "all"],
                )

                # Summarize the task details
                self.console.print("\n[bold]Task Configuration:[/bold]")

                details_table = Table(box=SIMPLE, show_header=False, expand=True)
                details_table.add_column("Property", style=Theme.TEXT_MUTED)
                details_table.add_column("Value", style=Theme.TEXT_BRIGHT)

                for key, value in task_config.items():
                    details_table.add_row(key.capitalize(), str(value))

                self.console.print(details_table)

            elif task == "Schedule Sync":
                # Create a form for sync configuration
                task_config["repos"] = Prompt.ask(
                    "Repositories to sync",
                    choices=["all", "selected", f"only {selected_repo.name}"],
                    default=f"only {selected_repo.name}",
                )
                task_config["target"] = Prompt.ask(
                    "Sync target",
                    choices=["origin", "upstream", "custom"],
                    default="origin",
                )
                if task_config["target"] == "custom":
                    task_config["remote_url"] = Prompt.ask("Remote URL")

                task_config["schedule"] = Prompt.ask(
                    "Schedule",
                    choices=["hourly", "daily", "on commit", "manual"],
                    default="daily",
                )
                task_config["push"] = Confirm.ask("Push local changes?", default=True)
                task_config["pull"] = Confirm.ask("Pull remote changes?", default=True)

                # Summarize the task details
                self.console.print("\n[bold]Task Configuration:[/bold]")

                details_table = Table(box=SIMPLE, show_header=False, expand=True)
                details_table.add_column("Property", style=Theme.TEXT_MUTED)
                details_table.add_column("Value", style=Theme.TEXT_BRIGHT)

                for key, value in task_config.items():
                    details_table.add_row(key.capitalize(), str(value))

                self.console.print(details_table)

            elif task == "Run Linters":
                # Create a form for linter configuration
                available_linters = ["flake8", "pylint", "black", "isort", "mypy"]

                self.console.print("\n[bold]Available Linters:[/bold]")
                for linter in available_linters:
                    self.console.print(f"- {linter}")

                task_config["linters"] = Prompt.ask(
                    "\nLinter type (comma-separated)",
                    default="flake8,black",
                )
                task_config["schedule"] = Prompt.ask(
                    "Schedule",
                    choices=["on commit", "daily", "manual"],
                    default="on commit",
                )
                task_config["auto_fix"] = Confirm.ask("Auto-fix issues?", default=True)
                task_config["fail_on_error"] = Confirm.ask(
                    "Fail on lint errors?", default=False
                )

                # Summarize the task details
                self.console.print("\n[bold]Task Configuration:[/bold]")

                details_table = Table(box=SIMPLE, show_header=False, expand=True)
                details_table.add_column("Property", style=Theme.TEXT_MUTED)
                details_table.add_column("Value", style=Theme.TEXT_BRIGHT)

                for key, value in task_config.items():
                    details_table.add_row(key.capitalize(), str(value))

                self.console.print(details_table)

            # Confirm task scheduling
            self.console.print()
            confirm = Confirm.ask("\nConfirm task scheduling?", default=True)

            if confirm:
                # Show a spinner while "scheduling"
                with Progress(
                    SpinnerColumn(spinner_name="dots"),
                    TextColumn("[bold green]Scheduling task..."),
                    TimeElapsedColumn(),
                ) as progress:
                    task_id = progress.add_task("Scheduling...", total=None)
                    # Simulate scheduling delay
                    time.sleep(2)

                # Show success message
                self.console.print()
                self.console.print(
                    f"[bold green]Task '{task}' scheduled successfully![/bold green]"
                )
                self.console.print(
                    f"[green]The task will run according to the specified schedule.[/green]"
                )
            else:
                self.console.print("\n[yellow]Task scheduling cancelled[/yellow]")

            self.console.print("\nPress any key to continue...")
            self.console.input()

        # Update status message
        if confirm:
            self.set_status(f"{task} scheduled for {selected_repo.name}", Theme.SUCCESS)
        else:
            self.set_status("Task scheduling cancelled", Theme.WARNING)

    def modify_setting(self, setting: str):
        """Modify a dashboard setting."""
        with self.console.screen():
            self.console.print(
                f"[bold {Theme.PRIMARY}]Modify Setting: {setting}[/bold {Theme.PRIMARY}]\n"
            )

            if setting == "Appearance":
                theme_options = [
                    "Default",
                    "Dark",
                    "Light",
                    "High Contrast",
                    "Terminal",
                ]
                selected_theme = Prompt.ask(
                    "Select theme", choices=theme_options, default="Default"
                )

                color_scheme = Prompt.ask(
                    "Color scheme",
                    choices=["Blue", "Green", "Purple", "Amber", "Custom"],
                    default="Blue",
                )

                font_size = Prompt.ask(
                    "Font size", choices=["Small", "Medium", "Large"], default="Medium"
                )

                layout_density = Prompt.ask(
                    "Layout density",
                    choices=["Compact", "Comfortable", "Spacious"],
                    default="Comfortable",
                )

                # Show preview
                self.console.print("\n[bold]Appearance Preview:[/bold]")

                # Create a sample panel to demonstrate the theme
                sample = Panel(
                    Text(
                        "Sample repository card with the selected theme",
                        justify="center",
                    ),
                    title="Repository Name",
                    subtitle="Last updated: 2 days ago",
                    border_style=Theme.PRIMARY,
                    padding=(1, 2),
                )

                self.console.print(sample)

                self.console.print(f"\nTheme: [bold]{selected_theme}[/bold]")
                self.console.print(f"Color Scheme: [bold]{color_scheme}[/bold]")
                self.console.print(f"Font Size: [bold]{font_size}[/bold]")
                self.console.print(f"Layout Density: [bold]{layout_density}[/bold]")

            elif setting == "Notifications":
                notify_commits = Confirm.ask("Notify on new commits?", default=True)

                notify_tasks = Confirm.ask("Notify on task completion?", default=True)

                notification_level = Prompt.ask(
                    "Notification level",
                    choices=["All", "Important only", "Critical only", "None"],
                    default="Important only",
                )

                notification_sound = Confirm.ask(
                    "Enable notification sounds?", default=False
                )

                # Summarize settings
                self.console.print("\n[bold]Notification Settings:[/bold]")
                self.console.print(
                    f"Commit Notifications: {'Enabled' if notify_commits else 'Disabled'}"
                )
                self.console.print(
                    f"Task Notifications: {'Enabled' if notify_tasks else 'Disabled'}"
                )
                self.console.print(f"Notification Level: {notification_level}")
                self.console.print(
                    f"Notification Sounds: {'Enabled' if notification_sound else 'Disabled'}"
                )

            elif setting == "Git Settings":
                git_user = Prompt.ask(
                    "Git username",
                    default=subprocess.run(
                        ["git", "config", "user.name"], capture_output=True, text=True
                    ).stdout.strip()
                    or "username",
                )

                git_email = Prompt.ask(
                    "Git email",
                    default=subprocess.run(
                        ["git", "config", "user.email"], capture_output=True, text=True
                    ).stdout.strip()
                    or "email@example.com",
                )

                git_editor = Prompt.ask(
                    "Git editor",
                    default=subprocess.run(
                        ["git", "config", "core.editor"], capture_output=True, text=True
                    ).stdout.strip()
                    or "vim",
                )

                git_default_branch = Prompt.ask(
                    "Default branch",
                    choices=["main", "master", "develop"],
                    default="main",
                )

                # Summarize settings
                self.console.print("\n[bold]Git Configuration:[/bold]")
                self.console.print(f"Username: {git_user}")
                self.console.print(f"Email: {git_email}")
                self.console.print(f"Editor: {git_editor}")
                self.console.print(f"Default Branch: {git_default_branch}")

            # Confirm settings change
            self.console.print()
            confirm = Confirm.ask("\nSave these settings?", default=True)

            if confirm:
                # Show a spinner while "saving"
                with Progress(
                    SpinnerColumn(spinner_name="dots"),
                    TextColumn("[bold green]Saving settings..."),
                    TimeElapsedColumn(),
                ) as progress:
                    task_id = progress.add_task("Saving...", total=None)
                    # Simulate saving delay
                    time.sleep(1.5)

                # Show success message
                self.console.print()
                self.console.print(
                    f"[bold green]Settings saved successfully![/bold green]"
                )
            else:
                self.console.print("\n[yellow]Settings change cancelled[/yellow]")

            self.console.print("\nPress any key to continue...")
            self.console.input()

        # Update status message
        if confirm:
            self.set_status(f"{setting} settings updated", Theme.SUCCESS)
        else:
            self.set_status("Settings change cancelled", Theme.WARNING)

    async def run(self):
        """Run the dashboard asynchronously."""
        self.console.print(
            f"[bold {Theme.SUCCESS}]Starting Repository Dashboard...[/bold {Theme.SUCCESS}]"
        )

        # Show a loading animation
        with Progress(
            SpinnerColumn(spinner_name="dots"),
            TextColumn("[bold green]Loading repositories..."),
            TimeElapsedColumn(),
        ) as progress:
            loading_task = progress.add_task("Loading...", total=None)
            await asyncio.sleep(1.5)  # Simulate loading

        try:
            # Using Python's built-in input handling instead of console.input
            import sys
            import tty
            import termios
            import select

            def get_key():
                """Get a single keypress without waiting for Enter."""
                fd = sys.stdin.fileno()
                old_settings = termios.tcgetattr(fd)
                try:
                    tty.setraw(sys.stdin.fileno())
                    # Check if there's input ready with a short timeout
                    if select.select([sys.stdin], [], [], 0.1)[0]:
                        ch = sys.stdin.read(1)
                        # Handle special keys (arrow keys, etc.)
                        if ch == "\x1b":  # Escape sequence
                            # Read the next two characters
                            if select.select([sys.stdin], [], [], 0.1)[0]:
                                ch1 = sys.stdin.read(1)
                                if ch1 == "[":
                                    if select.select([sys.stdin], [], [], 0.1)[0]:
                                        ch2 = sys.stdin.read(1)
                                        if ch2 == "A":
                                            return "up"
                                        elif ch2 == "B":
                                            return "down"
                                        elif ch2 == "C":
                                            return "right"
                                        elif ch2 == "D":
                                            return "left"
                        elif ch == "\r":  # Enter key
                            return "enter"
                        # Support for space bar
                        elif ch == " ":
                            return "space"
                        # Map vim keys
                        elif ch == "h":
                            return "h"  # left
                        elif ch == "j":
                            return "j"  # down
                        elif ch == "k":
                            return "k"  # up
                        elif ch == "l":
                            return "l"  # right
                        return ch
                    return None
                finally:
                    termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)

            with Live(self.layout, refresh_per_second=10, screen=True) as live:
                self._update_screen(live)
                self.set_status("Welcome to Repository Dashboard", Theme.SUCCESS)

                while self.running:
                    # Update screen with current state
                    self._update_screen(live)

                    # Get key input using our custom function
                    key = get_key()
                    if key:
                        self.handle_input(key)
                        # Small delay to prevent too rapid input
                        await asyncio.sleep(0.05)
                    else:
                        # Small sleep to prevent CPU hogging
                        await asyncio.sleep(0.01)

            # Show closing message
            self.console.print(
                f"[bold {Theme.SUCCESS}]Dashboard closed. Thank you for using Repository Manager![/bold {Theme.SUCCESS}]"
            )

        except Exception as e:
            # Restore terminal if something goes wrong
            self.console.print(f"[bold {Theme.ERROR}]Error: {e}[/bold {Theme.ERROR}]")
            import traceback

            self.console.print(traceback.format_exc())

        finally:
            # Make sure to restore terminal settings
            if "old_settings" in locals():
                termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)


def run_dashboard():
    """Run the dashboard."""
    dashboard = Dashboard()

    # Run with asyncio event loop
    try:
        import asyncio

        asyncio.run(dashboard.run())
    except ImportError:
        # Fallback for older Python versions
        dashboard.console.print(
            "[yellow]AsyncIO not available, running in synchronous mode[/yellow]"
        )
        # Create a simple event loop replacement
        while dashboard.running:
            key = input("Command (q to quit, up/down/enter): ")
            dashboard.handle_input(key)

            # Update display (simple version)
            dashboard.console.clear()
            dashboard.console.print(dashboard._render_header())
            dashboard.console.print(dashboard._render_sidebar())
            dashboard.console.print(dashboard._render_info())
            dashboard.console.print(dashboard._render_actions())
            dashboard.console.print(dashboard._render_footer())


if __name__ == "__main__":
    run_dashboard()
