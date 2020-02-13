// 2020, GUILLEUS Hugues <ghugues@netc.fr>
// BSD 3-Clause License

// The web socket connexion
var soc = null;
// reception zone #recep
var recep = null;
// Input message element #ms
var msArea = null;
// button to send a message #send
var sendBut = null;

document.addEventListener("DOMContentLoaded", () => {
	recep = document.querySelector('#recep');
	msArea = document.querySelector('#ms');
	sendBut = document.querySelector('#send');

	soc = new WebSocket(`${
		document.location.protocol == "http:" ? "ws" : "wss"
	}://${document.location.host}/ws/chatt`);

	soc.addEventListener("message", event => {
		let p = document.createElement("p");
		p.textContent = event.data;
		recep.appendChild(p);
	});
	soc.addEventListener("close", event => {
		console.log("fin connection webSocket", event);
		let p = document.createElement("p");
		p.textContent = "*** FIN DE CONNEXION ***";
		recep.appendChild(p);
	});

	sendBut.addEventListener('click', sending);
	msArea.addEventListener('keydown', event => {
		if (event.keyCode !== 13) return;
		sending();
	});
});

// Send a message to the server.
function sending() {
	if (!msArea.value) return;
	soc.send(msArea.value);
	msArea.value = "";
}
