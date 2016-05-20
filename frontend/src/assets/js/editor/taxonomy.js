var editorModule = angular.module("duck.editor");

/**
 * Manages the ISO 19944 Taxonomy including fuzzy lookup of data use statement element values.
 */
editorModule.service("TaxonomyService", function ($sce, $log) {
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
    this.populate("scope", ["this capability", "this application or this service","services listed in the service agreement","the CSP Services",
        "all the CSP Products and services", "third-party product and services"]);
    this.populate("qualifier", ["unlinked pseudonymized", "all"]);
    this.populate("dataCategory", ["email addresses", "telemetry data", "surfing habits"]);
    this.populate("action", ["provide", "inform"]);
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
            symbolTable.values.forEach(function (val) {
                vals.push({value: val, label: $sce.trustAsHtml(val)});
            });
            return vals;
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
    }

});