import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig({
    plugins: [react(), tailwindcss()],
    optimizeDeps: {
        include: ['@fluentui/react-components', '@fluentui/react-icons', 'scheduler'],
    },
    server: {
        allowedHosts: ['aegis.zhaoyun.wang'],
        proxy: {
            '/api': {
                target: 'https://127.0.0.1:8060',
                changeOrigin: true,
                secure: false,
            },
        },
    },
    resolve: {
        alias: {
            '@': path.resolve(import.meta.dirname, './src'),
        },
    },
    build: {
        outDir: '../static/dist',
        emptyOutDir: true,
    },
})
