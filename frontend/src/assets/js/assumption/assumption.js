/**
 * Assumption module.
 *
 * @type {angular.Module}
 */
var i18NModule = angular.module('duck.assumption');

i18NModule.service("AssumptionSetService", function () {
    var context = this;


    this.getAssumptionSets = function () {
        return [{id:"1", name:"Assumption Set 1", description:"The Default A Priori Assumption Set"}];
    }

});