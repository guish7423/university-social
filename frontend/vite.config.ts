import { defineConfig } from "vite"
import Uni from "@dcloudio/vite-plugin-uni"
import UnoCSS from "unocss/vite"
import AutoImport from "unplugin-auto-import/vite"

export default defineConfig({
  plugins: [
    Uni(),
    UnoCSS(),
    AutoImport({
      imports: ["vue", "pinia"],
      dts: "src/types/auto-imports.d.ts",
    }),
  ],
})
