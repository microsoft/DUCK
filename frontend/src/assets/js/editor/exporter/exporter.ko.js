var editorModule = angular.module("duck.editor");

/**
 * Exports documents in KO.
 */
editorModule.run(function (DocumentExporter) {

    /**
     * Exports KO text.
     */
    DocumentExporter.register("text/plain", "ko", function (document) {
        var text = "";
        document.statements.forEach(function (statement) {
            text = text + statement.useScope + " uses " + statement.qualifier + " " + statement.dataCategory + " from " + statement.sourceScope
                + " to " + statement.action + " the" + " " + statement.resultScope + ".\n\n";
        });
        return new Blob([text], {type: 'text/plain;charset=utf-16'});
    });

});