var i18NModule = angular.module('duck.i18n');

i18NModule.service("LocaleService", function () {
    var context = this;

    this.getLocales = function () {
        return ["English (US)", "Italian", "German", "Chineese"];
    }

});