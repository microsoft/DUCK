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
    }

});