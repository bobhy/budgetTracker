<script lang="ts">
	import {
		createTable,
		getCoreRowModel,
		FlexRender,
		type ColumnDef,
		type TableOptions,
        type Row,
        type SortingState
	} from '@tanstack/svelte-table';
	// import FlexRender from './FlexRender.svelte';
	import { createVirtualizer } from '@tanstack/svelte-virtual';
	import type { DatagridProps, DataGridConfig, DataGridColumn, SortKey } from './DatagridTypes';
	import { onMount, untrack } from 'svelte';
    import { cn } from '$lib/utils';
    import * as Dialog from '$lib/components/ui/dialog';
    import { ArrowUp, ArrowDown, ArrowUpDown } from '@lucide/svelte';
    import { writable, get } from 'svelte/store';

	let { config, dataSource, onEdit, onSelection, class: className }: DatagridProps = $props();

    // -- State --
	let data = $state<any[]>([]);
    let totalKnownRows = $state(0); 
    let hasMore = $state(true); 
    let isLoading = $state(false);
	let sorting = $state<SortingState>([]);
    let globalFilter = $state(""); 
    
    // Selection State
    let selectedRowIndices = $state<Set<number>>(new Set());
    let lastSelectedIndex = $state<number | null>(null);

    // Editing State
    let isEditOpen = $state(false);
    let editingRowIndex = $state<number | null>(null);
    let editingRowData = $state<any>({}); 
    
    // -- Columns --
    const columns = $derived(
        [
            ...config.columns.map((col) => {
                return {
                    accessorKey: col.name,
                    header: col.title || col.name.toLowerCase().replace(/_/g, " ").replace(/\b\w/g, l => l.toUpperCase()),
                    enableSorting: col.isSortable ?? false,
                    enableResizing: true, 
                    size: col.maxWidth ?? 150, 
                    meta: {
                        config: col 
                    },
                    cell: (info: any) => {
                        const val = info.getValue();
                        if (col.formatter) return col.formatter(val);
                        return val;
                    }
                } as ColumnDef<any>;
            })
        ]
    );

    // -- Table Instance --
	const options: TableOptions<any> = {
		get data() {
			return data;
		},
		get columns() {
			return columns;
		},
		getCoreRowModel: getCoreRowModel(),
        manualSorting: true, 
        manualPagination: true,
        enableColumnResizing: true,
        columnResizeMode: 'onChange',
        state: {
            get sorting() {
                return sorting;
            }
        },
        onSortingChange: (updater) => {
            if (typeof updater === 'function') {
                sorting = updater(sorting);
            } else {
                sorting = updater;
            }
        }
	};

	const table = createTable(options);

    // -- Virtualization --
	let tableContainer: HTMLDivElement | undefined = $state();
    
	const rowVirtualizer = createVirtualizer({
        count: 0,
        getScrollElement: () => tableContainer || null,
        estimateSize: () => 40,
        overscan: 5
    });

    $effect(() => {
        // Update virtualizer when dependencies change
        const count = hasMore ? data.length + 1 : data.length;
        const scrollElement = tableContainer || null;
        

        // Use get() to avoid subscription loop
        get(rowVirtualizer).setOptions({
            count,
            getScrollElement: () => scrollElement,
            estimateSize: () => 40,
            overscan: 5
        });
    });



    const virtualRows = $derived($rowVirtualizer.getVirtualItems());
    const totalSize = $derived($rowVirtualizer.getTotalSize());

    // -- Data Fetching --
    let backendFetchedCount = $state(0);
    
    async function performFetch(wantedGridRows: number) {

        if (isLoading || !hasMore) return;
        
        isLoading = true;
        
        try {
            const cols = [config.keyColumn, ...config.columns.map(c => c.name)];
            const batchSize = 100; 
            
            let addedRows = 0;
            let safety = 0;
            
            while (addedRows < wantedGridRows && hasMore && safety < 10) {
                safety++;
                
                // Convert TanStack SortingState to our SortKey[] format
                const sortKeys: SortKey[] = sorting.map(s => ({
                    key: s.id,
                    direction: s.desc ? 'desc' : 'asc'
                }));
                const newRawRows = await dataSource(cols, backendFetchedCount, batchSize, sortKeys);
                
                if (newRawRows.length < batchSize) {
                    hasMore = false;
                }
                backendFetchedCount += newRawRows.length;
                
                // Filter locally
                const filtered = newRawRows.filter(row => {
                   if (!globalFilter) return true;
                   return config.columns.some(col => {
                       const val = row[col.name];
                       if (val == null) return false;
                       return String(val).toLowerCase().includes(globalFilter.toLowerCase());
                   });
                });
                
                if (filtered.length > 0) {
                    data = [...data, ...filtered];
                    addedRows += filtered.length;
                }
                
                if (!hasMore) break;
            }
            
        } catch (e) {
             console.error("Datagrid fetch error", e);
        } finally {
            isLoading = false;
        }
    }
    
    // Effect: Watch scroll and fetch
    $effect(() => {

        if (!virtualRows.length) {
            if (data.length === 0 && hasMore && !isLoading) {

                 performFetch(config.maxVisibleRows || 50); 
            }
            return;
        }
        const lastVirtualRow = virtualRows[virtualRows.length - 1];
        if (lastVirtualRow.index >= data.length - 1 && hasMore && !isLoading) {
            performFetch(50); 
        }
    });

    function resetAndReload() {
        data = [];
        backendFetchedCount = 0;
        hasMore = true;
        isLoading = false;
        selectedRowIndices = new Set();
        lastSelectedIndex = null;
    }
    
    function handleSortOrFilterChange() {
        resetAndReload();
    }

    // -- Helpers --
    function getCellStyles(colConfig: DataGridColumn) {
        let classes = "";
        if (colConfig.justify === 'center') classes += " text-center";
        else if (colConfig.justify === 'right') classes += " text-right";
        else classes += " text-left";

        if (colConfig.wrappable === 'hard') classes += " break-all whitespace-normal";
        else if (colConfig.wrappable === 'word') classes += " break-words whitespace-normal";
        else classes += " truncate whitespace-nowrap";

        return classes;
    }

    // -- Selection Logic --
    function handleRowClick(index: number, e: MouseEvent) {
        if (index >= data.length) return;

        const newSet = new Set(selectedRowIndices);
        
        if (e.shiftKey && lastSelectedIndex !== null) {
            // Range
            const start = Math.min(lastSelectedIndex, index);
            const end = Math.max(lastSelectedIndex, index);
            if (!e.ctrlKey && !e.metaKey) {
                newSet.clear();
            }
            for (let i = start; i <= end; i++) {
                newSet.add(i);
            }
        } else if (e.ctrlKey || e.metaKey) {
            // Toggle
            if (newSet.has(index)) {
                newSet.delete(index);
            } else {
                newSet.add(index);
            }
            lastSelectedIndex = index;
        } else {
            // Single
            newSet.clear();
            newSet.add(index);
            lastSelectedIndex = index;
        }
        
        selectedRowIndices = newSet;
        
        if (onSelection) {
             // Pass simplistic notification
        }
    }

    // -- Copy / Clipboard --
    async function handleKeyDown(e: KeyboardEvent) {
        if ((e.ctrlKey || e.metaKey) && e.key === 'c') {
            if (selectedRowIndices.size === 0) return;
            
            e.preventDefault();
            
            const rows = Array.from(selectedRowIndices).sort((a,b) => a - b).map(i => data[i]);
            if (rows.length === 0) return;
            
            const header = config.columns.map(c => c.title || c.name).join(",");
            const lines = rows.map(r => {
                return config.columns.map(c => {
                    const val = r[c.name];
                    const s = String(val ?? "");
                    if (s.includes(",") || s.includes('"') || s.includes("\n")) {
                        return `"${s.replace(/"/g, '""')}"`;
                    }
                    return s;
                }).join(",");
            });
            
            const csv = [header, ...lines].join("\n");
            
            try {
                await navigator.clipboard.writeText(csv);
            } catch (err) {
                console.error("Failed to copy", err);
            }
            
            if (onSelection) {
                const cells = [];
                for(const r of rows) {
                    for(const c of config.columns) {
                        cells.push({ rowKey: r[config.keyColumn], columnKey: c.name });
                    }
                }
                const msg = await onSelection(cells);
                if (msg) alert(msg); 
            }
        }
    }

    // -- Editing --
    function handleRowDoubleClick(index: number) {
        if (index >= data.length) return;
        editingRowIndex = index;
        editingRowData = { ...data[index] }; 
        isEditOpen = true;
    }

    async function handleEditSave() {
        if (editingRowIndex === null) return;
        
        if (onEdit) {
            const original = data[editingRowIndex];
            const result = await onEdit([original], [editingRowData]);
            
            if (result.dbStatus) {
                alert("Error: " + result.dbStatus);
                return;
            }
            if (result.validationFailures.length > 0) {
                alert("Validation Failed: " + result.validationFailures.map(f => f.message).join("\n"));
                return;
            }
            
            const newData = [...data];
            newData[editingRowIndex] = editingRowData;
            data = newData;
            
            isEditOpen = false;
            editingRowIndex = null;
        }
    }
