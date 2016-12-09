"use strict";

import plugins from "gulp-load-plugins";
import yargs from "yargs";
import browser from "browser-sync";
import gulp from "gulp";
import rimraf from "rimraf";
import yaml from "js-yaml";
import fs from "fs";
import gulpgo from "gulp-go";
import gutil from "gulp-util";
import child from "child_process";
import addsrc from "gulp-add-src";
import zip from "gulp-zip";


// Load all Gulp plugins into one variable
const $ = plugins();

// Check for --production flag
const PRODUCTION = !!(yargs.argv.production);

// Load settings from settings.yml
const {COMPATIBILITY, PORT, UNCSS_OPTIONS, PATHS} = loadConfig();


if (process.env.GOPATH === undefined || process.env.GOPATH === null) {
    gutil.log(gutil.colors.red("GOPATH not set, aborting build"));
    process.exit(1);
}


function loadConfig() {
    let ymlFile = fs.readFileSync("config.yml", "utf8");
    return yaml.load(ymlFile);
}

// Build the "dist" folder by running all of the below tasks
gulp.task("build",
    gulp.series(clean, backend, gulp.parallel(pages, sass, vendorJS, javascript, images, config, partials, fonts, copy)));

// build the distribution
gulp.task("distro",
    gulp.series(clean, backendCompile, gulp.parallel(pages, sass, vendorJS, javascript, images, config, partials, fonts, copy), copyWeb, copyConfig, copyRuleBases, copyBinary, zipit));

// Build the site, run the server, and watch for file changes
gulp.task("default",
    gulp.series("build", server, watch));

// Delete the "dist" folder every time a build starts.
function clean(done) {
    rimraf(PATHS.dist, done);
    rimraf("./image", done);
}

function copyWeb() {
    return gulp.src(PATHS.dist + "/**")
        .pipe(gulp.dest("./image/stage/dist"));
}

function copyRuleBases() {
    return gulp.src("RuleBases/**")
        .pipe(gulp.dest("./image/stage/RuleBases"));
}

function copyBinary() {
    return gulp.src(process.env.GOPATH + "/bin/backend*")
        .pipe(gulp.dest("./image/stage"));
}

function copyConfig() {
    return gulp.src("configuration.json")
        .pipe(gulp.dest("./image/stage"));
}

function zipit() {
    return gulp.src("./image/stage/**")
        .pipe(zip('duck.zip'))
        .pipe(gulp.dest('./image'));
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
        // .pipe($.if(PRODUCTION, $.uncss(UNCSS_OPTIONS)))
        // .pipe($.if(PRODUCTION, $.cssnano()))
        .pipe(addsrc.append(PATHS.css))
        .pipe($.concat('app.css'))
        .pipe($.if(!PRODUCTION, $.sourcemaps.write()))
        .pipe(gulp.dest(PATHS.dist + "/assets/css"))
        .pipe(browser.reload({stream: true}));
}

// Combine JavaScript into one file. In production, the file is minified.
function vendorJS() {
    return gulp.src(PATHS.vendorJS)
        .pipe($.sourcemaps.init())
        .pipe($.concat("vendor.js"))
        // .pipe($.if(PRODUCTION, $.uglify()
        //     .on("error", e => {
        //         console.log(e);
        //     })
        // ))
        .pipe($.if(!PRODUCTION, $.sourcemaps.write()))
        .pipe(gulp.dest(PATHS.dist + "/assets/js"));
}

function javascript() {
    return gulp.src(PATHS.javascript)
        .pipe($.sourcemaps.init())
        .pipe($.babel())
        .pipe($.concat("app.js"))
        // .pipe($.if(PRODUCTION, $.uglify()
        //     .on("error", e => {
        //         console.log(e);
        //     })
        // ))
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

function partials() {
    return gulp.src("frontend/src/partials/**/*")
        .pipe(gulp.dest(PATHS.dist + "/partials"));
}

function config() {
    return gulp.src("frontend/src/assets/config/**/*")
        .pipe(gulp.dest(PATHS.dist + "/assets/config"));
}

function fonts() {
    return gulp.src(PATHS.fonts)
        .pipe(gulp.dest(PATHS.dist + "/assets/fonts"));
}


// proxy the backend server to support browser reloading
function server(done) {
    browser.init({
        proxy: "localhost:3000",
        port: PORT
    });
    done();
}

// function server(done) {
//     browser.init({
//         server: PATHS.dist, port: PORT
//     });
//     done();
// }

var go;

// launch the backend serving the web distribution directory
function backendCompile(done) {
    child.spawnSync('go', ['install'], {cwd: "backend", stdio: "inherit"});
    done();
}

// launch the backend serving the web distribution directory
function backend(done) {
    go = gulpgo.run("main.go", ["--webdir", "src/github.com/Microsoft/DUCK/" + PATHS.dist], {cwd: "backend", stdio: "inherit"});
    done();
}

function reloadBrowser(done) {
    browser.reload();
    done();
}

// Watch for changes to frontend assets and backend Go code
function watch() {
    gulp.watch(PATHS.assets, copy);
    gulp.watch("frontend/src/pages/**/*.html", gulp.series(pages, reloadBrowser));
    gulp.watch("frontend/src/{layouts,partials}/**/*.html", gulp.series(resetPages, pages, partials, reloadBrowser));
    gulp.watch("frontend/src/assets/scss/**/*.scss", sass);
    gulp.watch("frontend/src/assets/js/**/*.js", gulp.series(javascript, reloadBrowser));
    gulp.watch("frontend/src/assets/img/**/*", gulp.series(images, reloadBrowser));
    gulp.watch("frontend/src/assets/img/**/*", gulp.series(config, reloadBrowser));
    gulp.watch(["backend/**/*.go"]).on("change", function () {
        go.restart();
    });
}
