import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
    testDir: './tests',
    fullyParallel: false, // Serial for database tests
    forbidOnly: !!process.env.CI,
    retries: 0,
    workers: 1, // Single worker to prevent database conflicts
    reporter: 'html',

    use: {
        baseURL: 'http://localhost:34115',
        trace: 'on-first-retry',
        screenshot: 'only-on-failure',
    },

    projects: [
        {
            name: 'chromium',
            use: { ...devices['Desktop Chrome'] },
        },
    ],
});
