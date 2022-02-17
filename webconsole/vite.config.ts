import { fileURLToPath, URL } from "url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  base: '/console/ui/',
  plugins: [vue()],
  resolve: {
    alias: {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
});
