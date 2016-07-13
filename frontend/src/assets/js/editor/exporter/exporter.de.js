var editorModule = angular.module("duck.editor");

/**
 * Exports documents in DE.
 */
editorModule.run(function (DocumentExporter) {

    /**
     * Exports DE text.
     */
    DocumentExporter.register("text/plain", "de", function (document) {
        var text = "";
        document.statements.forEach(function (statement) {
            text = text + statement.useScope + " anwendungen " + statement.qualifier + " " + statement.dataCategory + " von " + statement.sourceScope
                + " nach " + statement.action + " das" + " " + statement.resultScope + ".\n\n";
        });
        return new Blob([text], {type: 'text/plain;charset=utf-8'});
    });

});