import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  optimizeDeps: {
    include: ['react-syntax-highlighter', 'react-syntax-highlighter/dist/cjs/styles/prism']
  },
  define: {
    // Make version and build time available at build time
    __VERSION__: JSON.stringify(process.env.VITE_VERSION || 'dev'),
    __BUILD_TIME__: JSON.stringify(process.env.VITE_BUILD_TIME || new Date().toISOString()),
  },
})
