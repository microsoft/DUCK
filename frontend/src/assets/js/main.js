/**
 * This module handles the protected (signed in) areas of the application.
 *
 * @type {angular.Module}
 */
var mainModule = angular.module("duck.main", ["ui.router"]);

mainModule.controller("MainController", function ($scope) {

    /**
     * Logs the user out of the system by clearing the local storage token.
     */
    this.logout = function () {
        alert("Logout");
    }
});
