const NounProject = require("the-noun-project")
const jsonfile = require("jsonfile")
const csv = require("csvtojson")
const flags = require("node-flags")
const { exec } = require("child_process");

let cantidad = flags.get("CANT")
let saltoArg = flags.get("SALTO")

cantidad = cantidad ? parseInt(cantidad) : 0
saltoArg = saltoArg ? parseInt(saltoArg) : 0

let np = new NounProject({
	key: "c262b78520734708a4b90c9f951f708e",
	secret: "75639560369340d6b9211315f41bdaad"
})

let promises = []
let tags = {}

csv()
	.fromFile("./tags.csv")

	.on("json", (tag) => {
       
		var salto = saltoArg

		while(salto <= (cantidad + saltoArg)) {
			let promiseGet = new Promise((resolve) => {
                
				//console.log("inicio "+tag.ENG+"#"+salto)
				np.getIconsByTerm(tag.ENG, {offset: salto, limit_to_public_domain: 1}, (err, data) => {

					if (!err && data) {

						resolve({data: data.icons, tag: tag.ENG})

						tags[tag.ENG] = tag

					} else {
                       
						resolve(666)

					}
					/*
					var d = new Date()

					console.log(`Respuesta: ${d.getHours()}:${d.getMinutes()}:${d.getSeconds()}:${d.getMilliseconds()}...`)*/

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


	})