var extend = require('lodash-node/modern/object/assign'),
    path = require('path'),
    webpack = require('webpack'),
    ExtractTextPlugin = require('extract-text-webpack-plugin'),
    HtmlWebpackPlugin = require('html-webpack-plugin');


// ----------------------------------------------------------------------
// Configuration
var options = {
    // Address to start webpack-dev-server on.
    devServerAddr: 'localhost',
    devServerPort: 3000,
};

var IS_PRODUCTION = ('production' === process.env.NODE_ENV);


// ----------------------------------------------------------------------
// Computed configuration (paths, etc.)

var sassOptions = [
  'precision=10',
  'outputStyle=expanded',
  'sourceMap=' + !IS_PRODUCTION,
  'includePaths[]=' + path.resolve(__dirname, './bower_components/bootstrap-sass-official/assets/stylesheets'),
].join('&');

var devServerHost = options.devServerAddr + ':' + options.devServerPort;

// ----------------------------------------------------------------------
// Webpack Config

var plugins = [
  new HtmlWebpackPlugin({
    template: 'app/index.html',
  })
];

// Add plugins for production vs. development.
if( IS_PRODUCTION ) {
  // Production configuration
  plugins.push(
    new webpack.optimize.OccurenceOrderPlugin({ preferEntry: true }),
    new webpack.optimize.UglifyJsPlugin(),
    new webpack.optimize.DedupePlugin(),
    new webpack.DefinePlugin({
        "process.env": {
            NODE_ENV:   JSON.stringify("production"),
            NODE_DEBUG: JSON.stringify(""),
        }
    })
  );

} else {
  // Development configuration
  plugins.push(
    new webpack.HotModuleReplacementPlugin(),
    new webpack.DefinePlugin({
        "process.env": {
            NODE_ENV:   JSON.stringify(process.env.NODE_ENV || ""),
            NODE_DEBUG: JSON.stringify(process.env.NODE_DEBUG || ""),
        }
    })
  );
}

// Javascript loaders
var jsLoaders = ['react-hot', 'jsx?harmony'];
var scriptModLoaders = [
    { test: /\.js$/,   loaders: jsLoaders, exclude: /node_modules/ },
    { test: /\.jsx$/,  loaders: jsLoaders, exclude: /node_modules/ },
    { test: /\.json$/, loaders: ['json'] },
];

// Styles loaders
var cssLoaders = [
  'style',

  // We apply autoprefixer-loader to files that are imported, in addition to
  // our files.  We also minimize in production.
  'css?importLoaders=1' + (IS_PRODUCTION ? '&minimize' : ''),

  // Autoprefixer settings
  'autoprefixer-loader?browsers=last 3 versions'
];
var styleModLoaders = [
    { test: /\.css$/, loader: cssLoaders },
    { test: /\.scss$/, loader: cssLoaders.concat(["sass-loader?" + sassOptions]) },
];

// If we're in production, we extend the style loaders to also extract text.
if( IS_PRODUCTION ) {
  styleModLoaders = styleModLoaders.map(function(x) {
    return {
      test: x.test,
      // Strip off the 'style' loader
      loader: ExtractTextPlugin.extract(x.loader.slice(1).join('!')),
    };
  });

  // Add the plugin to write the final CSS file.
  plugins.push(
    new ExtractTextPlugin("assets/styles-[hash].css", { allChunks: true })
  );
} else {
  // Just join them together, since arrays don't seem to work
  styleModLoaders = styleModLoaders.map(function(x) {
    return {
      test: x.test,
      loader: x.loader.join('!'),
    };
  });
}


// Static file loaders
var fileLoaderOptions = { name: 'assets/[hash].[ext]' },
    urlLoaderOptions = extend({}, fileLoaderOptions, {
      limit: 10000,
    });
