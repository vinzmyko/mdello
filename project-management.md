# Initial Starting Plan

1. Separation of Concerns
    - TrelloAPI - "Talks to Trello's servers" (Create a interface for this for testing later)
    - MarkdownParser - "Converts between markdown and board objects"  
    - FileHandler - "Reads and writes markdown files"
    - SyncCoordinator - "Orchestrates the sync process"
2. Error Handling Strategy
    Main things that can go wrong:
        - Network errors (API down) → retry 'x' amount of times and inform user
        - File errors → show clear message and exit
        - Parse error → Show which line is broken
        - Sync conflicts → For MVP just overwrite
3. Testing Strategy
    - MarkdownParser: has lots of edge cases 
    - SyncCoordinator: core logic

# MVP Task list

## Done
- Access API key and token ✅
- Create TrelloClient in trello.go ✅

## In Progress
- CRUD operations:
    - GET: boards ✅, lists ✅, cards ✅
    - CREATE: boards, lists, cards
    - UPDATE: boards, lists, cards
    - DELETE: boards, lists, cards

## Queued Up
- HealthCheck() error { return nil }, for use as a guard clause
- Error handling and robustness
    - Currently, using log.Fatal() everywhere
    - Make functions return ([]map[string]any, error)
    - Handle HTTP status codes (404, 401, etc.)
    - Add timeout handling?
- Basic CLI structure
- Data Models
    - Instead of map[string]any, create proper structs

## Backlog
- When trello.go gets too big separate into new dir/ with lists.go, cards.go, client.go, types.go etc
