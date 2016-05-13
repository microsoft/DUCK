var homeModule = angular.module("duck.editor");

homeModule.controller("EditorController", function (DocumentModel, $stateParams, ObjectUtils, $scope) {
    var controller = this;

    var documentId = ObjectUtils.notNull($stateParams.documentId) ? $stateParams.documentId : null;
    controller.noDocument = documentId === null;

    if (controller.noDocument) {
        return;
    }

    controller.dirty = function () {
        return DocumentModel.dirty;
    };

    controller.revert = function () {
        return DocumentModel.revert();
    };
    
    controller.editStatement = function (statement) {
        alert("TODO: Not Implemented");
    };

    controller.deleteStatement = function (statement) {
        DocumentModel.deleteStatement(statement);
    };

    controller.addStatement = function () {
        DocumentModel.addStatement({id: $scope.document.statements.length + 1, content: "Another statement " + ($scope.document.statements.length + 1)});
    };

    DocumentModel.initialize(documentId).then(function () {
        // ng-sortable requires the use of $scope
        $scope.document = DocumentModel.document;
    }, function (status) {
        // FIXME display error
        alert('Failed: ' + status);
    });


    // setup the sortable control listener
    controller.dragControlListeners = {
        allowDuplicates: true,
        orderChanged: function (event) {
            // var source = event.source.index;
            // var destination = event.dest.index;
        }
    };


});