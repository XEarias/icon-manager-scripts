const net = require('net');
const handlers = require('./handlers.js');
const config = require('./config.js');
const async = require('async');


const handler = (data) => {

	return new Promise((resolve, reject) => {

		let id = JSON.parse(data.toString()).id;
		let action = JSON.parse(data.toString()).action;
		let body = JSON.parse(data.toString()).data;

		handlers[action](body, res => {

			res.id = id;
			res.action = action;

			if(res.status == 200){
				resolve( JSON.stringify(res) );
			}else{
				reject( JSON.stringify(res) )
			}

			
		})
	})
}

var server = net.createServer(socket => {

	console.log('Servidor creado');

	socket.on('data', async data => {

		try {
			res = await handler(data);
			socket.write(res); 
		} catch (e) {
			socket.write(e.toString()); 
		}

	})

	socket.on('close', (had_error) => {
		console.log(had_error);
	});

	socket.on('error', (e) => {
		console.log(e.toString())
	})

});

server.listen(config.port, config.host, 1, () => {
	console.log("Servidor corriendo en el HOST:"+config.host+", PORT:"+config.port)
});
