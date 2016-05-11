var homeModule = angular.module("duck.home");

homeModule.controller("HomeController", function (DocumentService, $scope) {
    var home = this;
    home.summaries = [];
    DocumentService.getAuthoredDocumentSummaries().then(function (summaries) {
        home.summaries = summaries;
    }, function (status) {
        // FIXME display error
        alert('Failed: ' + status);
    })


});