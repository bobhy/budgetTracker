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
