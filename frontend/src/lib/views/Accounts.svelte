<script lang="ts">
   import { onMount } from 'svelte';
   import { Button } from "$lib/components/ui/button";
   import * as Table from "$lib/components/ui/table";
   import * as Dialog from "$lib/components/ui/dialog";
   import { Input } from "$lib/components/ui/input";
   import { Label } from "$lib/components/ui/label";
   import { Trash2, Pencil, Plus } from "@lucide/svelte";
   import { GetAccounts, AddAccount, UpdateAccount, DeleteAccount, GetBeneficiaries } from "../../../wailsjs/go/main/App";
   import { models } from '$wailsjs/go/models';

   let accounts: models.Account[] = $state([]);
   let beneficiaries: models.Beneficiary[] = $state([]);
   
   let isDialogOpen = $state(false);
   let isEditing = $state(false);
   
   // Form state
   let currentName = $state("");
   let currentDesc = $state("");
   let currentBenID = $state("");
   let originalName = $state("");

   async function load() {
       const [accs, bens] = await Promise.all([
           GetAccounts(),
           GetBeneficiaries()
       ]);
       accounts = accs || [];
       beneficiaries = bens || [];
   }

   onMount(load);

   function openAdd() {
       isEditing = false;
       currentName = "";
       currentDesc = "";
       currentBenID = beneficiaries.length > 0 ? beneficiaries[0].Name : "";
       isDialogOpen = true;
   }

   function openEdit(acc: any) {
       isEditing = true;
       currentName = acc.Name;
       currentDesc = acc.Description;
       currentBenID = acc.BeneficiaryID;
       originalName = acc.Name;
       isDialogOpen = true;
   }

   async function save() {
       try {
           if (isEditing) {
               await UpdateAccount(originalName, currentName, currentDesc, currentBenID);
           } else {
               await AddAccount(currentName, currentDesc, currentBenID);
           }
           isDialogOpen = false;
           load();
       } catch (e) {
           console.error(e);
           alert("Error saving: " + e);
       }
   }

   async function remove(name: string) {
       if (!confirm(`Delete account ${name}?`)) return;
       try {
           await DeleteAccount(name);
           load();
       } catch (e) {
           console.error(e);
           alert("Error deleting: " + e);
       }
   }


</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">Accounts</h2>
        <div class="flex gap-2">

            <Button onclick={openAdd}>
                <Plus class="w-4 h-4 mr-2" /> Add Account
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
                    <Table.Head class="text-right">Actions</Table.Head>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {#if accounts.length === 0}
                    <Table.Row>
                        <Table.Cell colspan={4} class="text-center h-24 text-muted-foreground">
                            No accounts found.
                        </Table.Cell>
                    </Table.Row>
                {:else}
                    {#each accounts as a}
                        <Table.Row>
                            <Table.Cell>{a.Name}</Table.Cell>
                            <Table.Cell>{a.Description}</Table.Cell>
                            <Table.Cell>{a.BeneficiaryID}</Table.Cell>
                            <Table.Cell class="text-right">
                                <Button variant="ghost" size="icon" onclick={() => openEdit(a)}>
                                    <Pencil class="w-4 h-4" />
                                </Button>
                                <Button variant="ghost" size="icon" class="text-destructive" onclick={() => remove(a.Name)}>
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
                <Dialog.Title>{isEditing ? 'Edit' : 'Add'} Account</Dialog.Title>
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
                    <!-- Simple Select for v0.1 -->
                    <select class="col-span-3 flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors" bind:value={currentBenID}>
                        <option value="" disabled>Select Beneficiary</option>
                        {#each beneficiaries as b}
                            <option value={b.Name}>{b.Name}</option>
                        {/each}
                    </select>
                </div>
            </div>
            <Dialog.Footer>
                <Button variant="outline" onclick={() => isDialogOpen = false}>Cancel</Button>
                <Button onclick={save}>Save</Button>
            </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Root>
</div>
