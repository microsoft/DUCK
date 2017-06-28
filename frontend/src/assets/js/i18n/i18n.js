// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

/**
 * Internationalization module.
 *
 * @type {angular.Module}
 */
var i18NModule = angular.module('duck.i18n');

/**
 * Constant that defines the default frontend language. This may be modified.
 */
i18NModule.constant("DEFAULT_LOCALE", "en");

i18NModule.service("LocaleService", function (DEFAULT_LOCALE) {
    var context = this;

    this.defaultLocale = DEFAULT_LOCALE;

    this.getLocales = function () {
        return [{description: "English (US)", id: "en"}, {description: "German", id: "de"}, {description: "Korean", id: "ko"}];
    }

});