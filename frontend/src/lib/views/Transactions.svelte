<script lang="ts">
    import { onMount } from 'svelte';
    import DataTable from '$lib/components/ui/datagrid/DataTable.svelte';
    import type { DataTableConfig, DataSourceCallback } from '$lib/components/ui/datagrid/DataTableTypes';
    import { GetTransactions } from '$wailsjs/go/main/App';

    let allTransactions: any[] = $state([]);

    const config: DataTableConfig = {
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
            const sk = sortKeys[0]; // DataTable only passes one sort key usually? Or array?
            result.sort((a, b) => {
                const valA = a[sk.key];
                const valB = b[sk.key];
                if (valA < valB) return sk.direction === 'asc' ? -1 : 1;
                if (valA > valB) return sk.direction === 'asc' ? 1 : -1;
                return 0;
            });
        } else {
            result.sort((a, b) => {
                if (a.PostedDate < b.PostedDate) return 1;
                if (a.PostedDate > b.PostedDate) return -1;
                return 0;
            });
        }

        const endRow = Math.min(startRow + numRows, result.length);
        const slice = result.slice(startRow, endRow);
        return slice;
    };
</script>

<div class="h-[calc(100vh-100px)] w-full p-4">
    {#if allTransactions.length > 0}
        <DataTable {config} {dataSource} />
    {:else}
        <div class="flex items-center justify-center h-full text-muted-foreground">
            Loading transactions...
        </div>
    {/if}
</div>
