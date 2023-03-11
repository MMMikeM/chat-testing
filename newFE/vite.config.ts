import { defineConfig } from "vite"
import react from "@vitejs/plugin-react-swc"
import tsconfigpaths from "vite-tsconfig-paths"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), tsconfigpaths()],
  server: {
    host: true,
  },
})
