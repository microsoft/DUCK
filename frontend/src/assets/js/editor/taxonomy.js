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
        });
        var localeCache = new Hashtable();
        localeCache.put("eng", {entries: entries, values: values, fuse: fuse});
        context.cache.put(type, localeCache);
    };

    // TODO the cache will be populated from a backend request
    this.populate("scope", [
        {value: "this capability", subtype: "capability", category: "1"},
        {value: "this application or this service", subtype: "service", category: "2"},
        {value: "services listed in the service agreement", subtype: "services-agreement", category: "3"},
        {value: "the CSP Services", subtype: "csp-services", category: "4"},
        {value: "all the CSP Products and services", subtype: "csp-products", category: "5"},
        {value: "third-party product and services", subtype: "third-party", category: "6"}]);

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
     * @param categories if true, include only terms that are categories
     * @return {Array} containing matching values in the form {value, label}
     */
    this.lookup = function (type, locale, term, categories) {
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
            return vals;
        }
        var ret = symbolTable.fuse
            .search(term)
            .slice(0, 10)
            .map(function (entry) {
                var label = context.formatLabel(entry);
                return {
                    value: entry.value,
                    label: label,
                    dictionary: entry.dictionary,
                    subtype: entry.subtype
                };
            });
        if (categories) {
            // filter terms that are not categories
            ret = ret.filter(function (term) {
                return !term.dictionary
            })
        }
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
     * Adds a new term to the taxonomy.
     * @param type the ISO type
     * @param subtype the subtype category
     * @param value the term value
     * @param dictionaryType the type of dictionary, e.g. global or document
     */
    this.addTerm = function (type, subtype, value, dictionaryType) {
        // deactivate all dictionaries, add the new term to the deactivated terms and reactivate the terms; this preserves sort order and may be faster 
        // than iterating over all dictionaries to determine the insertion point
        var entries = context.deactivateDictionaries();
        entries.push({type: type, subtype: subtype, value: value, dictionaryType: dictionaryType, dictionary: true});
        context.activate([entries]);

    };

    /**
     * Activates a set of dictionaries. Dictionaries are activated when a document is edited, including the global user dictionary and the document dictionary.
     * @param dictionaries an array of dictionaries containing term objects in the form {value, type, subtype, dictionaryType}.
     */
    this.activate = function (dictionaries) {
        // Sort all entries in reverse alphabetical order. This is because they must be inserted under their respective categories and the insertion
        // algorithm used below inserts at the first entry after the category is found. Hence, when iterating through the list, we need to do so in reverse
        // order so the highest entry is inserted last, at the top slot before the category. Otherwise, if the list were in alphabetical order, the
        // insertion algorithm would need to insert the entry at the last position after the category (and before the next category), which is less efficient.
        var terms = [];
        dictionaries.forEach(function (dictionary) {
            dictionary.forEach(function (term) {
                terms.push(term);
            })
        });
        terms.sort(function (entry1, entry2) {
            // note the comparison is correct
            return entry2.value.localeCompare(entry1.value);
        });

        //dictionaries.forEach(function (dictionary) {
        terms.forEach(function (term) {
            // put the term in each locale taxonomy
            var localeCache = context.cache.get(term.type);
            localeCache.values().forEach(function (symbolTable) {
                var inserted = false;
                // items must be inserted under their category type; iterate until the category is found and splice the entry in
                for (var i = 0; i < symbolTable.entries.length; i++) {
                    if (term.subtype === symbolTable.entries[i].subtype) {
                        symbolTable.entries.splice(i + 1, 0, {
                            value: term.value,
                            type: term.type,
                            subtype: term.subtype,
                            dictionary: true,
                            dictionaryType: term.dictionaryType
                        });
                        symbolTable.values.splice(i + 1, 0, term.value);
                        inserted = true;
                        break;
                    }
                }
                if (!inserted) {
                    // no category, add at the end
                    symbolTable.entries.push({value: term.value, subtype: term.subtype, dictionary: true, dictionaryType: term.dictionaryType});
                    symbolTable.values.push(term.value);
                }
            })
        });
        //});
        context.reindex();

    };

    /**
     * Deactivates previously registered dictionaries by iterating all symbol table entries and removing any marked as a dictionary type.
     * @return the deactivated terms
     */
    this.deactivateDictionaries = function () {
        var deactivatedEntries = [];
        context.cache.values().forEach(function (localeCache) {
            localeCache.values().forEach(function (symbolTable) {
                for (var i = symbolTable.entries.length - 1; i >= 0; i--) {
                    if (symbolTable.entries[i].dictionary) {
                        deactivatedEntries.push(symbolTable.entries[i]);
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
        return deactivatedEntries;
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
        if (!entry.dictionary) {
            return $sce.trustAsHtml("<div class='clearfix'><div class='float-left'><strong>" + entry.value + "</strong></div>" +
                "<div class='dark-gray float-right'  style='z-index:10000'</div>" +
                "<div class='dark-gray float-right'>ISO</div>" +
                "</div>");
        } else {
            return $sce.trustAsHtml("&nbsp;&nbsp;&nbsp;&nbsp;" + entry.value);

        }
    };


});