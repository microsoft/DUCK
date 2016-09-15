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
     * The original document being edited
     */
    this.originalDocument = null;

    this.alternativeVersions = [];

    /**
     * Tracks the local edit state of the document.
     */
    this.dirty = false;

    /**
     * The state of the current document: NOT_VALIDATED; NON_COMPLIANT; UNKNOWN; or COMPLIANT
     */
    this.state = "NOT_VALIDATED";

    this.currentStatement = null;

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
                context.originalDocument = useDocument;
                context.document = useDocument;
                context.document.dictionary = new Hashtable();
                // Create a fake term dictionary for testing
                // context.document.dictionary.put("Foo Service", {
                //     value: "Foo Service",
                //     type: "scope",
                //     code: "foo-service",
                //     category: "2",
                //     dictionaryType: "document"
                // });

                context.dirty = false;

                // configure the taxonomy service with the global and document dictionaries as the document will be edited.
                TaxonomyService.activate([GlobalDictionary.getDictionary(), context.document.dictionary.values()]);

                context.document.statements.forEach(function (statement) {
                    statement.errors = {
                        useScope: {active: false, level: null, action: false},
                        action: {active: false, level: null, action: false}
                    };
                    context.lookupAndSetTerms(context.document);
                });
                context.alternativeVersions.length = 0; // clear the array
                // context.alternativeVersions.push({id: context.document.id, name: context.document.name, locale: context.document.locale, statements: []});

                // callback
                resolve();
            });
        });
    };

    this.release = function () {
        TaxonomyService.deactivateDictionaries();
    };

    /**
     * Returns true if the current selected document can be edited. Only the original document can be edited, so if an alternative version is selected, this
     * method will return false.
     *
     * @return {boolean}
     */
    this.isEditable = function () {
        return context.document === context.originalDocument;
    };

    /**
     * Selects the original document.
     */
    this.selectOriginal = function () {
        context.document = context.originalDocument;
    };

    /**
     * Selects an alternative version of the document.
     *
     * @param document the document to select
     */
    this.selectAlternateVersion = function (document) {
        context.document = document;
        context.clearCurrentStatement();
    };

    this.setDocumentLocale = function (locale) {
        if (context.document) {
            context.document.locale = locale;
            context.lookupAndSetTerms(context.document);
            context.alternativeVersions.forEach(function (alternative) {
                alternative.locale = locale;
                context.lookupAndSetTerms(alternative);
            });
            context.markDirty();
        }

    };

    this.setAssumptionSet = function (assumptionSetId) {
        if (context.document) {
            context.document.assumptionSet = assumptionSetId;
            context.markDirty();
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

    this.toggleEdit = function (statement) {
        statement.$_edit = !statement.$_edit;
    };

    this.edit = function (statement) {
        statement.$_edit = true;
    };

    this.close = function (statement) {
        statement.$_edit = false;
    };

    this.editing = function (statement) {
        return statement.$_edit;
    };

    /**
     * Adds the statement in the local model (i.e. it is not synchronized to the backend.
     *
     * @param statement the statement
     */
    this.addStatement = function (statement) {
        context.document.statements.push(statement);
        statement.trackingId = UUID.next();
        context.dirty = true;
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

    this.complianceCheckWithAlternatives = function () {
        context.selectOriginal(); // make sure the original is selected and clear existing alternatives
        context.alternativeVersions.length = 0;
        // FIXME ruleset id
        DataUseDocumentService.complianceCheckWithAlternatives(context.document, "123").then(function (complianceResult) {
            context.state = complianceResult.compliant;
            context.alternativeVersions.length = 0; // clear array
            complianceResult.documents.forEach(function (alternative) {

                // decorate diff information
                for (var i = 0; i < context.originalDocument.statements.length; i++) {
                    var originalStatement = context.originalDocument.statements[i];
                    var diffStatement;
                    if (i >= alternative.statements.length) {
                        // statement was deleted at the end of the array
                        diffStatement = context.createDiffStatement(originalStatement);
                        alternative.statements.push(diffStatement);

                    } else if (originalStatement.trackingId !== alternative.statements[i].trackingId) {
                        diffStatement = context.createDiffStatement(originalStatement);
                        alternative.statements.splice(i, 0, diffStatement);
                    }
                }
                context.lookupAndSetTerms(alternative);
                context.alternativeVersions.push(alternative);

            });
        });
    };

    this.createDiffStatement = function (originalStatement) {
        return {
            _diff: true,
            trackingId: originalStatement.trackingId,
            actionCode: originalStatement.actionCode,
            dataCategoryCode: originalStatement.dataCategoryCode,
            qualifierCode: originalStatement.qualifierCode,
            resultScopeCode: originalStatement.resultScopeCode,
            sourceScopeCode: originalStatement.sourceScopeCode,
            useScopeCode: originalStatement.useScopeCode,
            passive: originalStatement.passive
        };
    };

    /**
     * Saves the local model to the backend.
     */
    this.save = function () {
        DataUseDocumentService.saveDocument(context.document);
        context.originalDocument = context.document;
        context.dirty = false;
    };

    /**
     * Reverts local changes to the document.
     */
    this.revert = function () {
        if (context.document === null) {
            return;
        }
        context.initialize(context.document.id);
    };

    /**
     * Marks the local state as edited.
     */
    this.markDirty = function () {
        context.dirty = true;
    };

    this.adoptAlternativeVersion = function () {
        context.clearCurrentStatement();
        context.document.statements.without(function (statement) {
            return statement._diff;     // clear out diff statements
        });
        context.originalDocument = context.document;
        context.markDirty();
        context.state = "COMPLIANT";
        context.alternativeVersions.length = 0; // clear alternatives
    };

    this.validateSyntax = function () {
        var errorNumber = 1;
        context.document.statements.forEach(function (statement) {
            if (ObjectUtils.isNull(statement.errors)) {
                return;
            }
            if (!ObjectUtils.isEmptyString(statement.useScope)) {
                if (TaxonomyService.contains("scope", context.document.locale, statement.useScope)) {
                    context.resetValidation(statement.errors.useScope);
                } else {
                    statement.errors.useScope.active = true;
                    statement.errors.useScope.level = "error";
                    statement.errors.useScope.errorNumber = errorNumber;
                    statement.errors.useScope.message = "Use scope is not recognized";
                    errorNumber++;
                }
            } else {
                context.resetValidation(statement.errors.useScope);
            }
            // statement.errors.action = {active: true, level: "warning", errorNumber: 2, message: "Action is not an ISO-defined term"};
        })
    };

    this.resetValidation = function (errorObject) {
        errorObject.active = false;
        errorObject.level = null;
        errorObject.errorNumber = 0;
        errorObject.message = null;
    };

    this.emptyStatement = function (statement) {
        return (ObjectUtils.isNull(statement.useScope) || statement.useScope.trim().length === 0) && (ObjectUtils.isNull(statement.action) || statement.action.trim().length === 0)
    }
}); 
    