<script lang="ts">
    import { DataTable } from "datatable";
    import type {
        DataTableConfig,
        DataSourceCallback,
        RowEditAction,
        RowEditResult,
    } from "datatable";
    import { models } from "$wailsjs/go/models";
    import * as Service from "$wailsjs/go/models/Service";

    const config: DataTableConfig = {
        name: "beneficiaries_grid",
        keyColumn: "name",
        title: "Beneficiaries",
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
        return await Service.GetBeneficiariesPaginated(
            startRow,
            numRows,
            goSortKeys,
        );
    };

    const handleRowEdit = async (
        action: RowEditAction,
        row: any,
        oldRow: any,
    ): Promise<RowEditResult> => {
        try {
            if (action === "update") {
                await Service.UpdateBeneficiary(oldRow, row);
            } else if (action === "create") {
                await Service.AddBeneficiary(row);
            } else if (action === "delete") {
                await Service.DeleteBeneficiary(oldRow);
            }
            return true;
        } catch (e) {
            console.error(`Beneficiary ${action} failed:`, e);
            return { error: String(e) };
        }
    };
</script>

<div class="h-[calc(100vh-100px)] w-full p-4">
    <DataTable {config} {dataSource} onRowEdit={handleRowEdit} />
</div>
