var gatewayModule = angular.module("duck.gateway");

/**
 * Manages synchronization of user statement documents with the backend.
 */
gatewayModule.service('DataUseDocumentService', function (CurrentUser, UUID, $http, $q) {

    var context = this;

    // FIXME - remove when server is enabled
    this.summaries = new Hashtable();
    this.summaries.put("1", {name: "Customer Document v1", id: "1"});
    // this.summaries.put("2", {name: "Third-Party Document v2", id: "2"});
    // this.summaries.put("3", {name: "Partner Document", id: "3"});

    this.documents = new Hashtable();
    this.documents.put("1", {
        name: "Customer Document v1", id: "1",
        statements: [{
            useScope: "this product", qualifier: "account", dataCategory: "data", sourceScope: "those cloud services",
            action: "provide", resultScope: "cloud services defined in the service agreement", trackingId: UUID.next()
        }]
    });
    // this.documents.put("2", {name: "Third-Party Document v2", id: "2", statements: [{content: "Statement 1"}, {content: "Statement 2"}]});
    // this.documents.put("3", {name: "Partner Document", id: "3", statements: [{content: "Statement 1"}, {content: "Statement 2"}]});
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
            if (true) {
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
            if (true) {
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
