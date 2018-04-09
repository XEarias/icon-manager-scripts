const Etiqueta = require('../Models/Etiqueta.js');
const Idioma = require('../Models/Idioma.js');
const async = require('async');

exports.ObtenerTodos = () => 
{
	Etiqueta.ObtenerTodos((err, data) => {
		if (data.length > 0) {
			return {
                status : true,
                data : data
            };
		} else {
			return {
                status : false,
                msg: 'No hay etiquetas en la base de datos',
            };
		}
	});
}

exports.GuardarEtiquetas = (etiquetas) => 
{
	let insertIds = [];

	// itera cada etiqueta
	async.forEachOf(etiquetas, (etiqueta, keyEtiqueta, callback) => {

		// itera las traducciones de la etiqueta actual
		async.forEachOf(etiqueta.traducciones, (traduccion, keyTraduccion, callback) => {

			let normalize = (() => {
				const from = "ÃÀÁÄÂÈÉËÊÌÍÏÎÒÓÖÔÙÚÜÛãàáäâèéëêìíïîòóöôùúüûÑñÇç",
					to = "AAAAAEEEEIIIIOOOOUUUUaaaaaeeeeiiiioooouuuunncc",
					mapping = {};

				for (let i = 0, j = from.length; i < j; i++)
					mapping[from.charAt(i)] = to.charAt(i);

				return (str) => {
					let ret = [];
					for (var i = 0, j = str.length; i < j; i++) {
						let c = str.charAt(i);
						if (mapping.hasOwnProperty(str.charAt(i)))
							ret.push(mapping[c]);
						else
							ret.push(c);
					}
					return ret.join('');
				}
			})();

			etiquetas[keyEtiqueta].traducciones[keyTraduccion].valor = normalize(etiquetas[keyEtiqueta].traducciones[keyTraduccion].valor.toLowerCase());

			// obtiene el idioma de la traduccion actual
			Idioma.ObtenerPorCodigo(traduccion.idioma, (err, data) => {
				if (data !== null) {
					// sobreescribe el campo idioma de la etiqueta actual
					etiquetas[keyEtiqueta].traducciones[keyTraduccion].idioma = data._id;
					callback();

				} else {
					return callback({
						msg: "No existe el idioma"
					});
				}
			})
		}, err => { // fin de each para las traducciones

			if (err) return callback(err);

			let etiquetaData = etiquetas[keyEtiqueta];
			etiquetaData.iconos = [];

			// guardamos la etiqueta sobreescrita despues que termine el loop de sus traducciones
			Etiqueta.Guardar(etiquetaData, (err, data) => {
				if (typeof data !== 'undefined' && data.insertId) {
					insertIds.push(data.insertId);
					callback();
				} else {
					return callback(err);
				}
			})
		})

	}, err => { // fin de each para las etiquetas

		if (err) {
			return {
                status : false,
                msg: err,
            };
		} else {
			return {
                status : true,
                data : insertIds
            };
		}

	})

}


exports.ObtenerPorIcono = (id) => 
{
	Etiqueta.ObtenerPorIcono(id, (err, data) => {
		if (data.length > 0) {
			return {
                status : true,
                data : data
            };
		} else {
            return {
                status : false,
                msg : 'No hay etiquetas en la base de datos'
            };
		}
	})
}

exports.Actualizar = (_id, etiqueta) => 
{
	const etiquetaData = etiqueta;

	Etiqueta.Actualizar(_id, etiquetaData, (err, data) => {
		if (data !== null && data.affectedRow) {
			return {
                status : true,
                data : data
            };
		} else {
            return {
                status : false,
                msg : 'Algo ocurrio'
            };
		}
	})
}

exports.AsignarIconos = (_ids, iconos) => 
{
	const idsIconos = req.body.iconos;

	let affectedRows = [];

	async.forEachOf(_ids, (id, key, callback) => {

		Etiqueta.AsignarIconos(id, idsIconos, (err, data) => {
			if (data !== null && data.affectedRow) {
				affectedRows.push(data.affectedRow);
				callback();
			} else {
				return callback(err);
			}
		})

	}, err => {
		if (err) {
            return {
                status : false,
                msg : 'Algo ocurrio'
            }
        };

		return {
            status : true,
            data : affectedRows
        };
    })
}

exports.DesasignarIcono = (_id, idIcono) => 
{
	Etiqueta.DesasignarIcono(_id, idIcono, (err, data) => {
		if (data !== null && data.affectedRow) {
            return {
                status : true,
                data : data
            };
		} else {
            return {
                status : false,
                msg : 'Algo ocurrio'
            };
		}
	})
}

exports.Borrar = (_id) => 
{
	Etiqueta.Borrar(_id, (err, data) => {
		if (data !== null && data.affectedRow) {
			return {
                status : true,
                data : data
            };
		} else {
            return {
                status : false,
                msg : 'Algo ocurrio'
            };
		}
	})
}