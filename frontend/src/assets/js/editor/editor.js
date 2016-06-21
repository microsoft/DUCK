var editorModule = angular.module("duck.editor");

editorModule.controller("EditorController", function (DocumentModel, TaxonomyService, EventBus, LocaleService,
                                                      $stateParams, AbandonComponent, ObjectUtils, $scope, $rootScope) {

    var controller = this;

    controller.locales = LocaleService.getLocales();

    var documentId = ObjectUtils.notNull($stateParams.documentId) ? $stateParams.documentId : null;
    controller.noDocument = documentId === null;

    if (controller.noDocument) {
        return;
    }

    controller.active = true;

    // dynamic watches setup during statement editing
    controller.watches = new Hashtable();

    var unregisterDirtyCheck = $rootScope.$on("$stateChangeStart", function (event, toState) {
        if (!DocumentModel.dirty) {
            return;
        }
        AbandonComponent.open(event, toState);
    });

    $scope.$on("$destroy", function () {
        unregisterDirtyCheck();
        DocumentModel.release();
    });

    initializeCompletions();

    controller.setDocumentLocale = function (locale) {
        if (DocumentModel.document) {
            DocumentModel.document.locale = locale;
            DocumentModel.markDirty();
        }

    };

    controller.getLocale = function () {
        return DocumentModel.document ? DocumentModel.document.locale : null;
    };

    controller.save = function () {
        DocumentModel.save();
    };

    controller.toggleEdit = function (statement) {
        DocumentModel.toggleEdit(statement);
    };

    controller.editing = function (statement) {
        return DocumentModel.editing(statement);
    };

    controller.editAll = function () {
        return DocumentModel.document.statements.forEach(function (statement) {
            DocumentModel.edit(statement);
        });
    };

    controller.closeAll = function () {
        return DocumentModel.document.statements.forEach(function (statement) {
            DocumentModel.close(statement);
        });
    };

    controller.dirty = function () {
        return DocumentModel.dirty;
    };

    controller.revert = function () {
        return DocumentModel.revert();
    };

    controller.deleteStatement = function (statement) {
        DocumentModel.deleteStatement(statement);
    };

    controller.getLocalePrefix = function () {
        return "editor/" + DocumentModel.document.locale + "/" + DocumentModel.document.locale;
    };

    controller.addStatement = function () {
        DocumentModel.addStatement({
            useScope: null,
            qualifier: null,
            dataCategory: null,
            sourceScope: null,
            action: null,
            resultScope: null
        });
    };

    controller.hasErrors = function (statement) {
        var errors = statement.errors;
        if (ObjectUtils.isNull(errors)) {
            return false;
        }
        return errors.useScope.active || errors.action.active;
    };

    controller.emptyStatement = function (statement) {
        return DocumentModel.emptyStatement(statement);
    };

    controller.makePassive = function (statement) {
        DocumentModel.makePassive(statement);
    };

    controller.makeActive = function (statement) {
        DocumentModel.makeActive(statement);
    };

    DocumentModel.initialize(documentId).then(function () {
        // ng-sortable and watch requires the use of $scope
        $scope.document = DocumentModel.document;

        // deep watch the collection of statements
        $scope.$watch("document.statements", function () {
            DocumentModel.reCalculateCodes();
        }, true)
    }, function (status) {
        // FIXME display error
        alert("Failed: " + status);
    });


    // setup the sortable control listener
    controller.dragControlListeners = {
        allowDuplicates: true,
        orderChanged: function (event) {
            DocumentModel.markDirty();
        }
    };

    /**
     * Sets up Mass autocomplete input fields.
     */
    function initializeCompletions() {

        /**
         * Suggestion function used by Mass autocomplete for ISO scope fields. Since these fields may be extended with end-user defined values, a new term
         * option is added.
         *
         * @param term the input text typed by the user; may be a partial term
         * @return {Array} suggestions matching the input text
         */
        var scopeSuggest = function (term) {
            var terms = TaxonomyService.lookup("scope", DocumentModel.document.locale, term);
            terms.push({value: "_new", label: "<span class='primary-text'>New term...</span>"});
            return terms
        };

        /**
         * Mass autocomplete callback for ISO scope fields. Called when a field is activated by the end user for editing.
         *
         * This function sets up a watch for the new term option selected by a user. If triggered, the new term event is published which will result in a
         * modal being activated.
         *
         * @param fieldName the specific scope field name being activated, e.g. useScope.
         */
        var scopeAttach = function (fieldName) {
            $scope.currentField = fieldName;
            $scope.currentFieldType = "scope";
            DocumentModel.document.statements.forEach(function (statement) {
                if (!DocumentModel.editing(statement)) {
                    return;
                }
                // Register a watch all all use scopes of statements being edited. The watches monitor for the new term option selected by the user.
                // If this occurs, an event to open the new term dialog is fired
                var unregister = $scope.$watch(function () {
                    return statement[fieldName]
                }, function (newValue) {
                    if (newValue === "_new") {   // new term entered
                        statement[fieldName] = "";
                        DocumentModel.setCurrentStatement(statement);
                        EventBus.publish("ui.newTerm");
                    }
                });
                var other = controller.watches.put(statement.trackingId, unregister);
                if (other !== null) {
                    // unregister a previous watch setup for the statement
                    other();
                }

            });
        };

        /**
         * Mass autocomplete callback for ISO scope fields. Called when a field is deactivated.
         *
         * This function de-registers watches setup by the scopAttach function.
         *
         * @param value the value passed from Mass autocomplete.
         */
        var scopeDetach = function (value) {
            DocumentModel.validateSyntax();
            controller.watches.values().forEach(function (unregister) {
                unregister();     // deregister watches
            });

        };

        // setup autocompletes - requires $scope
        $scope.useScopeCompletion = {
            suggest: scopeSuggest,
            on_attach: function (value) {
                scopeAttach("useScope");
            },
            on_detach: scopeDetach
        };

        $scope.sourceScopeCompletion = {
            suggest: scopeSuggest,
            on_attach: function (value) {
                scopeAttach("sourceScope");
            },
            on_detach: scopeDetach
        };

        $scope.resultScopeCompletion = {
            suggest: scopeSuggest,
            on_attach: function (value) {
                scopeAttach("resultScope");
            },
            on_detach: scopeDetach
        };

        $scope.qualifierCompletion = {
            suggest: function (term) {
                return TaxonomyService.lookup("qualifier", DocumentModel.document.locale, term)
            },
            on_detach: function (value) {
                DocumentModel.validateSyntax();
            }
        };

        $scope.dataCategoryCompletion = {
            suggest: function (term) {
                var terms = TaxonomyService.lookup("dataCategory", DocumentModel.document.locale, term);
                terms.push({value: "_new", label: "<span class='primary-text'>New term...</span>"});
                return terms
            },
            on_attach: function (value) {
                $scope.currentField = "dataCategory";
                $scope.currentFieldType = "dataCategory";
                DocumentModel.document.statements.forEach(function (statement) {
                    if (!DocumentModel.editing(statement)) {
                        return;
                    }
                    var unregister = $scope.$watch(function () {
                        return statement.dataCategory
                    }, function (newValue) {
                        if (newValue === "_new") {   // new term entered
                            statement.dataCategory = "";
                            DocumentModel.setCurrentStatement(statement);
                            EventBus.publish("ui.newTerm");
                        }
                    });
                    var other = controller.watches.put(statement.trackingId, unregister);
                    if (other !== null) {
                        // unregister a previous watch setup for the statement
                        other();
                    }


                });
            },
            on_detach: function (value) {
                DocumentModel.validateSyntax();
            }
        };

        $scope.actionCompletion = {
            on_attach: function (value) {
                $scope.currentFieldType = "action";
            },
            suggest: function (term) {
                return TaxonomyService.lookup("action", DocumentModel.document.locale, term)
            },
            on_detach: function (value) {
                DocumentModel.validateSyntax();
            }
        };

    }

});

