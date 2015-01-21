/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app'),
    es3 = require('broccoli-es3-safe-recast'),
    isProduction = ( process.env.EMBER_ENV || 'development' ) === 'production';

var app = new EmberApp({
  'ember-cli-selectize': {
    'theme': 'bootstrap3'
  },
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
// Moment.js (https://github.com/moment/moment)
app.import({
  development: 'bower_components/momentjs/moment.js',
  production:  'bower_components/momentjs/min/moment.min.js',
});

// ----------------------------------------------------------------------
// Humanize Plus (https://github.com/HubSpot/humanize)
app.import({
  development: 'bower_components/humanize-plus/public/src/humanize.js',
  production: 'bower_components/humanize-plus/public/dist/humanize.min.js',
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
