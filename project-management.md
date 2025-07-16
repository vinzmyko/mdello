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
- Update the operations.go to use the Params structs at input parameters ✅
- Refactor operations.go functions to use helper functions as similar code for each: createrequest, request, response ✅

## In Progress
- Basic CLI structure
    - `mdello init`
        - connect trello api key ✅
        - select default board ✅
            - under the hood chooses the boardId, uses trello.Board ✅
    - `mdello boards`
        - change config.currentBoards → to new selected one like in railway ✅
    - Feature like `git rebase -i` that opens a new window ✅
    - Refactor cli.go to have multiple files ✅
    - Get all the board data into the new file
        - Show current board name ✅, list names with cards ✅, cards with checklists ✅
            - What do I want information do I want to convey for the Card?
                - Labels ✅, due date ✅
        - The `[x]` and `[]` implemented on the cards ✅
    - Figure out how to parse the markdown into data that we can send to the trello api
        - Implement the movement of lists and cards when you move them in the editor
        - If there are not invalid commands ask the user to fix it

// How cards are structured example
card title:vinzmykodelrosario
        id:b293
        status:[x]
        position:0
        listID:b281
        due date:17-07-2025 15:47
        labels:
                test
                vinz


- Add id generation to ConvertToMarkdown() ✅
    - Add shortId helper func ✅
    - Edit all trello objects to include ids ✅
- Create Basic parser structure
    - scan line by line for board, lists, card ✅
    - handle edge cases for regex carefully ✅
- Extend parser.go to extract the ids ✅
- Create structure based on markdown file ✅
- Clean up the board.go file and parser.go file and commit ✅
- Compare the original markdown file struct vs the editted markdown file ✅
    - What are the changes from the original and the editted version ✅
- Simple board name changes ✅

TODO: Refactor code now that it works
- Figure out how to deal with the card ids
- Refactor and clean up all the trail code
    - parser.go
        - Which structs need to be private
        - Correct regex based on how we solved the shortID in convert.go
        - Should detect changes function be in here?
        - right now the file handles markdown → intermidiate structs and comparing the original and editted
        - origList and edList seem messy fix the variables I made here
    - actions.go
        - Think about the best layout to handle the structs and the interface implementation
            - all structs together then the implementation methods?
            - struct then related impl methods?
                - in order of board, list, cards

- Moving cards

## Queued Up
- Show all the labels the board has and when you add `- @labelName` you can then add it in the cards underneath
    - To avoid errors we need to process the label changes first and then apply them to the cards

## Backlog
- When trello/operations.go gets too big separate into new dir/ with lists.go, cards.go
