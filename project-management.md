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
- Error handling and robustness ✅
- Validate the trello token if it's incorrect, in TrelloClient ✅

## In Progress
- Basic CLI structure
    - `mdello init`
        - connect trello api key ✅
        - select default board ✅
            - under the hood chooses the boardId, uses trello.Board ✅
    - `mdello boards`
        - change config.currentBoards → to new selected one like in railway ✅
    - Feature like `git rebase -i` that opens a new window ✅
    - Refactor cli.go to have multiple files
    - Get all the board data into the new file
        - Show current board name, list names with cards, cards with checklists
        - Implement the movement of lists and cards when you move them in the editor
        - The `[x]` and `[]` implemented on the cards
        - If there are not invalid commands ask the user to fix it

## Queued Up
- When trello.go gets too big separate into new dir/ with lists.go, cards.go

## Backlog
- Update the operations.go to use the Params structs at input parameters
- Refactor operations.go functions to use helper functions as similar code for each: createrequest, request, response
