import { defineConfig } from "vite"
import Uni from "@dcloudio/vite-plugin-uni"
import UnoCSS from "unocss/vite"
import AutoImport from "unplugin-auto-import/vite"

export default defineConfig(({ mode }) => ({
  base: process.env.VITE_BASE || '/',
  plugins: [
    Uni(),
    UnoCSS(),
    AutoImport({
      imports: ['vue', 'pinia'],
      dts: 'src/types/auto-imports.d.ts',
    }),
  ],
}))
