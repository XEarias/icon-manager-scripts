const express = require("express");
const router = express.Router();
const controllers = require("./Controllers");

router.get('/etiquetas', controllers.EtiquetaController.ObtenerTodos);
router.post('/etiquetas', controllers.EtiquetaController.GuardarEtiquetas);
router.post('/etiquetas/iconos', controllers.EtiquetaController.AsignarIconos);

module.exports = router;