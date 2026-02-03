<script lang="ts">
    import { DataTable } from "datatable";
    import type {
        DataTableConfig,
        DataSourceCallback,
        RowAction,
        RowEditResult,
    } from "datatable";
    import {
        GetTransactionsPaginated,
        AddTransaction,
        UpdateTransaction,
        DeleteTransaction,
    } from "$wailsjs/go/main/App";
    import { models } from "$wailsjs/go/models";

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
        return await GetTransactionsPaginated(startRow, numRows, goSortKeys);
    };

    const handleRowEdit = async (
        action: RowAction,
        row: any,
    ): Promise<RowEditResult> => {
        try {
            if (action === "update") {
                await UpdateTransaction(
                    row.id,
                    row.posted_date,
                    row.account_id,
                    row.amount,
                    row.description,
                    row.tag,
                    row.beneficiary,
                    row.budget_line,
                    row.raw_hint || "",
                );
            } else if (action === "create") {
                await AddTransaction(
                    row.posted_date,
                    row.account_id,
                    row.amount,
                    row.description,
                    row.tag,
                    row.beneficiary,
                    row.budget_line,
                    row.raw_hint || "",
                );
            } else if (action === "delete") {
                await DeleteTransaction(row.id);
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
