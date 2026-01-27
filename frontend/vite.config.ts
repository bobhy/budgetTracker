import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { resolve } from "path";
import { getSveltekitMocksAlias, sveltekitOptimizeDepsExclude } from "datatable/vite";

export default defineConfig(({ mode }) => ({
  plugins: [
    tailwindcss(),
    svelte({

    }),
  ],
  resolve: {
    alias: {
      $lib: resolve("./src/lib"),
      $wailsjs: resolve("./wailsjs"),
      ...getSveltekitMocksAlias(),
    },
  },
  build: {
    sourcemap: mode === "development" ? "inline" : false,
    rollupOptions: {
      onwarn(warning, warn) {
        if (warning.message.includes('externalized for browser compatibility')) return;
        warn(warning);
      }
    }
  },
  server: {
    port: 5173,
    strictPort: true,
  },
  css: {
    devSourcemap: mode === "development",
  },
  optimizeDeps: {
    exclude: sveltekitOptimizeDepsExclude,
  },
}));

