$(document).foundation();

var app = angular.module('application', []);

app.factory('AppInfo', function () {

    return {
        test: 'value'
    };
});


//app.config(['$urlProvider', '$locationProvider', function ($urlProvider, $locationProvider) {
//
//    $urlProvider.otherwise('/');
//
//    $locationProvider.html5Mode({
//        enabled: false,
//        requireBase: false
//    });
//
//    $locationProvider.hashPrefix('!');
//}]);

app.controller('AppController', function ($scope) {
    $scope.test = "This is a test";
});
