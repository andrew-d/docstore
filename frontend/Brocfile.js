/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app'),
    es3 = require('broccoli-es3-safe-recast'),
    isProduction = ( process.env.EMBER_ENV || 'development' ) === 'production';

var app = new EmberApp();

// ----------------------------------------------------------------------
// Pure (https://github.com/yahoo/pure/)
app.import({
  development: 'bower_components/pure/pure.css',
  production:  'bower_components/pure/pure.min.css',
});

// ----------------------------------------------------------------------
// Font Awesome (https://fortawesome.github.io/Font-Awesome/)
app.import({
  development: 'bower_components/fontawesome/css/font-awesome.css',
  production:  'bower_components/fontawesome/css/font-awesome.min.css',
});
app.import('bower_components/fontawesome/fonts/FontAwesome.otf', {
  destDir: 'fonts',
});
['eot', 'svg', 'ttf', 'woff'].forEach(function(ext) {
  app.import('bower_components/fontawesome/fonts/fontawesome-webfont.' + ext, {
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

// ----------------------------------------------------------------------
// Lodash (https://lodash.com)
app.import({
  development: 'bower_components/lodash/dist/lodash.js',
  production:  'bower_components/lodash/dist/lodash.min.js'
}, {
  'lodash': [
    'default'
  ]
});



var appTree = app.toTree();

// Only apply the ES3 transform in production - it takes a long time.
if( isProduction ) {
  appTree = es3(appTree);
}

module.exports = appTree;
