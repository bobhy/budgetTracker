
export interface DataGridColumn {
    name: string;
    title?: string;
    isSortable?: boolean;
    wrappable?: 'none' | 'word' | 'hard'; // default 'none'
    maxLines?: number; // default 1
    justify?: 'left' | 'center' | 'right'; // default 'left'
    maxWidth?: number; // characters
    formatter?: (value: any) => string;
}

export interface DataGridConfig {
    name: string; // ID for callbacks
    keyColumn: string;
    title?: string;
    maxVisibleRows?: number;
    isFilterable?: boolean; // Default false
    isFindable?: boolean; // Default false
    columns: DataGridColumn[];
}

// Data Source Callback Types
export interface SortKey {
    key: string;
    direction: 'asc' | 'desc';
}

export type DataSourceCallback = (
    columnKeys: string[],
    startRow: number,
    numRows: number,
    sortKeys: SortKey[]
) => Promise<any[]>; // Returns array of row objects

// Edit Callback Types
export interface DataEditCallbackResult {
    validationFailures: { columnKey: string; message: string }[];
    dbStatus: string | null;
}

export type DataEditCallback = (
    beforeData: any[],
    afterData: any[]
) => Promise<DataEditCallbackResult>;

// Selection Callback Types
export interface SelectedCell {
    rowKey: any;
    columnKey: string;
}

export type SelectedCellsCallback = (selectedCells: SelectedCell[]) => Promise<string | null>;

export interface DatagridProps {
    config: DataGridConfig;
    dataSource: DataSourceCallback;
    onEdit?: DataEditCallback;
    onSelection?: SelectedCellsCallback;
    class?: string;
}
