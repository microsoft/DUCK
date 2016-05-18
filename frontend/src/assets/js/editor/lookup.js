var editorModule = angular.module("duck.editor");

/**
 * Provides fuzzy lookup of data use statement field values. A field value is used to populate an ISO-identified segment of the data use statement such as use
 * scope or qualifier.
 */
editorModule.service("ValueLookupService", function ($sce, $log) {
    var context = this;

    context.cache = new Hashtable();  // key: symbol type, value: Hashtable [key: locale, value: list of symbol values]

    this.populate = function (type, values) {
        var fuse = new Fuse(values, {
            shouldSort: true,
            caseSensitive: false,
            threshold: 0.4
        });
        var localCache = new Hashtable();
        localCache.put("eng", {values: values, fuse: fuse});
        context.cache.put(type, localCache);
    };

    // TODO the cache will be populated from a backend request
    this.populate("useScope", ["this company", "this product", "this site", "this application"]);
    this.populate("qualifier",["unlinked pseudonymized", "all"] );
    this.populate("dataCategory", ["email addresses", "telemetry data", "surfing habits"]);
    this.populate("sourceScope", ["this capability"]);
    this.populate("action",["provide", "inform"] );
    this.populate("resultScope", ["the services listed in this services agreement"]);
    // TODO end cache population


    /**
     * Performs a fuzzy lookup of a set of values matching the given term.
     * 
     * @param type the value type, e.g. use scope or qualifier
     * @param locale the language, e.g. "eng"
     * @param term the term to match
     * @return {Array} containing matching values
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

        return symbolTable.fuse
            .search(term)
            .slice(0, 10)
            .map(function (i) {
                var val = symbolTable.values[i];
                return {
                    value: val,
                    label: $sce.trustAsHtml(val)
                };
            });
    };


});