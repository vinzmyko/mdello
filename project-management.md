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
    - `mdello init` ✅
    - `mdello boards` ✅
    - `mdello board`
        - Finish all the actions for actions.go
            - CRUD for TrelloObjects:
                - Board: Only thing you should be able to do it change the board name ✅
                - List CreateList ✅, UpdateListName ✅, UpdateListPosition ✅, ArchiveList ✅,
                - For list/card creation need to think of a way to deal with the shortID since when you create one you don't know the id ✅
                - Card
                    - CREATE, UPDATE
                        - TaskCompleted ✅, Name ✅, labels, due date
        - Apply them to Diff() function

// How the get labels look like
label name='cool'
label name='label with space'
label name=''
label name=''
label name=''
label name=''
label name=''
label name='from main func'

- Get the actions for creating/deleting a label from the board
    FILES:
        - from_markdown.go - parsing the labels+labelid into data we send to the diff.go to compare them and then we need to implement the functions ✅
    - Just use short ids for the labels so we can detect changes instead of our current approach ✅
        - Add shortid in to_markdown() needs to be in this format: `@feature:blue {d3e4f}` and they are stacked vertically ✅
    - Test all the Create, Update, Delete for the actions
        - actions.go LABELS
            - Update label name ✅
            - Update label colour ✅
            - Create Label ✅
            - Delete Label ✅
    - Remember the subtle = _light, bold = _dark ✅
    - Might need to delete some functions when I was testing ✅

    - When moving list might change the inside data ✅
        - Right now when creating a new card we should check the shortID to see if it existed in the previous lists and if it did move that data ✅
            - Were gonna need a way to check all the lists for this somehow meaning maybe change the structure of how we find lists/cards in the map nested loop ✅

- get labels to be updated for card actions ✅
    - Implement add label to cards via markdown ✅
    - Implement delete label to cards via markdown ✅
    - Need to add the `~` to the boardID and get it working with the cards as well ✅
- Get the Duedate working in the actions ✅
    - think about the best way to represent the duedate ✅

- Now get the deleting and adding labels to cards ✅
    - TODO: create label and add it to card
    - TODO: when you create the card it should be able to accept labels and duedate as well


git add markdown/from_markdown.go
git commit -m "fix(markdown): remove label validation that conflicts with due dates"

git add markdown/to_markdown.go
git commit -m "feat(markdown): add parseMarkdownDate() and make formatDate() more modular"

git add markdown/actions.go
git commit -m "feat(markdown): add update and delete card duedate actions"

git add markdown/diff.go
git commit -m "feat(markdown): handle card due date creation, update, and deletion"

git add cli/board.go
git commit -m "fix(cli): inject userconfig to markdown conversion functions"

git add project-management.md


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

## Queued Up
- Add the `!` feature at the end in which you can add more detail to your changes not just the quick ones
- in actions.go update all the `return err` into `return fmt.Errorf()`

## Backlog
- If there are not invalid commands ask the user to fix it
- When I didn't have net and tried to do cli command I got a weird error. Error wasn't obvious that my internet was down so fix that
- Create a `board --view/web` in which it opens your prefered browser to look at your tasks
- When trello/operations.go gets too big separate into new dir/ with lists.go, cards.go

### BUG
- [Non Fatal] Missing ID detection logic creates false positives for new items
    - Currently, the system assumes any list or card without a {shortID} is a new item that needs to be created. 
      However, this creates false positives when users accidentally delete IDs from existing items.
        - The ResolveShortID("") method will generate sentinel values (NEW_ITEM_1, etc.) for any missing ID, without verifying if the item actually exists or is genuinely new.
            - No {shortID} → Generate sentinel ID → Diff treats as new item → Create action
