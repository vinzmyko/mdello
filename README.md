# mdello

A command-line tool that lets you manage Trello boards using markdown as the primary interface.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Basic Commands](#basic-commands)
  - [Markdown Structure](#markdown-structure)
  - [Working with Cards](#working-with-cards)
  - [Managing Labels](#managing-labels)
  - [Due Dates](#due-dates)
  - [Detailed Editing Mode](#detailed-editing-mode)
- [Examples](#examples)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [Licence](#licence)

## Features

- **Markdown-first workflow**: Edit Trello boards using familiar markdown syntax
- **Hierarchical structure**: Boards → Lists → Cards represented as `#` → `##` → `- [ ]`
- **Label management**: Create, assign, and manage labels directly in markdown
- **Due date support**: Set and manage card due dates with customisable formats
- **Bulk operations**: Move multiple cards and lists by reorganising markdown
- **Detailed editing**: Access advanced board, list, and card settings
- **Browser integration**: Open boards directly in Trello via your default browser
- **Offline-first**: Make changes in your preferred text editor

## Installation

### Prerequisites

- Go 1.18 or higher
- A Trello account and API token

### Install from releases

Download the latest binary from the [releases page](https://github.com/vinzmyko/mdello/releases) and add it to your PATH.

### Install from source

```bash
git clone https://github.com/vinzmyko/mdello.git
cd mdello
go build -o mdello
sudo mv mdello /usr/local/bin/
```

### Install via Go

```bash
go install github.com/vinzmyko/mdello@latest
```

## Quick Start

1. **Initialise mdello with your Trello token:**
   ```bash
   mdello init
   ```
   You'll be prompted to enter your Trello API token. Get your token from [https://trello.com/app-key](https://trello.com/app-key).

2. **List your boards:**
   ```bash
   mdello boards
   ```

3. **Edit a board:**
   ```bash
   mdello board
   ```
   This opens your current board as a markdown file in your default editor.

4. **Open board in browser:**
   ```bash
   mdello open
   ```

## Usage

### Basic Commands

```bash
mdello [command]

Available Commands:
  board       Edit current board via markdown file
  boards      Get all current user's boards
  help        Help about any command
  init        Initialise mdello with your Trello token
  open        Open current board in Trello via default browser

Flags:
  -h, --help   help for mdello

Use "mdello [command] --help" for more information about a command.
```

### Markdown Structure

mdello uses a hierarchical markdown structure to represent Trello boards:

```markdown
# Board Name {board_id}
@label_name:colour {label_id}

## List Name {list_id}
- [ ] Card name @label due:25-07-2025 21:34 {card_id}
- [x] Completed card {card_id}

## Another List {list_id}
- [ ] Another card {card_id}
```

**Structure breakdown:**
- `#` = Board name
- `@label:colour` = Board labels (below board name)
- `##` = List names
- `- [ ]` = Incomplete cards
- `- [x]` = Completed cards
- `{id}` = Unique identifiers (automatically managed)

### Working with Cards

**Creating cards:**
Add new lines with `- [ ]` under any list:
```markdown
## To Do
- [ ] New task
- [ ] Another new task
```

**Moving cards:**
Cut and paste cards between lists:
```markdown
## Backlog
- [ ] Move this card

## In Progress
- [ ] Move this card  # Moved from Backlog
```

**Completing cards:**
Change `[ ]` to `[x]`:
```markdown
- [x] This task is complete
```

### Managing Labels

**Creating labels:**
Add labels below the board name:
```markdown
# My Project {board_id}
@urgent:red {label_id}
@feature:blue {label_id}
@bug:yellow {label_id}
```

**Applying labels to cards:**
Reference labels with `@` in card text:
```markdown
- [ ] Fix critical bug @urgent @bug {card_id}
- [ ] Add new feature @feature {card_id}
```

### Due Dates

Set due dates using the format configured during `mdello init`:

```markdown
- [ ] Task with deadline due:25-07-2025 21:34 {card_id}
- [ ] Another deadline due:30-07-2025 12:00 {card_id}
```

### Detailed Editing Mode

For advanced editing of boards, lists, and cards, append `!` to any item:

```markdown
# My Board {board_id}!
```

This opens a detailed editing interface with additional options:

**Board settings:**
- Description
- Permissions (org/private/public)
- Voting settings
- Comment permissions
- Display preferences

**List settings:**
- Position
- Archive status
- Subscriptions

**Card settings:**
- Detailed descriptions
- Start and due dates
- Archive status
- Notifications

## Examples

### Basic Project Board

```markdown
# Website Redesign {8ee70}
@urgent:red {c5adc}
@design:blue {d4e8f}
@development:green {f9a2b}

## Backlog {0be10}
- [ ] User research @design {f5061}
- [ ] Competitor analysis @design {a3f84}

## In Progress {c9e65}
- [ ] Homepage mockup @design due:28-07-2025 17:00 {b7d92}

## Review {1f0d6}
- [ ] Logo concepts @design {e8c47}

## Done {383bd}
- [x] Project kickoff meeting {30e09}
- [x] Requirements gathering due:20-07-2025 12:00 {a59dd}
```

### Detailed Board Editing

```markdown
# Website Redesign {8ee70}!

# =============================================================================
# EDITING BOARD: Website Redesign {68826e536b5f666e87d80ca8}
# =============================================================================

## Description
=== DESCRIPTION START ===
Complete redesign of company website focusing on:
* Improved user experience
* Mobile responsiveness  
* Modern design language
* Accessibility compliance
=== DESCRIPTION END ===

## Board Settings
Name: Website Redesign
Closed: false

## Permissions
Permission Level: org
Voting: members
Comments: members
Invitations: members
```

## Configuration

mdello stores configuration in `~/.mdello/config.json`:

```json
{
  "trello_token": "your-api-token",
  "date_format": "02-01-2006 15:04",
  "editor": "vim",
  "current_board": "board-id"
}
```

**Customisable settings:**
- Date format for due dates
- Default text editor
- API credentials
- Current working board

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Development setup:**
```bash
git clone https://github.com/vinzmyko/mdello.git
cd mdello
go mod download
go build
```

**Running tests:**
```bash
go test ./...
```

## Licence

This project is licensed under the MIT License

---

**Need help?** Open an issue on [GitHub](https://github.com/vinzmyko/mdello/issues)
