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

- Get the actions for creating/deleting a label from the board ✅
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
    - create label and add it to card ✅
        - Created problems since the way I created cards before. I used the state when the user started editing the markdown. However, I couldn't create something then use it
        in the same markdown file. Therefore, I made it so the CreateLabelAction command got a fresh set of board labels each time. ✅
    - when you create the card it should be able to accept labels and duedate as well ✅
    - Check creating new list and new cards with card properties ✅

"As a user, I want to 'Add a `!` to the end of the markdown representation of boards, lists, and cards'. So I get more specific changes to the trello objects"
    - Add the `!` feature at the end in which you can add more detail to your changes not just the quick ones
        - Make all the board, list, card lines take `!$` or something at the end of the string find the `!` and put it in a group `()`
            - Check to see if they appear with something that notices it and does a print line
                - Board ✅, list ✅, card ✅
        - We then need to figure out how to put then in a list and apply the to markdown at the end of the first markdown file edit ✅
            - Do I want it to show the fmt.Printf("\n Change:...") lines between new markdown file and previous markdown file? I guess not in the new updates we should user
              the new API calls that get fresh data maybe. Or we can change the return type of Diff() to return like a struct to contain simpleActions []TrelloAction and
              detailedActions []TrelloAction. I am just giving ideas ✅
    - We would need a way to put a lot of the specific edits into one markdown file and then process each section as separate api requests after the first markdown file ✅
    - These api trello call changes should happen after the quick changes. We should get the new data based on the quick changes ✅
    - Fill in the GenerateDetailedXContent with things that the user can change for detailed information. ✅
        - Decide which params the user can edit and display then in the Generate function ✅
        - Might need to pass in the *trello.TrelloClient to get the data about it. definitely ✅
    - Create another function in from_markdown.go that analyses the string content with readers and such
        - Possibly spit diff.go into markdown/diff/quickactions.go and detailedactions.go ✅
        - I guess this should go into diff.go since we don't need to create two versions since everything will be there. We need a way to track the trello object ideas
          for the UpdateXParams. ✅ 
        - We should do the section thing but we will need to compare both the originalContent and the editedContent. ✅
    - Use the strconv I did with trello/operations.go to handle these fields since the trello object sections of the same object type are the same ✅
        - Create mappings for the markdown keys and the request param keys ✅
        - Create state machine parser ✅
            - Field level change detection and the trello api value conversion ✅
            - I think we need to create the updates in bulk and not use composition like I did when creating cards. This is because the sections are the same and not different
              each time like the quick actions editor ✅
            - Figure out a way to handle the boolean values with validation ✅
                - Need to do this with string and position field as well ✅

- Create a `open` in which it opens your prefered browser to look at your tasks ✅
- Update the `help` command ✅
- Need to allow the edited file to recognise which trello object id it is by putting the trello object id in the header ✅
    - Now able to insert the trello id key onto the bulk action for detailed actions ✅
- Implement README.txt for my cli tool ✅
- in actions.go update all the `return err` into `return fmt.Errorf()` ✅

## Queued Up
- Use cli tool and when you notice bugs write them down and fix

## Backlog
- Add the Card checklist implementation in the detailedActions markdown section
- If there are not invalid commands ask the user to fix it
- When I didn't have net and tried to do cli command I got a weird error. Error wasn't obvious that my internet was down so fix that
- Make a new command in which you can update user preferences. The only one I can think of would be if you want the board labels to show and put the date format here as well
- When trello/operations.go gets too big separate into new dir/ with lists.go, cards.go

### BUG
- [Non Fatal] Missing ID detection logic creates false positives for new items
    - Currently, the system assumes any list or card without a {shortID} is a new item that needs to be created. 
      However, this creates false positives when users accidentally delete IDs from existing items.
        - The ResolveShortID("") method will generate sentinel values (NEW_ITEM_1, etc.) for any missing ID, without verifying if the item actually exists or is genuinely new.
            - No {shortID} → Generate sentinel ID → Diff treats as new item → Create action
- [Medium] Missing the `due:20-03-2025 00:00` more specifically the `:` will result in this error 
"Error parsing markdown: invalid label format: labels cannot contain spaces. Use ~ for spaces (e.g., @front~end for 'front end')"
    - On this line `	if invalidLabelPattern := regexp.MustCompile(`@\w+\s+\w`).FindString(tempText); invalidLabelPattern != "" {`
        - This is because it didn't recognise that format. We need to update the error message because I did not know what the error was that I forgot the `:`
- [High] Seems like there are some positions bugs when moving lists
    - To reproduce I deleted a list called `Done` and then created another one and moved it to the end
        - Need more ways to reproduce the bug so I can fix it. Do not know the edge cases
