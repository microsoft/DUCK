// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
var editorModule = angular.module("duck.editor");

/**
 * Manages the current document being edited.
 */
editorModule.service("DocumentModel", function (CurrentUser, TaxonomyService, GlobalDictionary, DataUseDocumentService, $q, UUID, ObjectUtils) {
    /**
     * A local copy of the document.
     */
    this.document = null;

    /**
     * Tracks the local edit state of the document.
     */
    this.dirty = false;

    /**
     * The state of the current document: NOT_VALIDATED; NON_COMPLIANT; UNKNOWN; or COMPLIANT
     */
    this.state = "NOT_VALIDATED";

    this.currentStatement = null;

    this.state = "UNKNOWN";

    var context = this;

    /**
     * Retrieves the document from the backend and uses it to initialize the local model.
     *
     * @param documentId the document id
     * @return a promise that will be resolved after initialization has successfully completed or failed
     */
    this.initialize = function (documentId) {
        return $q(function (resolve) {
            DataUseDocumentService.getDocument(documentId).then(function (useDocument) {
                context.document = useDocument;
                context.dirty = false;

                // configure the taxonomy service with the global and document dictionaries as the document will be edited.
                TaxonomyService.activate([GlobalDictionary.getDictionary(), context.document.dictionary.values()]);

                context.document.statements.forEach(function (statement) {
                    context.addStatementErrorObject(statement);
                    context.lookupAndSetTerms(context.document);
                });

                // callback
                resolve();
            });
        });
    };

    this.release = function () {
        TaxonomyService.deactivateDictionaries();
        context.resetCompliance();
    };

    this.setDocumentLocale = function (locale) {
        if (context.document) {
            context.document.locale = locale;
            context.lookupAndSetTerms(context.document);
            context.markDirty();
        }

    };

    this.setAssumptionSet = function (assumptionSetId) {
        if (context.document) {
            context.document.assumptionSet = assumptionSetId;
            context.markDirty();
            context.resetCompliance();
        }
    };

    /**
     * Resets the statement field codes as when an ISO field value is edited.
     */
    this.reCalculateCodes = function () {
        context.document.statements.forEach(function (statement) {
            statement.useScopeCode = TaxonomyService.findCode("scope", statement.useScope, context.document.locale, statement.useScope);
            statement.qualifierCode = TaxonomyService.findCode("qualifier", statement.qualifier, context.document.locale, statement.qualifier);
            statement.dataCategoryCode = TaxonomyService.findCode("dataCategory", statement.dataCategory, context.document.locale, statement.dataCategory);
            statement.sourceScopeCode = TaxonomyService.findCode("scope", statement.sourceScope, context.document.locale, statement.sourceScope);
            statement.actionCode = TaxonomyService.findCode("action", statement.action, context.document.locale, statement.action);
            statement.resultScopeCode = TaxonomyService.findCode("scope", statement.resultScope, context.document.locale, statement.resultScope);
            // console.log(statement.useScopeCode + ", " + statement.qualifierCode + ", " + statement.dataCategoryCode + ", " + statement.sourceScopeCode
            //     + ", " + statement.actionCode + ", " + statement.resultScopeCode);
        });
    };

    /**
     * Sets the statement field terms based on their corresponding code.
     */
    this.lookupAndSetTerms = function (document) {
        document.statements.forEach(function (statement) {
            statement.useScope = TaxonomyService.findTerm("scope", statement.useScopeCode, document.locale, statement.useScopeCode);
            statement.qualifier = TaxonomyService.findTerm("qualifier", statement.qualifierCode, document.locale, statement.qualifierCode);
            statement.dataCategory = TaxonomyService.findTerm("dataCategory", statement.dataCategoryCode, document.locale, statement.dataCategoryCode);
            statement.sourceScope = TaxonomyService.findTerm("scope", statement.sourceScopeCode, document.locale, statement.sourceScopeCode);
            statement.action = TaxonomyService.findTerm("action", statement.actionCode, document.locale, statement.actionCode);
            statement.resultScope = TaxonomyService.findTerm("scope", statement.resultScopeCode, document.locale, statement.resultScopeCode);
            // console.log(statement.useScope + ", " + statement.qualifier + ", " + statement.dataCategory + ", " + statement.sourceScope
            //     + ", " + statement.action + ", " + statement.resultScope);
        });
    };

    /**
     * Deletes the statement in the local model (i.e. it is not synchronized to the backend).
     *
     * @param statement the statement
     */
    this.deleteStatement = function (statement) {
        context.document.statements.without(function (element) {
            return element === statement;
        });
        context.dirty = true;
        context.resetCompliance();
    };

    /**
     * Duplicates the statement in the local model (i.e. it is not synchronized to the backend).
     *
     * @param statement the statement
     */
    this.duplicateStatement = function (statement) {
        var pos = -1;
        var found = -1;
        context.document.statements.forEach(function (element) {
            pos++;

            if (element === statement) {
                found = pos;
            }
        });
        if (found >= 0) {
            var newStatement = {
                useScope: statement.useScope,
                qualifier: statement.qualifier,
                dataCategory: statement.dataCategory,
                sourceScope: statement.sourceScope,
                action: statement.action,
                resultScope: statement.resultScope,
                passive: statement.passive
            };
            context.addStatementErrorObject(newStatement);
            newStatement.trackingId = UUID.next();

            context.document.statements.splice(found, 0, newStatement);
            context.dirty = true;
            context.resetCompliance();
        }
    };

    /**
     * Sets the current statement for editing purposes.
     * @param statement the current statement
     */
    this.setCurrentStatement = function (statement) {
        this.currentStatement = statement;
    };

    /**
     * Clears the current statement.
     */
    this.clearCurrentStatement = function () {
        this.currentStatement = null;
    };

    this.getCurrentStatement = function () {
        return this.currentStatement;
    };

    /**
     * Adds the statement in the local model (i.e. it is not synchronized to the backend.
     *
     * @param statement the statement
     */
    this.addStatement = function (statement) {
        context.addStatementErrorObject(statement);
        context.document.statements.push(statement);
        statement.trackingId = UUID.next();
        context.dirty = true;
        context.resetCompliance();
    };

    this.addStatementErrorObject = function (statement) {
        statement.errors = {
            useScope: {active: false, level: null, action: false},
            qualifier: {active: false, level: null, action: false},
            dataCategory: {active: false, level: null, action: false},
            sourceScope: {active: false, level: null, action: false},
            action: {active: false, level: null, action: false},
            resultScope: {active: false, level: null, action: false}
        };

    };

    /**
     * Adds a new term to either the global or document dictionary.
     * @param type the ISO type
     * @param code the code
     * @param category category
     * @param value the term value
     * @param dictionaryType the type of dictionary, e.g. global or document
     */
    this.addTerm = function (type, code, category, value, dictionaryType) {
        if (dictionaryType === "document") {
            context.document.dictionary.put(value, {value: value, type: type, code: code, category: category, dictionaryType: "document"});
        } else {
            GlobalDictionary.addTerm(type, code, category, value);
        }
        TaxonomyService.addTerm(type, code, category, value, dictionaryType);
    };

    this.makePassive = function (statement) {
        if (ObjectUtils.notNull(statement)) {
            statement.passive = true;
        } else {
            context.document.statements.forEach(function (statement) {
                statement.passive = true
            });
        }
        context.markDirty();
    };

    this.makeActive = function (statement) {
        if (ObjectUtils.notNull(statement)) {
            statement.passive = false;
        } else {
            context.document.statements.forEach(function (statement) {
                statement.passive = false;
            });
        }
        context.markDirty();
    };

    this.complianceCheck = function () {
        return $q(function (resolve, reject) {
            if (!context.validateSyntax(true)) {
                resolve();
                return;
            }
            // FIXME ruleset id
            DataUseDocumentService.complianceCheck(context.document, "123").then(function (complianceResult) {
                context.state = complianceResult.compliant;
                context.document.statements.forEach(function (statement) {
                    var statementExplanation = complianceResult.map.get(statement.trackingId);
                    if (statementExplanation == null) {
                        return;
                    }
                    statement.$$statementExplanation = statementExplanation;    // note  $$ signals Angular to ignore this property during serialization

                });
                resolve();
            });
        });
    };

    /**
     * Saves the local model to the backend.
     */
    this.save = function () {
        var promise = DataUseDocumentService.saveDocument(context.document);
        context.dirty = false;
        return promise;
    };

    /**
     * Marks the local state as edited.
     */
    this.markDirty = function () {
        context.dirty = true;
    };

    this.validateSyntax = function (full) {
        var fullValidation = full === undefined ? false : full;
        var totalErrors = 0;
        context.document.statements.forEach(function (statement) {
            var errorNumber = 1;
            if (ObjectUtils.isNull(statement.errors) && !fullValidation) {
                return true;
            }

            // use scope
            if (!ObjectUtils.isEmptyString(statement.useScope)) {
                if (TaxonomyService.contains("scope", context.document.locale, statement.useScope)) {
                    context.resetValidation(statement.errors.useScope);
                } else {
                    context.createError(statement.errors.useScope, "Use scope is not recognized", errorNumber);
                    errorNumber++;
                }
            } else if (fullValidation) {
                context.createError(statement.errors.useScope, "Use scope is not specified", errorNumber);
                errorNumber++;
            } else {
                context.resetValidation(statement.errors.useScope);
            }

            // qualifier
            if (!ObjectUtils.isEmptyString(statement.qualifier)) {
                if (TaxonomyService.contains("qualifier", context.document.locale, statement.qualifier)) {
                    context.resetValidation(statement.errors.qualifier);
                } else {
                    context.createError(statement.errors.qualifier, "Qualifier is not recognized", errorNumber);
                    errorNumber++;
                }
            } else if (fullValidation) {
                context.createError(statement.errors.qualifier, "Qualifier is not specified", errorNumber);
                errorNumber++;
            } else {
                context.resetValidation(statement.errors.qualifier);
            }

            // data category
            if (!ObjectUtils.isEmptyString(statement.dataCategory)) {
                if (TaxonomyService.contains("dataCategory", context.document.locale, statement.dataCategory)) {
                    context.resetValidation(statement.errors.dataCategory);
                } else {
                    context.createError(statement.errors.dataCategory, "Data category is not recognized", errorNumber);
                    errorNumber++;
                }
            } else if (fullValidation) {
                context.createError(statement.errors.dataCategory, "Data category is not specified", errorNumber);
                errorNumber++;
            } else {
                context.resetValidation(statement.errors.dataCategory);
            }

            // source scope
            if (!ObjectUtils.isEmptyString(statement.sourceScope)) {
                if (TaxonomyService.contains("scope", context.document.locale, statement.sourceScope)) {
                    context.resetValidation(statement.errors.sourceScope);
                } else {
                    context.createError(statement.errors.sourceScope, "Source scope is not recognized", errorNumber);
                    errorNumber++;
                }
            } else if (fullValidation) {
                context.createError(statement.errors.sourceScope, "Source Scope is not specified", errorNumber);
                errorNumber++;
            } else {
                context.resetValidation(statement.errors.sourceScope);
            }

            // action
            if (!ObjectUtils.isEmptyString(statement.action)) {
                if (TaxonomyService.contains("action", context.document.locale, statement.action)) {
                    context.resetValidation(statement.errors.action);
                } else {
                    context.createError(statement.errors.action, "Action is not recognized", errorNumber);
                    errorNumber++;
                }
            } else if (fullValidation) {
                context.createError(statement.errors.action, "Action is not specified", errorNumber);
                errorNumber++;
            } else {
                context.resetValidation(statement.errors.action);
            }

            // result scope
            if (!ObjectUtils.isEmptyString(statement.resultScope)) {
                if (TaxonomyService.contains("scope", context.document.locale, statement.resultScope)) {
                    context.resetValidation(statement.errors.resultScope);
                } else {
                    context.createError(statement.errors.resultScope, "Result scope is not recognized", errorNumber);
                    errorNumber++;
                }
            } else if (fullValidation) {
                context.createError(statement.errors.resultScope, "Result Scope is not specified", errorNumber);
                errorNumber++;
            } else {
                context.resetValidation(statement.errors.resultScope);
            }
            totalErrors = totalErrors + errorNumber - 1;

        });
        return totalErrors === 0;
    };

    this.createError = function (errorObject, message, errorNumber) {
        errorObject.active = true;
        errorObject.level = "error";
        errorObject.errorNumber = errorNumber;
        errorObject.message = message;

    };

    this.resetValidation = function (errorObject) {
        errorObject.active = false;
        errorObject.level = null;
        errorObject.errorNumber = 0;
        errorObject.message = null;
    };

    this.emptyStatement = function (statement) {
        return (ObjectUtils.isNull(statement.useScope) || statement.useScope.trim().length === 0) && (ObjectUtils.isNull(statement.action) || statement.action.trim().length === 0)
    };

    this.resetCompliance = function () {
        context.state = "UNKNOWN";
        context.document.statements.forEach(function (statement) {
            delete statement.$$statementExplanation;
        });
    }
});
    