# form_component_plan
The goal is to implement a highly modular data grid which can be used for many purposes within budgetTracker and also be contributed to the community as a standalone component.  The grid displays the currently visible portion of what may be a larger table of data.

The datagrid uses svelte, shadcn-svelte, and components from [shadcn-svelte](https://www.shadcn-svelte.com/docs/components).

## Grid Capabilities
- The grid displays a window on a potentially larger table of data.
- The schema of data is configured by caller, the grid adapts its display to match the configuration.
- Grid displays rows and columns based on their "natural" size, and does not  attempt to pad to completely fill the parent container. Individual columns can be configured to wrap to a specified number of lines if they would otherwise exceed a configured number of characters width.
- When size of data would exceed the bounds of the parent container, the grid displays scroll bars in that direction.
- it implements infinite scroll by invoking "data source" callbacks to provide chunks of row data. [[#data_source]]
- it implements manual horizontal scroll if the rows are too wide.
- it allows the user to sort the current grid
- it allows the user to filter the data, showing only rows containing a match.
- it allows the user to find a match in the data, scrolling the grid to make the match visible.
- allows editing the data and communicating the changed row(s) through a callback
- User can select one or more rows and process them through a callback
  
## User interface
### Concepts
- current  position -- the row and column of the data (not the grid) which is the current focus of operations.  
- visible location -- the row  and column of the grid where the current location in the data is displayed.  

When an operation updates the current position, the new visible location in the grid is briefly highlighted.  The grid will scroll data as necessary to ensure the visible location is indeed visible.
The grid will minimize movement of the visible location.  If an operation is changing the row, the visible location column will not change and vice versa.
### Initial appearance
- Grid fills its container from upper left corner to bottom right, displaying cells at their "natural" (or explicitly configured) size.  Parent container can shrink if there's not enough data to fill it.  
- If there is more data than the grid can display, grid will show scroll bars in the direction(s) where there is more data.
- Current row is set to first row of data in the grid.
- The title and column headers are always visible.  The data grid scrolls under the column headers, but the column witdths stay in sync.

### Moving within the grid
The user can move the current position by keystrokes:
- arrow keys -- one row or column in the desired direction
- home -- first column of the current row in the grid
- end -- last column of the current row in the grid
- page up -- scrolls the data up by the  number of rows currently visible in the grid, leaving the visible location unchanged within the grid.  But if there is not enough data, the visible location will move to the top of the grid.
- page down -- scrolls the data down by the number of rows currently visible in the grid, leaving the visible location unchanged within the grid, or moving it to the bottom of the grid if there is not enough data to fill the grid.
- shift-home and shift-page up -- first row of data
- shift-end and shift-page down -- last row of data

### Data editing
 When the user double-clicks on any cell in a row, a popup form opens allowing the use to edit all fields of the row.  The form has OK and Cancel, and if the user clicks OK, the grid invokes a data callback to update that row.  A future version will support editing the data  table cells in place.

### Column resizing
If the user double-clicks on the column separator between columns, the grid will display a vertical splitter and allow the user to manually resize the column to the left, shifting the subsequent columns without changing their sizes. The datagrid will remember the column width and apply it to future views.

### Sorting
initial implementation: if a column is configured for sorting, user can click on the header for that column and cycle through ascending, descending, and no sorting.
Sorting resorts the whole data source, but leaves the visible location unchanged (if there is enough data to fill the grid).
Sorting is sensitive to the datatype of the column, so numeric columns are sorted numerically and text columns are sorted lexicographically

### Finding.and filtering
If configured to enable it, the grid shows a "Search" field and also a "Filter by" field.  
- The find and filter operations are incremental and asynchronous: as the user types, character by character, the find or filter operation is repeated and the view updated asynchronously.  If the user happens to type multiple characters before the view is updated, the next update will use all the characters currently in the find/filter input area.
- Find and filter are both case insensitive.  They operate on all columns.

- Find scrolls the grid to make the found item visible. It does not move the current location if no match is found.
- Filter filters the data source, hiding any row that does not contain a substring match in any column.  The filter is applied to the data source as a whole, not to the visible portion of the grid.  

implementation note: Datasource does not implement find or filter, this is done by the datagrid.

### Selecting cells
The user can select one or more rows in the datagrid which can be processed through a callback, or sent to the clipboard.

To start a selection the user:
1. clicks on the row to start a selection
2. scrolls as desired, then either:
   1. shift-click on antother row to select all rows from the start to end, inclusive.
   2. ctrl-click on a row to add or remove that row from the selection. and repeats as desired
3. completes the selection by typing ctrl-c.

At this point, a CSV representation of the data is sent to the system clipboard and the data selection callback (if any) is invoked.

It is also possible to select a single row.  The user:
1. navigates to the desired row
2. types ctrl-c.

Typing ctrl-c with no previously selected row initiates the same selection processing (clipboard and callback).

## configuration
### General datagrid attributes
- name -- the name of the datagrid.  This is used to identify the datagrid  in callbacks
- keyColumn -- the name of a column in the datasource which is unique for each row.  This column name is shared with the data source callback and is not assumed to be string datatype.
- Title -- a visible title for the datagrid, displayed in the header
- maxVisibleRows -- the maximum number of rows to display at once.  If not specified, the datagrid will fill its container.
- isFilterable -- whether the datagrid should allow the user to filter on this datagrid.  
- isFindable -- whether the datagrid should allow the user to find on this datagrid..  

!! Open question: is the data source reponsible for converting binary data to character form, or does the datagrid do it?   Affects sort, find and filter. Agent should provide design feedback before committing to a full plan.   For now, assume datasource does it and the datagrid just does centering/justification.

### Columns
An array of column descriptions with the following attributes:
- name -- the field name in the data source for data in this column.  For a database data source, this is the column name, e.g.
- title -- the label for the column header.  If not specified, the name is used, in Title Case.
- isSortable -- whether the datgrid should allow the user to sort on this column.  The sort is always lexocographic.  If caller wants numbers or dates to sort "nicely", s/he must choose a format that suports it.
- wrappable -- 3 options:
    - none -- field is not wrappable at all.
    - word -- field can be wrapped at word separator characters
    - hard -- field can be wrapped at any character
- maxLines -- defaults to 1.  If wrappable is not none, this is the maximum number of lines to display in a cell.  If the cell contains more lines, it is truncated. Note that all lines are considered part of one row in the datagrid.
- justify -- 3 options:
    - left -- field is left-justified
    - center -- field is centered
    - right -- field is right-justified
- maxWidth -- maximum number of characters to display in a cell.  Default is no limit.  If set and if the cell contains more characters, it is truncated. Note this is number of characters in the data, not number of pixels of screen.

### Data Source Callback
A function invoked by the datagrid to get the data to display.  It is passed the following arguments:
- columnKeys -- list of keys of columns to query.  Datasource is expected to provide results back with columns in same order.  Note that the first element of this list is the keyColumn, the rest are the configured columns.
- startRow -- the first row to return
- numRows -- the minimum number of rows to return
- sortColumnKeys -- a list of objects {key, sortDirection (asc or desc) }

It returns an array of row objects with attributes corresponding to keyColumn and columnKeys.  The keyColumn may be numeric or string, the columnKeys are always string.

### Data Edit Callback
A function invoked by the datagrid to edit the data.  It is passed the following arguments:
- beforeData -- the row(s) of data before editing.  This is an array of objects with attributes corresponding to keyColumn and columnKeys.
- afterData -- the row(s) of data after user editing.  This is an array of objects with attributes corresponding to keyColumn and columnKeys.

It returns an object as follows:
- validationFailures -- an array of objects {columnKey, message} indicating which columns failed validation and why.  If all columns validated successfully, this array is empty.
- dbStatus -- a string indicating any data source error completing the edit.  If null, the edit was successful.  If non-empty, it is a single string describing the error.

If validationFailures is non-empty or if dbStatus is non-null, no changes were made to the data source.  The datagrid should display the error(s) and allow the user to retry the edit operation.
If validationFailures is non-empty, the datagrid should display the error(s) associated with the particular column in columnKey.  
If dbStatus is non-null, it should be displayed as an error at the bottom of the whole edit operation form.

Note that this interface implies the user can edit multiple rows in one operation.  However, the current datagrid only allows edit of one row at a time.  Future versions will support editing multiple rows at once.

### Selected cells callback
A function invoked by the datagrid to notify the caller of the user's selection of cells.  It is passed the following arguments:
- selectedCells -- an array of objects {rowKey, columnKey} indicating the selected cells.

It returns a status string, which is displayed at the bottom of the datagrid.  If null, the operation on the selected cells, whatever it was, succeedeed and the user shoujld see a standard OK message.

## Implementation notes
Implement unit tests to verify the callbacks in isolation.  Use table driven subtests to cover  normal and error inputs and error returns.
Implement unit tests and subtests to verify basic navigation, sort, filter, find and selection.  Use table driven subtesst to cover scenarios when the datasource does not and then does exceed the visible data grid.
