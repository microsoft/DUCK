"use strict";

import plugins  from "gulp-load-plugins";
import yargs    from "yargs";
import browser  from "browser-sync";
import gulp     from "gulp";
import rimraf   from "rimraf";
import yaml     from "js-yaml";
import fs       from "fs";
import gulpgo   from "gulp-go";

// Load all Gulp plugins into one variable
const $ = plugins();

// Check for --production flag
const PRODUCTION = !!(yargs.argv.production);

// Load settings from settings.yml
const { COMPATIBILITY, PORT, UNCSS_OPTIONS, PATHS } = loadConfig();

function loadConfig() {
    let ymlFile = fs.readFileSync("config.yml", "utf8");
    return yaml.load(ymlFile);
}

// Build the "dist" folder by running all of the below tasks
gulp.task("build",
    gulp.series(clean, backend, gulp.parallel(pages, sass, vendorJS, javascript, images, copy)));

// Build the site, run the server, and watch for file changes
gulp.task("default",
    gulp.series("build", server, watch));

// Delete the "dist" folder every time a build starts.
function clean(done) {
    rimraf(PATHS.dist, done);
}

// Copy files out of the assets folder. This task skips over the "img", "js", and "scss" folders, which are parsed separately
function copy() {
    return gulp.src(PATHS.assets)
        .pipe(gulp.dest(PATHS.dist + "/assets"));
}

// Copy page templates into finished HTML files
function pages() {
    return gulp.src("frontend/src/pages/**/*")
        .pipe(gulp.dest(PATHS.dist));
}

function resetPages(done) {
    done();
}

// Compile Sass into CSS. In production, the CSS is compressed
function sass() {
    return gulp.src("frontend/src/assets/scss/app.scss")
        .pipe($.sourcemaps.init())
        .pipe($.sass({
                includePaths: PATHS.sass
            })
            .on("error", $.sass.logError))
        .pipe($.autoprefixer({
            browsers: COMPATIBILITY
        }))
        .pipe($.if(PRODUCTION, $.uncss(UNCSS_OPTIONS)))
        .pipe($.if(PRODUCTION, $.cssnano()))
        .pipe($.if(!PRODUCTION, $.sourcemaps.write()))
        .pipe(gulp.dest(PATHS.dist + "/assets/css"))
        .pipe(browser.reload({stream: true}));
}

// Combine JavaScript into one file. In production, the file is minified.
function vendorJS() {
    return gulp.src(PATHS.vendorJS)
        .pipe($.sourcemaps.init())
        .pipe($.concat("vendor.js"))
        .pipe($.if(PRODUCTION, $.uglify()
            .on("error", e => {
                console.log(e);
            })
        ))
        .pipe($.if(!PRODUCTION, $.sourcemaps.write()))
        .pipe(gulp.dest(PATHS.dist + "/assets/js"));
}

function javascript() {
    return gulp.src(PATHS.javascript)
        .pipe($.sourcemaps.init())
        .pipe($.babel())
        .pipe($.concat("app.js"))
        .pipe($.if(PRODUCTION, $.uglify()
            .on("error", e => {
                console.log(e);
            })
        ))
        .pipe($.if(!PRODUCTION, $.sourcemaps.write()))
        .pipe(gulp.dest(PATHS.dist + "/assets/js"));
}

// Copy images to the "dist" folder. In production, the images are compressed.
function images() {
    return gulp.src("frontend/src/assets/img/**/*")
        .pipe($.if(PRODUCTION, $.imagemin({
            progressive: true
        })))
        .pipe(gulp.dest(PATHS.dist + "/assets/img"));
}

// proxy the backend server to support browser reloading
function server(done) {
    browser.init({
        proxy: "localhost:3000",
        port: PORT
    });
    done();
}

var go;

// launch the backend serving the web distribution directory
function backend(done) {
    go = gulpgo.run("main.go", ["--webdir", __dirname + "/" + PATHS.dist], {cwd: "backend", stdio: "inherit"});
    done();
}


// Watch for changes to frontend assets and backend Go code
function watch() {
    gulp.watch(PATHS.assets, copy);
    gulp.watch("frontend/src/pages/**/*.html", gulp.series(pages, browser.reload));
    gulp.watch("frontend/src/{layouts,partials}/**/*.html", gulp.series(resetPages, pages, browser.reload));
    gulp.watch("frontend/src/assets/scss/**/*.scss", sass);
    gulp.watch("frontend/src/assets/js/**/*.js", gulp.series(javascript, browser.reload));
    gulp.watch("frontend/src/assets/img/**/*", gulp.series(images, browser.reload));
    gulp.watch(["backend/**/*.go"]).on("change", function () {
        go.restart();
    });
}
