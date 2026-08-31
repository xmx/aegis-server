import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import {wallpaperListPlugin} from './vite-plugin-wallpapers.ts'

export default defineConfig({
  plugins: [react(), wallpaperListPlugin()],
  server: {
    host: '0.0.0.0',
    port: 8061,
    allowedHosts: ['aegis.zhaoyun.wang'],
    proxy: {
      '/api': {
        target: 'https://127.0.0.1:8060',
        changeOrigin: true,
        secure: false,
        ws: true,
      },
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(import.meta.dirname, './src'),
      '@icons': path.resolve(import.meta.dirname, './fluentui-system-icons/assets'),
      '@extract': path.resolve(import.meta.dirname, './fluentui-system-icons/extract'),
      '@win11': path.resolve(import.meta.dirname, './win11react'),
    },
  },
  build: {
    outDir: '../static/dist',
    emptyOutDir: true,
  },
})