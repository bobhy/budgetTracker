// Shim for svelte/internal to support Svelte 4 libraries (like @tanstack/svelte-table v8) in Svelte 5.
// This is a "best effort" shim to prevent build errors. Runtime behavior of components relying on this (like FlexRender from lib) might fail,
// but we are using our own local FlexRender.

export class SvelteComponent {
    $$ = {};
    $destroy() { }
    $on() { }
    $set() { }
}

export function init() { }
export function safe_not_equal() { return false; }
export function text() { return document.createTextNode(''); }
export function claim_text() { }
export function insert_hydration() { }
export function set_data() { }
export function noop() { }
export function detach() { }
export function create_ssr_component() { return { render: () => '' }; }
export function escape() { return ''; }
export function validate_component() { return { $$render: () => '' }; }
export function create_component() { }
export function claim_component() { }
export function mount_component() { }
export function transition_in() { }
export function transition_out() { }
export function destroy_component() { }
export function run_all() { }
export function binding_callbacks() { }
