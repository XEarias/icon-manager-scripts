const MongoClient = require("mongodb").MongoClient
const ObjectId = require("mongodb").ObjectID
const config = require("../config.js")

// Connection URL
const datos = config.mongo

// Use connect method to connect to the server
const mongo = {
	connection: callback => { 
		MongoClient.connect(datos.url, (err, client) => {
			if (err) {console.log(err); return}
			const db = client.db(datos.database)
			callback(db)
			client.close()
		})
	},
	objectId: id => {
		return new ObjectId(id)
	}
}

module.exports = mongo