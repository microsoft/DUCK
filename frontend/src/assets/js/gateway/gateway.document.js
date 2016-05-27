var gatewayModule = angular.module("duck.gateway");

/**
 * Manages synchronization of user statement documents with the backend.
 */
gatewayModule.service('DataUseDocumentService', function (CurrentUser, UUID, $http, $q) {

    var context = this;
    context.runServer = false;

    // FIXME - remove when server is enabled
    if (!context.runServer) {
        context.summaries = new Hashtable();
        context.summaries.put("1", {name: "Customer Document v1", id: "1"});

        context.documents = new Hashtable();
        context.documents.put("1", {
            name: "Customer Document v1", id: "1",
            statements: [{
                useScope: "this product", qualifier: "account", dataCategory: "data", sourceScope: "those cloud services",
                action: "provide", resultScope: "cloud services defined in the service agreement", trackingId: UUID.next(),
                passive:false
            }]
        });
    }
    // FIXME end remove

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

            // FIXME disable server call until implemented
            if (!context.runServer) {
                resolve(context.summaries.values());
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
            if (!context.runServer) {
                resolve(context.documents.get(documentId));
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
