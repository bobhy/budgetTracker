<script lang="ts">
   import { onMount } from 'svelte';
   import { Button } from "$lib/components/ui/button";
   import * as Table from "$lib/components/ui/table";
   import * as Dialog from "$lib/components/ui/dialog";
   import { Input } from "$lib/components/ui/input";
   import { Label } from "$lib/components/ui/label";
   import { Trash2, Pencil, Plus } from "@lucide/svelte";
   import { GetTransactions, AddTransaction, UpdateTransaction, DeleteTransaction, GetAccounts, GenerateTransactions } from "../../../wailsjs/go/main/App";

   let transactions = $state([]);
   let accounts = $state([]);
   
   let isDialogOpen = $state(false);
   let isEditing = $state(false);
   let currentID = $state(0); // transactions have uint IDs
   
   let currentDate = $state(new Date().toISOString().split('T')[0]);
   let currentAccountID = $state("");
   let currentAmount = $state(0);
   let currentDesc = $state("");
   let currentTag = $state("");

   async function load() {
       const [txs, accs] = await Promise.all([
           GetTransactions(),
           GetAccounts()
       ]);
       transactions = txs || [];
       accounts = accs || [];
   }

   onMount(load);

   function openAdd() {
       isEditing = false;
       currentDate = new Date().toISOString().split('T')[0];
       currentAccountID = accounts.length > 0 ? accounts[0].Name : "";
       currentAmount = 0;
       currentDesc = "";
       currentTag = "";
       isDialogOpen = true;
   }

   function openEdit(t: any) {
       isEditing = true;
       currentID = t.ID;
       // Assuming PostedDate comes as string or we might need formatting
       currentDate = t.PostedDate; 
       currentAccountID = t.AccountID;
       currentAmount = t.Amount;
       currentDesc = t.Description;
       currentTag = t.Tag;
       isDialogOpen = true;
   }

   async function save() {
       try {
           const amount = Number(currentAmount);
           if (isEditing) {
               await UpdateTransaction(currentID, currentDate, currentAccountID, amount, currentDesc, currentTag);
           } else {
               await AddTransaction(currentDate, currentAccountID, amount, currentDesc, currentTag);
           }
           isDialogOpen = false;
           load();
       } catch (e) {
           console.error(e);
           alert("Error saving: " + e);
       }
   }

   async function remove(id: number) {
       if (!confirm(`Delete transaction?`)) return;
       try {
           await DeleteTransaction(id);
           load();
       } catch (e) {
           console.error(e);
           alert("Error deleting: " + e);
       }
   }

   async function generate() {
       try {
           await GenerateTransactions(10);
           load();
       } catch (e) {
           console.error(e);
           alert("Error generating: " + e);
       }
   }
</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">Transactions</h2>
        <div class="flex gap-2">
            <Button variant="outline" onclick={generate}>
                Generate 10
            </Button>
            <Button onclick={openAdd}>
                <Plus class="w-4 h-4 mr-2" /> Add Transaction
            </Button>
        </div>
    </div>

    <div class="border rounded-md">
        <Table.Root>
            <Table.Header>
                <Table.Row>
                    <Table.Head>Date</Table.Head>
                    <Table.Head>Account</Table.Head>
                    <Table.Head>Amount</Table.Head>
                    <Table.Head>Description</Table.Head>
                    <Table.Head>Tag</Table.Head>
                    <Table.Head class="text-right">Actions</Table.Head>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {#if transactions.length === 0}
                    <Table.Row>
                        <Table.Cell colspan={6} class="text-center h-24 text-muted-foreground">
                            No transactions found.
                        </Table.Cell>
                    </Table.Row>
                {:else}
                    {#each transactions as t}
                        <Table.Row>
                            <Table.Cell>{t.PostedDate}</Table.Cell>
                            <Table.Cell>{t.AccountID}</Table.Cell>
                            <Table.Cell>{t.Amount}</Table.Cell>
                            <Table.Cell>{t.Description}</Table.Cell>
                            <Table.Cell>{t.Tag}</Table.Cell>
                            <Table.Cell class="text-right">
                                <Button variant="ghost" size="icon" onclick={() => openEdit(t)}>
                                    <Pencil class="w-4 h-4" />
                                </Button>
                                <Button variant="ghost" size="icon" class="text-destructive" onclick={() => remove(t.ID)}>
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
                <Dialog.Title>{isEditing ? 'Edit' : 'Add'} Transaction</Dialog.Title>
            </Dialog.Header>
            <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Date</Label>
                    <Input type="date" class="col-span-3" bind:value={currentDate} />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Account</Label>
                    <select class="col-span-3 flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors" bind:value={currentAccountID}>
                        <option value="" disabled>Select Account</option>
                        {#each accounts as a}
                            <option value={a.Name}>{a.Name}</option>
                        {/each}
                    </select>
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Amount (cents)</Label>
                    <Input type="number" class="col-span-3" bind:value={currentAmount} />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Description</Label>
                    <Input class="col-span-3" bind:value={currentDesc} />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Tag</Label>
                    <Input class="col-span-3" bind:value={currentTag} />
                </div>
            </div>
            <Dialog.Footer>
                <Button variant="outline" onclick={() => isDialogOpen = false}>Cancel</Button>
                <Button onclick={save}>Save</Button>
            </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Root>
</div>
