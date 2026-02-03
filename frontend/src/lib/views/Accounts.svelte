<script lang="ts">
    import { onMount } from "svelte";
    import { DataTable } from "datatable";
    import type {
        DataTableConfig,
        DataSourceCallback,
        RowAction,
        RowEditResult,
    } from "datatable";
    import {
        GetAccountsPaginated,
        GetAccounts,
        AddAccount,
        UpdateAccount,
        DeleteAccount,
        GetBeneficiaries,
    } from "$wailsjs/go/main/App";
    import { models } from "$wailsjs/go/models";

    let beneficiaries = $state<string[]>([]);

    onMount(async () => {
        const fetched = await GetBeneficiaries();
        beneficiaries = fetched.map((b) => b.name);
    });

    const config: DataTableConfig = {
        name: "accounts_grid",
        keyColumn: "name",
        title: "Accounts",
        maxVisibleRows: 20,
        isFilterable: true,
        isFindable: true,
        isEditable: true,
        columns: [
            {
                name: "name",
                title: "Name",
                isSortable: true,
                justify: "center",
            },
            {
                name: "description",
                isSortable: true,
                justify: "left",
                wrappable: "word",
                maxLines: 3,
                maxChars: 40,
            },
            {
                name: "beneficiary_id",
                title: "Beneficiary",
                isSortable: true,
                justify: "center",
                enumValues: () => beneficiaries,
            },
        ],
    };

    const dataSource: DataSourceCallback = async (
        columnKeys,
        startRow,
        numRows,
        sortKeys,
    ) => {
        const goSortKeys: models.SortOption[] = sortKeys.map(
            (k) =>
                ({ key: k.key, direction: k.direction }) as models.SortOption,
        );
        return await GetAccountsPaginated(startRow, numRows, goSortKeys);
    };

    const handleRowEdit = async (
        action: RowAction,
        row: any,
    ): Promise<RowEditResult> => {
        try {
            if (action === "update") {
                await UpdateAccount(
                    row.name,
                    row.new_name || row.name,
                    row.description,
                    row.beneficiary_id,
                );
            } else if (action === "create") {
                await AddAccount(row.name, row.description, row.beneficiary_id);
            } else if (action === "delete") {
                await DeleteAccount(row.name);
            }
            return true;
        } catch (e) {
            console.error(`Account ${action} failed:`, e);
            return { error: String(e) };
        }
    };
</script>

<div class="h-[calc(100vh-100px)] w-full p-4">
    <DataTable {config} {dataSource} onRowEdit={handleRowEdit} />
</div>
