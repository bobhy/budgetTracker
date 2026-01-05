<script lang="ts">
    import { onMount } from 'svelte';
    import Datagrid from '$lib/components/ui/datagrid/Datagrid.svelte';
    import type { DataGridConfig, DataSourceCallback } from '$lib/components/ui/datagrid/DatagridTypes';
    import { GetTransactions } from '$wailsjs/go/main/App';

    let allTransactions: any[] = $state([]);

    const config: DataGridConfig = {
        name: 'transactions_grid',
        keyColumn: 'ID',
        title: 'Transactions',
        maxVisibleRows: 20, // Or whatever fits
        isFilterable: true,
        isFindable: true,
        columns: [
            { name: 'PostedDate', title: 'Date', isSortable: true, justify: 'left' },
            { name: 'AccountID', title: 'Account', isSortable: true },
            { name: 'Amount', title: 'Amount', isSortable: true, formatter: (v) => (v/100).toFixed(2) },
            { name: 'Description', title: 'Description', isSortable: true, wrappable: 'word' },
            { name: 'Beneficiary', title: 'Beneficiary', isSortable: true },
            { name: 'BudgetLine', title: 'Budget Line', isSortable: true },
            { name: 'Tag', title: 'Tag', isSortable: true }
        ]
    };

    onMount(async () => {
         // Load *all* transactions for now, as GetTransactions returns []models.Transaction
         // Ideally backend should support pagination, but for now we simulate it in dataSource
         allTransactions = await GetTransactions() || [];
    });

    const dataSource: DataSourceCallback = async (columnKeys, startRow, numRows, sortKeys) => {
        // Since backend doesn't seem to have paginated API exposed yet [based on limited exploration],
        // we filter/sort 'allTransactions' in memory.
        
        let result = [...allTransactions];

        // 1. Sort
        if (sortKeys.length > 0) {
            const sk = sortKeys[0]; // Datagrid only passes one sort key usually? Or array?
            // Types say array.
            result.sort((a, b) => {
                const valA = a[sk.key];
                const valB = b[sk.key];
                
                if (valA < valB) return sk.direction === 'asc' ? -1 : 1;
                if (valA > valB) return sk.direction === 'asc' ? 1 : -1;
                return 0;
            });
        } 
        // Default sort by Date desc if no sort?
        else {
             result.sort((a, b) => {
                if (a.PostedDate < b.PostedDate) return 1;
                if (a.PostedDate > b.PostedDate) return -1;
                return 0;
            });
        }

        // 2. Pagination
        // The Datagrid asks for [startRow ... startRow+numRows]
        // But wait, the Datagrid logic (client side filtering) *also* exists?
        // No, `Datagrid.svelte` line 116 "Filter locally".
        // The Datagrid assumes dataSource returns "raw rows" matching the query
        // BUT the Datagrid handles filtering "locally" on the batch it got?
        // Actually line 116 `const filtered = newRawRows.filter(...)` suggests Datagrid 
        // filters what it RECEIVED from dataSource.
        // If the dataSource returns paginated data (e.g. rows 0-100), and I type "Walmart",
        // and "Walmart" is in row 500, Datagrid will never see it if I don't scroll?
        // Wait, the requirements said: "Filter filters the data source... masking any row".
        // "Find scrolls the grid".
        // The Current Datagrid Implementation (Datagrid.svelte):
        // It fetches chunks.
        // It filters those chunks *locally* (L116).
        // If `hasMore` is true, it keeps fetching?
        // No, `performFetch` loop: `while (addedRows < wantedGridRows && hasMore...)`
        // So if I filter "Walmart", it will keep asking dataSource for more rows until it finds enough matches OR dataSource runs out.
        // THIS IS GOOD. It means dataSource just needs to provide "next N rows".
        // However, `dataSource` signature does NOT take a "filter string". 
        // So the backend (or this callback) doesn't know about the filter.
        // So `Datagrid` responsibilty is to scan the stream provided by dataSource.
        
        // So this callback just returns the slice of *sorted* data.
        
        const endRow = Math.min(startRow + numRows, result.length);
        const slice = result.slice(startRow, endRow);
        
        return slice;
    };

</script>

<div class="h-[calc(100vh-100px)] w-full p-4">
    {#if allTransactions.length > 0}
        <Datagrid {config} {dataSource} />
    {:else}
        <div class="flex items-center justify-center h-full text-muted-foreground">
            Loading transactions...
        </div>
    {/if}
</div>
