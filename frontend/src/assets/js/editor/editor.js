var homeModule = angular.module("duck.editor");

homeModule.controller("EditorController", function (DataUseDocumentService, $stateParams, ObjectUtils, $scope) {
    var controller = this;

    var documentId = ObjectUtils.notNull($stateParams.documentId) ? $stateParams.documentId : null;
    controller.noDocument = documentId === null;

    if (controller.noDocument) {
        return;
    }

    controller.deleteStatement = function (statement) {
        $scope.statements.without(function(element){
            return element === statement;
        });
        controller.listContents();
    };

    controller.editStatement = function (statement) {
        alert("TODO: Not Implemented");
    };

    controller.addStatement = function () {
        $scope.statements.push({id: $scope.statements.length + 1, content: "Another statement " + ($scope.statements.length + 1)});
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

    controller.listContents = function () {
        $scope.statements.forEach(function (statement) {
            console.log(statement.content);
        });
        console.log("-----------------------");
    };

    controller.dragControlListeners = {
        allowDuplicates: true,
        orderChanged: function (event) {
            var source = event.source.index;
            var destination = event.dest.index;
            listContents();
        }
    };


});