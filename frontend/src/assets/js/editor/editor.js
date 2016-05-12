var homeModule = angular.module("duck.editor");

homeModule.controller("EditorController", function (DataUseDocumentService, $stateParams, ObjectUtils, $scope) {
    var controller = this;

    var documentId = ObjectUtils.notNull($stateParams.documentId) ? $stateParams.documentId : null;
    controller.noDocument = documentId === null;

    if (controller.noDocument) {
        return;
    }

    controller.deleteStatement = function (statement) {
        alert("TODO: Not Implemented");
    };

    controller.editStatement = function (statement) {
        alert("TODO: Not Implemented");
    };

    DataUseDocumentService.getDocument(documentId).then(function (document) {
        controller.document = document;
    }, function (status) {
        // FIXME display error
        alert('Failed: ' + status);
    });


    // ng-sortable requires the use of $scope
    $scope.statements = [{id: 1, content: "Your data will not be shared with any third-party."}, {
        id: 2,
        content: "Your data may be used for advertising purposes."
    }];

    $scope.dragControlListeners = {
        containment: '#board',//optional param.
        allowDuplicates: true //optional param allows duplicates to be dropped.
    };


});