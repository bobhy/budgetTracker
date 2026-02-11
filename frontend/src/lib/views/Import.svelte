<script lang="ts">
    import { onMount, tick } from "svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card";
    import { ImportFile, SelectFile, ApplyTags } from "$wailsjs/go/main/App";
    import { models } from "$wailsjs/go/models";
    import * as Service from "$wailsjs/go/models/Service";
    import { toast } from "svelte-sonner";
    import { parseMoney } from "$lib/money";
    import { DataTable, DataTableContainer } from "datatable";
    import type {
        DataTableConfig,
        DataSourceCallback,
        RowEditAction,
        RowEditResult,
    } from "datatable";

    // State using Runes
    let accounts = $state<any[]>([]);
    let selectedAccount = $state("");
    let filePath = $state("");
    let totalRawTransactions = $state(0);
    let loading = $state(false);
    let importStats = $state(""); // Keeping to support logic, but will not display bottom div. Actually, used below for finalize reporting?
    // Wait, requirement said "The view currently has a last div that's conditioned on importStats, that should be removed".
    // I removed the display block in previous edit.
    // Now implementing status logic.

    let statusMessage = $state("");
    let statusType = $state<"success" | "error">("success");

    // Options for Edit Form
    let accountOptions = $state<string[]>([]);
    let budgetOptions = $state<string[]>([]);
    let beneficiaryOptions = $state<string[]>([]);

    // Component ref for DataTable to force refresh
    let dataTableRef = $state<any>();

    onMount(async () => {
        await loadAccounts();
        await loadOptions();
        await loadRawTransactionCount();
    });

    async function loadAccounts() {
        try {
            accounts = (await Service.GetAccounts()) || [];
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

    async function loadRawTransactionCount() {
        try {
            totalRawTransactions = await Service.GetRawTransactionCount();
        } catch (err) {
            toast.error("Failed to load raw transactions count: " + err);
        }
    }

    async function handleSelectFile() {
        try {
            const path = await SelectFile();
            if (path) {
                filePath = path;
            }
        } catch (err) {
            toast.error("Error selecting file: " + err);
        }
    }

    async function handleImport() {
        console.log("[Import] Starting import process");
        if (!selectedAccount || !filePath) {
            console.warn("[Import] Aborted: Missing account or file path");
            return;
        }
        loading = true;
        await tick();
        try {
            console.log(
                `[Import] Calling backend ImportFile. Account: ${selectedAccount}, Path: ${filePath}`,
            );
            const msg = await ImportFile(selectedAccount, filePath);
            console.log("[Import] Backend response:", msg);
            toast.success(msg);
            await loadRawTransactionCount();
            dataTableRef?.refresh();
        } catch (err) {
            console.error("[Import] Error calling ImportFile:", err);
            toast.error("Import failed: " + err);
        } finally {
            loading = false;
        }
    }

    async function handleFinalize() {
        if (totalRawTransactions === 0) return;

        loading = true;
        statusMessage = ""; // clear previous
        await tick();
        try {
            const msg = await Service.FinalizeImport();
            statusMessage = msg;
            statusType = "success";
            await loadRawTransactionCount();
            dataTableRef?.refresh();
        } catch (err) {
            statusMessage = "Finalize failed: " + err;
            statusType = "error";
        } finally {
            loading = false;
        }
    }

    async function handleApplyTags() {
        loading = true;
        statusMessage = ""; // clear previous
        await tick();
        try {
            const msg = await ApplyTags();
            statusMessage = msg;
            statusType = "success";
            dataTableRef?.refresh();
        } catch (err) {
            statusMessage = "Auto-Tag failed: " + err;
            statusType = "error";
        } finally {
            loading = false;
        }
    }

    // DataTable Configuration
    const tableConfig: DataTableConfig = {
        name: "raw-transactions-table",
        title: "Staging Area",
        keyColumn: "ID",
        isEditable: true,
        isFilterable: true,
        columns: [
            {
                name: "PostedDate",
                title: "Date",
                isSortable: true,
                justify: "left",
            },
            {
                name: "Description",
                isSortable: true,
                justify: "left",
                wrappable: "word",
                maxChars: 10,
                maxLines: 1,
            },
            {
                name: "Amount",
                isSortable: true,
                justify: "right",
                formatter: (val: number) => (val / 100).toFixed(2),
            },
            {
                name: "Beneficiary",
                isSortable: true,
                justify: "left",
                enumValues: () => beneficiaryOptions,
            },
            {
                name: "Tag",
                isSortable: true,
                justify: "left",
                maxChars: 20,
            },
            {
                name: "Budget",
                title: "Budget",
                isSortable: true,
                justify: "left",
                enumValues: () => budgetOptions,
            },
            {
                name: "Action",
                title: "Status",
                isSortable: true,
                justify: "center",
                // Status is read-only in this view conceptually, but if edited, it's just text
            },
        ],
    };

    // Server-side data source wrapper
    //todo: handle error return from GetRawTransactions.
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
        return await Service.GetRawTransactionsPaginated(
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
                await Service.UpdateRawTransaction(oldRow, row);
            } else if (action === "delete") {
                await Service.DeleteRawTransaction(oldRow);
                // Count changes on delete
                await loadRawTransactionCount();
            } else if (action === "create") {
                // Not really creating raw transactions manually in this flow usually,
                // but if we did:
                await Service.AddRawTransaction(row);
                // Count changes on create
                await loadRawTransactionCount();
            }

            // Refresh local data to reflect changes
            // For pagination, usually we rely on the component re-fetching or we manually invalidate.
            // Returning true usually triggers a re-fetch if configured effectively?
            // Actually currently DataTable might just update the local row if we return true on update,
            // but for add/delete we definitely need a re-fetch.
            return true;
        } catch (e) {
            console.error(`RawTransaction ${action} failed:`, e);
            return { error: String(e) };
        }
    };
