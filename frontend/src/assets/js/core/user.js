// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

/**
 * This module defines the current user service.
 */
var coreModule = angular.module("duck.core");

coreModule.service("CurrentUser", function ($log, LocaleService, AssumptionSetService, ObjectUtils) {
    this.loggedIn = false;
    this.firstName = "anonymous";
    this.lastName = "anonymous";
    this.token = null;  // the JSON Web Token provided by the server
    this.id = "";

    /**
     * Initializes the current user from local storage if present; otherwise initializes an anonymous user.
     */
    this.initialize = function () {
        if (localStorage.getItem("duck.token") === null) {
            $log.debug("Current user initialized as anonymous");
            return;
        }
        this.firstName = localStorage.getItem("duck.firstName");
        this.lastName = localStorage.getItem("duck.lastName");
        this.id = localStorage.getItem("duck.id");
        this.token = localStorage.getItem("duck.token");
        var locale = this.locale = localStorage.getItem("duck.locale");
        if (ObjectUtils.isNull(locale)) {
            this.locale = LocaleService.defaultLocale;
        } else {
            this.locale = locale;
        }

        var assumptionSet = this.assumptionSet = localStorage.getItem("duck.assumptionSet");
        if (ObjectUtils.isNull(assumptionSet)) {
            this.assumptionSet = AssumptionSetService.getAssumptionSets()[0].id;
        } else {
            this.locale = assumptionSet;
        }

        this.loggedIn = true;
        $log.debug("Current user initialized");
    };

    /**
     * Initializes the current user from the given data.
     * @param data the data
     */
    this.initializeWith = function (data) {
        this.firstName = data.firstName;
        this.lastName = data.lastName;
        this.id = data.id;
        this.token = data.token;

        if (data.locale) {
            this.locale = data.locale;
        } else {
            this.locale = LocaleService.defaultLocale;
        }

        if (data.assumptionSet) {
            this.assumptionSet = data.assumptionSet;
        } else {
            this.assumptionSet = AssumptionSetService.getAssumptionSets()[0].id;
        }


        this.loggedIn = true;
        this.save();
    };

    /**
     * Saves the current user to local storage.
     */
    this.save = function () {
        localStorage.setItem("duck.firstName", this.firstName);
        localStorage.setItem("duck.lastName", this.lastName);
        localStorage.setItem("duck.token", this.token);
        localStorage.setItem("duck.id", this.id);
        localStorage.setItem("duck.locale", this.locale);
    };

    /**
     * Resets the current user to anonymous.
     */
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
        $log.debug("Current user signed out and set to anonymous");
    }

});
