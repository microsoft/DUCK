// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
'use strict';

/**
 * Event bus implementation that wraps Angular and potentially other implementations. Avoids proliferation of $scope references.
 */
var coreModule = angular.module('duck.core');

coreModule.service("EventBus", function ($rootScope) {
    var context = this;
    
    /**
     * Publishes a message to a topic.
     * @param topic the topic
     * @param message the message
     */
    this.publish = function (topic, message) {
        if (context.startsWith(topic, "ui.")) {  // modal event - open the Foundation modal
            $("#" + topic.substring(3)).foundation("open");
            return;
        }
        $rootScope.$broadcast(topic, message)
    };

    /**
     * Subscribes the function to a topic and returns a function that can be called to cancel the subscription.
     * @param topic the topic
     * @param subscriber the subscriber function
     * @returns {*} a function that can be invoked to cancel the subscription
     */
    this.subscribe = function (topic, subscriber) {
        // wrap the sender in a function so only the message is sent
        return $rootScope.$on(topic, function (sender, message) {
            subscriber(message);
        });
    };

    this.startsWith = function (value, expr) {
        return value.lastIndexOf(expr, 0) === 0;
    }

});
