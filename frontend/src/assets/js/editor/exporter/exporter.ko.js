// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

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
            text = text + statement.useScope + " uses " + statement.qualifier + " " + statement.dataCategory;

            statement.dataCategories.forEach(function(category){
                text = text + " " + category.operator + " " + category.qualifier + " " + category.dataCategory;
            });
            
            text = text + " from " + statement.sourceScope
                + " to " + statement.action + " the" + " " + statement.resultScope + ".\n\n";
        });
        return new Blob([text], {type: 'text/plain;charset=utf-16'});
    });

});