var homeModule = angular.module("duck.home");

homeModule.controller("HomeController", function (DataUseDocumentService, $state) {
    var home = this;
    home.summaries = [];

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


    }

});