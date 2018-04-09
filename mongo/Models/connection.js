const MongoClient = require('mongodb').MongoClient;
const ObjectId = require('mongodb').ObjectID;

// Connection URL
const datos = {
    url: "mongodb+srv://cis:unaclavemuysegura5000@bazam-cgzwr.mongodb.net/admin",
    database: "disenador"
};

// Use connect method to connect to the server
const mongo = {
    connection: callback => { 
        MongoClient.connect(datos.url, (err, client) => {
            if (err) throw err;
            const db = client.db(datos.database);
            callback(db)
            client.close();
        })
    },
    objectId: id => {
        return new ObjectId(id);
    }
}

module.exports = mongo;