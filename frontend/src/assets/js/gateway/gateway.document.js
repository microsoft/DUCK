var gatewayModule = angular.module("duck.gateway");

/**
 * Manages synchronization of user statement documents with the backend.
 */
gatewayModule.service('DocumentService', function (CurrentUser, $http, $q) {

    var context = this;

    this.getAuthoredDocumentSummaries = function () {
        return $q(function (resolve, reject) {
            // make sure the user is signed in is set
            if (!CurrentUser.loggedIn) {
                $state.go('signin');
            }
            var url = "/v1/documents/" + CurrentUser.id + "/summary";

            // disable server call until implemented
            if (true) {
                resolve([{name: "Document 1"}, {name: "Document 2"}, {name: "Document 3"}]);
                return;
            }

            //noinspection JSUnusedLocalSymbols
            $http.get(url).success(function (data, status, headers, config) {
                var documents = angular.fromJson(data);
                resolve(documents);
            }).error(function (data, status, headers, config) {
                reject(status);
            });
        });

    }

});