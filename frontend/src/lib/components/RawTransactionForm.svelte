<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { createEventDispatcher } from "svelte";

  // Svelte 5 Props using Runes
  let { 
    open = $bindable(false),
    item = null,
    budgetOptions = [],
    beneficiaryOptions = []
  } = $props();

  const dispatch = createEventDispatcher();

  let formData = $state({ ...item });

  // Update formData when item changes
  $effect(() => {
    if (item) {
        formData = { ...item };
    }
  });

  function handleSave() {
    dispatch("save", formData);
    open = false;
  }

  function handleDelete() {
    if (confirm("Are you sure you want to delete this raw transaction?")) {
        dispatch("delete", formData.ID);
        open = false;
    }
  }

</script>

<Dialog.Root bind:open={open}>
  <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
      <Dialog.Title>Edit Raw Transaction</Dialog.Title>
      <Dialog.Description>
        Make changes to the imported transaction before finalizing.
      </Dialog.Description>
    </Dialog.Header>
    
    <div class="grid gap-4 py-4">
      <!-- Date -->
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="postedDate" class="text-right">Date</Label>
        <Input id="postedDate" bind:value={formData.PostedDate} class="col-span-3" />
      </div>

      <!-- Description -->
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="description" class="text-right">Description</Label>
        <Input id="description" bind:value={formData.Description} class="col-span-3" />
      </div>

       <!-- Amount (Cents) -->
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="amount" class="text-right">Amount (Cents)</Label>
        <Input id="amount" type="number" bind:value={formData.Amount} class="col-span-3" />
      </div>

       <!-- Beneficiary -->
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="beneficiary" class="text-right">Beneficiary</Label>
        <div class="col-span-3">
             <select class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" 
                bind:value={formData.Beneficiary}>
                {#each beneficiaryOptions as b}
                    <option value={b}>{b}</option>
                {/each}
             </select>
        </div>
      </div>

       <!-- Budget Line -->
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="budgetLine" class="text-right">Budget Line</Label>
          <div class="col-span-3">
             <select class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" 
                bind:value={formData.BudgetLine}>
                <option value="">(None)</option>
                {#each budgetOptions as b}
                    <option value={b}>{b}</option>
                {/each}
             </select>
        </div>
      </div>

      <!-- Raw Hint -->
       <div class="grid grid-cols-4 items-center gap-4">
        <Label for="rawHint" class="text-right">Hint</Label>
        <Input id="rawHint" bind:value={formData.RawHint} class="col-span-3" />
      </div>

    </div>

    <Dialog.Footer>
      <Button variant="destructive" onclick={handleDelete}>Delete</Button>
      <Button onclick={handleSave}>Save changes</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
