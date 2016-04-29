/**
 * This module bootstraps the application, including defining URLs and pages.
 *
 * @type {angular.Module}
 */

var app = angular.module("application", ["duck.main", "duck.signin", "ui.router"]);

app.factory("AppInfo", function () {

    return {
        name: "DUCK Application",
        version: "1.0.0"
    };
});


app.config(["$urlRouterProvider", "$locationProvider", "$stateProvider", function ($urlRouterProvider, $locationProvider, $stateProvider) {

    // setup URL structure
    $urlRouterProvider.otherwise("/");

    $locationProvider.html5Mode({
        enabled: true,
        requireBase: false
    });

    $locationProvider.hashPrefix("!");

    // define the application states
    $stateProvider
        .state("main", {  // the top-level state for protected (signed in) areas of the application
            url: "/",
            templateUrl: "../../main.html"
        })

        .state("signin", {   // signin and registration process
            url: "/signin",
            templateUrl: "../../signin.html"
        })
}]);

app.controller("AppController", function ($scope) {
});

app.run(function ($rootScope) {
    // load Foundation
    $rootScope.$on('$viewContentLoaded', function () {
        $(document).foundation();
    });
});
