<script lang="ts">
  import { onMount } from "svelte";
  import { navigation } from "$lib/stores/navigation.svelte.ts";
  import Navbar from "$lib/components/Navbar.svelte";
  import Beneficiaries from "$lib/views/Beneficiaries.svelte";
  import Accounts from "$lib/views/Accounts.svelte";
  import Budgets from "$lib/views/Budgets.svelte";

  import Transactions from "$lib/views/Transactions.svelte";
  import Database from "$lib/views/Database.svelte";
  import Import from "$lib/views/Import.svelte";

  let error = $state<string | null>(null);
  let errorDetails = $state<string>("");

  onMount(() => {
    window.addEventListener("error", (e) => {
      console.error("Global error caught:", e);
      error = e.message;
      errorDetails = e.error?.stack || "";
    });

    window.addEventListener("unhandledrejection", (e) => {
      console.error("Unhandled promise rejection:", e);
      error = e.reason?.message || String(e.reason);
      errorDetails = e.reason?.stack || "";
    });
  });
</script>

{#if error}
  <div class="fixed inset-0 bg-red-100 flex items-center justify-center p-4">
    <div class="bg-white p-8 rounded shadow-lg max-w-2xl w-full">
      <h1 class="text-2xl font-bold text-red-600 mb-4">Application Error</h1>
      <p class="text-gray-700 mb-4 font-mono text-sm">{error}</p>
      {#if errorDetails}
        <details class="mb-4">
          <summary
            class="cursor-pointer text-sm text-gray-600 hover:text-gray-800"
            >Stack trace</summary
          >
          <pre
            class="mt-2 text-xs bg-gray-100 p-2 rounded overflow-auto max-h-64">{errorDetails}</pre>
        </details>
      {/if}
      <p class="text-sm text-gray-500 mb-4">
        Check browser console for more details
      </p>
      <button
        onclick={() => {
          error = null;
          errorDetails = "";
        }}
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        Dismiss
      </button>
    </div>
  </div>
{:else}
  <main class="min-h-screen bg-background text-foreground flex flex-col">
    <Navbar />

    <div class="flex-1 container mx-auto p-4">
      {#if navigation.currentView === "beneficiaries"}
        <Beneficiaries />
      {:else if navigation.currentView === "accounts"}
        <Accounts />
      {:else if navigation.currentView === "budgets"}
        <Budgets />
      {:else if navigation.currentView === "transactions"}
        <Transactions />
      {:else if navigation.currentView === "import"}
        <Import />
      {:else if navigation.currentView === "database"}
        <Database />
      {/if}
    </div>
  </main>
{/if}
