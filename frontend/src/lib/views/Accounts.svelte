<script lang="ts">
    import { onMount } from "svelte";
    import { DataTable } from "datatable";
    import type {
        DataTableConfig,
        DataSourceCallback,
        RowEditAction,
        RowEditResult,
    } from "datatable";
    import { models } from "$wailsjs/go/models";
    import * as Service from "$wailsjs/go/models/Service";

    let beneficiaries = $state<string[]>([]);

    onMount(async () => {
        const fetched = await Service.GetBeneficiaries();
        beneficiaries = fetched.map((b: any) => b.Name);
    });

    const config: DataTableConfig = {
        name: "accounts_grid",
        keyColumn: "Name",
        title: "Accounts",
        maxVisibleRows: 20,
        isFilterable: true,
        isFindable: true,
        isEditable: true,
        columns: [
            {
                name: "Name",
                isSortable: true,
                justify: "center",
            },
            {
                name: "Description",
                isSortable: true,
                justify: "left",
                wrappable: "word",
                maxLines: 3,
                maxChars: 20,
            },
            {
                name: "Beneficiary",
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
        return await Service.GetAccountsPaginated(
            startRow,
            numRows,
            goSortKeys,
        );
    };

    const handleRowEdit = async (
        action: RowEditAction,
        row: any,
        oldRow?: any,
        keyColumn?: string,
    ): Promise<RowEditResult> => {
        try {
            if (action === "update") {
                await Service.UpdateAccount(oldRow, row);
            } else if (action === "create") {
                await Service.AddAccount(row);
            } else if (action === "delete") {
                await Service.DeleteAccount(oldRow);
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
