import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  //adding the server proxy
  server: {
    proxy: {
      '/v1': 'http://localhost:3000'
    }
  }
})
