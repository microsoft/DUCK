/**
 * Base services.
 */
var coreModule = angular.module('duck.core');

coreModule.service("ObjectUtils", function () {
    this.isNull = function (obj) {
        return obj === undefined || obj === null;
    };

    this.notNull = function (obj) {
        return !this.isNull(obj);
    };

    this.isEmpty = function (array) {
        return this.isNull(array) || array.length < 1;
    };

    this.isEmptyString = function (value) {
        return this.isNull(value) || value.trim().length === 0;
    };

    /**
     * Safely evaluates the expression on the object; if null or undefined, returns the default value.
     * @param root the object to evaluate on
     * @param expression the expression
     * @param defaultValue the default value
     * @returns {Object} the result
     */
    this.safeGet = function (root, expression, defaultValue) {
        if (this.isNull(root)) {
            return defaultValue;
        }
        var tokens = expression.split('.');
        var o = root;
        for (var i = 0; i < tokens.length; i++) {
            if (this.notNull(o[tokens[i]])) {
                o = o[tokens[i]];
            } else {
                return defaultValue;
            }
        }
        return o;
    };

    this.getOrDefault = function (value, defaultValue) {
        return this.notNull(value) ? value : defaultValue;
    };


});

/**
 * Capitalizes the initial character of an expression.
 */
coreModule.filter("capitalizeInitial", function () {
    return function (expression) {
        return expression.charAt(0).toUpperCase() + expression.slice(1);
    }
});

/**
 * Removes elements matched by the provided function. Unlike Array.filter(), this function modifies the original array rather than returning a new one.
 *
 * @param func the matching function
 */
Array.prototype.without = function (func) {
    for (var i = this.length - 1; i >= 0; i--) {
        if (func(this[i])) {
            this.splice(i, 1);
        }
    }
};

