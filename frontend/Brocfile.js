/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app');

var app = new EmberApp();

// Bootstrap
app.import({
    development: 'bower_components/bootstrap/dist/css/bootstrap.css',
    production:  'bower_components/bootstrap/dist/css/bootstrap.min.css',
});
app.import({
    development: 'bower_components/bootstrap/dist/js/bootstrap.js',
    production:  'bower_components/bootstrap/dist/js/bootstrap.min.js',
});
[
    'glyphicons-halflings-regular.eot',
    'glyphicons-halflings-regular.svg',
    'glyphicons-halflings-regular.ttf',
    'glyphicons-halflings-regular.woff',
].forEach(function(n) {
    app.import('bower_components/bootstrap/dist/fonts/' + n, {
        destDir: 'fonts',
    });
});

// Moment.js
app.import({
    development: 'bower_components/momentjs/moment.js',
    production:  'bower_components/momentjs/min/moment.min.js',
});

// Humanize Plus (https://github.com/HubSpot/humanize)
app.import({
    development: 'bower_components/humanize-plus/public/src/humanize.js',
    production:  'bower_components/humanize-plus/public/dist/humanize.min.js',
});

// mOxie (https://github.com/moxiecode/moxie)
app.import({
    development: 'bower_components/moxie/bin/js/moxie.js',
    production:  'bower_components/moxie/bin/js/moxie.min.js',
});
app.import('bower_components/moxie/bin/flash/Moxie.min.swf', {
    destDir: 'assets',
});
app.import('bower_components/moxie/bin/silverlight/Moxie.min.xap', {
    destDir: 'assets',
});


module.exports = app.toTree();
