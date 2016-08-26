/**
 * Assumption module.
 *
 * @type {angular.Module}
 */
var i18NModule = angular.module('duck.assumption');

i18NModule.service("AssumptionSetService", function ($http, $q) {
    this.assumptionSets = [{id: "1", name: "Assumption Set 1", description: "The Default A Priori Assumption Set"}];
    var context = this;


    this.getAssumptionSets = function () {
        return context.assumptionSets;
    };

    this.initialize = function () {
        return $q(function (resolve, reject) {
            resolve();  // remove this line when assumption-sets handler is implemented in backend

            // uncomment when assumption-sets handler is implemented in backend:

            // $http.get('assumption-sets').success(function (data) {
            //     context.assumptionSets = data;
            //     resolve();
            // }).error(function (data, status) {
            //     reject(status);
            // });

            // end assumption-sets handler comment
        });
    }
});