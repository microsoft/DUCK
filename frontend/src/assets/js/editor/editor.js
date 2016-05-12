var homeModule = angular.module("duck.editor");

homeModule.controller("EditorController", function (DataUseDocumentService, $stateParams, ObjectUtils) {
    var controller = this;

    var documentId = ObjectUtils.notNull($stateParams.documentId) ? $stateParams.documentId : null;
    controller.noDocument = documentId === null;

    if (controller.noDocument) {
        return;
    }

    DataUseDocumentService.getDocument(documentId).then(function (document) {
        controller.document = document;
    }, function (status) {
        // FIXME display error
        alert('Failed: ' + status);
    });
});