var editorModule = angular.module("duck.editor");

/**
 * Exports documents in US-EN.
 */
editorModule.run(function (DocumentExporter) {

    /**
     * Exports US-EN text.
     */
    DocumentExporter.register("text/plain", "json", function (document) {
        var exportObject = {statements:[]};

        document.statements.forEach(function (statement) {
            exportObject.statements.push(statement);
            // text = text + statement.useScope + " uses " + statement.qualifier + " " + statement.dataCategory + " from " + statement.sourceScope
            //     + " to " + statement.action + " the" + " " + statement.resultScope + ".\n\n";
        });
        var text = angular.toJson(exportObject);
        return new Blob([text], {type: 'text/plain;charset=utf-8'});
    });

});