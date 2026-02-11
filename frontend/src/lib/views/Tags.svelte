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

    let budgets = $state<string[]>([]);

    onMount(async () => {
        const fetched = await Service.GetBudgets();
        budgets = fetched.map((b: any) => b.Name);
    });

    const config: DataTableConfig = {
        name: "tags_grid",
        keyColumn: "Name",
        title: "Tags",
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
                name: "Budget",
                isSortable: true,
                justify: "center",
                enumValues: () => budgets,
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
        return await Service.GetTagsPaginated(startRow, numRows, goSortKeys);
        //tags, _, err := models.GetPage[models.Tag](s.DB, startRow, numRows, goSortKeys, nil)
        //return tags, err
    };

    const handleRowEdit = async (
        action: RowEditAction,
        row: any,
        oldRow?: any,
        keyColumn?: string,
    ): Promise<RowEditResult> => {
        try {
            if (action === "update") {
                await Service.UpdateTag(oldRow, row);
            } else if (action === "create") {
                await Service.AddTag(row);
            } else if (action === "delete") {
                await Service.DeleteTag(oldRow);
            }
            return true;
        } catch (e) {
            console.error(`Tag ${action} failed:`, e);
            return { error: String(e) };
        }
    };
</script>

<div class="h-[calc(100vh-100px)] w-full p-4">
    <DataTable {config} {dataSource} onRowEdit={handleRowEdit} />
</div>
