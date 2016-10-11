var homeModule = angular.module("duck.home");

homeModule.controller("HomeController", function (DataUseDocumentService, DocumentModel, $state) {
    var home = this;
    home.summaries = [];
    home.testDataType = "DEMO_DOCUMENT1";
    home.testData = null;

    // load document summaries for the current user
    DataUseDocumentService.getAuthoredDocumentSummaries().then(function (summaries) {
        home.summaries = summaries;
    }, function (status) {
        // FIXME display error
        alert('Failed: ' + status);
    });

    home.editDocument = function (documentId) {
        $state.go('main.editor', {documentId: documentId});
    };

    home.deleteDocument = function (documentId) {
        DataUseDocumentService.deleteDocument(documentId);
        home.summaries.without(function (summary) {
            return summary.id === documentId;
        });
    };

    home.createDocument = function (name) {
        DataUseDocumentService.createDocument(name).then(function (document) {
            $state.go('main.editor', {documentId: document.id});
        }, function (error) {
            //FIXME
        });


    };

    home.copyDocument = function (name, documentId) {
        DataUseDocumentService.copyDocument(name, documentId).then(function (document) {
            $state.go('main.editor', {documentId: document.id});
        }, function (error) {
            //FIXME
        });


    };

    home.createTestDocument = function (name, type) {
        DataUseDocumentService.createDocument(name).then(function (document) {
            DocumentModel.initialize(document.id).then(function () {
                home.testData[type].statements.forEach(function (statement) {
                    DocumentModel.addStatement(statement);
                });
                DocumentModel.lookupAndSetTerms(DocumentModel.document);

                DocumentModel.save().then(function () {
                    $state.go('main.editor', {documentId: document.id});
                });
            });
        }, function (error) {
            //FIXME
        });

    };

    DataUseDocumentService.getTestData().then(function (data) {
        home.testData = data;
    });

});