</script>

<div 
    class={cn("flex flex-col h-full w-full border rounded-md overflow-hidden bg-background outline-none", className)}
    tabindex="0"
    onkeydown={handleKeyDown}
    role="grid"
>
    <!-- Toolbar (Filter) -->
    {#if config.isFilterable || config.isFindable}
    <div class="flex-none p-2 border-b flex gap-2">
        {#if config.isFilterable}
             <input 
                type="text" 
                placeholder="Filter..." 
                class="px-2 py-1 text-sm border rounded bg-transparent focus:outline-none focus:ring-1 focus:ring-primary"
                bind:value={globalFilter}
                oninput={handleSortOrFilterChange}
             />
        {/if}
    </div>
    {/if}

    <!-- Header -->
    <div class="flex-none border-b bg-muted/40 font-medium text-sm">
        <div class="flex w-full min-w-max">
             {#each table.getHeaderGroups() as headerGroup}
                {#each headerGroup.headers as header}
                    <div class="relative px-4 py-3 text-left font-medium text-muted-foreground select-none flex items-center gap-1 group border-r border-transparent hover:border-border/50 transition-colors"
                         style="width: {header.getSize()}px; flex: {header.column.getIsResizing() ? '0 0 auto' : '1 1 0%'};">
                        {#if !header.isPlaceholder}
                            <button
                                class="flex items-center gap-1 hover:text-foreground transition-colors w-full overflow-hidden"
                                class:cursor-pointer={header.column.getCanSort()}
                                class:select-none={true}
                                onclick={(e) => {
                                    header.column.getToggleSortingHandler()?.(e); 
                                    handleSortOrFilterChange(); 
                                }}
                            >
                                <div class="truncate font-semibold text-foreground/80">
                                    <FlexRender
                                        content={header.column.columnDef.header}
                                        context={header.getContext()}
                                    />
                                </div>
                                
                                {#if header.column.getIsSorted() === 'asc'}
                                    <ArrowUp class="h-3 w-3" />
                                {:else if header.column.getIsSorted() === 'desc'}
                                    <ArrowDown class="h-3 w-3" />
                                {:else if header.column.getCanSort()}
                                    <ArrowUpDown class="h-3 w-3 opacity-0 group-hover:opacity-50" />
                                {/if}
                            </button>
                            
                            <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
                            <div
                                role="separator"
                                class="absolute right-0 top-0 h-full w-1 cursor-col-resize touch-none select-none hover:bg-primary opacity-0 group-hover:opacity-100"
                                onmousedown={header.getResizeHandler()}
                                ontouchstart={header.getResizeHandler()}
                            ></div>
                        {/if}
                    </div>
                {/each}
             {/each}
        </div>
    </div>

    <!-- Body (Virtual Scroll Container) -->
    <div 
        bind:this={tableContainer} 
        class="flex-1 overflow-auto relative w-full"
    >
        <div style="height: {totalSize}px; width: 100%; position: relative;">
            
            {#each virtualRows as virtualRow (virtualRow.index)}
                {@const row = table.getRowModel().rows[virtualRow.index]}
                {@const isSelected = selectedRowIndices.has(virtualRow.index)}
                
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <div
                    class={cn(
                        "absolute top-0 left-0 w-full flex min-w-max border-b transition-colors cursor-default",
                        isSelected ? "bg-primary/20 hover:bg-primary/30" : "bg-background hover:bg-muted/50"
                    )}
                    style="transform: translateY({virtualRow.start}px); height: {virtualRow.size}px;"
                    onclick={(e) => handleRowClick(virtualRow.index, e)}
                    ondblclick={() => handleRowDoubleClick(virtualRow.index)}
                    role="row"
                    aria-selected={isSelected}
                    tabindex={row ? row.original[config.keyColumn] : 0}
                >
                     {#if row}
                        {#each row.getVisibleCells() as cell}
                            {@const colConfig = (cell.column.columnDef.meta as any)?.config as DataGridColumn}
                            <div 
                                class={cn("px-4 py-2 text-sm flex items-center border-r border-transparent last:border-0", colConfig ? getCellStyles(colConfig) : "")}
                                style="width: {cell.column.getSize()}px; flex: 1 1 0%;"
                            >
                                <FlexRender
                                    content={cell.column.columnDef.cell}
                                    context={cell.getContext()}
                                />
                            </div>
                        {/each}
                     {:else}
                         {#if isLoading}
                            <div class="w-full h-full flex items-center px-4 text-muted-foreground/50 text-xs text-center justify-center italic">
                                Loading more data...
                            </div>
                         {/if}
                     {/if}
                </div>
            {/each}

        </div>
    </div>
    
    <!-- Footer / Status -->
    <div class="flex-none border-t bg-muted/40 p-2 text-xs text-muted-foreground flex justify-between select-none">
        <div>
            {data.length} rows loaded. {selectedRowIndices.size > 0 ? `${selectedRowIndices.size} selected` : ''}
        </div>
        <div>
            {isLoading ? 'Fetching...' : 'Idle'}
        </div>
    </div>
    
    <!-- Edit Dialog -->
    <Dialog.Root bind:open={isEditOpen}>
        <Dialog.Content class="sm:max-w-[425px]">
          <Dialog.Header>
            <Dialog.Title>Edit Row</Dialog.Title>
            <Dialog.Description>
              Make changes to the row data here. Click save when you're done.
            </Dialog.Description>
          </Dialog.Header>
          <div class="grid gap-4 py-4 max-h-[60vh] overflow-y-auto">
             {#each config.columns as col}
                 <div class="grid grid-cols-4 items-center gap-4">
                    <label for="edit-{col.name}" class="text-right text-sm font-medium">
                      {col.title || col.name}
                    </label>
                    <input
                      id="edit-{col.name}"
                      bind:value={editingRowData[col.name]}
                      class="col-span-3 flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                    />
                 </div>
             {/each}
          </div>
          <Dialog.Footer>
            <button class="px-4 py-2 bg-primary text-primary-foreground rounded-md text-sm font-medium hover:bg-primary/90" onclick={handleEditSave}>Save changes</button>
          </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Root>

</div>
