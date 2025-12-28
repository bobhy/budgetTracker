# Transaction Import Implementation Plan

This plan outlines the implementation of the CSV import feature, following the logic defined in `.agent/rules/import_plan.md` (Update 2).

## 1. Backend Implementation

### 1.1. Data Models (`models/models.go`)
- [ ] **Transaction & RawTransaction Structs**:
    - Ensure `Amount` is `Money` (int64 cents). Rules: Positive = Expense, Negative = Income.
    - Fields:
        - `PostedDate`, `AccountID`, `Amount`, `Description`, `Tag`.
        - `Beneficiary` (string, FK to Beneficiary.Name).
        - `BudgetLine` (string, FK to Budget.Name).
        - `RawHint` (string).
    - `RawTransaction` adds `Action` ("add", "update").

### 1.2. Database Configuration (`models/service.go`)
- [ ] **Seeding Logic**:
    - Modify `Clean()` or add `Seed()` to insert the required production data:
        - **Beneficiaries**: "Bob", "Jessie", "Us".
        - **Accounts**:
            - "CapitalOne" (Format: "CapitalOne")
            - "WfChecking" (Format: "WFChecking")
            - "WfVisa" (Format: "WFVisa")
    - Remove `GenerateBeneficiaries`, `GenerateAccounts` methods if no longer needed, or just deprioritize them (Plan says "Remove").

### 1.3. Import Logic (`transactionImport/`)
- [ ] **Parsers (`parsers.go`)**:
    - `CapitalOneParser`:
        - Header mapping.
        - Beneficiary logic: Card 3028 -> Bob, 6539 -> Jessie, Else -> Us.
        - Amount logic: `Date * 100 - Credit * 100`. (Positive for Debit/Expense).
    - `WFCheckingParser`:
        - No header.
        - Amount logic: `-1 * value * 100`.
    - `WFVisaParser`:
        - No header.
        - Amount logic: `-1 * value * 100`.
- [ ] **Processor (`processor.go`)**:
    - Idempotency logic using `RawTransaction` table as staging.
    - Matching logic against `Transaction` table to determine "add" vs "update".

## 2. Frontend Implementation

### 2.1. Views
- [ ] **Import View**:
    - Dropdown for Account (populated from DB).
    - File Upload.
    - "Import" button.
    - `RawTransactionsTable` with sorting (Asc/Desc/None).
- [ ] **Forms**:
    - `RawTransactionForm` (Modal) for editing staging data.

## 3. Execution Steps
1. **Models**: Update structs in `models.go`.
2. **Service**: Update `Clean`/Seed logic in `service.go` / `app.go`.
3. **Parsers**: Implement specific bank parsers in `transactionImport/`.
4. **Processor**: Implement matching/upsert logic.
5. **Frontend**: Build the Import UI.
