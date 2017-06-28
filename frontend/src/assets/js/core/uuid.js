// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

var coreModule = angular.module("duck.core");

/**
 * Creates UUIDs that are as unique as Math.random() provides truly random number generation.
 */
coreModule.service("UUID", function () {

    this.next = function () {
        return 'axxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            var r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }
});