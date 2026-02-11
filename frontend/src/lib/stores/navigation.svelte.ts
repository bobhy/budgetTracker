export type View = 'beneficiaries' | 'accounts' | 'budgets' | 'tags' | 'transactions' | 'database' | 'import';

class NavigationState {
    currentView = $state<View>('beneficiaries');

    setView(view: View) {
        this.currentView = view;
    }
}

export const navigation = new NavigationState();
