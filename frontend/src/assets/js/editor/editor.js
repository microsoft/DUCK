var editorModule = angular.module("duck.editor");

editorModule.controller("EditorController", function (DocumentModel, ValueLookupService,
                                                      $stateParams, ObjectUtils, AbandonComponent, $scope, $rootScope) {

    var controller = this;

    var documentId = ObjectUtils.notNull($stateParams.documentId) ? $stateParams.documentId : null;
    controller.noDocument = documentId === null;

    if (controller.noDocument) {
        return;
    }

    controller.active = true;

    var unregisterDirtyCheck = $rootScope.$on("$stateChangeStart", function (event, toState) {
        if (!DocumentModel.dirty) {
            return;
        }
        AbandonComponent.open(event, toState);
    });

    $scope.$on("$destroy", function () {
        unregisterDirtyCheck();
    });

    // setup autocompletes - requires $scope
    $scope.useScopeCompletion = {
        suggest: function (term) {
            return ValueLookupService.lookup("useScope", "eng", term)
        }
    };
    $scope.qualifierCompletion = {
        suggest: function (term) {
            return ValueLookupService.lookup("qualifier", "eng", term)
        }
    };

    $scope.dataCategoryCompletion = {
        suggest: function (term) {
            return ValueLookupService.lookup("dataCategory", "eng", term)
        }
    };

    $scope.sourceScopeCompletion = {
        suggest: function (term) {
            return ValueLookupService.lookup("sourceScope", "eng", term)
        }
    };

    $scope.actionCompletion = {
        suggest: function (term) {
            return ValueLookupService.lookup("action", "eng", term)
        }
    };

    $scope.resultScopeCompletion = {
        suggest: function (term) {
            return ValueLookupService.lookup("resultScope", "eng", term)
        }
    };



    controller.toggleEdit = function (statement) {
        DocumentModel.toggleEdit(statement);
    };

    controller.editing = function (statement) {
        return DocumentModel.editing(statement);
    };

    controller.dirty = function () {
        return DocumentModel.dirty;
    };

    controller.revert = function () {
        return DocumentModel.revert();
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
        alert("Failed: " + status);
    });


    // setup the sortable control listener
    controller.dragControlListeners = {
        allowDuplicates: true,
        orderChanged: function (event) {
            DocumentModel.markDirty();
        }
    };


});