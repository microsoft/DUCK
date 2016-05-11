var gatewayModule = angular.module("duck.gateway");

/**
 * Manages synchronization of user statement documents with the backend.
 */
gatewayModule.service('DataUseDocumentService', function (CurrentUser, $http, $q) {

    var context = this;

    /**
     * Retrieves summaries for data use statement documents authored by the current user.
     *
     * @return an array of summaries
     */
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

    };

    /**
     * Retrieves a data use statement document authored by the current user.
     *
     * @param documentId the document id
     * @return the document
     */
    this.getDocument = function (documentId) {
        return $q(function (resolve, reject) {
            // make sure the user is signed in is set
            if (!CurrentUser.loggedIn) {
                $state.go('signin');
            }
            var url = "/v1/documents/" + CurrentUser.id + "/" + documentId;

            // disable server call until implemented
            if (true) {
                resolve({name: "Document 1", statements: [{order: 1, content: "Statement 1"}, {order: 2, content: "Statement 2"}]});
                return;
            }

            //noinspection JSUnusedLocalSymbols
            $http.get(url).success(function (data, status, headers, config) {
                var document = angular.fromJson(data);
                resolve(document);
            }).error(function (data, status, headers, config) {
                reject(status);
            });
        });
    };
});
