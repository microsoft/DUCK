var editorModule = angular.module("duck.editor");

/**
 * Manages the current document being edited.
 */
editorModule.service("DocumentModel", function (TaxonomyService, GlobalDictionary, DataUseDocumentService, $q, UUID, ObjectUtils) {
    /**
     * A local copy of the document.
     */
    this.document = null;

    /**
     * Tracks the local edit state of the document.
     */
    this.dirty = false;

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
            DataUseDocumentService.getDocument(documentId).then(function (document) {
                document.statements.forEach(function (statement) {
                    statement.errors = {
                        useScope: {active: false, level: null, action: false},
                        action: {active: false, level: null, action: false}
                    };
                });
                // FIXME: hardcode locale for now
                document.locale = "eng";

                // FIXME create a fake document dictionary for testing
                document.dictionary = new Hashtable();
                document.dictionary.put("Foo Service", {value: "Foo Service", type: "scope", subtype: "service", dictionaryType: "document"});

                context.document = document;
                context.dirty = false;

                // configure the taxonomy service with the global and document dictionaries as the document will be edited.
                TaxonomyService.activate([GlobalDictionary.getDictionary(), context.document.dictionary.values()]);

                // callback
                resolve();
            });
        });
    };

    this.release = function () {
        TaxonomyService.deactivateDictionaries();
    };

    /**
     * Deletes the statement in the local model (i.e. it is not synchronized to the backend.
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
     * @param subtype the subtype category
     * @param value the term value
     * @param dictionaryType the type of dictionary, e.g. global or document
     */
    this.addTerm = function (type, subtype, value, dictionaryType) {
        if (dictionaryType === "document") {
            context.document.dictionary.put(value, {value: value, type: type, subtype: subtype, dictionaryType: "document"});
        } else {
            GlobalDictionary.addTerm("scope", $scope.newCategory.subtype, $scope.newTermValue);
        }
        TaxonomyService.addTerm(type, subtype, value, dictionaryType);
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


    /**
     * Saves the local model to the backend.
     */
    this.save = function () {
        // TODO implement
        context.false = false;
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
    