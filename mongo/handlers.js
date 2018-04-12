const Controllers = require('./Controllers/index.js');

module.exports = {
    obtenerTodos: Controllers.EtiquetaController.ObtenerTodos,
    guardar: Controllers.EtiquetaController.Guardar,
    asignarIconos: Controllers.EtiquetaController.AsignarIconos
}