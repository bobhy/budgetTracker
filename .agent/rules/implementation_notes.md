# Wails & GORM Implementation Notes

## Frontend Module Resolution
When using packages like `lucide-svelte` with Svelte 5 and Vite, always import from the standard package name:
- **Correct**: `import { Icon } from "@lucide/svelte"`
- **Incorrect**: `import { Icon } from "lucide-svelte"`

## Wails & GORM Bindings
Wails has trouble generating TypeScript bindings for complex GORM structs like `gorm.Model` (specifically `gorm.DeletedAt`).

### Best Practice for Models
Avoid embedding `gorm.Model` directly in structs exposed to the frontend. Instead:
1. Explicitly define the fields: `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`.
2. Use `*time.Time` for `DeletedAt` instead of `gorm.DeletedAt` to maintain soft-delete functionality while being JSON-friendly.
3. Use `json:"-"` tags for internal GORM fields (`CreatedAt`, `UpdatedAt`, `DeletedAt`) if they aren't needed in the frontend, preventing "Not found" warnings during binding generation.

**Example Pattern:**
```go
type MyEntity struct {
    ID        uint       `gorm:"primarykey" json:"ID"`
    CreatedAt time.Time  `json:"-"`
    UpdatedAt time.Time  `json:"-"`
    DeletedAt *time.Time `gorm:"index" json:"-"`
    // ... other fields
}
```

## Application Icon
To successfully update the application icon on Linux/Desktop:
1. Place the icon in `build/appicon.png`.
2. Remove old generated icons (`build/windows/icon.ico`, etc.) to force regeneration.
3. Explicitly embed the icon in `main.go` and pass it to the Linux options:
```go
//go:embed build/appicon.png
var icon []byte

// ... inside wails.Run options:
Linux: &linux.Options{
    Icon: icon,
},
```

## Svelte 5 Stores
For global state management (like navigation), use Svelte 5 Runes in `.svelte.ts` files:
```typescript
// navigation.svelte.ts
export const navigationState = $state({ currentView: 'Dashboard' });
```
This is simpler and more reactive than traditional Svelte stores.

## Frontend Development Standards

### Svelte 5 & Runes
The project relies on Svelte 5. All new components and views **MUST** utilize **Runes** (`$state`, `$props`, `$derived`, `$effect`, etc.) instead of legacy Svelte 4 syntax (`export let`, `let` without $state for reactivity).

- **Strict Consistency**: Do not mix legacy syntax and Runes within the same file. While technically supported during migration, it leads to confusion and potential compiler edge-cases (e.g. blank screens).
- **Imports**: Ensure all UI component imports are actually used or removed. Importing complex UI libraries (like shadcn-svelte) without usage can sometimes cause issues if the underlying library has dependencies.

**Example of Component using Runes:**
```html
<script lang="ts">
  let { title = "Default" } = $props();
  let count = $state(0);

  function increment() {
    count += 1;
  }
</script>

<button onclick={increment}>{title}: {count}</button>
```
