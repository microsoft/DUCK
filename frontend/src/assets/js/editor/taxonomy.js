var editorModule = angular.module("duck.editor");

/**
 * Manages the ISO 19944 taxonomy including fuzzy lookup of data use statement element values.
 *
 * This service hosts a dictionary of standard ISO terms. When a document is edited, additional dictionaries will be activated, namely the global user
 * dictionary and the dictionary associated with the document. This allows end-users to create their own terms as subtypes of standard ISO terms, for example,
 * a product name that is specific to the organization hosting the DUCK application.
 */
editorModule.service("TaxonomyService", function ($sce, $log) {
    var context = this;

    context.cache = new Hashtable();  // key: symbol type, value: Hashtable [key: locale, value: list of symbol values]

    this.populate = function (type, entries) {
        var values = [];
        entries.forEach(function (entry) {
            values.push(entry.value);
        });
        var fuse = new Fuse(entries, {
            shouldSort: true,
            caseSensitive: false,
            threshold: 0.4,
            keys: ["value", "category"]
            // id: "value"
        });
        var localeCache = new Hashtable();
        localeCache.put("eng", {entries: entries, values: values, fuse: fuse});
        context.cache.put(type, localeCache);
    };

    // TODO the cache will be populated from a backend request
    this.populate("scope", [
        {value: "this capability", category: "1"},
        {value: "this application or this service", category: "2"},
        {value: "services listed in the service agreement", category: "3"},
        {value: "the CSP Services", category: "4"},
        {value: "all the CSP Products and services", category: "5"},
        {value: "third-party product and services", category: "6"}]);

    this.populate("qualifier", [{value: "unlinked pseudonymized"}, {value: "all"}]);
    this.populate("dataCategory", [{value: "email addresses"}, {value: "telemetry data"}, {value: "surfing habits"}]);
    this.populate("action", [{value: "provide"}, {value: "inform"}]);
    // TODO end cache population

    /**
     * Performs a fuzzy lookup of a set of values matching the given term.
     *
     * @param type the value type, e.g. use scope or qualifier
     * @param locale the language, e.g. "eng"
     * @param term the term to match
     * @return {Array} containing matching values in the form {value, label}
     */
    this.lookup = function (type, locale, term) {
        if (!term) {
            return [];
        }
        var typeCache = context.cache.get(type);
        if (typeCache == null) {
            $log.error("Unknown symbol type: " + type);
            return;
        }

        var symbolTable = typeCache.get(locale);
        if (symbolTable == null) {
            $log.error("Unknown locale when looking up symbol type '" + type + "': " + locale);
            return;
        }
        if (term.trim().length === 0) {
            // return all options
            var vals = [];
            symbolTable.entries.forEach(function (entry) {
                vals.push({value: entry.value, label: context.formatLabel(entry)});
            });
            vals.push({value: "_new", label: "New term..."});
            return vals;
        }
        var ret = symbolTable.fuse
            .search(term)
            .slice(0, 10)
            .map(function (entry) {
                var label = context.formatLabel(entry);
                return {
                    value: entry.value,
                    label: label
                };
            });
        ret.push({value: "_new", label: "New term..."});
        return ret;
    };

    this.contains = function (type, locale, term) {
        if (!term || term.trim().length === 0) {
            return false;
        }

        var typeCache = context.cache.get(type);
        if (typeCache == null) {
            $log.error("Unknown symbol type: " + type);
            return false;
        }

        var symbolTable = typeCache.get(locale);
        if (symbolTable == null) {
            $log.error("Unknown locale when looking up symbol type '" + type + "': " + locale);
            return false;
        }
        return symbolTable.values.includes(term);
    };

    /**
     * Activates a set of dictionaries. Dictionaries are activated when a document is edited, including the global user dictionary and the document dictionary.
     * @param dictionaries an array of dictionaries containing term objects in the form {value, type, subtype, dictionaryType}.
     */
    this.activate = function (dictionaries) {
        dictionaries.forEach(function (dictionary) {
            dictionary.forEach(function (term) {
                // put the term in each locale taxonomy
                var localeCache = context.cache.get(term.type);
                localeCache.values().forEach(function (symbolTable) {
                    symbolTable.entries.push({value: term.value, subtype: term.subtype, dictionary: true, dictionaryType: term.dictionaryType});
                    symbolTable.values.push(term.value);
                })
            });
        });
        context.reindex();

    };

    /**
     * Deactivates previously registered dictionaries by iterating all symbol table entries and removing any marked as a dictionary type.
     */
    this.deactivateDictionaries = function () {
        context.cache.values().forEach(function (localeCache) {
            localeCache.values().forEach(function (symbolTable) {
                for (var i = symbolTable.entries.length - 1; i >= 0; i--) {
                    if (symbolTable.entries[i].dictionary) {
                        symbolTable.values.without(function (entry) {
                            //noinspection JSReferencingMutableVariableFromClosure
                            return entry === symbolTable.entries[i].value;
                        });
                        symbolTable.entries.splice(i, 1);

                    }
                }
            });
        });
        context.reindex();
    };

    /**
     * Reindexes all symbol tables.
     */
    this.reindex = function () {
        // reset the search indexes
        context.cache.values().forEach(function (localeCache) {
            localeCache.values().forEach(function (symbolTable) {
                symbolTable.fuse = new Fuse(symbolTable.entries, {
                    shouldSort: true,
                    caseSensitive: false,
                    threshold: 0.4,
                    keys: ["value", "category"]
                    // id: "value"
                });
            });
        });

    };

    this.formatLabel = function (entry) {
        if (entry.category) {
            return $sce.trustAsHtml("<div class='clearfix'><div class='float-left'><strong>" + entry.value + "</strong></div>" +
                "<div class='dark-gray float-right'  style='z-index:10000'</div>" +
                "<div class='dark-gray float-right'>ISO</div>" +
                "</div>");
        } else {
            return $sce.trustAsHtml("    " + entry.value);

        }
    };


});