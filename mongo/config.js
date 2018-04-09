const node_env = process.env.NODE_ENV || "desarrollo";

const configuracion = {
    produccion : {
        servidor : "Produccion",
        mongo : {
            url: "mongodb+srv://cis:unaclavemuysegura5000@bazam-cgzwr.mongodb.net/admin",
            database : "disenador"
        }
    },

    desarrollo : {
        servidor : "Desarrollo",
        secret :"unaclavemuysegura",
    	seguridad :true,
    	puerto : "666",
    	url : "http://localhost:666",
        mongo : {
            url : "mongodb+srv://cis:unaclavemuysegura5000@bazam-cgzwr.mongodb.net/admin",
            database : "disenador"
        }
    }
}


module.exports = configuracion[node_env];