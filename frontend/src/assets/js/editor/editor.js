var homeModule = angular.module("duck.editor");

homeModule.controller("EditorController", function (DataUseDocumentService, $stateParams, ObjectUtils) {
    var controller = this;

    controller.documentId = ObjectUtils.notNull($stateParams.documentId) ? $stateParams.documentId : null;
    controller.noDocument = controller.documentId === null;

});