</script>

<div class="flex-1 w-full h-full flex flex-col gap-4 p-4">
    <Card.Root>
        <Card.Header>
            <Card.Title>Import Transactions</Card.Title>
            <Card.Description
                >Select an account and CSV file to import.</Card.Description
            >
        </Card.Header>
        <Card.Content class="space-y-4">
            <div class="flex items-end gap-4">
                <!-- Account Select -->
                <div class="grid gap-2">
                    <span class="text-sm font-medium">Account</span>
                    <select
                        class="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        bind:value={selectedAccount}
                    >
                        {#each accounts as acc}
                            <option value={acc.Name}>{acc.Name}</option>
                        {/each}
                    </select>
                </div>

                <!-- File Select -->
                <div class="grid gap-2 flex-1">
                    <span class="text-sm font-medium">File</span>
                    <div class="flex gap-2">
                        <input
                            type="text"
                            readonly
                            value={filePath}
                            class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                            placeholder="No file selected"
                        />
                        <Button variant="outline" onclick={handleSelectFile}
                            >Select File</Button
                        >
                    </div>
                </div>

                <Button
                    onclick={handleImport}
                    disabled={!selectedAccount || !filePath || loading}
                >
                    {loading ? "Importing..." : "Import"}
                </Button>
            </div>
        </Card.Content>
    </Card.Root>

    {#if totalRawTransactions > 0}
        <Card.Root class="flex-1 flex flex-col min-h-0">
            <Card.Header>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full">
                    <!-- Column 1: Info -->
                    <div class="flex flex-col justify-start">
                        <Card.Title
                            >Staging Area ({totalRawTransactions})</Card.Title
                        >
                        <Card.Description
                            >Review, edit and assign budget to finalize.</Card.Description
                        >
                    </div>

                    <!-- Column 2: Auto-Tag -->
                    <div class="flex flex-col gap-2 items-start">
                        <Button
                            onclick={handleApplyTags}
                            disabled={loading}
                            variant="outline"
                        >
                            Auto-Tag
                        </Button>
                        <p class="text-sm text-muted-foreground">
                            Use the auto-tag button (or manually edit staged
                            transactions) to assign a budget.
                        </p>
                    </div>

                    <!-- Column 3: Finalize -->
                    <div class="flex flex-col gap-2 items-start">
                        <Button
                            onclick={handleFinalize}
                            disabled={loading}
                            variant="outline"
                        >
                            Finalize
                        </Button>
                        <p class="text-sm text-muted-foreground">
                            Use the Finalize button to move transactions out of
                            the staging area.
                        </p>
                    </div>
                </div>
            </Card.Header>
            <Card.Content class="flex-1 min-h-0 flex flex-col relative p-0">
                {#if statusMessage}
                    <div
                        class="mx-4 my-2 p-3 rounded-md text-sm font-medium border {statusType ===
                        'error'
                            ? 'bg-destructive/15 text-destructive border-destructive/20'
                            : 'bg-green-500/15 text-green-700 border-green-500/20'}"
                    >
                        {statusMessage}
                    </div>
                {/if}
                <!-- 
                    Use DataTableContainer to handle responsive layout.
                    It provides flex-1 min-h-0 relative context.
                  -->
                <DataTableContainer>
                    <DataTable
                        bind:this={dataTableRef}
                        config={tableConfig}
                        {dataSource}
                        onRowEdit={handleRowEdit}
                        class="absolute inset-0 border-0"
                    />
                </DataTableContainer>
            </Card.Content>
        </Card.Root>
    {/if}
</div>
