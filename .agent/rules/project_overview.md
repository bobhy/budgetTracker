---
trigger: always_on
---

# Project Overview

This document contains notes on the project structure and configuration for the Wails + Svelte 5 + Tailwind 4 template.

## Tech Stack

- **Backend framework**: [Wails v2](https://wails.io/) (Go)
- **Frontend framework**: [Svelte 5](https://svelte.dev/) (Vite)
- **Styling**: [Tailwind CSS v4](https://tailwindcss.com/) (configured via `@tailwindcss/vite`)
- **UI Components**: [shadcn-svelte](https://www.shadcn-svelte.com/)
- **Icons**: [lucide-svelte](https://lucide.dev/guide/packages/lucide-svelte)
- **Forms**: [sveltekit-superforms](https://superforms.rocks/)
- **Build Tool**: Vite

## Project Structure

```
.
├── main.go             # Wails application entry point
├── app.go              # App struct and startup logic (bound to frontend)
├── wails.json          # Wails configuration
├── go.mod              # Go dependencies
├── frontend/           # Frontend source
│   ├── package.json    # Frontend dependencies and scripts
│   ├── vite.config.ts  # Vite config (Tailwind & Svelte plugins)
│   ├── svelte.config.js# Svelte config
│   ├── src/
│   │   ├── App.svelte  # Main application component (SPA entry)
│   │   ├── main.ts     # Entry point
│   │   ├── vite-env.d.ts # Vite types
│   │   ├── style.css   # Global styles (Tailwind 4 import + Theme variables)
│   │   ├── lib/        # Shared code and aliases ($lib)
│   │   │   └── components/ui/ # shadcn-svelte components
```

## Key Configurations

### Backend (`main.go`)
- Embeds `frontend/dist`.
- Binds `App` struct to frontend.

### Frontend (`frontend/vite.config.ts`)
- Uses `@tailwindcss/vite` for Tailwind 4 integration.
- Uses `@sveltejs/vite-plugin-svelte` for Svelte 5.
- Sets up alias `$lib` -> `./src/lib`.

### Frontend (`frontend/package.json`)
- Scripts: `dev`, `build`, `check`.
- Dependencies include standard Svelte ecosystem packages. See `frontend/package.json` for specific versions.

### Styling (`frontend/src/style.css`)
- Imports Tailwind 4: `@import "tailwindcss";`
- Defines CSS variables for key colors (foreground, background, primary, etc.) supporting light/dark modes via `.dark` class.

## Integration Notes
- Components are imported from `$lib/components/ui/...`.
- Tailwind classes are available globally.
