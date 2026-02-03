import { test, expect } from '@playwright/test';

test.describe('Error Boundary', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/');
        await page.waitForLoadState('networkidle');
    });

    test('catches and displays global JavaScript errors', async ({ page }) => {
        // Inject an error by dispatching an ErrorEvent to window
        await page.evaluate(() => {
            const errorEvent = new ErrorEvent('error', {
                message: 'Test error message',
                filename: 'test.js',
                lineno: 1,
                colno: 1,
                error: new Error('Test error message'),
            });
            window.dispatchEvent(errorEvent);
        });

        // Wait for error modal to appear
        await page.waitForSelector('text="Application Error"', { timeout: 3000 });

        // Verify error message is displayed
        await expect(page.locator('text="Test error message"')).toBeVisible();

        // Verify "Check browser console" message exists
        await expect(page.locator('text="Check browser console for more details"')).toBeVisible();

        // Verify dismiss button exists
        const dismissBtn = page.locator('button:has-text("Dismiss")');
        await expect(dismissBtn).toBeVisible();

        // Click dismiss and verify error modal closes
        await dismissBtn.click();
        await expect(page.locator('text="Application Error"')).not.toBeVisible();
    });

    test('catches unhandled promise rejections', async ({ page }) => {
        // Dispatch an unhandledrejection event to window
        await page.evaluate(() => {
            const rejectionEvent = new PromiseRejectionEvent('unhandledrejection', {
                promise: Promise.reject(new Error('Async test error')),
                reason: new Error('Async test error'),
            });
            window.dispatchEvent(rejectionEvent);
        });

        // Wait for error boundary to catch it
        await page.waitForSelector('text="Application Error"', { timeout: 3000 });

        // Verify error content
        await expect(page.locator('text="Async test error"')).toBeVisible();

        // Verify stack trace details section is available
        await expect(page.locator('summary:has-text("Stack trace")')).toBeVisible();
    });

    test('shows stack trace in expandable details', async ({ page }) => {
        // Inject error with stack
        await page.evaluate(() => {
            const err = new Error('Error with stack');
            const errorEvent = new ErrorEvent('error', {
                message: 'Error with stack',
                error: err,
            });
            window.dispatchEvent(errorEvent);
        });

        // Wait for error modal
        await page.waitForSelector('text="Application Error"');

        // Verify stack trace summary exists
        const stackSummary = page.locator('summary:has-text("Stack trace")');
        await expect(stackSummary).toBeVisible();

        // Expand stack trace
        await stackSummary.click();

        // Verify stack trace content is visible (looks for <pre> tag)
        const stackContent = page.locator('details pre');
        await expect(stackContent).toBeVisible();
    });

    test('allows dismissing error and continuing to use app', async ({ page }) => {
        // Cause an error via ErrorEvent
        await page.evaluate(() => {
            const errorEvent = new ErrorEvent('error', {
                message: 'Recoverable error',
                error: new Error('Recoverable error'),
            });
            window.dispatchEvent(errorEvent);
        });

        // Wait for error
        await page.waitForSelector('text="Application Error"');

        // Dismiss error
        await page.click('button:has-text("Dismiss")');

        // Verify error is gone
        await expect(page.locator('text="Application Error"')).not.toBeVisible();

        // Verify can navigate to different views
        await page.click('button:has-text("Beneficiaries")');
        await page.waitForTimeout(500);

        // App should still be functional - look for navbar or main content
        const navbar = page.locator('nav, header, [role="navigation"]');
        await expect(navbar).toBeVisible();
    });
});
