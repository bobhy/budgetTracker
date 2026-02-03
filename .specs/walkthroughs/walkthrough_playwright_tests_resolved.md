# Add Playwright E2E Tests to BudgetTracker

## Summary

Successfully added Playwright end-to-end tests to budgetTracker with test database seeding infrastructure and two comprehensive test suites covering:
1. Accounts view data loading
2. Error boundary functionality

## Implementation

### 1. Installed Playwright

```bash
npm install -D @playwright/test
npx playwright install chromium
```

---

### 2. Created Test Infrastructure

#### [playwright.config.ts](file:///home/bobhy/worktrees/budgetTracker/playwright.config.ts)

Configured Playwright with:
- Single worker to prevent database conflicts
- Test database via `budgetTracker_config` env var
- Wails dev server integration
- HTML reporter for results

#### [tests/config/test_config1.toml](file:///home/bobhy/worktrees/budgetTracker/tests/config/test_config1.toml)

Test-specific configuration pointing to `/tmp/budgetTracker_test.db`.

---

### 3. Created Test Database Seeding Helper

#### [tests/helpers/seedDatabase.ts](file:///home/bobhy/worktrees/budgetTracker/tests/helpers/seedDatabase.ts)

**Features:**
- [seedTestDatabase()](file:///home/bobhy/worktrees/budgetTracker/tests/helpers/seedDatabase.ts#51-101) - Populates test database with known data
- [clearTestDatabase()](file:///home/bobhy/worktrees/budgetTracker/tests/helpers/seedDatabase.ts#102-114) - Removes test database between runs
- `DEFAULT_TEST_DATA` - Predefined test dataset

**Test Data:**
- **3 Beneficiaries**: Alice, Bob, Us
- **5 Accounts**: CapitalOne, WfChecking, WfVisa, AliceCard, BobSavings
- **6 Budgets**: groceries, travel_us, travel_alice, boat, utilities, entertainment
- **12 Transactions**: Varied dates, amounts, and categories

---

### 4. Accounts View Tests

#### [tests/accounts-view.spec.ts](file:///home/bobhy/worktrees/budgetTracker/tests/accounts-view.spec.ts)

**4 test cases:**

1. **loads and displays account data**
   - Verifies all 5 seeded accounts appear in DataTable
   - Checks for CapitalOne, WfChecking, WfVisa, AliceCard, BobSavings

2. **displays correct number of account rows**
   - Counts rows in DataTable
   - Ensures at least 5 data rows present

3. **shows account descriptions**
   - Verifies descriptions visible
   - Checks "Capital One Rewards Credit", "Wells Fargo Checking"

4. **shows beneficiary assignments**
   - Confirms beneficiary column displays
   - Verifies "Us", "Alice", "Bob" appear

---

### 5. Error Boundary Tests

#### [tests/error-boundary.spec.ts](file:///home/bobhy/worktrees/budgetTracker/tests/error-boundary.spec.ts)

**4 test scenarios:**

1. **catches and displays global JavaScript errors**
   - Throws error via `page.evaluate()`
   - Verifies "Application Error" modal appears
   - Checks error message displays
   - Tests dismiss button functionality

2. **catches unhandled promise rejections**
   - Creates rejected promise
   - Confirms error boundary catches it
   - Verifies async error message shows

3. **shows stack trace in expandable details**
   - Verifies stack trace summary exists
   - Tests expand/collapse functionality
   - Checks stack content visibility

4. **allows dismissing error and continuing to use app**
   - Dismisses error
   - Navigates to different view
   - Confirms app remains functional

---

### 6. Added npm Scripts

Updated [package.json](file:///home/bobhy/worktrees/budgetTracker/package.json):

```json
"scripts": {
  "test:e2e": "playwright test",
  "test:e2e:ui": "playwright test --ui",
  "test:e2e:headed": "playwright test --headed",
  "test:e2e:report": "playwright show-report"
}
```

## Running Tests

```bash
# Run all tests
npm run test:e2e

# Run with UI mode (interactive)
npm run test:e2e:ui

# Run in headed mode (see browser)
npm run test:e2e:headed

# View last test report
npm run test:e2e:report
```

## Test Coverage

✅ **Database Seeding**: Reusable helper with rich test data  
✅ **Accounts View**: 4 tests covering data display and accuracy  
✅ **Error Boundary**: 4 tests for error catching and recovery  
✅ **Total**: 8 end-to-end tests

## Files Created

- **Config**: [playwright.config.ts](file:///home/bobhy/worktrees/budgetTracker/playwright.config.ts), `~/budgetTrackerTest/c2_test.toml`
- **Helpers**: [tests/helpers/seedDatabase.ts](file:///home/bobhy/worktrees/budgetTracker/tests/helpers/seedDatabase.ts)
- **Tests**: [tests/accounts-view.spec.ts](file:///home/bobhy/worktrees/budgetTracker/tests/accounts-view.spec.ts), [tests/error-boundary.spec.ts](file:///home/bobhy/worktrees/budgetTracker/tests/error-boundary.spec.ts)
- **Scripts**: Updated [package.json](file:///home/bobhy/worktrees/budgetTracker/package.json)

## Next Steps

Run the tests to verify everything works:
```bash
npm run test:e2e
```

Note: Tests require Wails dev server to start, which may take 30-60 seconds on first run.
