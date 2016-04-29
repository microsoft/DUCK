/**
 * This module defines the current user service.
 *
 * @type {angular.Module}
 */
var userModule = angular.module("duck.user", []);

userModule.service("CurrentUser", function ($log) {
    this.loggedIn = false;
    this.firstName = "anonymous";
    this.lastName = "anonymous";
    this.token = null;
    this.id = "";

    this.initialize = function () {

        if (localStorage.getItem("duck.token") === null) {
            $log.debug("Current user initialized");
            return;
        }
        this.firstName = localStorage.getItem("duck.firstName");
        this.lastName = localStorage.getItem("duck.lastName");
        this.id = localStorage.getItem("duck.id");
        this.token = localStorage.getItem("duck.token");
        this.loggedIn = true;
        $log.debug("Current user initialized");
    };

    this.initializeWith = function (data) {
        this.firstName = data.firstName;
        this.lastName = data.lastName;
        this.id = data.id;
        this.token = data.token;

        this.loggedIn = true;
        this.save();
    };

    this.save = function () {
        localStorage.setItem("duck.firstName", this.firstName);
        localStorage.setItem("duck.lastName", this.lastName);
        localStorage.setItem("duck.token", this.token);
        localStorage.setItem("duck.id", this.id);
    };

    this.reset = function () {
        localStorage.removeItem("duck.firstName");
        localStorage.removeItem("duck.lastName");
        localStorage.removeItem("duck.token");
        localStorage.removeItem("duck.id");

        this.loggedIn = false;
        this.firstName = "anonymous";
        this.lastName = "anonymous";
        this.token = null;
        this.id = "";
    }

});
