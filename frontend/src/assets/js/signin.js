/**
 * This module handles user sign-in and registration.
 *
 * @type {angular.Module}
 */
var signinModule = angular.module("duck.signin", ["ui.router"]);

signinModule.controller("SigninController", function ($state) {

    this.signin = function () {
        $state.go("main");
    }

});
