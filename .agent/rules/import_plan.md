---
trigger: manual
---

This rule applies to importing transaction data from external accounts to populate transactions table.

The import process allows the user to import new, raw transactions from banks and credit card accounts, to review, edit and categorize them accourding to budget items they relate to, then to merge the updated raw transactions into the main transactions table

Most of the code for doing this is contained in a new package "transactionImport".

The import process proceeds as follows:
1. ask the user to select a CSV file and specify the account into which the transactions shold be loaded.  The CSV format is different for each account, and each account will need custom logic to transform a CSV record into the corresponding raw transaction record.  See table below for sample data for each account.
2. load the specified CSV file into the raw transactions table (which may include raw transactions from other accounts).  Inserting the same csv file multiple times will update matching records in the raw transaction table and create new ones.
3. run an automated data cleanup and categorization process over the imported raw transactions  The categorization process sets the budgetLine field to link the transaction with a budget line item.
3.1.  The user runs steps 1-3 for all the accounts hving new CSV data available.

When the user has completed the CSV upload, the import process:
4. allows the user to review and manually edit the raw transactions table, creating, deleting or updating raw transactions.
5. the user can sort and filter a view of the raw transactions table to assist with the manual review process.
1. The user selects a "finalize imports" button, which upserts the main transactions table from the raw transactions table and then empties the raw transactions table.  The user is informed how many new records were inserted and how many were updated.

## Schema changes
The existing Transactions table has field "BudgetLineID".  Replace this with field called "BudgetLine" which is a foreign key link to Budget table field "name"

Both Transaction and RawTransaction have a field Amount.  This will always be positive for an expense and negative for income.  It will always be a number of cents, as noted previously.

The raw transactions table has the postedDate, Account, Amount, Description, Tag and BudgetLine fields that the main Transaction table has.  In addition it has a field "Action" which has values "update" or "add".  Action is "update" if this raw transaction will change an existing record already present in Transactions.

## new data views and components
Raw Transaction form: a data form for a raw transaction that shows all the fields and allows the user to CRUD a single raw transaction.  The view should have buttons to Update the existing raw transaction, to delete it or to create a new one (gives an error if it would overwrite an existing raw transaction already in the raw transactions table).

Raw Transactions table: a table view for the raw transactions table.  The table can be sorted by clicking on any column header, where clicking cycles through sort ascending, sort descending or don't sort by this column. The user can click on a row to edit that record. When the user clicks on the row, the Raw Transaction form pops up, prepopulated with the data from that row.

Also write unit tests to exercise all the service methods of the raw transactions table and its data views.

Each month, banks and credit card companies provide CSV files containing raw transactions.
The user must download these CSV files, then run the import process in budgetTracker to update transactions.
The overall process must be idempotent, meaning running the import multiple times for the same month and account will not create duplicate transactions.  
Instead, it will update data in existing transactions or create new transactions as necessary.  The idea is that user may manually fix up the raw data and run the import multiple times and does not want to create duplicate transactions.  This also allows the user to run the import for overlapping date ranges without creating duplicate data.

In the course of coding this, ask for clarification of anything that's inconsistent when presenting the plan for implementation.

## table of CSV and sample data

For production use, we will have these Beneficiaries:
- Bob
- Jessie
- Us

and these accounts

- CapitalOne
    - name: CapitalOne
    - description: Capital One rewards Credit Account
    - beneficiary: us
    - CSV sample file:  import/2025-12-27_transaction_download.csv
    - custom field processing.
        - first line of file is a header, naming the fields below.
        - Transaction Date: ignored
        - Posted Date: maps to models.Transaction.PostedDate
        - Card No.: if == 3028, models.Transaction.Beneficiary is "bob", if == 6539, models.Transaction.Beneficiary is "jessie" if anything else, beneficiary is "us".
        - Description: maps to .Description
        - Category: maps to models.Transaction.RawHint (a new string)
        - Debit,Credit: the value (Debit*100 - Credit*100) maps to models.Transaction.Amount

- WFChecking
    - name: WfChecking
    - description: Wells Fargo checking
    - beneficiary: us
    - CSV sample file: import/Checking1.csv
    - custom field processing
        - no header row.  
        - field 1: maps to .PostedDate
        - field 2: -1 *(value* 100) maps to .Amount
        - field 3: ignored
        - field 4: ignored (check number)
        - field 5: maps to .Description
        - models.Transaction.Beneficiary defaults to "us"

-- WFVisa
- name WfVisa
- description: Wells Fargo Visa
- beneficiary: us
- CSV  sample file: import/CreditCard3.csv
- custom field processing:
    - no header row
        - field 1: maps to .PostedDate
        - field 2: -1 *(value* 100) maps to .Amount
        - field 3: ignored
        - field 4: ignored
        - field 5: maps to .Description
        - models.Transaction.Beneficiary defaults to "us"

Remove the functions that generate random accounts and random beneficiaries, add code in the "clean database" function that initializes Accounts and Beneficiaries with these values.
