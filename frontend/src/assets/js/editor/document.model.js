var homeModule = angular.module("duck.editor");

/**
 * Manages the current document being edited.
 */
homeModule.service("DocumentModel", function (DataUseDocumentService, $q) {
    /**
     * A local copy of the document.
     */
    this.document = null;

    /**
     * Tracks the local edit state of the document.
     */
    this.dirty = false;

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
                context.document = document;
                context.dirty = false;
                resolve();
            });
        });
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
     * Adds the statement in the local model (i.e. it is not synchronized to the backend.
     *
     * @param statement the statement
     */
    this.addStatement = function (statement) {
        context.document.statements.push(statement);
        context.dirty = true;
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
    }

}); 
    