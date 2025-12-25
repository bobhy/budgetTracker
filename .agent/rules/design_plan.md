# Design Plan for budgetTracker
This is guidance for how to structure the project and which APIs to use, intended for use when writing new code.  It does not (yet) describe the code currently in the project!

## Back end
The back end uses a sqlite database stored in a configurable location.

Project uses gorm ORM to integrate backend go code with database schema.  

Packages named `model<xxx>` contain the go description of the database objects.  There might be multiple independent classes of objects in the same database, each will have a separate package, such as "models" for the budget Tracker objects or "models_test" for objects used for testing.

In addition to the types describing each object, each model package has a "Service" type with methods for each of the database queries to be performed.  gorm API calls are only used within the models package and external callers rely instead on methods of Service for  all database operations that the app can perform.

## Front end
Embedded WebView using SvelteKit, tailwind, shadcn-svelte, etc. from the starter template.

The main window has a top menubar with each of the main functions, such as "Import transactions", "Reports" and "Adjust Budget".  Each main menuitem will have multiple sub-items for specific operations.
