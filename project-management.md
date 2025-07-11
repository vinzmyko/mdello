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
- Trello CRUD operation functions ✅
- HealthCheck() error { return nil }, for use as a guard clause ✅
- Create, Read, Update Data Models ✅
- Make functions return created structs not `[]map[string]any,` error ✅
- Currently, using log.Fatal() everywhere, replace it ✅

## In Progress
- Error handling and robustness
    - Handle HTTP status codes (404, 401, etc.)
    - Add timeout handling?

## Queued Up
- Basic CLI structure

## Backlog
- When trello.go gets too big separate into new dir/ with lists.go, cards.go
- Update the operations.go to use the Params structs at input parameters
