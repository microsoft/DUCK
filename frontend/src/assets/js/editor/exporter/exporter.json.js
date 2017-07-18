// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
var editorModule = angular.module("duck.editor");

/**
 * Exports documents in US-EN.
 */
editorModule.run(function (DocumentExporter) {

    /**
     * Exports US-EN text.
     */
    DocumentExporter.register("text/plain", "json", function (document) {
        var exportObject = {statements: []};

        document.statements.forEach(function (statement) {
            exportObject.statements.push({
                "useScopeCode": statement.useScopeCode,
                "qualifierCode": statement.qualifierCode,
                "dataCategoryCode": statement.dataCategoryCode,
                "sourceScopeCode": statement.sourceScopeCode,
                "actionCode": statement.actionCode,
                "resultScopeCode": statement.resultScopeCode

            });
        });
        var text = angular.toJson(exportObject);
        return new Blob([text], {type: 'text/plain;charset=utf-8'});
    });

});