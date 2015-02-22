var express = require('express'),
    nomnom = require('nomnom'),
    httpProxy = require('http-proxy'),
    morgan = require('morgan'),
    url = require('url'),
    webpack = require('webpack'),
    when = require('when'),
    WebpackDevServer = require('webpack-dev-server');

var config = require('./webpack.config');


// This little script helps with running webpack-dev-server and proxying
// requests to a backend API.  It will do the following:
//
//  1. Start webpack-dev-server on the address that's specified in
//     webpack.config.js
//  2. Start a regular HTTP server on a configurable port that serves the
//     application's static assets, and redirects all unmatched requests to
//     the application's index page, to support the HTML5 history API.
//  3. If configured, proxies requests with a given prefix to a backend URL.


// ----------------------------------------------------------------------
// Config parsing
var devServerHost = 'http://' +
                    config.userOptions.devServerAddr + ':' +
                    config.userOptions.devServerPort;

var opts = nomnom
  .option('port', {
    abbr: 'p',
    default: 3001,
    help: 'Port to listen on',
    callback: function(port) {
      var parsed = parseInt(port);

      if( port != parsed ) {
        return "port must be an integer";
      }

      if( port <= 0 || port > 65536 ) {
        return "port must be in the range 1-65535";
      }
    },
  })
  .option('upstream', {
    default: null,
    help: "If specified, proxy all requests matching 'upstream-prefix' here",
  })
  .option('upstream-prefix', {
    default: '/api',
    help: 'Any request matching this prefix gets proxied upstream',
  })
  .parse();


// ----------------------------------------------------------------------
// Set up proxy
var app = express();

var devServerProxy = httpProxy.createProxyServer({
  target: devServerHost + '/assets',
});

devServerProxy.on('error', function(e) {
  console.log('Error proxying to webpack-dev-server');
  console.log(e);
});

// Proxy request for static assets to the webpack dev server.
app.use('/assets', function(req, res, next) {
  return devServerProxy.web(req, res);
});

if( opts.upstream ) {
  var upstreamTarget = opts.upstream + opts['upstream-prefix'];
  var upstreamProxy = httpProxy.createProxyServer({
    target: upstreamTarget,
    changeOrigin: true,
  });

  upstreamProxy.on('error', function(e) {
    console.log('Error proxying to ' + upstreamTarget);
    console.log(e);
  });

  console.log('Proxying to ' + upstreamTarget);

  // Print the status of requests to/from the upstream API.
  app.use(opts['upstream-prefix'], morgan('dev'));

  // Proxy things in this prefix to the API
  app.use(opts['upstream-prefix'], function(req, res, next) {
    return upstreamProxy.web(req, res);
  });
}

// Send the index page for all remaining requests.
app.get('/*', function(req, res) {
  res.sendFile(__dirname + '/build/index.html');
});


// ----------------------------------------------------------------------
// Set up webpack dev server
var devServer = new WebpackDevServer(webpack(config), {
  publicPath: config.output.publicPath,
  hot: true,
  stats: {
    colors: true,
    hash: true,
    timings: true,
    assets: true,
    chunks: true,

    // Hide some verbose information that we don't need.
    modules: false,
    cached: false,
  },
});

// ----------------------------------------------------------------------
// Run both servers
var p1 = when.promise(function(resolve, reject) {
  devServer.listen(
    config.userOptions.devServerPort,
    config.userOptions.devServerAddr,
    function (err, result) {
      if( err ) {
        console.log('Error starting webpack-dev-server: ' + err);
        reject(err);
        return;
      }

      resolve();
    });
});

var p2 = when.promise(function(resolve, reject) {
  app.listen(opts.port, function(err, result) {
    if( err ) {
      console.log('Error starting server: ' + err);
      reject(err);
      return;
    }

    resolve();
  });
});


when.all([p1, p2]).then(function() {
  console.log('Server listening at ' + 'localhost' + ':' + opts.port);
  console.log('Dev server listening at ' +
              config.userOptions.devServerAddr + ':' +
              config.userOptions.devServerPort);
  console.log('--------------------------------------------------');
  console.log('');

  // Should sleep forever.
}).catch(function(err) {
  console.log('Error starting');
});
