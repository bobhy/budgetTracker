<script lang="ts">
    import { onMount, tick } from "svelte";
    import { DataTable } from "datatable";
    import type {
        DataTableConfig,
        DataSourceCallback,
        RowEditAction,
        RowEditResult,
    } from "datatable";
    import { models } from "$wailsjs/go/models";
    import * as Service from "$wailsjs/go/models/Service";
    import { toast } from "svelte-sonner";
    import { parseMoney } from "$lib/money";

    // Options for Edit Form
    let accountOptions = $state<string[]>([]);
    let budgetOptions = $state<string[]>([]);
    let beneficiaryOptions = $state<string[]>([]);

    onMount(async () => {
        await loadAccounts();
        await loadOptions();
    });

    async function loadAccounts() {
        try {
            const accounts = (await Service.GetAccounts()) || [];
            accountOptions = accounts.map((a: any) => a.Name);
        } catch (err) {
            toast.error("Failed to load accounts: " + err);
        }
    }

    async function loadOptions() {
        try {
            const buds = (await Service.GetBudgets()) || [];
            budgetOptions = buds.map((b: any) => b.Name);
            const bens = (await Service.GetBeneficiaries()) || [];
            beneficiaryOptions = bens.map((b: any) => b.Name);
        } catch (err) {
            toast.error("Failed to load budgets, beneficiaries: " + err);
        }
    }

    const config: DataTableConfig = {
        name: "transactions_grid",
        keyColumn: "ID",
        title: "Transactions",
        maxVisibleRows: 20,
        isFilterable: true,
        isFindable: true,
        isEditable: true,
        columns: [
            {
                name: "PostedDate",
                title: "Date",
                isSortable: true,
                justify: "center",
            },
            {
                name: "Account",
                isSortable: true,
                justify: "center",
                enumValues: () => accountOptions,
            },
            {
                name: "Amount",
                isSortable: true,
                justify: "right",
                formatter: (v) => (v / 100).toFixed(2),
            },
            {
                name: "Description",
                isSortable: true,
                wrappable: "word",
                maxLines: 3,
                maxChars: 20,
            },
            {
                name: "Beneficiary",
                isSortable: true,
                justify: "center",
                enumValues: () => beneficiaryOptions,
            },
            {
                name: "Budget",
                title: "Budget Line",
                enumValues: () => budgetOptions,
                isSortable: true,
            },
            { name: "Tag", isSortable: true },
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
            // Ensure numeric fields are numbers
            if (row.Amount) row.Amount = parseMoney(row.Amount);

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
