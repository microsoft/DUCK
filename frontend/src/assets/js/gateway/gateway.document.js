var gatewayModule = angular.module("duck.gateway");

/**
 * Manages synchronization of user statement documents with the backend.
 */
gatewayModule.service('DataUseDocumentService', function (CurrentUser, UUID, $http, $q) {

    var context = this;
    context.runServer = true;

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
                document.revision = newDocument.revison;
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

    /**
     * Performs a compliance check on a document and returns possible alternatives.
     *
     * @param document the document
     * @param rulebaseId the rulebase id
     * @return the request promise
     */
    this.complianceCheckWithAlternatives = function (document, rulebaseId) {
        return $q(function (resolve, reject) {
            var url = "/v1/rulebases/" + rulebaseId+"/documents";
            var complianceResult;
            // stub for testing
            // compliant values: NON_COMPLIANT; UNKNOWN; or COMPLIANT
/*
            if (document.statements.length <= 2) {
                complianceResult = {
                    compliant: "COMPLIANT",
                    documents: []
                };
                resolve(complianceResult);
                return;
            }

            complianceResult = {
                compliant: "NON_COMPLIANT",
                documents: [{
                    id: document.id,
                    locale: document.locale,
                    name: document.name,
                    owner: document.owner,
                    statements: []
                }]
            };

            // remove the second and the last statements for testing
            for (var i = 0; i < document.statements.length; i++) {
                if (i == 1 || i == document.statements.length -1) {
                    continue;
                }
                var statement = document.statements[i];
                complianceResult.documents[0].statements.push({
                    trackingId: statement.trackingId,
                    actionCode: statement.actionCode,
                    dataCategoryCode: statement.dataCategoryCode,
                    qualifierCode: statement.qualifierCode,
                    resultScopeCode: statement.resultScopeCode,
                    sourceScopeCode: statement.sourceScopeCode,
                    useScopeCode: statement.useScopeCode,
                    passive: statement.passive
                });
            }

            resolve(complianceResult);
            // end stub */

             var documentData = context.createDocumentData(document);
             $http.put(url, documentData).success(function () {
                 var complianceResult = angular.fromJson(data);
                 resolve(complianceResult);
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
                actionCode: statement.actionCode,
                dataCategoryCode: statement.dataCategoryCode,
                qualifierCode: statement.qualifierCode,
                resultScopeCode: statement.resultScopeCode,
                sourceScopeCode: statement.sourceScopeCode,
                useScopeCode: statement.useScopeCode,
                passive: statement.passive
            })
        });
        data.revision = document.revision;
        return data;
    }

});
