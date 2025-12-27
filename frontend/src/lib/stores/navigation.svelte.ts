export type View = 'beneficiaries' | 'accounts' | 'budgets' | 'transactions' | 'database';

class NavigationState {
    currentView = $state<View>('beneficiaries');

    setView(view: View) {
        this.currentView = view;
    }
}

export const navigation = new NavigationState();
