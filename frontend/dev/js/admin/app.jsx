import React from "react";
import { createRoot } from "react-dom";

function App() {
	return (
		<>
			<div>hej 48</div>
			<div>Kage 3</div>
		</>
	);
}

const domContainer = document.getElementById("application");
const root = createRoot(domContainer);
root.render(<App />);
