const esbuild = require("esbuild");

// main.min.css
esbuild.build({
	entryPoints: ['./frontend/src/css/main.css'],
	sourcemap: false,
	bundle: true,
	outfile: './frontend/out/css/main.min.css', // Corrected file extension
	minify: true,
	target: 'es2020', // This is fine if you're targeting specific CSS features
}).then(() => {
	console.log("main.min.css build complete!");
}).catch(() => {
	process.exit(1);
});

// main.bundle.min.js
esbuild.build({
	entryPoints: ['./frontend/src/js/main.js'],
	sourcemap: false,
	bundle: true,
	outfile: './frontend/out/js/main.bundle.min.js',
	target: 'es2020', /* es2016 */
	minify: true,
	platform: 'browser',
}).then(()=>{
	console.log("main.bundle.min.js build complete!")
}).catch(()=>{
	process.exit(1)
});

// admin.bundle.min.js
esbuild.build({
	entryPoints: ['./frontend/src/js/admin/app.jsx'],
	sourcemap: false,
	bundle: true,
	outfile: './frontend/out/js/admin.bundle.min.js',
	target: 'es2020',
	jsxFactory: 'React.createElement',
	jsxFragment: 'React.Fragment',
	minify: true,
	plugins: [],
	platform: 'browser',
}).then(()=>{
	console.log("admin.bundle.min.js build complete!")
}).catch(()=>{
	process.exit(1)
});