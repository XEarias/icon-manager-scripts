const net = require("net")
const handlers = require("./handlers.js")
const config = require("./config.js")
const async = require("async")


const handler = (data) => {

	return new Promise((resolve, reject) => {

		let id = JSON.parse(data.toString()).id
		let action = JSON.parse(data.toString()).action
		let body = JSON.parse(data.toString()).data

		handlers[action](body, res => {

			res.id = id
			res.action = action

			if(res.status == 200){
				resolve( JSON.stringify(res) )
			}else{
				reject( JSON.stringify(res) )
			}

			
		})
	})
}

var server = net.createServer(socket => {

	console.log("----Nueva conexion establecida----")

	socket.on("data", data => {
	
		handler(data)
			.then((res) => {
				socket.write(res)
				socket.end()
			})
			.catch((res) => {
				socket.write(res) 
				socket.end()
			})
			


		

	})

	socket.on("close", (error) => {
		if(error) console.log("error:" + error)
	})

	socket.on("error", (e) => {
		if(e) console.log(e.toString())
	})

})

server.listen(config.port, config.host, 1, () => {
	console.log("Servidor corriendo en el HOST:"+config.host+", PORT:"+config.port)
})
