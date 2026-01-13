<script lang="ts">
    import { onMount, untrack } from 'svelte';
    import { DataTable } from 'datatable';
    import type { DataTableConfig, DataSourceCallback } from 'datatable';
    // Ideally this import exists after Wails rebuild
    import { GetTransactionsPaginated } from '$wailsjs/go/main/App';
    import { models } from '$wailsjs/go/models';

    const config: DataTableConfig = {
        name: 'transactions_grid',
        keyColumn: 'ID', // Maps to Go struct field ID (json:"ID")
        title: 'Transactions',
        maxVisibleRows: 20,
        isFilterable: true,
        isFindable: true,
        columns: [
            { name: 'PostedDate', title: 'Date', isSortable: true, justify: 'left' },
            { name: 'AccountID', title: 'Account', isSortable: true, justify: 'center' },
            { name: 'Amount', title: 'Amount', isSortable: true, justify: 'right', formatter: (v) => (v/100).toFixed(2) },
            { name: 'Description', title: 'Description', isSortable: true, wrappable: 'word', maxLines: 3, maxWidth:1000},
            { name: 'Beneficiary', title: 'Beneficiary', isSortable: true, justify: 'center' },
            { name: 'BudgetLine', title: 'Budget Line', isSortable: true },
            { name: 'Tag', title: 'Tag', isSortable: true }
        ]
    };

    // State for client-side filtering + backend pagination
    let cachedMatches: any[] = [];
    let backendOffset = 0;
    let backendHasMore = true;
    
    // Change detection
    let lastSortJSON = "";
    let lastFilterJSON = "";

    const dataSource: DataSourceCallback = async (columnKeys, startRow, numRows, sortKeys, filters) => {
        const sortJSON = JSON.stringify(sortKeys);
        const filterJSON = JSON.stringify(filters);

        // Reset if context changes
        if (sortJSON !== lastSortJSON || filterJSON !== lastFilterJSON) {
            cachedMatches = [];
            backendOffset = 0;
            backendHasMore = true;
            lastSortJSON = sortJSON;
            lastFilterJSON = filterJSON;
        }

        // Map SortKeys
        const goSortKeys: models.SortOption[] = sortKeys.map(k => {
             return { key: k.key, direction: k.direction } as models.SortOption; 
        });

        // Filter Function
        const globalTerm = filters?.global?.toLowerCase() || "";
        const isMatch = (row: any) => {
             if (!globalTerm) return true;
             // Search common text fields
             // Note: Depending on row structure (Go struct), fields might be Capitalized.
             // We can check values.
             return Object.values(row).some(v => String(v).toLowerCase().includes(globalTerm));
        };

        // Fetch Loop
        const BATCH_SIZE = 100;
        const neededEnd = startRow + numRows;

        while (cachedMatches.length < neededEnd && backendHasMore) {
            try {
                // Fetch next batch from backend (sorted, but NOT filtered)
                const batch = await GetTransactionsPaginated(backendOffset, BATCH_SIZE, goSortKeys);
                
                if (batch && batch.length > 0) {
                    // Filter in memory
                    const matches = batch.filter(isMatch);
                    cachedMatches.push(...matches);
                    
                    backendOffset += batch.length;
                    
                    // If batch was full, we assume there MIGHT be more. 
                    // If batch < BATCH_SIZE, we are definitely done.
                    if (batch.length < BATCH_SIZE) {
                        backendHasMore = false;
                    }
                } else {
                    backendHasMore = false;
                }
            } catch (e) {
                console.error("Backend fetch error:", e);
                backendHasMore = false; 
                break;
            }
        }

        // Return slice of what we found
        return cachedMatches.slice(startRow, neededEnd);
    };

</script>


<div class="h-[calc(100vh-100px)] w-full p-4">
    <DataTable {config} {dataSource} />
</div>
