<script lang="ts">
  import { navigation } from '$lib/stores/navigation.svelte.ts';
  import Navbar from '$lib/components/Navbar.svelte';
  import Beneficiaries from '$lib/views/Beneficiaries.svelte';
  import Accounts from '$lib/views/Accounts.svelte';
  import Budgets from '$lib/views/Budgets.svelte';

  import Transactions from '$lib/views/Transactions.svelte';
  import Database from '$lib/views/Database.svelte';
  import Import from '$lib/views/Import.svelte';
</script>

<main class="min-h-screen bg-background text-foreground flex flex-col">
  <Navbar />

  <div class="flex-1 container mx-auto p-4">
    {#if navigation.currentView === 'beneficiaries'}
      <Beneficiaries />
    {:else if navigation.currentView === 'accounts'}
      <Accounts />
    {:else if navigation.currentView === 'budgets'}
      <Budgets />

    {:else if navigation.currentView === 'transactions'}
      <Transactions />
    {:else if navigation.currentView === 'import'}
      <Import />
    {:else if navigation.currentView === 'database'}
      <Database />
    {/if}
  </div>
</main>
