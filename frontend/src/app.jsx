import React from "react"
import ReactDOM from "react-dom/client"
function App(){
	return (
		<>
			<div>hej 28</div>
		</>
	)
}

const domContainer = document.querySelector("#application");
const root = ReactDOM.createRoot(domContainer)
root.render(<App />)