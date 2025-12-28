<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import * as Table from "$lib/components/ui/table";
  import { Button } from "$lib/components/ui/button";
  import { ArrowUp, ArrowDown } from "@lucide/svelte";
  
  let { data = [] } = $props();

  const dispatch = createEventDispatcher();

  let sortColumn = $state<string | null>(null);
  let sortDirection = $state<"asc" | "desc" | null>(null);

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

    // Mutating the prop is generally not recommended unless binding, but for simple sorting handling:
    // We should ideally copy. BUT 'data' is passed by reference.
    // In Runes, to mutate, it should be a $state proxy.
    // If 'data' passed from Import.svelte is $state(rawTransactions), modifying it here updates parent?
    // Arrays in Svelte 5 state are proxies. So yes.
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

</script>

<div class="rounded-md border">
  <Table.Root>
    <Table.Header>
      <Table.Row>
        <Table.Head class="cursor-pointer select-none" onclick={() => cycleSort("PostedDate")}>
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
        <Table.Head class="cursor-pointer select-none" onclick={() => cycleSort("Amount")}>
            Amount
             {#if sortColumn === "Amount"}
             {#if sortDirection === "asc"}<ArrowUp class="inline h-4 w-4"/>{/if}
             {#if sortDirection === "desc"}<ArrowDown class="inline h-4 w-4"/>{/if}
          {/if}
        </Table.Head>
        <Table.Head>Beneficiary</Table.Head>
        <Table.Head>BudgetLine</Table.Head>
        <Table.Head>Action</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#each data as row}
        <!-- Row Click Event -->
        <Table.Row 
            class="cursor-pointer hover:bg-muted/50"
            onclick={() => dispatch('edit', row)}
        >
          <Table.Cell>{row.PostedDate}</Table.Cell>
          <Table.Cell>{row.Description}</Table.Cell>
          <Table.Cell>{formatMoney(row.Amount)}</Table.Cell>
          <Table.Cell>{row.Beneficiary}</Table.Cell>
          <Table.Cell>{row.BudgetLine}</Table.Cell>
          <Table.Cell>
            <span class={row.Action === "add" ? "text-green-600 font-bold" : "text-blue-600 font-bold"}>
                {row.Action}
            </span>
          </Table.Cell>
        </Table.Row>
      {/each}
    </Table.Body>
  </Table.Root>
</div>
