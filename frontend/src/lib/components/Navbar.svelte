<script lang="ts">
	import { navigation, type View } from '$lib/stores/navigation.svelte.ts';
	import { Database, Users, CreditCard, PiggyBank, Receipt, Upload } from '@lucide/svelte';
	import * as Menubar from '$lib/components/ui/menubar';
    import { Button } from '$lib/components/ui/button';

    const items: { label: string; view: View; icon: any }[] = [
        { label: 'Beneficiaries', view: 'beneficiaries', icon: Users },
        { label: 'Accounts', view: 'accounts', icon: CreditCard },
        { label: 'Budgets', view: 'budgets', icon: PiggyBank },

        { label: 'Transactions', view: 'transactions', icon: Receipt },
        { label: 'Import', view: 'import', icon: Upload },
        { label: 'Database', view: 'database', icon: Database },
    ];
</script>

<div class="flex items-center gap-2 p-4 border-b">
    <h1 class="text-xl font-bold mr-4">BudgetTracker</h1>
    <nav class="flex gap-2">
        {#each items as item}
            <Button
                variant={navigation.currentView === item.view ? "default" : "ghost"}
                onclick={() => navigation.setView(item.view)}
                size="sm"
                class="gap-2"
            >
                <item.icon class="w-4 h-4" />
                {item.label}
            </Button>
        {/each}
    </nav>
</div>
