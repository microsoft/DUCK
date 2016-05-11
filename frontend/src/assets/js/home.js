var homeModule = angular.module("duck.home");

homeModule.controller("HomeController", function (DataUseDocumentService) {
    var home = this;
    home.summaries = [];
    DataUseDocumentService.getAuthoredDocumentSummaries().then(function (summaries) {
        home.summaries = summaries;
    }, function (status) {
        // FIXME display error
        alert('Failed: ' + status);
    })


});