// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

var gatewayModule = angular.module("duck.gateway");

/**
 * Manages synchronization of user statement documents with the backend.
 */
gatewayModule.service('DataUseDocumentService', function (CurrentUser, NotificationService, UUID, $http, $q, ObjectUtils) {

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
                // replace dictionary with hashtable
                var dictionaryObj = document.dictionary;
                document.dictionary = new Hashtable();
                angular.forEach(dictionaryObj, function (value, key) {
                    document.dictionary.put(key, value);
                });
                resolve(document);
            }).error(function (data, status, headers, config) {
                reject(status);
            });
        });
    };

    this.getTestData = function () {
        return $q(function (resolve, reject) {
            // make sure the user is signed in is set
            if (!CurrentUser.loggedIn) {
                $state.go('signin');
            }
            var url = "assets/config/test-data.json";

            $http.get(url).success(function (data) {
                var testData = angular.fromJson(data);
                resolve(testData);
            }).error(function (data, status) {
                reject(status);
            });
        });
    };

    this.createDocument = function (name) {
        return $q(function (resolve, reject) {
            var url = "/v1/documents";
            var data = {};
            data.locale = CurrentUser.locale;
            data.assumptionSet = CurrentUser.assumptionSet;
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

    this.copyDocument = function (name, documentId) {
        return $q(function (resolve, reject) {
            var url = "/v1/documents/copy/" + documentId;
            var data = {};
            data.locale = CurrentUser.locale;
            data.name = name;
            data.owner = CurrentUser.id;
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
        NotificationService.display("document_saving");
        return $q(function (resolve, reject) {
            var url = "/v1/documents";
            var documentData = context.createDocumentData(document);
            //noinspection JSUnusedLocalSymbols
            $http.put(url, documentData).success(function (data, status, headers, config) {
                var newDocument = angular.fromJson(data);
                // update the document revision
                document.revision = newDocument.revision;
                NotificationService.display("document_saved", 1000);
                resolve(document);
            }).error(function (data, status, headers, config) {
                NotificationService.clear();
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
     * Performs a document compliance check, returning a compliance result. Note the compliance result has an additional 'map' property of statements
     * explanations keyed by statement tracking id.
     *
     * @param document the document to check
     * @param rulebaseId the rule base to use
     * @return {*} the result
     */
    this.complianceCheck = function (document, rulebaseId) {
        NotificationService.display("document_validating");
        return $q(function (resolve, reject) {
            var url = "/v1/rulebases/" + rulebaseId + "/documents";
            var documentData = context.createDocumentData(document);

            // compliant values: NON_COMPLIANT; UNKNOWN; or COMPLIANT

            // stub for testing

            if (false) {
                if (document.statements.length != 4) {
                    resolve({ compliant: "COMPLIANT" });
                    return;
                }
                var complianceResult = {
                    compliant: "COMPLIANT",
                    explanation: {}
                };

                complianceResult.explanation[document.statements[0].trackingId] = {
                    consentRequired: { value: true, assumed: false },
                    pii: { value: true, assumed: false },
                    li: { value: true, assumed: false },
                    compatiblePurpose: [document.statements[1].trackingId]
                };
                complianceResult.explanation[document.statements[1].trackingId] = {
                    consentRequired: { value: true, assumed: false },
                    pii: { value: true, assumed: true },
                    li: { value: true, assumed: true },
                    compatiblePurpose: [document.statements[0].trackingId]
                };
                complianceResult.explanation[document.statements[2].trackingId] = {
                    consentRequired: { value: true, assumed: false },
                    pii: { value: false, assumed: false },
                    li: { value: false, assumed: false },
                    compatiblePurpose: []
                };
                complianceResult.explanation[document.statements[3].trackingId] = {
                    consentRequired: { value: false, assumed: false },
                    pii: { value: false, assumed: false },
                    li: { value: true, assumed: false },
                    compatiblePurpose: []
                };
                NotificationService.clear();
                resolve(context.mapExplanation(complianceResult));
                return;
            }
            // end testing stub

            $http.put(url, documentData).success(function (data) {
                var complianceResult = angular.fromJson(data);
                resolve(context.mapExplanation(complianceResult));
                // FIXME handle errors
            }).error(function (data, status) {
                reject(status);
            }).finally(function () {
                NotificationService.clear();
            });
        });
    };

    this.mapExplanation = function (complianceResult) {
        var map = new Hashtable();
        angular.forEach(complianceResult.explanation, function (statementExplanation, trackingId) {
            map.put(trackingId, statementExplanation);
        });
        complianceResult.map = map;
        return complianceResult;
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
            var url = "/v1/rulebases/" + rulebaseId + "/documents";
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
            $http.put(url, documentData).success(function (data) {
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
        data.assumptionSet = document.assumptionSet;
        data.name = document.name;
        data.owner = document.owner;
        data.statements = [];
        context.copyStatements(document, data);
        data.revision = document.revision;
        data.dictionary = {};
        var entries = document.dictionary.entries();
        entries.forEach(function (entry) {
            var term = entry[1];
            data.dictionary[entry[0]] = { value: term.value, type: term.type, code: term.code, category: term.category, dictionaryType: "document" };
        });
        return data;
    };

    this.copyStatements = function (document, data) {
        document.statements.forEach(function (statement) {
            var newStatement = {
                trackingId: statement.trackingId,
                actionCode: statement.actionCode,
                dataCategoryCode: statement.dataCategoryCode,
                dataCategories: [],
                qualifierCode: statement.qualifierCode,
                resultScopeCode: statement.resultScopeCode,
                sourceScopeCode: statement.sourceScopeCode,
                useScopeCode: statement.useScopeCode,
                passive: statement.passive

            };
            if (ObjectUtils.notNull(statement.tag)) {
                newStatement.tag = statement.tag;
            }
            statement.dataCategories.forEach(function (dataCategory) {
                var newCategory = {
                    dataCategoryCode: dataCategory.dataCategoryCode,
                    qualifierCode: dataCategory.qualifierCode,
                    operator: dataCategory.operator
                }
                newStatement.dataCategories.push(newCategory);
            });

            data.statements.push(newStatement)


        });
    }

});
