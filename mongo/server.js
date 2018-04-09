const net = require('net');
const Etiqueta = require('./Controllers/EtiquetaController.js');


const server = net.createServer(socket => {
	socket.write('Echo server\r\n');
	socket.pipe(socket);
});
server.listen(666, '127.0.0.1');


const client = new net.Socket();
client.connect(666, '127.0.0.1', () => {
	console.log('Connected');
	client.write('Hello, server! Love, Client.');
});

client.on('etiqueta:guardar', data => {
    console.log(data)    
});

client.on('close', () => {
	console.log('Connection closed');
});