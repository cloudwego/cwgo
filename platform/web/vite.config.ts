import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import viteEslint from "vite-plugin-eslint";

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [react(), viteEslint()],
	server: {
		proxy: {
			"/api": {
				target: "http://localhost:8089",
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, "/api")
			}
		}
	},
	logLevel: "info"
});
