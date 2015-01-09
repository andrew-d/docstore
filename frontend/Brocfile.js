/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app');

var app = new EmberApp();

// ----------------------------------------------------------------------
// Bootstrap
app.import({
  development: 'bower_components/bootstrap/dist/css/bootstrap.css',
  production:  'bower_components/bootstrap/dist/css/bootstrap.min.css',
});
app.import({
  development: 'bower_components/bootstrap/dist/js/bootstrap.js',
  production:  'bower_components/bootstrap/dist/js/bootstrap.min.js',
});

['eot', 'svg', 'ttf', 'woff'].forEach(function(ext) {
  app.import('bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.' + ext);
});


module.exports = app.toTree();
