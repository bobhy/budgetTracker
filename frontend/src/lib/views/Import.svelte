<script lang="ts">
  import { onMount } from "svelte";
  import { Button } from "$lib/components/ui/button";
  // Removed unused Select imports
  import * as Card from "$lib/components/ui/card";
  import RawTransactionsTable from "$lib/components/RawTransactionsTable.svelte";
  import { 
    ImportFile, 
    SelectFile, 
    GetAccounts, 
    GetRawTransactions, 
    FinalizeImport, 
    UpdateRawTransaction, 
    DeleteRawTransaction,
    GetBeneficiaries,
    GetBudgets
  } from "$wailsjs/go/main/App"; 
  import { toast } from "svelte-sonner";

  // State using Runes
  let accounts = $state<any[]>([]);
  let selectedAccount = $state("");
  let filePath = $state("");
  let rawTransactions = $state<any[]>([]);
  let loading = $state(false);
  let importStats = $state("");
  
  // Options for Edit Form
  let budgetOptions = $state<string[]>([]);
  let beneficiaryOptions = $state<string[]>([]);

  onMount(async () => {
    await loadAccounts();
    await loadOptions();
    await loadRawTransactions(); 
  });

  async function loadAccounts() {
    try {
        accounts = await GetAccounts() || [];
        if (accounts.length > 0) {
            selectedAccount = accounts[0].Name;
        }
    } catch (err) {
        toast.error("Failed to load accounts: " + err);
    }
  }

  async function loadOptions() {
      try {
          const buds = await GetBudgets() || [];
          budgetOptions = buds.map((b: any) => b.Name);
          const bens = await GetBeneficiaries() || [];
          beneficiaryOptions = bens.map((b: any) => b.Name);
      } catch (err) {
          console.error(err);
      }
  }

  async function loadRawTransactions() {
      try {
          rawTransactions = await GetRawTransactions() || [];
      } catch (err) {
          console.error(err);
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
      try {
          console.log(`[Import] Calling backend ImportFile. Account: ${selectedAccount}, Path: ${filePath}`);
          const msg = await ImportFile(selectedAccount, filePath);
          console.log("[Import] Backend response:", msg);
          toast.success(msg);
          await loadRawTransactions();
      } catch (err) {
          console.error("[Import] Error calling ImportFile:", err);
          toast.error("Import failed: " + err);
      } finally {
          loading = false;
      }
  }

  async function handleFinalize() {
      if (rawTransactions.length === 0) return;
      if (!confirm("Are you sure you want to finalize imports? This will clear the staging table.")) return;
       loading = true;
      try {
          const msg = await FinalizeImport();
          toast.success(msg);
          importStats = msg;
          await loadRawTransactions();
      } catch (err) {
          toast.error("Finalize failed: " + err);
      } finally {
          loading = false;
      }
  }

  async function saveEdit(event: CustomEvent) {
      const item = event.detail;
      try {
          // item.Amount is already cents from the Table's logic
          // But ensure type safety
          const amount = Number(item.Amount);
          await UpdateRawTransaction(
              item.ID, 
              item.PostedDate, 
              amount, 
              item.Description, 
              item.Beneficiary, 
              item.BudgetLine, 
              item.RawHint
          );
          toast.success("Transaction updated");
          await loadRawTransactions();
      } catch (err) {
          toast.error("Update failed: " + err);
      }
  }
  
  async function deleteEdit(event: CustomEvent) {
      const id = event.detail;
      try {
          await DeleteRawTransaction(id);
          toast.success("Transaction deleted");
          await loadRawTransactions();
      } catch (err) {
          toast.error("Delete failed: " + err);
      }
  }

</script>

<div class="h-full flex flex-col gap-4 p-4">
  <Card.Root>
    <Card.Header>
      <Card.Title>Import Transactions</Card.Title>
      <Card.Description>Select an account and CSV file to import.</Card.Description>
    </Card.Header>
    <Card.Content class="space-y-4">
      <div class="flex items-end gap-4">
        <!-- Account Select -->
        <div class="grid gap-2">
            <span class="text-sm font-medium">Account</span>
             <select class="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              bind:value={selectedAccount}>
                {#each accounts as acc}
                    <option value={acc.Name}>{acc.Name}</option>
                {/each}
             </select>
        </div>

        <!-- File Select -->
        <div class="grid gap-2 flex-1">
             <span class="text-sm font-medium">File</span>
             <div class="flex gap-2">
                 <input type="text" readonly value={filePath} class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" placeholder="No file selected" />
                 <Button variant="outline" onclick={handleSelectFile}>Select File</Button>
             </div>
        </div>

        <Button onclick={handleImport} disabled={!selectedAccount || !filePath || loading}>
            {loading ? "Importing..." : "Import"}
        </Button>
      </div>
    </Card.Content>
  </Card.Root>

  {#if rawTransactions.length > 0}
    <Card.Root class="flex-1 flex flex-col min-h-0">
        <Card.Header class="flex flex-row items-center justify-between">
            <div>
                <Card.Title>Staging Area ({rawTransactions.length})</Card.Title>
                <Card.Description>Review and edit before finalizing.</Card.Description>
            </div>
            <Button onclick={handleFinalize} disabled={loading}>Finalize Import</Button>
        </Card.Header>
        <Card.Content class="flex-1 overflow-auto min-h-0">
            <RawTransactionsTable 
                data={rawTransactions} 
                budgetOptions={budgetOptions}
                beneficiaryOptions={beneficiaryOptions}
                on:update={saveEdit}
                on:delete={deleteEdit}
            />
        </Card.Content>
    </Card.Root>
  {/if}

   {#if importStats}
      <div class="text-green-600 font-bold p-2 border rounded bg-green-50">
          Last Import: {importStats}
      </div>
   {/if}
</div>
