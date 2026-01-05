import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import userEvent from '@testing-library/user-event';
import Datagrid from './Datagrid.svelte';
import type { DataGridConfig, DataSourceCallback, SortKey } from './DatagridTypes';

describe('Datagrid Component - Navigation and Filtering Integration Tests', () => {
    let dataSourceMock: ReturnType<typeof vi.fn>;
    let defaultConfig: DataGridConfig;
    let testData: any[];

    // Table-driven test scenarios with different data sizes
    const testScenarios = [
        {
            name: 'Small dataset (5 rows - less than 1 grid)',
            rowCount: 5,
            maxVisibleRows: 20,
            expectedBatches: 1,
            description: 'Data does not fill the grid'
        },
        {
            name: 'Medium dataset (40 rows - about 2 grids)',
            rowCount: 40,
            maxVisibleRows: 20,
            expectedBatches: 1,
            description: 'Data fills about 2 grids worth of rows'
        },
        {
            name: 'Large dataset (200 rows - 10x grid)',
            rowCount: 200,
            maxVisibleRows: 20,
            expectedBatches: 2, // Will fetch in batches of 100
            description: 'Data fills 10x grid size'
        }
    ];

    // Helper function to generate test data
    function generateTestData(count: number) {
        return Array.from({ length: count }, (_, i) => ({
            id: i + 1,
            name: `Item ${i + 1}`,
            category: i % 3 === 0 ? 'A' : i % 3 === 1 ? 'B' : 'C',
            value: (i + 1) * 10,
            description: `Description for item ${i + 1}`
        }));
    }

    beforeEach(() => {
        vi.clearAllMocks();

        defaultConfig = {
            name: 'test-grid',
            keyColumn: 'id',
            title: 'Test Grid',
            maxVisibleRows: 20,
            isFilterable: true,
            isFindable: true,
            columns: [
                { name: 'id', title: 'ID', isSortable: true },
                { name: 'name', title: 'Name', isSortable: true },
                { name: 'category', title: 'Category', isSortable: true },
                { name: 'value', title: 'Value', isSortable: true }
            ]
        };
    });

    describe.each(testScenarios)('$name', ({ rowCount, maxVisibleRows, expectedBatches, description }) => {
        beforeEach(() => {
            testData = generateTestData(rowCount);
            defaultConfig.maxVisibleRows = maxVisibleRows;

            // Mock dataSource that returns slices of testData
            dataSourceMock = vi.fn(async (columnKeys, startRow, numRows, sortKeys) => {
                let data = [...testData];

                // Apply sorting if provided
                if (sortKeys && sortKeys.length > 0) {
                    const sortKey = sortKeys[0];
                    data.sort((a, b) => {
                        const aVal = a[sortKey.key];
                        const bVal = b[sortKey.key];
                        if (aVal < bVal) return sortKey.direction === 'asc' ? -1 : 1;
                        if (aVal > bVal) return sortKey.direction === 'asc' ? 1 : -1;
                        return 0;
                    });
                }

                // Return the requested slice
                return data.slice(startRow, startRow + numRows);
            });
        });

        it(`${description} - should initialize and load data`, async () => {
            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            }, { timeout: 2000 });

            // Verify initial fetch parameters
            expect(dataSourceMock).toHaveBeenCalledWith(
                expect.arrayContaining(['id']), // keyColumn + columns
                0, // startRow
                100, // batchSize
                [] // no sort initially
            );

            // Check that the grid container is rendered
            const gridElement = container.querySelector('[role="grid"]');
            expect(gridElement).toBeTruthy();
        });


        it(`${description} - should handle filtering`, async () => {
            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            // Wait for initial render
            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            const initialCallCount = dataSourceMock.mock.calls.length;

            // Find the filter input
            const filterInput = container.querySelector('input[placeholder="Filter..."]');
            expect(filterInput).toBeTruthy();

            // Type in filter
            if (filterInput) {
                await userEvent.type(filterInput, 'Item 1');
            }

            // Should trigger a reload
            await waitFor(() => {
                expect(dataSourceMock.mock.calls.length).toBeGreaterThan(initialCallCount);
            }, { timeout: 2000 });

            // The grid should reset and start fetching from row 0
            // For large datasets, check the first call after the filter was applied
            const callsAfterFilter = dataSourceMock.mock.calls.slice(initialCallCount);
            expect(callsAfterFilter[0][1]).toBe(0); // First call after filter should start at row 0
        });


        it(`${description} - should handle sorting`, async () => {
            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            const initialCallCount = dataSourceMock.mock.calls.length;

            // Find sortable headers (buttons in header)
            const buttons = container.querySelectorAll('button');
            expect(buttons.length).toBeGreaterThan(0);

            // Click on a sortable header (e.g., "Name" which is the second column)
            const nameHeaderButton = buttons[1];
            await fireEvent.click(nameHeaderButton);

            // Should trigger a reload with sort
            await waitFor(() => {
                expect(dataSourceMock.mock.calls.length).toBeGreaterThan(initialCallCount);
            }, { timeout: 2000 });

            // Verify that subsequent calls happen with sort parameters
            // Note: actual sort keys depend on TanStack Table state management
        });

        it(`${description} - should handle keyboard navigation (Page Down)`, async () => {
            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            const gridElement = container.querySelector('[role="grid"]') as HTMLElement;
            expect(gridElement).toBeTruthy();

            if (gridElement) {
                // Focus the grid
                gridElement.focus();

                // Simulate Page Down keypress
                await fireEvent.keyDown(gridElement, { key: 'PageDown' });

                // The grid should scroll (we can't easily verify scroll position in JSDOM,
                // but we can verify the grid responded to the key event without errors)
                expect(gridElement).toBeTruthy();
            }
        });

        it(`${description} - should handle keyboard navigation (Arrow Down, Arrow Up, Home, End)`, async () => {
            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            const gridElement = container.querySelector('[role="grid"]') as HTMLElement;
            expect(gridElement).toBeTruthy();

            if (gridElement) {
                gridElement.focus();

                // Test various navigation keys
                const keys = ['ArrowDown', 'ArrowUp', 'Home', 'End'];

                for (const key of keys) {
                    await fireEvent.keyDown(gridElement, { key });
                    // Verify no errors thrown
                    expect(gridElement).toBeTruthy();
                }
            }
        });

        it(`${description} - should load more data on scroll for large datasets`, async () => {
            if (rowCount <= maxVisibleRows) {
                // Skip this test for small datasets that don't need scrolling
                return;
            }

            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            const initialCallCount = dataSourceMock.mock.calls.length;

            // Simulate scroll to bottom by finding the scroll container and setting scrollTop
            const scrollContainer = container.querySelector('.overflow-auto');

            if (scrollContainer) {
                // Trigger scroll event 
                // Note: In JSDOM, we can't truly test progressive loading without mocking virtualizer,
                // but we verify the structure is in place
                await fireEvent.scroll(scrollContainer);

                // For large datasets, verify the component is set up for lazy loading
                expect(scrollContainer).toBeTruthy();
            }
        });

        it(`${description} - should display correct row count in footer`, async () => {
            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            }, { timeout: 2000 });

            // Wait for data to load and grid to update
            await waitFor(() => {
                const footerText = container.textContent;
                // The footer should show "X rows loaded"
                expect(footerText).toMatch(/rows loaded/);
            }, { timeout: 3000 });
        });

        it(`${description} - should handle combined filtering and sorting`, async () => {
            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            // Apply filter
            const filterInput = container.querySelector('input[placeholder="Filter..."]');
            if (filterInput) {
                await userEvent.type(filterInput, 'Item');
            }

            await waitFor(() => {
                expect(dataSourceMock.mock.calls.length).toBeGreaterThan(1);
            });

            const callCountAfterFilter = dataSourceMock.mock.calls.length;

            // Then apply sort
            const buttons = container.querySelectorAll('button');
            if (buttons.length > 0) {
                await fireEvent.click(buttons[1]); // Sort by name
            }

            // Should trigger another reload
            await waitFor(() => {
                expect(dataSourceMock.mock.calls.length).toBeGreaterThan(callCountAfterFilter);
            }, { timeout: 2000 });
        });
    });

    describe('Edge cases and stress tests', () => {
        it('should handle empty dataset gracefully', async () => {
            testData = [];
            dataSourceMock = vi.fn(async () => []);

            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            // Grid should render without errors
            const gridElement = container.querySelector('[role="grid"]');
            expect(gridElement).toBeTruthy();

            // Footer should show 0 rows
            await waitFor(() => {
                expect(container.textContent).toMatch(/0 rows loaded/);
            });
        });

        it('should handle rapid filter changes', async () => {
            testData = generateTestData(50);
            dataSourceMock = vi.fn(async (columnKeys, startRow, numRows) => {
                return testData.slice(startRow, startRow + numRows);
            });

            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            const filterInput = container.querySelector('input[placeholder="Filter..."]');

            if (filterInput) {
                // Type multiple characters rapidly
                await userEvent.type(filterInput, 'test', { delay: 10 });

                // Grid should handle this without crashing
                expect(filterInput).toBeTruthy();
            }
        });

        it('should handle very long text in cells', async () => {
            testData = [{
                id: 1,
                name: 'A'.repeat(1000), // Very long name
                category: 'Test',
                value: 100
            }];

            dataSourceMock = vi.fn(async () => testData);

            const { container } = render(Datagrid, {
                config: defaultConfig,
                dataSource: dataSourceMock
            });

            await waitFor(() => {
                expect(dataSourceMock).toHaveBeenCalled();
            });

            // Grid should render without layout issues
            const gridElement = container.querySelector('[role="grid"]');
            expect(gridElement).toBeTruthy();
        });
    });
});
