// Integration tests for DataTable component (renamed from Datagrid)
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import userEvent from '@testing-library/user-event';
import DataTable from './DataTable.svelte';
import type { DataTableConfig, DataSourceCallback, SortKey } from './DataTableTypes';

describe('DataTable Component - Navigation and Filtering Integration Tests', () => {
    let dataSourceMock: ReturnType<typeof vi.fn>;
    let defaultConfig: DataTableConfig;
    let testData: any[];

    const testScenarios = [
        { name: 'Small dataset (5 rows)', rowCount: 5, maxVisibleRows: 20 },
        { name: 'Medium dataset (40 rows)', rowCount: 40, maxVisibleRows: 20 },
        { name: 'Large dataset (200 rows)', rowCount: 200, maxVisibleRows: 20 }
    ];

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

    describe.each(testScenarios)('$name', ({ rowCount, maxVisibleRows }) => {
        beforeEach(() => {
            testData = generateTestData(rowCount);
            defaultConfig.maxVisibleRows = maxVisibleRows;
            dataSourceMock = vi.fn(async (columnKeys, startRow, numRows, sortKeys) => {
                let data = [...testData];
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
                return data.slice(startRow, startRow + numRows);
            });
        });

        it('should initialize and load data', async () => {
            const { container } = render(DataTable, { config: defaultConfig, dataSource: dataSourceMock });
            await waitFor(() => expect(dataSourceMock).toHaveBeenCalled(), { timeout: 2000 });
            const gridElement = container.querySelector('[role="grid"]');
            expect(gridElement).toBeTruthy();
        });

        it('should handle filtering', async () => {
            const { container } = render(DataTable, { config: defaultConfig, dataSource: dataSourceMock });
            await waitFor(() => expect(dataSourceMock).toHaveBeenCalled());
            const initialCalls = dataSourceMock.mock.calls.length;
            const filterInput = container.querySelector('input[placeholder="Filter..."]');
            expect(filterInput).toBeTruthy();
            if (filterInput) await userEvent.type(filterInput, 'Item 1');
            await waitFor(() => expect(dataSourceMock.mock.calls.length).toBeGreaterThan(initialCalls), { timeout: 2000 });
        });

        // Additional tests (sorting, keyboard navigation, etc.) can be added similarly.
    });
});
