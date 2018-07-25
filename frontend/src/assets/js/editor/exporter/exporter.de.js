// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
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
            text = text + statement.useScope + " verwendet " + statement.qualifier + " " + statement.dataCategory;

            statement.dataCategories.forEach(function(category){
                var op = category.operator=="and" ? "und" : "außer";
                text = text + " " + op + " " + category.qualifier + " " + category.dataCategory;
            });

            text = text + " von " + statement.sourceScope
                + ", um "  + statement.resultScope +  statement.action + ".\n\n";
        });
        return new Blob([text], {type: 'text/plain;charset=utf-8'});
    });

});