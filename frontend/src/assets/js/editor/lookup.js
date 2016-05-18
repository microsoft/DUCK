var editorModule = angular.module("duck.editor");
editorModule.service("UseScopeLookupService", function ($sce, $log) {
    this.scopeCache = new Hashtable();

    // TODO this will be created from a backend request
    var useScopes = ["this company", "this product", "this site", "this application"];
    var fuse = new Fuse(useScopes, {
        shouldSort: true,
        caseSensitive: false,
        threshold: 0.4
    });
    this.scopeCache.put("eng", {scopes: useScopes, fuse: fuse});

    var context = this;
    
    this.lookup = function (locale, term) {
        if (!term) {
            return [];
        }
        var scopeTable = context.scopeCache.get(locale);
        if (scopeTable == null) {
            $log.error("Unknown locale when looking up use scope: " + locale);
        }
        return scopeTable.fuse
            .search(term)
            .slice(0, 10)
            .map(function (i) {
                var val = scopeTable.scopes[i];
                return {
                    value: val,
                    label: $sce.trustAsHtml(val)
                };
            });
    }


});