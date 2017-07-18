// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

var editorModule = angular.module("duck.editor");

/**
 * Handles document export.
 *
 * This service delegates to a registered exporter function based on the export content type (e.g. text/plain) and locale (e.g. "en", "de"). An exporter
 * renders a data use statement document into a representation based on the document locale and returns the result as an HTML 5 Blob.
 */
editorModule.service("DocumentExporter", function (FileSaver, $log) {
    this.exporters = new Hashtable();
    this.extensions = new Hashtable();

    this.extensions.put("text/plain", "txt");

    var context = this;

    this.register = function (type, locale, exporter) {
        context.exporters.put(type + ":" + locale, exporter);
    };

    this.export = function (type, document, localeOverride) {

        var exporter = localeOverride ? context.exporters.get(type + ":" + localeOverride) : context.exporters.get(type + ":" + document.locale);
        if (exporter === null) {
            $log.error("Document exporter not found: " + type + "," + document.locale);
            return "Error exporting document: Exporter not registered."
        }
        var data = exporter(document);

        // strip the name of invalid characters
        var name = document.name.replace("[\\~#%&*{}/:<>?|\"-]");

        // save the file
        FileSaver.saveAs(data, name + "." + context.extensions.get(type));
    };

});