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

    // Filter State
    let filterTerm = $state("");
    let findTerm = $state("");
    
    // Change detection
    let lastSortJSON = "";
    let lastFilterHelper = ""; // Combine filters to detect change

    const dataSource: DataSourceCallback = async (columnKeys, startRow, numRows, sortKeys) => {
        const sortJSON = JSON.stringify(sortKeys);
        const currentFilterHelper = filterTerm + "||" + findTerm;

        // Reset if context changes
        if (sortJSON !== lastSortJSON || currentFilterHelper !== lastFilterHelper) {
            cachedMatches = [];
            backendOffset = 0;
            backendHasMore = true;
            lastSortJSON = sortJSON;
            lastFilterHelper = currentFilterHelper;
        }

        // Map SortKeys
        const goSortKeys: models.SortOption[] = sortKeys.map(k => {
             return { key: k.key, direction: k.direction } as models.SortOption; 
        });

        // Filter Function
        const globalT = filterTerm.toLowerCase();
        const findT = findTerm.toLowerCase();
        
        const isMatch = (row: any) => {
             const rowValues = Object.values(row).map(v => String(v).toLowerCase());
             
        // Check Global Filter
             if (globalT && !rowValues.some(v => v.includes(globalT))) {
                 return false;
             }
             
             // Check Find Filter (REMOVED - Find is now navigation only)
             // if (findT && !rowValues.some(v => v.includes(findT))) {
             //     return false;
             // }
             
             return true;
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

    // Incremental Find Logic
    const handleFind = async (term: string, direction: 'next' | 'previous', currentIndex: number): Promise<{rowIndex: number, columnName: string} | null> => {
        if (!term) return null;
        const lowerTerm = term.toLowerCase();

        // Helper to find match in a row
        const findInRow = (index: number): string | null => {
            const row = cachedMatches[index];
            if (!row) return null;
            
            // Search only visible columns
            for (const col of config.columns) {
                const val = row[col.name];
                if (val !== undefined && val !== null && String(val).toLowerCase().includes(lowerTerm)) {
                    return col.name;
                }
            }
            return null;
        };

        if (direction === 'next') {
            // 1. Search forward in currently cached data
            for (let i = currentIndex + 1; i < cachedMatches.length; i++) {
                const col = findInRow(i);
                if (col) return { rowIndex: i, columnName: col };
            }
            
            // 2. If not found, fetch more from backend until we find match or end
            while (backendHasMore) {
                 const BATCH_SIZE = 100;
                 let sortKeys = [];
                 try { sortKeys = JSON.parse(lastSortJSON || "[]"); } catch {}
                 const goSortKeys = sortKeys.map((k: any) => ({ key: k.key, direction: k.direction }));
                 
                 try {
                     const batch = await GetTransactionsPaginated(backendOffset, BATCH_SIZE, goSortKeys as any);
                     if (batch && batch.length > 0) {
                         // Filter via global filter ONLY
                         const globalT = filterTerm.toLowerCase();
                         const isMatch = (row: any) => {
                             if (!globalT) return true;
                             const rowValues = Object.values(row).map(v => String(v).toLowerCase());
                             return rowValues.some(v => v.includes(globalT));
                         };
                         
                         const matches = batch.filter(isMatch);
                         const startIdx = cachedMatches.length;
                         cachedMatches.push(...matches);
                         backendOffset += batch.length;
                         if (batch.length < BATCH_SIZE) backendHasMore = false;
                         
                         // Search in the added batch
                         for (let i = startIdx; i < cachedMatches.length; i++) {
                             const col = findInRow(i);
                             if (col) return { rowIndex: i, columnName: col };
                         }
                     } else {
                         backendHasMore = false;
                     }
                 } catch (e) {
                     console.error("Find fetch error", e);
                     break;
                 }
            }
        } else {
            // Previous: Just search backwards in cache
            for (let i = currentIndex - 1; i >= 0; i--) {
                const col = findInRow(i);
                if (col) return { rowIndex: i, columnName: col };
            }
        }
        
        return null; // Not found
    };

</script>


<div class="h-[calc(100vh-100px)] w-full p-4">
    <DataTable {config} {dataSource} bind:globalFilter={filterTerm} bind:findTerm={findTerm} onFind={handleFind} />
</div>
