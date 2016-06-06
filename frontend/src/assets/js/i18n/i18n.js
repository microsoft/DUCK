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
        return ["English (US)", "German"];
    }

});