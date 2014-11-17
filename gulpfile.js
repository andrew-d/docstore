var gulp = require('gulp'),
    concat = require('gulp-concat'),
    filter = require('gulp-filter'),
    mainBowerFiles = require('main-bower-files'),
    minifyCSS = require('gulp-minify-css'),
    print = require('gulp-print'),
    uglify = require('gulp-uglify');


var filterByExtension = function(ext) {
    return filter(function(file) {
        return file.path.match(new RegExp('\\.' + ext + '$'));
    });
};


gulp.task('concat', function() {
    var mainFiles = mainBowerFiles();

    if( !mainFiles.length ) return;

    var jsFilter = filterByExtension('js');

    return gulp.src(mainFiles)
        .pipe(jsFilter)
        .pipe(print())
        .pipe(concat('vendor.js'))
        .pipe(gulp.dest('./app/static/js'))
        .pipe(jsFilter.restore())
        .pipe(filterByExtension('css'))
        .pipe(print())
        .pipe(concat('vendor.css'))
        .pipe(gulp.dest('./app/static/css'));
});

gulp.task('fonts', function() {
    return gulp.src("bower_components/bootstrap/dist/fonts/**/*")
        .pipe(gulp.dest('./app/static/fonts'));
});

gulp.task('minify-css', ['concat'], function() {
    return gulp.src("./app/static/css/vendor.css")
        .pipe(minifyCSS())
        .pipe(concat("vendor.min.css"))
        .pipe(gulp.dest('./app/static/css'));
});

gulp.task('minify-js', ['concat'], function() {
    return gulp.src("./app/static/js/vendor.js")
        .pipe(uglify())
        .pipe(concat("vendor.min.js"))
        .pipe(gulp.dest('./app/static/js'));
});

gulp.task('default', ['concat', 'fonts', 'minify-css', 'minify-js']);
