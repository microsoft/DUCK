// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
var componentModule = angular.module('duck.component');

/**
 * Displays and error label if a statement fragment is in the error state. The error object for the statement fragment is passed and is of the form:
 *
 * {active: boolean, level: error|warning, errorNumber: number}
 */
componentModule.directive('errorLabel', function () {
    return {
        restrict: 'E',
        transclude: true,
        scope: {fragment: "="},
        template: "<div ng-show='fragment.active' class='badge' " +
        "ng-class='{\"alert\":fragment.level === \"error\", \"warning\":fragment.level === \"warning\"}' ng-transclude>{{fragment.errorNumber}}</div>",
        replace: true
    };
});

/**
 * Highlights the statement fragment if it is in error.
 */
componentModule.directive('errorMarker', function (ObjectUtils) {
    return {
        restrict: 'A',
        link: function (scope, element, attrs) {
            scope.$watchCollection(attrs.errorMarker, function (value) {
                if (ObjectUtils.isNull(value) || ObjectUtils.isNull(value.active) || ObjectUtils.isNull(value.level)) {
                    element.removeClass('error-highlight');
                    element.removeClass('warning-highlight');
                    return;
                }
                if (value.active && value.level === "error") {
                    element.addClass('error-highlight');
                } else if (value.active && value.level === "warning") {
                    element.addClass('warning-highlight');
                }
            });
        }
    };
});
