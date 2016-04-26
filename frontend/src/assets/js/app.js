$(document).foundation();

var booter = angular.module('booter', []);

booter.factory('AppInfo', function () {

    return {
        test: 'value'
    };
});