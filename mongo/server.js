const config = require('./config.js');

const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const bodyParser = require('body-parser');
const compression = require('compression');

const app = express();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));
app.use(compression());
app.enable('trust proxy');

const rutas = require('./routes.js');

app.use('/app', rutas);

app.use(function(req, res, next) {
  var err = new Error('No se encuentra');
  err.status = 404;
  next(err);
});

app.disable('x-powered-by');

app.use(function(err, req, res, next) {
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  res.status(err.status || 500);
  res.send({ error: err.message });
});

app.listen(config.puerto, function () {
  console.log('Servidor corriendo en : modo('+config.servidor+'), puerto('+config.puerto+')');
});

module.exports = app;