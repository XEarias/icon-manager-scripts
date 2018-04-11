const node_env = process.env.NODE_ENV || "desarrollo";

const configuracion = {
    produccion : {
        servidor : "Produccion"
    },

    desarrollo : {
        servidor : "Desarrollo",
        secret :"unaclavemuysegura",
    	seguridad :true,
    	puerto : "666",
    	url : "http://127.0.0.1:666",
    }
}


module.exports = configuracion[node_env];