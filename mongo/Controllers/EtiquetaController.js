const Etiqueta = require("../Models/Etiqueta.js")
const Idioma = require("../Models/Idioma.js")
const async = require("async")

exports.ObtenerTodos = (data, res) => 
{
	Etiqueta.ObtenerTodos((err, data) => {
		if (data.length > 0) {
			return {
				status: 200,
				data: data
			}
		} else {
			return {
				status: 404,
				data: "No hay etiquetas en la base de datos",
			}
		}
	})
}

exports.Guardar = (data, res) => 
{
	const etiquetas = data.etiquetas

	let insertIds = []

	// itera cada etiqueta
	async.forEachOf(etiquetas, (etiqueta, keyEtiqueta, callback) => {

		Etiqueta.ObtenerPorTraduccion(etiqueta["ENG"], (err, data) => {

			if (data.length > 0) {

				data.forEach( el => {
					insertIds.push(el._id)
				})

				callback()

			} else {

				let keys = Object.keys(etiqueta).filter( (val) => {
					return val != "categoria"
				})

				let eitquetaOriginal = keys.map(function(val) {
					return {
						idioma: val,
						valor: etiqueta[val]
					}
				})

				let et = { traducciones: eitquetaOriginal }

				etiquetas[keyEtiqueta] = et

				// itera las traducciones de la etiqueta actual
				async.forEachOf(et.traducciones, (traduccion, keyTraduccion, callback) => {

					let normalize = (() => {
						const from = "ÃÀÁÄÂÈÉËÊÌÍÏÎÒÓÖÔÙÚÜÛãàáäâèéëêìíïîòóöôùúüûÑñÇç",
							to = "AAAAAEEEEIIIIOOOOUUUUaaaaaeeeeiiiioooouuuunncc",
							mapping = {}

						for (let i = 0, j = from.length; i < j; i++)
							mapping[from.charAt(i)] = to.charAt(i)

						return (str) => {
							let ret = []
							for (var i = 0, j = str.length; i < j; i++) {
								let c = str.charAt(i)
								if (mapping.hasOwnProperty(str.charAt(i)))
									ret.push(mapping[c])
								else
									ret.push(c)
							}
							return ret.join("")
						}
					})()
					
					etiquetas[keyEtiqueta].traducciones[keyTraduccion].valor = normalize(etiquetas[keyEtiqueta].traducciones[keyTraduccion].valor.toLowerCase())

					// obtiene el idioma de la traduccion actual
					Idioma.ObtenerPorCodigo(traduccion.idioma, (err, data) => {
						if (data !== null) {
							// sobreescribe el campo idioma de la etiqueta actual
							etiquetas[keyEtiqueta].traducciones[keyTraduccion].idioma = data._id
							callback()

						} else {
							return callback("No existe el idioma")
						}
					})

				}, err => { // fin de each para las traducciones

					if (err) return callback(err)

					let etiquetaData = etiquetas[keyEtiqueta]
					etiquetaData.iconos = []

					// guardamos la etiqueta sobreescrita despues que termine el loop de sus traducciones
					Etiqueta.Guardar(etiquetaData, (err, data) => {
						if (typeof data !== "undefined" && data.insertId) {
							insertIds.push(data.insertId)
							callback()
						} else {
							return callback(err)
						}
					})
				})
			}
		})

	}, err => { // fin de each para las etiquetas

		if (err) {
			res({
				status: 500,
				error: err
			})
		} else {

			if(insertIds.length){
				console.log("Etiquetas devueltas!")
				res({
					status: 200,
					data: insertIds
				})
			}else{
				res({
					status: 500,
					error: "Algo ocurrio"
				})
			}
		}

	})

}

/*
exports.ObtenerPorIcono = (data, res) => 
{
	const id = req.params.id;

	Etiqueta.ObtenerPorIcono(id, (err, data) => {
		if (data.length > 0) {
			res.status(200).json(data);
		} else {
			res.status(404).json({
				'msg': 'No hay etiquetas en la base de datos'
			});
		}
	})
}

exports.Actualizar = (data, res) => 
{
	const _id = req.body._id;
	const etiquetaData = req.body.etiqueta;

	Etiqueta.Actualizar(_id, etiquetaData, (err, data) => {
		if (data !== null && data.affectedRow) {
			res.status(200).json(data);
		} else {
			res.status(500).json({
				'msg': 'Algo ocurrio'
			});
		}
	})
}
*/

exports.AsignarIconos = (data, res) => 
{
	const _ids = data._ids
	const idsIconos = data.iconos

	let affectedRows = []

	
	async.forEachOf(_ids, (id, key, callback) => {

		Etiqueta.AsignarIconos(id, idsIconos, (err, data) => {
			if (data !== null && data.affectedRow) {
				affectedRows.push(data.affectedRow._id)
				callback()
			} else {
				return callback(err)
			}
		})

	}, err => {
		if (err) {
			res({
				status: 500,
				error: "Algo ocurrio"
			})
		}
		console.log("Iconos asignados!")
		res({
			status: 200,
			data: affectedRows
		})

		
	})
}

/*exports.DesasignarIcono = (data, res) => 
{
	const _id = req.params._id;
	const idIcono = req.body.idIcono;

	Etiqueta.DesasignarIcono(_id, idIcono, (err, data) => {
		if (data !== null && data.affectedRow) {
			res.status(200).json(data);
		} else {
			res.status(500).json({
				'msg': 'Algo ocurrio'
			});
		}
	})
}

exports.Borrar = (data, res) => 
{
	const _id = req.params._id;

	Etiqueta.Borrar(_id, (err, data) => {
		if (data !== null && data.affectedRow) {
			res.status(200).json(data);
		} else {
			res.status(500).json({
				'msg': 'Algo ocurrio'
			});
		}
	})
}¨*/