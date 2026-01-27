import { writable } from 'svelte/store';

// $app/environment
export const browser = true;
export const dev = true;
export const building = false;
export const version = 'mock';

// $app/stores
export const page = writable({ status: 200, error: null, params: {}, route: { id: null }, state: {}, url: new URL('http://localhost') });
export const navigating = writable(null);
export const updated = writable(false);

// $app/navigation
export const goto = () => Promise.resolve();
export const invalidate = () => Promise.resolve();
export const invalidateAll = () => Promise.resolve();
export const afterNavigate = () => { };
export const beforeNavigate = () => { };
export const preloadData = () => Promise.resolve();
export const pushState = () => { };
export const replaceState = () => { };

// $app/forms
export const applyAction = () => { };
export const deserialize = () => { };
export const enhance = () => { };

export { };
