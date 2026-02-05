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
        beneficiaries = fetched.map((b: any) => b.name);
    });

    const config: DataTableConfig = {
        name: "budgets_grid",
        keyColumn: "name",
        title: "Budget Lines",
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
                name: "Beneficiary",
                isSortable: true,
                justify: "center",
                enumValues: () => beneficiaries,
            },
            {
                name: "Amount",
                isSortable: true,
                justify: "right",
            },
            {
                name: "Description",
                isSortable: true,
                justify: "left",
                wrappable: "word",
                maxLines: 2,
                maxChars: 40,
            },
            {
                name: "IntervalMonths",
                title: "Budget Period (mos)",
                isSortable: true,
                justify: "center",
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
        return await Service.GetBudgetsPaginated(startRow, numRows, goSortKeys);
    };

    const handleRowEdit = async (
        action: RowEditAction,
        row: any,
        oldRow: any,
    ): Promise<RowEditResult> => {
        try {
            if (action === "update") {
                await Service.UpdateBudget(oldRow, row);
            } else if (action === "create") {
                await Service.AddBudget(row);
            } else if (action === "delete") {
                await Service.DeleteBudget(oldRow);
            }
            return true;
        } catch (e) {
            console.error(`Budget ${action} failed:`, e);
            return { error: String(e) };
        }
    };
</script>

<div class="h-[calc(100vh-100px)] w-full p-4">
    <DataTable {config} {dataSource} onRowEdit={handleRowEdit} />
</div>
