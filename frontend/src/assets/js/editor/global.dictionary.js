// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

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

    this.addTerm = function (type, code, category, location, value, case_1, case_2) {
        var item = {value: value, case_1: case_1, case_2: case_2, type: type, code: code, category: category, location:location, dictionaryType: "global"};

        context.dictionary.put(value, item);
        return $q(function (resolve, reject) {
            $http.put('/v1/users/'+CurrentUser.id + "/dictionary/" + code, item).success(function (data) {
                resolve(item)

            }).error(function (data, status) {
                reject(status);
            });
        });
    };

    this.removeTerm = function (type, code, value) {
        context.dictionary.remove(value);
        // FIXME implement server delete
    };

    this.initialize = function () {
        $http.get('/v1/users/' + CurrentUser.id + "/dictionary").success(function (data) {
            angular.forEach(data, function(term){
               context.dictionary.put(term.value,{value: term.value, case_1: term.case_1, case_2: term.case_2, type: term.type, code: term.code, location: term.location, category: term.category, dictionaryType: "global"});
            });
        }).error(function (data, status) {
            reject(status);
        });
    }

});