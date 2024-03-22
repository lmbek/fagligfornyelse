const socket = new WebSocket('ws://127.0.0.1:8080/ws');

socket.onmessage = function(event) {
	if (event.data === 'fileChanged') {
		console.log('File changed, refreshing the browser.');
		setTimeout(()=>{
			window.location.reload(true);
		},150)

	} else {
		document.body.innerHTML = event.data;
	}
};


socket.onerror = function(error) {
	console.log('WebSocket Error: ' + error);
};

export default socket;

