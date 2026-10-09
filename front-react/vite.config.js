import {defineConfig, transformWithOxc} from 'vite'
import react from '@vitejs/plugin-react'

const jsxInJs = {
    name: 'treat-js-as-jsx',
    enforce: 'pre',
    async transform(code, id) {
        if (!/node_modules/.test(id) && /\.js$/.test(id)) {
            return transformWithOxc(code, id, {lang: 'jsx'})
        }
    },
}

export default defineConfig({
    plugins: [jsxInJs, react()],
    optimizeDeps: {rolldownOptions: {moduleTypes: {'.js': 'jsx'}}},
    base: './',
    publicDir: false,
    build: {
        outDir: 'static/dist',
        manifest: true,
        rolldownOptions: {
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
