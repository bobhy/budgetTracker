import { exec } from 'child_process';
import { promisify } from 'util';
import * as fs from 'fs/promises';

const execAsync = promisify(exec);

export interface TestData {
    beneficiaries: Array<{ name: string }>;
    accounts: Array<{ name: string; description: string; beneficiary: string }>;
    budgets: Array<{ name: string; description: string; beneficiary: string; amount: number; interval_months: number }>;
    transactions: Array<{ posted_date: string; account_id: string; amount: number; description: string; beneficiary: string; budget_line: string }>;
}

export const DEFAULT_TEST_DATA: TestData = {
    beneficiaries: [
        { name: 'Alice' },
        { name: 'Bob' },
        { name: 'Us' },
    ],
    accounts: [
        { name: 'CapitalOne', description: 'Capital One Rewards Credit', beneficiary: 'Us' },
        { name: 'WfChecking', description: 'Wells Fargo Checking', beneficiary: 'Us' },
        { name: 'WfVisa', description: 'Wells Fargo Visa', beneficiary: 'Us' },
        { name: 'AliceCard', description: 'Alice Personal Card', beneficiary: 'Alice' },
        { name: 'BobSavings', description: 'Bob Savings Account', beneficiary: 'Bob' },
    ],
    budgets: [
        { name: 'groceries', description: 'Food and household', beneficiary: 'Us', amount: 50000, interval_months: 1 },
        { name: 'travel_us', description: 'Shared travel fund', beneficiary: 'Us', amount: 200000, interval_months: 1 },
        { name: 'travel_alice', description: 'Alice personal travel', beneficiary: 'Alice', amount: 50000, interval_months: 1 },
        { name: 'boat', description: 'Bob boat expenses', beneficiary: 'Bob', amount: 100000, interval_months: 1 },
        { name: 'utilities', description: 'Monthly utilities', beneficiary: 'Us', amount: 30000, interval_months: 1 },
        { name: 'entertainment', description: 'Movies, dining out', beneficiary: 'Us', amount: 20000, interval_months: 1 },
    ],
    transactions: [
        { posted_date: '2024-01-05', account_id: 'CapitalOne', amount: -5432, description: 'Whole Foods Market', beneficiary: 'Us', budget_line: 'groceries' },
        { posted_date: '2024-01-07', account_id: 'WfChecking', amount: -12500, description: 'Electric Company', beneficiary: 'Us', budget_line: 'utilities' },
        { posted_date: '2024-01-10', account_id: 'WfVisa', amount: -8900, description: 'Restaurant Downtown', beneficiary: 'Us', budget_line: 'entertainment' },
        { posted_date: '2024-01-12', account_id: 'AliceCard', amount: -45000, description: 'Flight to Paris', beneficiary: 'Alice', budget_line: 'travel_alice' },
        { posted_date: '2024-01-15', account_id: 'CapitalOne', amount: -6780, description: 'Trader Joes', beneficiary: 'Us', budget_line: 'groceries' },
        { posted_date: '2024-01-18', account_id: 'BobSavings', amount: -23000, description: 'Boat maintenance', beneficiary: 'Bob', budget_line: 'boat' },
        { posted_date: '2024-01-20', account_id: 'WfVisa', amount: -3200, description: 'Movie theater', beneficiary: 'Us', budget_line: 'entertainment' },
        { posted_date: '2024-01-22', account_id: 'CapitalOne', amount: -15000, description: 'Gas bill', beneficiary: 'Us', budget_line: 'utilities' },
        { posted_date: '2024-01-25', account_id: 'WfChecking', amount: -5600, description: 'Grocery store', beneficiary: 'Us', budget_line: 'groceries' },
        { posted_date: '2024-01-28', account_id: 'WfVisa', amount: -120000, description: 'Hotel booking', beneficiary: 'Us', budget_line: 'travel_us' },
        { posted_date: '2024-02-01', account_id: 'CapitalOne', amount: -4300, description: 'Coffee shop', beneficiary: 'Us', budget_line: 'entertainment' },
        { posted_date: '2024-02-03', account_id: 'BobSavings', amount: -18000, description: 'Boat insurance', beneficiary: 'Bob', budget_line: 'boat' },
    ],
};

/**
 * Seeds the test database using Go's database seeding
 */
export async function seedTestDatabase(data: TestData = DEFAULT_TEST_DATA): Promise<void> {
    const dbPath = '/tmp/budgetTracker_test.db';

    // Ensure we start with a clean state
    await clearTestDatabase();

    // Create SQL to seed database
    const sql: string[] = [];

    // Insert beneficiaries
    for (const ben of data.beneficiaries) {
        sql.push(`INSERT INTO beneficiaries (name) VALUES ('${ben.name}');`);
    }

    // Insert accounts
    for (const acc of data.accounts) {
        sql.push(`INSERT INTO accounts (name, description, beneficiary) VALUES ('${acc.name}', '${acc.description}', '${acc.beneficiary}');`);
    }

    // Insert budgets
    for (const budget of data.budgets) {
        sql.push(`INSERT INTO budgets (name, description, beneficiary, amount, interval_months) VALUES ('${budget.name}', '${budget.description}', '${budget.beneficiary}', ${budget.amount}, ${budget.interval_months});`);
    }

    // Insert transactions
    for (const txn of data.transactions) {
        sql.push(`INSERT INTO transactions (posted_date, account_id, amount, description, beneficiary, budget_line, tag, raw_hint) VALUES ('${txn.posted_date}', '${txn.account_id}', ${txn.amount}, '${txn.description}', '${txn.beneficiary}', '${txn.budget_line}', '', '');`);
    }

    // Write SQL to temp file
    const sqlFile = '/tmp/seed_budget_test.sql';
    await fs.writeFile(sqlFile, sql.join('\n'));

    // Use sqlite3 to execute the SQL
    // First, we need to let the app create the schema by running briefly
    // This is handled by the webServer in playwright.config.ts which will start the app
    // For now, we'll use the Go app to initialize the DB via a brief connection

    // Alternative: directly create and seed using sqlite3
    await execAsync(`sqlite3 ${dbPath} < ${sqlFile}`);

    console.log(`Test database seeded at ${dbPath}`);
}

/**
 * Clears the test database
 */
export async function clearTestDatabase(): Promise<void> {
    const dbPath = '/tmp/budgetTracker_test.db';

    try {
        // Clear all tables
        // Order matters due to foreign keys if they are enabled
        const sql = `
            PRAGMA foreign_keys = OFF;
            DELETE FROM transactions;
            DELETE FROM budgets;
            DELETE FROM accounts;
            DELETE FROM beneficiaries;
            DELETE FROM raw_transactions;
            DELETE FROM sqlite_sequence; -- Reset autoincrement counters
            PRAGMA foreign_keys = ON;
        `;

        await execAsync(`sqlite3 ${dbPath} "${sql}"`);
        console.log('Test database cleared via SQL');
    } catch (err) {
        console.error('Error clearing database:', err);
        // Fallback or ignore if DB doesn't exist yet (first run)
    }
}