editorModule.controller("NewTermController", function (DocumentModel, TaxonomyService, GlobalDictionary, ObjectUtils, $scope) {
    var controller = this;
    controller.newTerm = {
        value: null,
        category: null,
        categoryValue: null,
        dictionary: "document"
    };

    controller.clear = function () {
        controller.newTerm.value = null;
        controller.newTerm.category = null;
        controller.newTerm.categoryValue = null;
        controller.newTerm.dictionary = "document"
    };

    $scope.newCategoryCompletion = {
        suggest: function (term) {
            return TaxonomyService.lookup($scope.currentFieldType, DocumentModel.document.locale, term, true);
        },
        on_select: function (category) {
            controller.newTerm.category = category;
        },
        auto_select_first: true
    };

    controller.addTerm = function () {
        var dictionaryType = controller.newTerm.dictionary ? "document" : "global";
        DocumentModel.addTerm($scope.currentFieldType, controller.newTerm.value, controller.newTerm.category.category, controller.newTerm.value, dictionaryType);
        var statement = DocumentModel.getCurrentStatement();
        statement[$scope.currentField] = controller.newTerm.value;
        DocumentModel.clearCurrentStatement();
        controller.clear();
    };

    controller.cancel = function () {
        DocumentModel.clearCurrentStatement();
        controller.clear();
    };

});


