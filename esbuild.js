const esbuild = require("esbuild");

esbuild.build({
	entryPoints: ['./frontend/dev/js/admin/app.jsx','./frontend/dev/js/web/htmx.js', './frontend/dev/js/dev/autorefresher.js'],
	sourcemap: false,
	bundle: true,
	outdir: './frontend/live/public/js/',
	target: 'es2016',
	jsxFactory: 'React.createElement',
	jsxFragment: 'React.Fragment',
	minify: true,
	plugins: [],
	platform: 'browser',
}).then(()=>{
	console.log("Build Complete!")
}).catch(()=>{
	process.exit(1)
});