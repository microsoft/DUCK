var homeModule = angular.module("duck.editor");

/**
 * Manages the current document being edited.
 */
homeModule.service("DocumentModel", function (DataUseDocumentService, $q) {
    this.document = null;

    var context = this;

    /**
     * Retrieves the document from the backend and uses it to initialize the model.
     *
     * @param documentId the document id
     * @return a promise that will be resolved after initialization has successfully completed or failed
     */
    this.initialize = function (documentId) {
        return $q(function (resolve) {
            DataUseDocumentService.getDocument(documentId).then(function (document) {
                context.document = document;
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
    };

    /**
     * Adds the statement in the local model (i.e. it is not synchronized to the backend.
     *
     * @param statement the statement
     */
    this.addStatement = function (statement) {
        context.document.statements.push(statement);
    };

    /**
     * Saves the document to the backend.
     */
    this.save = function () {
        // TODO implement
    }

}); 
    