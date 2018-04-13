const Mongo = require('./connection.js');
const Connection = Mongo.connection;
const objectId = Mongo.objectId;

let Idioma = {}

Idioma.ObtenerTodos = callback =>
{ 
	Connection(db => {
		const collection = db.collection('idiomas');
		collection.find({}).toArray((err, docs) => {
			if(err)	{console.log(err); return;}
		  	callback(null, docs);
		});
	})
}

Idioma.Obtener = (_id, callback) =>
{ 
	Connection(db => {
		const collection = db.collection('idiomas');
		collection.findOne({'_id': objectId(_id) }, (err, doc) => {
            if(err)	{console.log(err); return;}
		  	callback(null, doc);
		});
	})
}

Idioma.ObtenerPorCodigo = (codigo, callback) =>
{ 
	Connection(db => {
		const collection = db.collection('idiomas');
		collection.findOne({'codigo': codigo}, (err, doc) => {
            if(err)	{console.log(err); return;}
		  	callback(null, doc);
		});
	})
}
 
Idioma.Guardar = (idiomaData, callback) =>
{
	Connection(db => {
		const collection = db.collection('idiomas');
		collection.insertOne(idiomaData, (err, doc) => {
			if(err)	{console.log(err); return;}
		  	callback(null, { 'insertId': doc.insertedId });
		});
	})
}

Idioma.Actualizar = (_id, idiomaData, callback) =>
{
	Connection(db => {
		const collection = db.collection('idiomas');
		collection.findOneAndUpdate({ '_id': objectId(_id) }, { $set: idiomaData }, (err, doc) => {
			if (err) {console.log(err); return;}
            callback(null, { 'affectedRow': doc.value });
        });
	})
}

Idioma.Borrar = (_id, callback) =>
{
	Connection(db => {
		const collection = db.collection('idiomas');
		collection.findOneAndDelete({ '_id': objectId(_id) }, (err, doc) => {
            if (err) {console.log(err); return;}
            callback(null, { 'affectedRow': doc.value });
        });
	})
}

module.exports = Idioma;