const esbuild = require("esbuild");

esbuild.build({
	entryPoints: ['src/app.jsx'],
	bundle: true,
	outfile: 'public/assets/app.js',
	target: 'es2016',
	jsxFactory: 'React.createElement',
	jsxFragment: 'React.Fragment',
	minify: true,
	plugins: [],
}).then(()=>{
	console.log("Build Complete!")
}).catch(()=>{
	process.exit(1)
});