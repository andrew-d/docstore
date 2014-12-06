var path      = require('path'),
    gulp      = require('gulp'),
    gutil     = require('gulp-util'),
    express   = require('express'),
    sass      = require('gulp-sass'),
    minifyCSS = require('gulp-minify-css'),
    clean     = require('gulp-clean'),
    watch     = require('gulp-watch'),
    rev       = require('gulp-rev'),
    tiny_lr   = require('tiny-lr'),
    webpack   = require("webpack");


// ----------------------------------------------------------------------
// CONFIGURATION
// ----------------------------------------------------------------------
var httpPort = 4000;

var webpackConfig = require("./webpack.config.js");
if( gulp.env.production ) {  // i.e. we were executed with a --production option
    webpackConfig.plugins = webpackConfig.plugins.concat(
        new webpack.optimize.UglifyJsPlugin());
    webpackConfig.output.filename = "main-[hash].js";
}

// Paths to files in bower_components that should be copied to dist/assets/vendor
var vendorPaths = [
    'es5-shim/es5-sham.js',
    'es5-shim/es5-shim.js',
    'bootstrap/dist/css/bootstrap.min.css'
];

// Configuraton for SASS
var sassConfig = {
    includePaths: ['src/styles'],
}

// ----------------------------------------------------------------------
// TASKS
// ----------------------------------------------------------------------

gulp.task('clean', function() {
    gulp.src([
        'dist/**/*',
        '!dist/.gitignore',
    ], {read: false}).pipe(clean());
});

gulp.task('sass', function() {
    gulp.src('src/styles/main.scss')
        .pipe(sass(sassConfig).on('error', gutil.log))
        .pipe(gulp.env.production ? minifyCSS() : gutil.noop())
        .pipe(gulp.env.production ? rev()       : gutil.noop())
        .pipe(gulp.dest('dist/assets'))
});


// Some JS and CSS files we want to grab from Bower and put them in a
// dist/assets/vendor directory.  For example, the es5-sham.js is loaded in the
// HTML only for IE via a conditional comment.
gulp.task('vendor', function() {
    var paths = vendorPaths.map(function(p) {
        return path.resolve("./bower_components", p);
    });

    gulp.src(paths)
        .pipe(gulp.dest('dist/assets/vendor'));
});

// Just copy over remaining assets to dist. Exclude the styles and scripts as
// we process those elsewhere
gulp.task('copy', function() {
    gulp.src([
        'src/**/*',
        '!src/Bootstrap.jsx',
        '!src/scripts',
        '!src/scripts/**/*',
        '!src/styles',
        '!src/styles/**/*'
    ]).pipe(gulp.dest('dist'));
});

// This task lets Webpack take care of all the coffeescript and JSX
// transformations, defined in webpack.config.js.  Webpack also does its own
// uglification if we are in --production mode
gulp.task('webpack', function(callback) {
    execWebpack(webpackConfig);
    return callback()
});

gulp.task('dev', ['build'], function() {
    servers = createServers(httpPort, 35729);

    // When /src changes, fire off a rebuild
    gulp.watch(['./src/**/*'], function(evt) {
        gulp.run('build');
    });

    // When /dist changes, tell the browser to reload
    gulp.watch(['./dist/**/*'], function(evt) {
        gutil.log(gutil.colors.cyan(evt.path), 'changed')
        servers.lr.changed({
            body: {
                files: [evt.path],
            },
        });
    });
});



gulp.task('build', ['webpack', 'sass', 'copy', 'vendor'], function() {});

gulp.task('default', ['build'], function() {
    // Give first-time users a little help
    setTimeout(function() {
        gutil.log("**********************************************");
        gutil.log("* gulp              (development build)");
        gutil.log("* gulp clean        (rm /dist)");
        gutil.log("* gulp --production (production build)");
        gutil.log("* gulp dev          (build and run dev server)");
        gutil.log("**********************************************");
    }, 3000);
});


// ----------------------------------------------------------------------
// HELPERS
// ----------------------------------------------------------------------
// Create both http server and livereload server
var createServers = function(port, lrport) {
    lr = tiny_lr();
    lr.listen(lrport, function() {
        gutil.log("LiveReload listening on", lrport);
    });

    app = express();
    app.use(express.static(path.resolve("./dist")));
    app.listen(port, function() {
        gutil.log("HTTP server listening on", port);
    });

    return {
        lr: lr,
        app: app
    };
};


// Run webpack and handle errors.
var execWebpack = function(config) {
    webpack(config, function(err, stats) {
        if( err ) {
            throw new gutil.PluginError("execWebpack", err);
        }

        gutil.log("[execWebpack]", stats.toString({colors: true}))
    });
};
