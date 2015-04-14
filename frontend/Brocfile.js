/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app'),
    es3 = require('broccoli-es3-safe-recast'),
    pickFiles = require('broccoli-static-compiler');

var isProduction = ( process.env.EMBER_ENV || 'development' ) === 'production';

var app = new EmberApp();

// ----------------------------------------------------------------------
// Moment.js (https://github.com/moment/moment)
app.import('bower_components/moment/min/moment.min.js');

// ----------------------------------------------------------------------
// Bootstrap collapse component
app.import('bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/collapse.js');

// ----------------------------------------------------------------------
// Font Awesome (https://fortawesome.github.io/Font-Awesome/)
app.import({
  development: 'bower_components/fontawesome/css/font-awesome.css',
  production: 'bower_components/fontawesome/css/font-awesome.min.css',
});
app.import('bower_components/fontawesome/fonts/FontAwesome.otf', {
  destDir: 'fonts',
});
['eot', 'svg', 'ttf', 'woff', 'woff2'].forEach(function(ext) {
  app.import('bower_components/fontawesome/fonts/fontawesome-webfont.' + ext, {
    destDir: 'fonts',
  });
});

// ----------------------------------------------------------------------
// mOxie (https://github.com/moxiecode/moxie)
app.import({
  development: 'bower_components/moxie/bin/js/moxie.js',
  production: 'bower_components/moxie/bin/js/moxie.min.js',
});
app.import('bower_components/moxie/bin/flash/Moxie.min.swf', {
  destDir: 'assets',
});
app.import('bower_components/moxie/bin/silverlight/Moxie.min.xap', {
  destDir: 'assets',
});





var appTree = app.toTree();

// Only apply the ES3 transform in production - it takes a long time.
if( isProduction ) {
  appTree = es3(appTree);
}

module.exports = appTree;
