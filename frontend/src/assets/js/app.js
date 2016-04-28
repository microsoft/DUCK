$(document).foundation();

var app = angular.module('application', ['ui.router']);

app.factory('AppInfo', function () {

    return {
        test: 'value'
    };
});


app.config(['$urlRouterProvider', '$locationProvider', function ($urlRouterProvider, $locationProvider) {

    $urlRouterProvider.otherwise('/');

    $locationProvider.html5Mode({
        enabled: false,
        requireBase: false
    });

    $locationProvider.hashPrefix('!');
}]);

app.controller('AppController', function ($scope) {
    $scope.test = "This is a test";
});
