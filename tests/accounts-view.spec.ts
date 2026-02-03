import { test, expect } from '@playwright/test';
import { seedTestDatabase, clearTestDatabase } from './helpers/seedDatabase';

test.describe('Accounts View', () => {
    test.beforeEach(async ({ page }) => {
        // Clear and seed database before each test
        await clearTestDatabase();
        await seedTestDatabase();

        // Navigate to app
        await page.goto('/');
        await page.waitForLoadState('networkidle');

        // Click Accounts button in navbar
        await page.click('button:has-text("Accounts")');
    });

    test('loads and displays account data', async ({ page }) => {
        // Wait for DataTable grid to appear
        await page.waitForSelector('[role="grid"]', { timeout: 10000 });

        // Wait a moment for data to load
        await page.waitForTimeout(1000);

        // Verify expected accounts are visible
        await expect(page.locator('text="CapitalOne"')).toBeVisible();
        await expect(page.locator('text="WfChecking"')).toBeVisible();
        await expect(page.locator('text="WfVisa"')).toBeVisible();
        await expect(page.locator('text="AliceCard"')).toBeVisible();
        await expect(page.locator('text="BobSavings"')).toBeVisible();
    });

    test('displays correct number of account rows', async ({ page }) => {
        // Wait for grid
        await page.waitForSelector('[role="grid"]');
        await page.waitForTimeout(1000);

        // Count data rows (excluding header)
        const rows = page.locator('[role="row"]');
        const rowCount = await rows.count();

        // Should have header row + 5 data rows
        expect(rowCount).toBeGreaterThanOrEqual(5);
    });

    test('shows account descriptions', async ({ page }) => {
        // Wait for grid
        await page.waitForSelector('[role="grid"]');
        await page.waitForTimeout(1000);

        // Verify descriptions are visible
        await expect(page.locator('text="Capital One Rewards Credit"')).toBeVisible();
        await expect(page.locator('text="Wells Fargo Checking"')).toBeVisible();
    });

    test('shows beneficiary assignments', async ({ page }) => {
        // Wait for grid
        await page.waitForSelector('[role="grid"]');
        await page.waitForTimeout(1000);

        // Verify beneficiaries are displayed
        await expect(page.locator('text="Us"').first()).toBeVisible();
        await expect(page.locator('text="Alice"')).toBeVisible();
        await expect(page.locator('text="Bob"')).toBeVisible();
    });
});
