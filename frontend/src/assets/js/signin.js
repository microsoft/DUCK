// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

/**
 * This module handles user sign-in and registration.
 *
 * @type {angular.Module}
 */
var signinModule = angular.module("duck.signin");

signinModule.controller("SigninController", function ($state, SigninService, Validator, GlobalDictionary) {
    this.showRegister = false;
    this.username = "";
    this.password = "";
    this.locale = "en";
    this.firstname = "";
    this.lastname = "";

    this.Validator = Validator;

    this.showSigninForm = function (form) {
        form.$setPristine();
        form.$setUntouched();
        this.showRegister = false;
    };

    this.showRegisterForm = function (form) {
        form.$setPristine();
        form.$setUntouched();
        this.showRegister = true;
    };

    this.signin = function () {
        if (!(Validator.validateEmail(this.username) && Validator.validatePassword(this.password))) {
            return;
        }
        var promise = SigninService.signin(this.username, this.password);
        promise.then(function () {
                GlobalDictionary.initialize();
                $state.go("main.home");
            }, function (code) {
                // FIXME
                alert('Login Failed: ' + code);
            }
        );


    };

    this.register = function () {
        if (!(Validator.validateEmail(this.username)
            && Validator.validatePassword(this.password)
            && Validator.validateRequired(this.locale)
            && Validator.validateRequired(this.firstname)
            && Validator.validateRequired(this.lastname))) {
            return;
        }
        var promise = SigninService.register(
            {
                email: this.username,
                password: this.password,
                locale: this.locale,
                firstname: this.firstname,
                lastname: this.lastname
            }
        );
        promise.then(function () {
                $state.go("main.home");
            }, function (code) {
                // FIXME
                alert('Registration Failed: ' + code);
            }
        );


    };

    this.validate = function () {
        return Validator.validateEmail(this.username) && Validator.validatePassword(this.password);

    }


});
