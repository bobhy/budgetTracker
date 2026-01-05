# Datagrid Component Implementation Plan

This plan details the implementation of a reusable, virtualized Datagrid component for BudgetTracker, using Svelte 5 runes, Tailwind 4, and `@tanstack` libraries.

## Objectives
- Create a highly modular `Datagrid` component.
- Implement efficient "Lazy" virtualization (load-as-you-scroll).
- Support standard features: Sort, Filter, Find, Selection, Resize, Edit.
- Use shadcn-svelte aesthetics.

## Technology Stack
- **Framework**: Svelte 5 (Runes)
- **Styling**: Tailwind CSS v4
- **Table Logic**: `@tanstack/svelte-table` (Headless UI logic)
- **Virtualization**: `@tanstack/svelte-virtual` (Scroll virtualization)
- **Icons**: `lucide-svelte`

## Implementation Steps

### Phase 1: Foundation and Virtualization
**Goal**: A basic Grid component that can scroll through a large list of locally generated dummy data using virtualization.

1. **Dependencies**
    - Install `@tanstack/svelte-table` and `@tanstack/svelte-virtual`.
    - Check/Ensure `lucide-svelte` is available.
2. **Basic Component Structure**
    - Create `frontend/src/lib/components/ui/datagrid/Datagrid.svelte`.
    - Define Props: `config` (DatagridConfig), `data` (Initial data / or simple array for now).
    - Setup `createSvelteTable` with basic columns.
3. **Virtualization Integration**
    - Setup `createVirtualizer` attached to the table container.
    - Implement the rendering loop: Render only virtual rows.
    - Validate scrolling performance with 10k dummy rows.

### Phase 2: DataSource and Lazy Loading
**Goal**: Connect the grid to the `DataSource` callback and implement the "Lazy Cache" strategy.

1. **State Management**
    - Create a state manager (using Runes) to hold:
        - `cachedRows`: The array of accumulated data.
        - `totalKnownRows`: Estimate or exact count.
        - `isLoading`: Loading state.
2. **DataSource Integration**
    - Implement `fetchMore(startIndex, count)` method.
    - Hook `fetchMore` into the Virtualizer's `onRangeUpdate` or similar event (or check scroll position).
    - Update `cachedRows` when data arrives.
3. **Handling "End of Data"**
    - Logic to detect when DataSource is exhausted.

### Phase 3: Columns and Formatting
**Goal**: Respect the detailed Column configuration (Align, Wrap, Resizing, Formatting).

1. **Column Definitions**
    - Map `DatagridConfig.columns` to Tanstack Column definitions.
    - Implement `cell` renderers based on configuration:
        - `formatter`: Call the user-provided formatter.
        - `justify`: Apply Tailwind classes (`text-left`, `text-center`, `text-right`).
        - `wrappable`: Apply line-clamp or whitespace classes.
        - `maxWidth`: Apply CSS styles.
2. **Resizing**
    - Enable `enableColumnResizing` in TanStack Table.
    - Implement the UI handle (splitter) in the header.

### Phase 4: Sorting, Filtering, and Finding
**Goal**: Implement the "Lazy In-Memory" logic.

1. **Sorting**
    - Enable standard TanStack sorting.
    - Map `isSortable` config.
    - Note: Since we are "Lazy Loading", frontend sorting only sorts *what we have*.
    - *Wait*, Requirement update: "Sorting resorts the whole data source".
    - **Refinement**: If sorting resorts the *data source*, we must pass sort keys to the `DataSource` callback and *clear the cache* to reload in the new order.
    - Action: Update `fetchMore` to include sort state. On sort change -> Clear Cache -> Fetch form index 0.
2. **Filtering & Finding (The Hybrid Approach)**
    - Implement a `filter` input.
    - Logic:
        - When filter changes: Clear Cache.
        - Start fetching chunks.
        - For each chunk, apply filter *locally*.
        - If matching rows < `maxVisibleRows` (or screen fill), fetch NEXT chunk immediately.
        - Repeat until screen is full or source exhausted.
        - *Note*: This effectively scans the backend via the `DataSource` API until it finds enough data to show.
3. **Find (Search)**
    - "Find" highlights and scrolls to a match.
    - Logic: Iterate through cache first. If not found, continue fetching/expanding cache until found or end.
    - Scroll virtualizer to the index of the found item.

### Phase 5: Interaction (Selection & Editing)
**Goal**: User interaction features.

1. **Selection**
    - Implement Row Selection state (Set of IDs).
    - Handle Click, Shift+Click, Ctrl+Click logic.
    - Copy to Clipboard (CSV generation).
2. **Deep Editing**
    - Double-click handler on Row.
    - Open simple Dialog/Modal form (using shadcn `Dialog`).
    - Generate form fields based on Column Config.
    - On Submit -> Call `DataEditCallback`.
    - Handle validation errors.

### Phase 6: Styling and Polish
- Ensure it looks "Premium" (Hover effects, smooth scrolling, nice scrollbars).
- Dark mode compatibility.
- Empty states.

## File Structure Plan
```
frontend/src/lib/components/ui/datagrid/
├── Datagrid.svelte       // Main component
├── DatagridTypes.ts      // TS Interfaces for Config/Callbacks
├── DatagridFilter.svelte // Filter/Search Toolbar
└── formatting.ts         // Default formatters/utils
```

## Review feedback
Concerns from reading this doc that I want to review in code before making any revisions...
1. does cachedRows handle scrolling from the end and stopping at the first row?  Spec sounds like it might only handle scrolling to the *end*.
2. I hadn't realized that sorting required changing the data provider callback interface (Adding current sort order).  Given that, maybe it does make sense to add filter and find to the interface as well.
3. It's not clear that this design preserves user's visual position in the way I intend.  Try it and see, then think about a V2.
