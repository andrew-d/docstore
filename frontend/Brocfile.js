/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app'),
    es3 = require('broccoli-es3-safe-recast'),
    isProduction = ( process.env.EMBER_ENV || 'development' ) === 'production';

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
  app.import('bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.' + ext, {
    destDir: 'fonts',
  });
});

// ----------------------------------------------------------------------
// Holder.js (https://github.com/imsky/holder)
app.import('bower_components/holderjs/holder.js');

// ----------------------------------------------------------------------
// Moment.js (https://github.com/moment/moment)
app.import({
  development: 'bower_components/momentjs/moment.js',
  production:  'bower_components/momentjs/min/moment.min.js',
});



var appTree = app.toTree();

// Only apply the ES3 transform in production - it takes a long time.
if( isProduction ) {
  appTree = es3(appTree);
}

module.exports = appTree;
