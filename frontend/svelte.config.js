import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			// default options are fine
			pages: 'build',
			assets: 'build',
			fallback: 'index.html', // SPA mode
			precompress: false
		})
	}
};

export default config;
