import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { resolve } from "path";

export default defineConfig(({ mode }) => ({
  plugins: [
    tailwindcss(),
    svelte({

    }),
  ],
  resolve: {
    alias: {
      $lib: resolve("./src/lib"),
    },
  },
  build: {
    sourcemap: mode === "development" ? "inline" : false,
  },
  server: {
    port: 5173,
    strictPort: true,
  },
  css: {
    devSourcemap: mode === "development",
  },
}));
