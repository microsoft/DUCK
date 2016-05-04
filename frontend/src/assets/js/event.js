'use strict';

/**
 * Event bus implementation that wraps Angular and potentially other implementations. Avoids proliferation of $scope references.
 */
var module = angular.module('duck.event', []);

module.service("EventBus", function ($rootScope) {

    /**
     * Publishes a message to a topic.
     * @param topic the topic
     * @param message the message
     */
    this.publish = function (topic, message) {
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
