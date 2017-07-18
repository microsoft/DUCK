// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

var coreModule = angular.module('duck.core');


/**
 * Manages notification messages. Messages are translated via the Angular Translate service.
 */
coreModule.service("NotificationService", function ($translate, $timeout) {
    this.message = null;

    var context = this;

    this.hasMessage = function () {
        return context.message !== null;
    };

    this.display = function (message, time) {
        $translate(message).then(function (message) {
            context.setMessageInternal(message, time);
        }, function (translationId) {
            context.setMessageInternal(translationId, time);
        });
    };

    this.getMessage = function () {
        return context.message;
    };

    this.clear = function () {
        context.message = null;
    };

    this.setMessageInternal = function (message, time) {
        context.message = message;
        if (time) {
            $timeout(function () {
                context.clear();
            }, time);
        }

    }
});