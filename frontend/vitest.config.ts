/// <reference types="vitest" />
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import path from 'path';

export default defineConfig({
    plugins: [
        svelte({ hot: !process.env.VITEST }),
        tailwindcss(),
    ],
    resolve: {
        conditions: ['browser'],
        alias: {
            $lib: path.resolve('./src/lib'),
            $wailsjs: path.resolve('./wailsjs'),
        },
    },
    test: {
        globals: true,
        environment: 'jsdom',
        include: ['src/**/*.{test,spec}.{js,ts}'],
        setupFiles: ['src/setupTests.ts'],
        deps: {
            optimizer: {
                web: {
                    include: ['svelte', '@tanstack/svelte-table']
                }
            }
        }
    },
});
