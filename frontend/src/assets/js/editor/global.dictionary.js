var editorModule = angular.module("duck.editor");

/**
 * Maintains the global dictionary owned by an author.
 */
editorModule.service("GlobalDictionary", function (CurrentUser, $http, $q) {
    this.dictionary = new Hashtable();
    var context = this;

    this.getDictionary = function () {
        return context.dictionary.values();
    };

    this.getTerm = function (value) {
        return dictionary.get(value);
    };

    this.addTerm = function (type, code, category, value) {
        var item = {value: value, type: type, code: code, category: category, dictionaryType: "global"};
        context.dictionary.put(value, item);
        return $q(function (resolve, reject) {
            // create the user and then log them in
            $http.put('/v1/users/'+CurrentUser.id + "/dictionary/" + code, item).success(function (data) {
                resolve(item)

            }).error(function (data, status) {
                reject(status);
            });
        });

        // FIXME update server
    };

    this.removeTerm = function (type, code, value) {
        context.dictionary.remove(value);
        // FIXME implement server delete
    };

    this.initialize = function () {
        $http.get('/v1/users/' + CurrentUser.id + "/dictionary").success(function (data) {
            var i = 1;

        }).error(function (data, status) {
            reject(status);
        });
        // context.dictionary = new Hashtable();
        context.dictionary.put("Microsoft Azure", {value: "Microsoft Azure", type: "scope", code: "microsoft_azure", category: "2", dictionaryType: "global"})
    }

});