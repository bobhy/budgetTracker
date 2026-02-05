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
        name: "transactions_grid",
        keyColumn: "id",
        title: "Transactions",
        maxVisibleRows: 20,
        isFilterable: true,
        isFindable: true,
        isEditable: true,
        columns: [
            {
                name: "posted_date",
                title: "Date",
                isSortable: true,
                justify: "center",
            },
            {
                name: "account_id",
                title: "Account",
                isSortable: true,
                justify: "center",
            },
            {
                name: "amount",
                title: "Amount",
                isSortable: true,
                justify: "right",
                formatter: (v) => (v / 100).toFixed(2),
            },
            {
                name: "description",
                title: "Description",
                isSortable: true,
                wrappable: "word",
                maxLines: 3,
                maxChars: 50,
            },
            {
                name: "beneficiary",
                title: "Beneficiary",
                isSortable: true,
                justify: "center",
            },
            { name: "budget_line", title: "Budget Line", isSortable: true },
            { name: "tag", title: "Tag", isSortable: true },
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
        return await Service.GetTransactionsPaginated(
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
                await Service.UpdateTransaction(oldRow, row);
            } else if (action === "create") {
                await Service.AddTransaction(row);
            } else if (action === "delete") {
                await Service.DeleteTransaction(oldRow);
            }
            return true;
        } catch (e) {
            console.error(`Transaction ${action} failed:`, e);
            return { error: String(e) };
        }
    };
</script>

<div class="h-[calc(100vh-100px)] w-full p-4">
    <DataTable {config} {dataSource} onRowEdit={handleRowEdit} />
</div>
