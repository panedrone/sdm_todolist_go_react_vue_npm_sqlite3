import {defineConfig, transformWithEsbuild} from 'vite'
import react from '@vitejs/plugin-react'

const jsxInJs = {
    name: 'treat-js-as-jsx',
    enforce: 'pre',
    async transform(code, id) {
        if (!/node_modules/.test(id) && /\.js$/.test(id)) {
            return transformWithEsbuild(code, id, {loader: 'jsx'})
        }
    },
}

export default defineConfig({
    plugins: [jsxInJs, react()],
    optimizeDeps: {esbuild: {loader: {'.js': 'jsx'}}},
    publicDir: false,
    build: {
        outDir: 'static/dist',
        manifest: true,
        rollupOptions: {
            input: 'static/App.js',
            output: {
                entryFileNames: 'assets/[name].js',
                chunkFileNames: 'assets/[name].js',
                assetFileNames: 'assets/[name][extname]',
            },
        },
        emptyOutDir: true,
    },
})
