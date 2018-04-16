const NounProject = require("the-noun-project")
const jsonfile = require("jsonfile")
//const csv = require("csvtojson")
const flags = require("node-flags")
const { exec } = require("child_process")

let saltoTags = flags.get("TAGS")
let cantidad = flags.get("CANT")
let saltoArg = flags.get("SALTO")

cantidad = cantidad ? parseInt(cantidad) : 0
saltoArg = saltoArg ? parseInt(saltoArg) : 0
saltoTags = saltoTags ? parseInt(saltoTags) : 0

let np = new NounProject({
	key: "049d89ea99f6415c837e7f0de9040b96",
	secret: "58dd191936204d42885a5ab4de7a6663"
})

let promises = []
let tags = {}

var mysql      = require("mysql")
var connection = mysql.createConnection({
	host     : "127.0.0.1",
	user     : "logoPro",
	password : "&rJ-fZ:1uZ24",
	database : "disenadorlogodb_tags_prov"
})

connection.connect(function(err) {
	if (err) {
		console.error("Error conectando: " + err.stack)
		return
	}

	let queryInterval = (limpiarIntervalo) => {

		connection.query("SELECT * FROM tags LIMIT 100 OFFSET ?", [offsetQuery], function (error, tagRows) {
		
			if (error) {
				console.log("error seleccionando tags "+error)
				if(limpiarIntervalo){
					clearInterval(invervalo)
				}
				return
			}

			if(!tagRows || !tagRows.length){
				if(limpiarIntervalo){
					clearInterval(invervalo)
				}
				return
			}
			
			var salto = saltoArg

			while(salto <= (cantidad + saltoArg)) {
				tagRows.map(function(tag){
					
					var promiseGet = new Promise((resolve) => {
						np.getIconsByTerm(tag.ENG, {offset: salto, limit_to_public_domain: 0}, (err, data) => {
							
							if (!err && data) {
	
								resolve({data: data.icons, tag: tag.ENG})
	
								tags[tag.ENG] = tag
	
							} else {
								resolve(666)
	
							}
							
						})
	
	
					})

					
					salto = salto + 50
					promises.push(promiseGet)
				
				})
			}
	
			let icons = {}
			console.log(promises)
			Promise.all(promises)
	
				.then((res) => {
	
					res.map((iconCollection) => {
	
						if (iconCollection && iconCollection !== 666) {
	
							iconCollection.data.forEach(icon => {
	
								if (icon.icon_url) {
	
									if (!icons[icon.id]) { //se registra por primera vez
										icons[icon.id] = {
											url: icon.icon_url,
											tags: [iconCollection.tag],
											category: tags[iconCollection.tag].categoria
										}
									} else if (icons[icon.id] && icons[icon.id].tags.indexOf(iconCollection.tag) == -1) {

										icons[icon.id].tags.push(iconCollection.tag)
									
									}
	
								}
	
							})
						}
	
					})
					
					console.log(`Cantidad de logos encontrados: ${Object.keys(icons).length}`)
	
					let promiseLogos = new Promise((resolve) => {
						//iconos => json file
						jsonfile.writeFile("./icons.json", icons, (err) => {
							if (!err) {
								console.log("Archivo de Logos listo")
								resolve()
							}
							if(err) {
								console.log(err)
							}
						})
					})
	
					let promiseTags = new Promise((resolve) => {
						//tags => json file
						jsonfile.writeFile("./tags.json", tags, (err) => {
							if (!err) {
								console.log("Archivo de tags listo")
								resolve()
							}
							if(err) console.log(err)
						})
					})
	
					Promise.all([promiseLogos, promiseTags])
						.then(() => {
	
							console.log("Todo preparado")
	
							exec("./main", {cwd: "../upload"}, (error, stdout) => {
								if (error) {
									console.error(` ${error}`)
									if(limpiarIntervalo){
										clearInterval(invervalo)
									}
									return
								}
								console.log(`${stdout}`)
							})
						})
						.catch(() => {
							console.log("Error previo a ejecucion de script")
						})
	
	
				})
				.catch((err) => {
					if(err) console.log(err)
				})

		})

		offsetQuery += 100

	}


	let offsetQuery = saltoTags
	
	queryInterval()

	let invervalo = setInterval(() => {
	
		queryInterval(true)
	}, 3600000)
	
})

/*
csv()
	.fromFile("./tags.csv")

	.on("json", (tag) => {
       
		var salto = saltoArg

		while(salto <= (cantidad + saltoArg)) {
			let promiseGet = new Promise((resolve) => {
                
				//console.log("inicio "+tag.ENG+"#"+salto)
				np.getIconsByTerm(tag.ENG, {offset: salto, limit_to_public_domain: 0}, (err, data) => {

					if (!err && data) {

						resolve({data: data.icons, tag: tag.ENG})

						tags[tag.ENG] = tag

					} else {
                       
						resolve(666)

					}
					/*
					var d = new Date()

					console.log(`Respuesta: ${d.getHours()}:${d.getMinutes()}:${d.getSeconds()}:${d.getMilliseconds()}...`)

				})


			})
			salto = salto + 50
			promises.push(promiseGet)

		}

	})
	.on("done", (err) => {
		if(err) return console.log(err)
        
		let icons = {}
        
		Promise.all(promises)

			.then((res) => {

				res.map((iconCollection) => {

					if (iconCollection && iconCollection !== 666) {

						iconCollection.data.forEach(icon => {

							if (icon.icon_url) {

								if (!icons[icon.id]) { //se registra por primera vez
									icons[icon.id] = {
										url: icon.icon_url,
										tags: [iconCollection.tag],
										category: tags[iconCollection.tag].categoria
									}
								} else if (icons[icon.id] && icons[icon.id].tags.indexOf(iconCollection.tag) == -1) {
                             
									icons[icon.id].tags.push(iconCollection.tag)
                                
								}

							}

						})
					}

				})
                
				console.log(`Cantidad de logos encontrados: ${Object.keys(icons).length}`)

				let promiseLogos = new Promise((resolve) => {
					//iconos => json file
					jsonfile.writeFile("./icons.json", icons, (err) => {
						if (!err) {
							console.log("Archivo de Logos listo")
							resolve()
						}
						if(err) console.log(err)
					})
				})

				let promiseTags = new Promise((resolve) => {
					//tags => json file
					jsonfile.writeFile("./tags.json", tags, (err) => {
						if (!err) {
							console.log("Archivo de tags listo")
							resolve()
						}
						if(err) console.log(err)
					})
				})

				Promise.all([promiseLogos, promiseTags]).then(() => {

					console.log("Todo preparado")

					exec("go run main.go", {cwd: "../upload"}, (error, stdout) => {
						if (error) {
							console.error(` ${error}`)
							return
						}
						console.log(`${stdout}`)
					})
				})


			})
			.catch((err) => {
				if(err) console.log(err)
			})


	})*/
