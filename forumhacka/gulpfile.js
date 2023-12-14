const gulp = require('gulp');
const concat = require('gulp-concat');
const autoprefixer = require('gulp-autoprefixer');
const cleanCSS = require('gulp-clean-css');
const uglify = require('gulp-uglify');
const del = require('del');
const browserSync = require('browser-sync').create();
const sass = require('gulp-sass');
const sourcemaps = require('gulp-sourcemaps');
const plumber = require('gulp-plumber');
const notify = require("gulp-notify");
const fileinclude = require('gulp-file-include');
const tiny = require('gulp-tinypng-nokey');

function styles() {
   return gulp.src('./src/scss/style.scss')
   .pipe(plumber({errorHandler: notify.onError("Error: <%= error.message %>")}))
   .pipe(sass().on('error', sass.logError))
   .pipe(autoprefixer({
      browsers: ['last 2 versions'],
      cascade: false
   }))
   .pipe(cleanCSS({
      level: 2
   }))
   .pipe(gulp.dest('./build/css'))
   .pipe(browserSync.stream());
}

const jsFiles = [
   './node_modules/jquery/dist/jquery.min.js',
   './node_modules/bootstrap/dist/js/bootstrap.min.js',
   './node_modules/magnific-popup/dist/jquery.magnific-popup.min.js',
   './src/external/perfect-scrollbar/perfect-scrollbar.min.js',
   './src/js/**/*.js'
];
function scripts() {
   return gulp.src(jsFiles)
   .pipe(concat('bundle.js'))
   .pipe(uglify({
      toplevel: true
   }))
   .pipe(gulp.dest('./build/js'))
   .pipe(browserSync.stream());
}

function clean() {
   return del(['build/*'])
}

function html() {
   return gulp.src('./src/*.html')
   .pipe(plumber({errorHandler: notify.onError("Error: <%= error.message %>")}))
   .pipe(fileinclude({
      prefix: '@@',
      basepath: '@file'
    }))
   .pipe(gulp.dest('./build/'))
   .pipe(browserSync.stream());
}

function favicon() {
   return gulp.src('./src/favicon/*')
   .pipe(gulp.dest('./build/favicon/'))
   .pipe(browserSync.stream());
}

function img() {
   return gulp.src('./src/images/*')
   // .pipe(tiny())
   .pipe(gulp.dest('./build/images/'))
   .pipe(browserSync.stream());
}



function watch() {
   browserSync.init({
      server: {
          baseDir: "./build/"
      }
  });
  gulp.watch('./src/scss/**/*.scss', styles)
  gulp.watch('./src/js/**/*.js', scripts)
  gulp.watch('./src/**/*.html', gulp.series(html)).on('change', browserSync.reload);
}

gulp.task('styles', styles);
gulp.task('html', html);
gulp.task('scripts', scripts);
gulp.task('del', clean);
gulp.task('img', img);
gulp.task('favicon', favicon);

gulp.task('watch', watch);
gulp.task('build', gulp.series(clean, gulp.parallel(html,styles,scripts,favicon,img)));
gulp.task('dev', gulp.series('build','watch'));