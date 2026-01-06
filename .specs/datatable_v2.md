# datatable v2

## integration tests
- integration tests that verify actual rendered result.  (playwright?)
- integration tests verify all scrolling operations leave current location *and* displayed datatable correct.

## Paginated datasource
- modify datasource callback api to handle pagination and sort spec.  
- update actual db to use pagination signature for all main tables (even though it's really only needed for transactions)
- fetch-ahead hueristic -- when it needs to fetch more data, datatable asks for twice as many rows as it needs to fill in the rest of the datatable and caches the extra.  So if datatable is 10 rows and user did full page down, datatable would request 20 rows.  But if datatable is doing a client-side filter and has filled 5 rows (with filtered data), it would request 10 more rows.  
- limiting hueristic -- when data datatable is about to fetch more data, it will throw away cached rows that are more than 10 screens away from the current location
- unit tests  verify correct operation when paginating by 1 row and by 3 full screens, and a few numbers of rows in between, all tests trying both directions.  Unit tests also verify correct operation when fewer rows are actually returned than were requested.

## New sort
Rather than allowing only sort by one column at a time, new datatable sort operation allows user to select up to N=3 columns to sort by.
Initial sort can be specified in config.
When user clicks on any sortable column header, datatable will expose a sort selection form (design tbd). In the form looks like:
```
Sort <drop down: uparrow - ascending, downarrow -descending> by: <drop down of sortable columns, including "none">  
Then <drop down: uparrow - ascending, downarrow -descending> by: <drop down of sortable columns, including "none", but excluding the previous chosen column>  
Then <drop down: uparrow - ascending, downarrow -descending> by: <drop down of sortable columns, including "none", but excluding the previous 2 chosen columns>  

<cancel> <ok>
```
If the user confirms the sort, the column headers will reflect the new sort order, showing a small number 1,2,3 followed by the sort direction arrow
And the datatable will do a full requery and re-render with the new sort order.

- unit tests verify correct operation when sort requested on none, one, two, and three columns in a selection of ascending and descending. Unit tests should use a dummy data source that returns a known data in each row and column whose contents can be used to verify the correct cell is populated with the expected data.  Dummy data set can be configured to provide as many or as few rows and columns as needed for the particular test scenario.
## new integration tests
Since pagination and sort both affect the data source API, need new integration tests which exercise combinations of both sorting and pagination.  These tests can use the same dummy data set as above, but should "peek" at the rendered data to verify the cell contents are as expected.
