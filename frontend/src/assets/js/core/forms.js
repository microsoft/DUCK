/**
 * Form services.
 */
var coreModule = angular.module('duck.core');

coreModule.service("Validator", function () {
    var context = this;

    this.validateEmail = function (email) {
        var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
        return re.test(email);
    };

    this.validateRequired = function (value) {
        return value.trim().length > 0;
    };
    
    this.validatePassword = function (password) {
        return password.length >= 4;
    };

    this.showRequiredError = function (value, field, form) {
        return (form.$submitted || field.$touched) && !context.validateRequired(value)
    };

    this.showEmailError = function (email, field, form) {
        return (form.$submitted || field.$touched) && !context.validateEmail(email)
    };

    this.showPasswordError = function (password, field, form) {
        return (form.$submitted || field.$touched) && !context.validatePassword(password)
    }


});


