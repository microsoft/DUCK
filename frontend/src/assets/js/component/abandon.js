// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

/**
 * Defines UI components.
 */
var componentModule = angular.module("duck.component");

/**
 * Manages the abandon changes modal, e.g. opened when a user navigates from a page with edits that have not been saved.
 */
componentModule.service("AbandonComponent", function ($state, $rootScope) {

    var context = this;

    /**
     * Opens the modal.
     * @param event the state transition event.
     * @param toState the state to transition to
     */
    this.open = function (event, toState) {
        context.event = event;
        context.toState = toState;
        event.preventDefault();
        $("#abandonChanges").foundation("open");
    };

    /**
     * Triggered when the user confirms an action in the modal.
     */
    this.confirm = function () {
        $state.go(context.toState.name, null, {notify: false}).then(function () {
            // workaround for angular router issue: https://github.com/angular-ui/ui-router/issues/178
            $rootScope.$broadcast("$stateChangeSuccess", context.toState);
        });

    }

});


componentModule.controller("AbandonController", function (AbandonComponent) {

    this.confirm = function () {
        AbandonComponent.confirm();
    }

});
