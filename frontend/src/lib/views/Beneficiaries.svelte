<script lang="ts">
    import { onMount } from "svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Table from "$lib/components/ui/table";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Trash2, Pencil, Plus } from "@lucide/svelte";
    import {
        GetBeneficiaries,
        AddBeneficiary,
        UpdateBeneficiary,
        DeleteBeneficiary,
    } from "../../../wailsjs/go/main/App";
    import { models } from "$wailsjs/go/models";

    let beneficiaries: models.Beneficiary[] = $state([]);
    let isDialogOpen = $state(false);
    let isEditing = $state(false);
    let currentName = $state("");
    let originalName = $state(""); // For updates

    async function load() {
        beneficiaries = (await GetBeneficiaries()) || [];
    }

    onMount(load);

    function openAdd() {
        isEditing = false;
        currentName = "";
        isDialogOpen = true;
    }

    function openEdit(name: string) {
        isEditing = true;
        currentName = name;
        originalName = name;
        isDialogOpen = true;
    }

    async function save() {
        try {
            if (isEditing) {
                await UpdateBeneficiary(originalName, currentName);
            } else {
                await AddBeneficiary(currentName);
            }
            isDialogOpen = false;
            load();
        } catch (e) {
            console.error(e);
            alert("Error saving: " + e);
        }
    }

    async function remove(name: string) {
        if (!confirm(`Delete beneficiary ${name}?`)) return;
        try {
            await DeleteBeneficiary(name);
            load();
        } catch (e) {
            console.error(e);
            alert("Error deleting: " + e);
        }
    }
</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">Beneficiaries</h2>
        <div class="flex gap-2">
            <Button onclick={openAdd}>
                <Plus class="w-4 h-4 mr-2" /> Add Beneficiary
            </Button>
        </div>
    </div>

    <div class="border rounded-md">
        <Table.Root>
            <Table.Header>
                <Table.Row>
                    <Table.Head>Name</Table.Head>
                    <Table.Head class="text-right">Actions</Table.Head>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {#if beneficiaries.length === 0}
                    <Table.Row>
                        <Table.Cell
                            colspan={2}
                            class="text-center h-24 text-muted-foreground"
                        >
                            No beneficiaries found.
                        </Table.Cell>
                    </Table.Row>
                {:else}
                    {#each beneficiaries as b}
                        <Table.Row>
                            <Table.Cell>{b.name}</Table.Cell>
                            <Table.Cell class="text-right">
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    onclick={() => openEdit(b.name)}
                                >
                                    <Pencil class="w-4 h-4" />
                                </Button>
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    class="text-destructive"
                                    onclick={() => remove(b.name)}
                                >
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
                <Dialog.Title
                    >{isEditing ? "Edit" : "Add"} Beneficiary</Dialog.Title
                >
            </Dialog.Header>
            <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label class="text-right">Name</Label>
                    <Input class="col-span-3" bind:value={currentName} />
                </div>
            </div>
            <Dialog.Footer>
                <Button variant="outline" onclick={() => (isDialogOpen = false)}
                    >Cancel</Button
                >
                <Button onclick={save}>Save</Button>
            </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Root>
</div>
