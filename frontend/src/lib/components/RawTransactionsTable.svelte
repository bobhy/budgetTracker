<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import * as Table from "$lib/components/ui/table";
  import { Button } from "$lib/components/ui/button";
  import { ArrowUp, ArrowDown, Pencil, Trash2, Check, X } from "@lucide/svelte";
  import { Input } from "$lib/components/ui/input";
  
  let { 
    data = $bindable([]),
    budgetOptions = [],
    beneficiaryOptions = []
  } = $props();

  const dispatch = createEventDispatcher();

  let sortColumn = $state<string | null>(null);
  let sortDirection = $state<"asc" | "desc" | null>(null);

  let editingId = $state<number | null>(null);
  let editForm = $state<any>({});
  let editAmountFloat = $state<number>(0);

  function cycleSort(column: string) {
    if (sortColumn !== column) {
      sortColumn = column;
      sortDirection = "asc";
    } else if (sortDirection === "asc") {
      sortDirection = "desc";
    } else if (sortDirection === "desc") {
      sortColumn = null;
      sortDirection = null;
    } else {
        sortColumn = column;
        sortDirection = "asc";
    }
    sortData();
  }

  function sortData() {
    if (!sortColumn || !sortDirection) {
        return;
    }
    data.sort((a: any, b: any) => {
        let valA = a[sortColumn!];
        let valB = b[sortColumn!];
        if (valA < valB) return sortDirection === "asc" ? -1 : 1;
        if (valA > valB) return sortDirection === "asc" ? 1 : -1;
        return 0;
    });
  }

  function formatMoney(cents: number) {
     return (cents / 100).toFixed(2);
  }

  function startEdit(row: any) {
    if (editingId !== null && editingId !== row.ID) {
        saveEdit();
    }
    editingId = row.ID;
    editForm = { ...row };
    editAmountFloat = row.Amount / 100;
  }

  function cancelEdit() {
      editingId = null;
      editForm = {};
  }

  function saveEdit() {
      // update amount from float
      editForm.Amount = Math.round(editAmountFloat * 100);
      dispatch('update', editForm);
      // Optimistic update? Parent will reload data usually. 
      // But we can close edit mode immediately.
      editingId = null;
  }

  function deleteRow(id: number) {
      if(confirm("Delete this transaction?")) {
          dispatch('delete', id);
      }
  }

</script>

<div class="rounded-md border">
  <Table.Root>
    <Table.Header>
      <Table.Row>
        <Table.Head class="w-[120px] cursor-pointer select-none" onclick={() => cycleSort("PostedDate")}>
          Date 
          {#if sortColumn === "PostedDate"}
             {#if sortDirection === "asc"}<ArrowUp class="inline h-4 w-4"/>{/if}
             {#if sortDirection === "desc"}<ArrowDown class="inline h-4 w-4"/>{/if}
          {/if}
        </Table.Head>
        <Table.Head class="cursor-pointer select-none" onclick={() => cycleSort("Description")}>
            Description
            {#if sortColumn === "Description"}
             {#if sortDirection === "asc"}<ArrowUp class="inline h-4 w-4"/>{/if}
             {#if sortDirection === "desc"}<ArrowDown class="inline h-4 w-4"/>{/if}
          {/if}
        </Table.Head>
        <Table.Head class="w-[100px] cursor-pointer select-none" onclick={() => cycleSort("Amount")}>
            Amount
             {#if sortColumn === "Amount"}
             {#if sortDirection === "asc"}<ArrowUp class="inline h-4 w-4"/>{/if}
             {#if sortDirection === "desc"}<ArrowDown class="inline h-4 w-4"/>{/if}
          {/if}
        </Table.Head>
        <Table.Head class="w-[150px]">Beneficiary</Table.Head>
        <Table.Head class="w-[150px]">BudgetLine</Table.Head>
         <Table.Head class="w-[150px]">Raw Hint</Table.Head>
        <Table.Head class="w-[100px] text-right">Action</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#each data as row (row.ID)}
        <Table.Row ondblclick={() => editingId !== row.ID && startEdit(row)}>
            {#if editingId === row.ID}
                <!-- Edit Mode -->
                <Table.Cell>
                    <Input type="date" bind:value={editForm.PostedDate} class="h-8" />
                </Table.Cell>
                <Table.Cell>
                    <Input bind:value={editForm.Description} class="h-8" />
                </Table.Cell>
                <Table.Cell>
                    <Input type="number" step="0.01" bind:value={editAmountFloat} class="h-8 text-right" />
                </Table.Cell>
                <Table.Cell>
                     <select class="flex h-8 w-full rounded-md border border-input bg-background px-2 text-xs" bind:value={editForm.Beneficiary}>
                        {#each beneficiaryOptions as b}
                            <option value={b}>{b}</option>
                        {/each}
                        {#if !beneficiaryOptions.includes(editForm.Beneficiary) && editForm.Beneficiary}
                             <option value={editForm.Beneficiary}>{editForm.Beneficiary}</option>
                        {/if}
                     </select>
                </Table.Cell>
                <Table.Cell>
                     <select class="flex h-8 w-full rounded-md border border-input bg-background px-2 text-xs" bind:value={editForm.BudgetLine}>
                        <option value="">(None)</option>
                        {#each budgetOptions as b}
                            <option value={b}>{b}</option>
                        {/each}
                     </select>
                </Table.Cell>
                 <Table.Cell>
                     <Input bind:value={editForm.RawHint} class="h-8" />
                </Table.Cell>
                <Table.Cell class="text-right">
                    <div class="flex justify-end gap-1">
                        <Button variant="ghost" size="icon" class="h-8 w-8 text-green-600" onclick={saveEdit}>
                            <Check class="h-4 w-4" />
                        </Button>
                        <Button variant="ghost" size="icon" class="h-8 w-8 text-muted-foreground" onclick={cancelEdit}>
                            <X class="h-4 w-4" />
                        </Button>
                    </div>
                </Table.Cell>
            {:else}
                <!-- Display Mode -->
                <Table.Cell>{row.PostedDate}</Table.Cell>
                <Table.Cell>{row.Description}</Table.Cell>
                <Table.Cell class="text-right">{formatMoney(row.Amount)}</Table.Cell>
                <Table.Cell>{row.Beneficiary}</Table.Cell>
                <Table.Cell>{row.BudgetLine}</Table.Cell>
                 <Table.Cell class="text-xs text-muted-foreground truncate max-w-[100px]">{row.RawHint}</Table.Cell>
                <Table.Cell class="text-right">
                    <div class="flex justify-end gap-1">
                        <Button variant="ghost" size="icon" class="h-8 w-8" onclick={() => startEdit(row)}>
                            <Pencil class="h-4 w-4" />
                        </Button>
                         <Button variant="ghost" size="icon" class="h-8 w-8 text-destructive" onclick={() => deleteRow(row.ID)}>
                            <Trash2 class="h-4 w-4" />
                        </Button>
                    </div>
                </Table.Cell>
            {/if}
        </Table.Row>
      {/each}
    </Table.Body>
  </Table.Root>
</div>
