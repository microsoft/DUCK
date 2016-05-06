/**
 * This module handles user sign-in and registration.
 *
 * @type {angular.Module}
 */
var signinModule = angular.module("duck.signin");

signinModule.controller("SigninController", function ($state, SigninService) {
    this.username = "";
    this.password = "";

    this.signin = function () {
        if (this.username.length < 1 || this.password.length < 1) {
            // FIXME
            alert("Invalid username/password");
        }
        var promise = SigninService.signin(this.username, this.password);
        promise.then(function () {
                $state.go("main");
            }, function (code) {
                // FIXME
                alert('Login Failed: ' + code);
            }
        );


    }

});
