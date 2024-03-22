import htmx from "./../main/third-party/htmx.js";
import autorefresher from "./../main/dev/autorefresher";
import React from "react";
import { createRoot } from "react-dom";

function App() {
	return (
		<>
			<div>hej 627</div>
			<div>Kage 3</div>
		</>
	);
}

const domContainer = document.getElementById("application");
const root = createRoot(domContainer);
root.render(<App />);


export default App;
