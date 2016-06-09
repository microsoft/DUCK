var gatewayModule = angular.module("duck.gateway");

/**
 * Manages synchronization of user statement documents with the backend.
 */
gatewayModule.service('DataUseDocumentService', function (CurrentUser, UUID, $http, $q) {

    var context = this;
    context.runServer = true;

    // local testing 
    // if (!context.runServer) {
    //     context.summaries = new Hashtable();
    //     context.summaries.put("1", {name: "Customer Document v1", id: "1"});
    //
    //     context.documents = new Hashtable();
    //     context.documents.put("1", {
    //         name: "Customer Document v1", id: "1",
    //         statements: [{
    //             useScope: "the CSP Services", qualifier: "identified", dataCategory: "credentials", sourceScope: "this capability",
    //             action: "provide", resultScope: "cloud services defined in the service agreement", trackingId: UUID.next(),
    //             passive: false
    //         }]
    //     });
    // }


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

            // local testing
            // if (!context.runServer) {
            //     resolve(context.summaries.values());
            //     return;
            // }

            //noinspection JSUnusedLocalSymbols
            $http.get(url).success(function (data, status, headers, config) {
                var documents = angular.fromJson(data);
                resolve(documents);
            }).error(function (data, status, headers, config) {
                if (404 === status) {
                    // no doc summaries available, resolve an empty array
                    resolve([]);
                    return;
                }
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
            var url = "/v1/documents/" + documentId;

            // local testing
            // if (!context.runServer) {
            //     resolve(context.documents.get(documentId));
            //     return;
            // }

            //noinspection JSUnusedLocalSymbols
            $http.get(url).success(function (data, status, headers, config) {
                var document = angular.fromJson(data);
                resolve(document);
            }).error(function (data, status, headers, config) {
                reject(status);
            });
        });
    };

    this.createDocument = function (name) {
        return $q(function (resolve, reject) {
            var url = "/v1/documents";
            var data = {};
            data.locale = CurrentUser.locale;
            data.name = name;
            data.owner = CurrentUser.id;
            data.statements = [];
            $http.post(url, data).success(function (data) {
                var newDocument = angular.fromJson(data);
                resolve(newDocument);
            }).error(function (data, status) {
                reject(status);
            });
        });
    };

    /**
     * Retrieves a data use statement document authored by the current user.
     *
     * @param document the document id
     * @return the request promise
     */
    this.saveDocument = function (document) {
        return $q(function (resolve, reject) {
            var url = "/v1/documents";
            var documentData = context.createDocumentData(document);
            //noinspection JSUnusedLocalSymbols
            $http.put(url, documentData).success(function (data, status, headers, config) {
                var newDocument = angular.fromJson(data);
                // update the document revision
                document._rev = newDocument._rev;
                resolve(document);
            }).error(function (data, status, headers, config) {
                reject(status);
            });
        });
    };

    /**
     * Deletes a document.
     *
     * @param id the document id
     * @return the request promise
     */
    this.deleteDocument = function (id) {
        return $q(function (resolve, reject) {
            var url = "/v1/documents/" + id;
            $http.delete(url).success(function () {
                resolve();
                // FIXME handle errors
            }).error(function (data, status) {
                reject(status);
            });
        });
    };


    this.createDocumentData = function (document) {
        var data = {};
        data.id = document.id;
        data.locale = document.locale;
        data.name = document.name;
        data.owner = document.owner;
        data.statements = [];
        document.statements.forEach(function (statement) {
            data.statements.push({
                trackingId: statement.trackingId,
                action: statement.action,
                dataCategory: statement.dataCategory,
                qualifier: statement.qualifier,
                resultScope: statement.resultScope,
                sourceScope: statement.sourceScope,
                useScope: statement.useScope,
                passive: statement.passive
            })
        });
        data._rev = document._rev;
        return data;
    }

});
