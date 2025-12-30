---
trigger: manual
---

# Design Plan for budgetTracker

## Back end
### database -- package "db"
The back end uses a sqlite database stored in a configurable location.

database types:
- text -- a character string (utf8 encoded)
- money -- an amount of currency.  Store as an integer number of cents.  
- date -- a year, month, day.  No hour, minute, second.  Store as a string format: yyyy-mm-dd
- timestamp -- a full date/time, with time zone.

items:

- beneficiary -- person owning an account
    - name -- type string, is primary key

- account -- bank or credit account
    - name -- type string, is primary key
    - description -- type string
    - beneficiary -- link to beneficiary record which owns the account  

- transaction -- transaction on an account
    - postedDate -- type "date"
    - account -- link to account record the transaction posted in
    - amount -- type "money".  Positive value increases the balance in the account, negative reduces it.
    - description -- type "text".  
    - tag -- type "text".  A classifier which is relatable to a budget item (computed later)
    - budgetItem -- link to a budget line item (when the transaction is finally assigned to a budget item)

- budget -- budget for a beneficiary
    - name -- type string, is primary key
    - description -- type string
    - beneficiary -- link to beneficiary record which owns the budget item
    - amount -- type money, planned amount for this budget item
    - intervalMonths -- type integer, number of months the budget amount covers.

### database administration -- package "db_admin"
The database will have additional administrative tables and functions which are not included in the budget application.  These shoould be implemented in db_admin.

## data models - package "models"

Project uses gorm ORM to integrate backend go code with database schema.  

Package "models" contains the gorm definitions for the data objects defined in the database.

Later, we might also create a package "models_test" for objects used for testing.

In addition to the gorm declarations, each model package has a "Service" type with methods initialize the database, perform automigration and access the database.  The application does not use gorm API calls directly, instead it uses methods on the Service type.

## Front end - module "frontend", multiple pacages as appropriate
Embedded WebView using SvelteKit, tailwind, shadcn-svelte, etc. from the starter template.

The main window has a top menubar with each of the main functions, such as "Import transactions", "Reports" and "Adjust Budget".  Each main menuitem will have multiple sub-items for specific operations.

## version 0.1
For the first verions of the application, implement packages db and models as specified above.
Implement a front end which has top menubar navigation as follows:
- database
    - clean -- initializes a database with empty tables for package db
- beneficiaries
    - add -- add a new beneficiary
    - list -- list all beneficiaries
    - update -- select an existing record by matching a key entered in an input field, show the current values in a form, then offer buttons to cancel (do nothing), update (write changed data back to the record) or delete (delete that record).

- accounts
    - add -- add a new account
    - list -- list all accounts
    - update -- select an existing record by matching a key entered in an input field, show the current values in a form, then offer buttons to cancel (do nothing), update (write changed data back to the record) or delete (delete that record).

- budgets
    - add -- add a new budget
    - list -- list all budgets
    - update -- select an existing record by matching a key entered in an input field, show the current values in a form, then offer buttons to cancel (do nothing), update (write changed data back to the record) or delete (delete that record).

- transactions
    - add -- add a new transaction
    - list -- list all transactions
    - update -- select an existing record by matching a key entered in an input field, show the current values in a form, then offer buttons to cancel (do nothing), update (write changed data back to the record) or delete (delete that record).
