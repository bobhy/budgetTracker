---
description: troubleshoot blank screen in wails app
---
1. Check if the frontend dev server (Vite) is running and accessible.
2. Verify that the frontend build output `dist/index.html` exists and has correct script tags.
3. Inject `console.log` at the very top of `frontend/src/main.ts` to check if JavaScript execution starts.
4. Check Wails terminal output for any frontend console errors (Wails forwards them in dev mode).
5. If JS is running but screen is blank, check for Svelte mounting errors (e.g. `mount` failing).
6. Verify `wails.json` asset configuration matches the Vite output directory.
7. If suspecting a specific component (like Datagrid), temporarily comment it out in `App.svelte` to see if the rest of the app renders.
