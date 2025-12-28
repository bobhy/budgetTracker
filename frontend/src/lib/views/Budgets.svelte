<script lang="ts">
   import { onMount } from 'svelte';
   import { Button } from "$lib/components/ui/button";
   import * as Table from "$lib/components/ui/table";
   import * as Dialog from "$lib/components/ui/dialog";
   import { Input } from "$lib/components/ui/input";
   import { Label } from "$lib/components/ui/label";
   import { Trash2, Pencil, Plus } from "@lucide/svelte";
   import { GetBudgets, AddBudget, UpdateBudget, DeleteBudget, GetBeneficiaries } from "../../../wailsjs/go/main/App";

   let budgets = $state([]);
   let beneficiaries = $state([]);
   
   let isDialogOpen = $state(false);
   let isEditing = $state(false);
   
   let currentName = $state("");
   let currentDesc = $state("");
   let currentBenID = $state("");
   let currentAmount = $state(0);
   let currentInterval = $state(1);
   let originalName = $state("");

   async function load() {
       const [bds, bens] = await Promise.all([
           GetBudgets(),
           GetBeneficiaries()
       ]);
       budgets = bds || [];
       beneficiaries = bens || [];
   }

   onMount(load);

   function openAdd() {
       isEditing = false;
       currentName = "";
       currentDesc = "";
       currentBenID = beneficiaries.length > 0 ? beneficiaries[0].Name : "";
       currentAmount = 0;
       currentInterval = 1;
       isDialogOpen = true;
   }

   function openEdit(b: any) {
       isEditing = true;
       currentName = b.Name;
       currentDesc = b.Description;
       currentBenID = b.BeneficiaryID;
       currentAmount = b.Amount;
       currentInterval = b.IntervalMonths;
       originalName = b.Name;
       isDialogOpen = true;
   }

   async function save() {
       try {
           const amount = Number(currentAmount);
           const interval = Number(currentInterval);
           if (isEditing) {
               await UpdateBudget(originalName, currentName, currentDesc, currentBenID, amount, interval);
           } else {
               await AddBudget(currentName, currentDesc, currentBenID, amount, interval);
           }
           isDialogOpen = false;
           load();
       } catch (e) {
           console.error(e);
           alert("Error saving: " + e);
       }
   }

   async function remove(name: string) {
       if (!confirm(`Delete budget ${name}?`)) return;
       try {
           await DeleteBudget(name);
           load();
       } catch (e) {
           console.error(e);
           alert("Error deleting: " + e);
       }
   }


</script>


<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">Budgets</h2>
        <div class="flex gap-2">

            <Button onclick={openAdd}>
                <Plus class="w-4 h-4 mr-2" /> Add Budget
            </Button>
        </div>
    </div>

    <div class="border rounded-md">
        <Table.Root>
            <Table.Header>
                <Table.Row>
                    <Table.Head>Name</Table.Head>
                    <Table.Head>Description</Table.Head>
                    <Table.Head>Beneficiary</Table.Head>
                    <Table.Head>Amount (cents)</Table.Head>
                    <Table.Head>Interval (Months)</Table.Head>
                    <Table.Head class="text-right">Actions</Table.Head>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {#if budgets.length === 0}
                    <Table.Row>
                        <Table.Cell colspan={6} class="text-center h-24 text-muted-foreground">
                            No budgets found.
                        </Table.Cell>
                    </Table.Row>
                {:else}
                    {#each budgets as b}
                        <Table.Row>
                            <Table.Cell>{b.Name}</Table.Cell>
                            <Table.Cell>{b.Description}</Table.Cell>
                            <Table.Cell>{b.BeneficiaryID}</Table.Cell>
                            <Table.Cell>{b.Amount}</Table.Cell>
                            <Table.Cell>{b.IntervalMonths}</Table.Cell>
                            <Table.Cell class="text-right">
                                <Button variant="ghost" size="icon" onclick={() => openEdit(b)}>
                                    <Pencil class="w-4 h-4" />
                                </Button>
                                <Button variant="ghost" size="icon" class="text-destructive" onclick={() => remove(b.Name)}>
                                    <Trash2 class="w-4 h-4" />
                                </Button>
                            </Table.Cell>
                        </Table.Row>
                    {/each}
                {/if}
            </Table.Body>
        </Table.Root>
    </div>

    <Dialog.Root bind:open={isDialogOpen}>
        <Dialog.Content>
            <Dialog.Header>
                <Dialog.Title>{isEditing ? 'Edit' : 'Add'} Budget</Dialog.Title>
            </Dialog.Header>
            <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Name</Label>
                    <Input class="col-span-3" bind:value={currentName} />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Description</Label>
                    <Input class="col-span-3" bind:value={currentDesc} />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Beneficiary</Label>
                    <select class="col-span-3 flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors" bind:value={currentBenID}>
                        <option value="" disabled>Select Beneficiary</option>
                        {#each beneficiaries as b}
                            <option value={b.Name}>{b.Name}</option>
                        {/each}
                    </select>
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Amount (cents)</Label>
                    <Input type="number" class="col-span-3" bind:value={currentAmount} />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Interval (Mo)</Label>
                    <Input type="number" class="col-span-3" bind:value={currentInterval} />
                </div>
            </div>
            <Dialog.Footer>
                <Button variant="outline" onclick={() => isDialogOpen = false}>Cancel</Button>
                <Button onclick={save}>Save</Button>
            </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Root>
</div>
