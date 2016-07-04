var editorModule = angular.module("duck.editor");

/**
 * Maintains the global dictionary owned by an author.
 */
editorModule.service("GlobalDictionary", function (CurrentUser, $q, ObjectUtils) {
    this.dictionary = null;
    var context = this;

    this.getDictionary = function () {
        context.initialize();
        return context.dictionary.values();
    };

    this.getTerm = function (value) {
        context.initialize();
        return dictionary.get(value);
    };

    this.addTerm = function (type, code, category, value) {
        context.initialize();
        context.dictionary.put(value, {value: value, type: type, code: code, category: category, dictionaryType: "global"});
        // FIXME update server
    };

    this.removeTerm = function (type, code, value) {
        context.initialize();
        context.dictionary.remove(value);
        // FIXME implement server delete
    };

    this.initialize = function () {
        // FIXME populate from server
        if (context.dictionary !== null) {
            return;
        }
        context.dictionary = new Hashtable();
        context.dictionary.put("Microsoft Azure", {value: "Microsoft Azure", type: "scope", code: "microsoft_azure", category: "2", dictionaryType: "global"})
    }

});