var staticModLoaders = [
  // Various font types from Bootstrap
  {
      test: /\.(woff|woff2)$/,
      loader: "url-loader",
      query: extend({}, urlLoaderOptions, { mimetype: 'application/font-woff' }),
  },
  {
    test: /\.ttf$/,
    loader: "file-loader",
    query: extend({}, fileLoaderOptions, { mimetype: 'application/vnd.ms-fontobject' }),
  },
  {
    test: /\.eot$/,
    loader: "file-loader",
    query: extend({}, fileLoaderOptions, { mimetype: 'application/x-font-ttf' }),
  },
  {
    test: /\.svg$/,
    loader: "file-loader",
    query: extend({}, fileLoaderOptions, { mimetype: 'image/svg+xml' }),
  },

  // Inline base64 URLs for small images, direct URLs for the rest
  {
      test: /\.png$/,
      loader: "url-loader",
      query: extend({}, urlLoaderOptions, { mimetype: 'image/png' }),
  },
  {
      test: /\.gif$/,
      loader: "url-loader",
      query: extend({}, urlLoaderOptions, { mimetype: 'image/gif' }),
  },
  {
      test: /\.(jpg|jpeg)$/,
      loader: "url-loader",
      query: extend({}, urlLoaderOptions, { mimetype: 'image/jpg' }),
  },
];

// This prevents any assets from being emitted if there are errors.
plugins.push(new webpack.NoErrorsPlugin());

// Entry points
var entries = ['./app/index'];

if( !IS_PRODUCTION ) {
  // Add entries for hot module reload.
  entries.push(
    'webpack-dev-server/client?http://' + devServerHost,
    'webpack/hot/only-dev-server'
  );
}

// Configuration for JSHint. Much of this taken from:
//    https://gist.github.com/connor/1597131
var jsHintOptions = {
  // Environment
  "esnext"        : true,   // We're using ES6 syntax
  "browser"       : true,   // We're in the browser
  "browserify"    : true,   // Browserify-style globals are allowed
  "debug"         : false,  // Disallow debugger statements - e.g. browser breakpoints
  "devel"         : true,   // Allow development statements - e.g. `console.log()`

  // Preventing errors
  "asi"           : false,  // Do not allow automatic semicolon insertion.
  "bitwise"       : true,   // Prohibit bitwise operators (&, |, ^, etc.).
  "boss"          : false,  // Tolerate assignments inside if, for & while.
                            // Usually conditions & loops are for comparison, not assignments.
  "curly"         : true,   // Require {} for every new block or scope.
  "eqeqeq"        : true,   // Require triple equals i.e. `===`.
  "eqnull"        : false,  // Tolerate use of `== null`.
  "evil"          : false,  // Tolerate use of `eval`.
  "expr"          : false,  // Tolerate `ExpressionStatement` as programs.
  "forin"         : false,  // Tolerate `for in` loops without `hasOwnPrototype`.
  "loopfunc"      : false,  // Allow functions to be defined within loops.
  "noarg"         : true,   // Prohibit use of `arguments.caller` and `arguments.callee`.
  "regexp"        : true,   // Prohibit `.` and `[^...]` in regular expressions.
  "supernew"      : false,  // Tolerate `new function () { ... };` and `new Object;`.
  "undef"         : true,   // Require all non-global variables be declared
                            // before they are used.
  "unused"        : true,   // Prevent defined and unused variables.

  // Styling preferences
  "newcap"        : true,   // Require capitalization of all constructor
                            // functions - e.g. `new F()`.
  "noempty"       : true,   // Prohibit use of empty blocks.
  "trailing"      : true,   // Prohibit trailing whitespaces.
  "indent"        : 2,      // Specify indentation spacing
  "strict"        : false,  // Require all functions to run in ES5's strict mode.
  "sub"           : true,   // Warn about using `[]` notation instead of using
                            // dot notation - but only where possible.
};

// Only enable some things if we're not in prod.
if( !IS_PRODUCTION ) {
  jsHintOptions.debug = true;   // Permit debugger statements - e.g. browser breakpoints
}


var config = {
  target: 'web',
  devtool: IS_PRODUCTION ? null : 'eval',

  entry: entries,

  output: {
    path:       __dirname + '/build/',
    publicPath: 'http://' + devServerHost + '/',

    // Cache busting in production
    filename:   'assets/' + (IS_PRODUCTION ? 'bundle-[hash].js' : 'bundle.js'),
  },

  resolve: {
    extensions: ['', '.js', '.jsx', '.json'],
  },

  // Various config options set above
  plugins: plugins,
  module: {
    preLoaders: [
      { test: /\.jsx?$/, exclude: /node_modules/,
        loader: 'jsxhint-loader?harmony' },
    ],

    loaders: [].concat(scriptModLoaders)
               .concat(styleModLoaders)
               .concat(staticModLoaders),
  },

  // Options for JSHint
  jshint: jsHintOptions,

  // For use from server.js
  userOptions: options,
};

// All done!
module.exports = config;
