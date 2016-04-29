/**
 * This module handles the protected (signed in) areas of the application.
 *
 * @type {angular.Module}
 */
var mainModule = angular.module("duck.main", ["ui.router"]);

mainModule.controller("MainController", function ($scope) {

    // signal that the main controller has been loaded and Foundation should be initialized
    if (!$scope.initFoundation) {
        $scope.initFoundation = true;
    }

});

mainModule.controller("SignoutController", function ($state) {

    /**
     * Logs the user out of the system by clearing the local storage token.
     */
    this.signout = function () {
        $state.go("signin");
    }

});


