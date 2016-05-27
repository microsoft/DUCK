var editorModule = angular.module("duck.editor");

editorModule.controller("EditorController", function (DocumentModel, TaxonomyService, EventBus,
                                                      $stateParams, AbandonComponent, ObjectUtils, $scope, $rootScope) {

    var controller = this;

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

    controller.toggleEdit = function (statement) {
        DocumentModel.toggleEdit(statement);
    };

    controller.editing = function (statement) {
        return DocumentModel.editing(statement);
    };

    controller.edit = function (statement) {
        DocumentModel.edit(statement);
    };

    controller.close = function (statement) {
        DocumentModel.close(statement);
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
        // ng-sortable requires the use of $scope
        $scope.document = DocumentModel.document;
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

    function initializeCompletions() {

        // setup autocompletes - requires $scope
        $scope.useScopeCompletion = {
            suggest: function (term) {
                var terms = TaxonomyService.lookup("scope", "eng", term);
                terms.push({value: "_new", label: "<span class='primary-text'>New term...</span>"});
                return terms
            },
            on_attach: function (value) {
                DocumentModel.document.statements.forEach(function (statement) {
                    if (!DocumentModel.editing(statement)) {
                        return;
                    }

                    // Register a watch all all use scopes of statements being edited. The watches monitor for the new term option selected by the user.
                    // If this occurs, an event to open the new term dialog is fired
                    var unregister = $scope.$watch(function () {
                        return statement.useScope
                    }, function (newValue) {
                        if (newValue === "_new") {   // new term entered
                            statement.useScope = "";
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
                controller.watches.values().forEach(function (unregister) {
                    unregister();     // deregister watches
                });
            }
        };

        $scope.qualifierCompletion = {
            suggest: function (term) {
                return TaxonomyService.lookup("qualifier", "eng", term)
            },
            on_detach: function (value) {
                DocumentModel.validateSyntax();
            }
        };

        $scope.dataCategoryCompletion = {
            suggest: function (term) {
                return TaxonomyService.lookup("dataCategory", "eng", term)
            },
            on_detach: function (value) {
                DocumentModel.validateSyntax();
            }
        };

        $scope.sourceScopeCompletion = {
            suggest: function (term) {
                return TaxonomyService.lookup("scope", "eng", term)
            },
            on_detach: function (value) {
                DocumentModel.validateSyntax();
            }
        };

        $scope.actionCompletion = {
            suggest: function (term) {
                return TaxonomyService.lookup("action", "eng", term)
            },
            on_detach: function (value) {
                DocumentModel.validateSyntax();
            }
        };

        $scope.resultScopeCompletion = {
            suggest: function (term) {
                return TaxonomyService.lookup("scope", "eng", term)
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
            console.log("suggest: " + term);
            return TaxonomyService.lookup("scope", "eng", term, true);
        },
        on_select: function (category) {
            controller.newTerm.category = category;
        },
        auto_select_first: true
    };

    controller.addTerm = function () {
        DocumentModel.addTerm("scope", controller.newTerm.category.subtype, controller.newTerm.value, controller.newTerm.dictionary ? "document" : "global");
        var statement = DocumentModel.getCurrentStatement();
        statement.useScope = controller.newTerm.value;
        DocumentModel.clearCurrentStatement();
        controller.clear();
    };

    controller.cancel = function () {
        DocumentModel.clearCurrentStatement();
        controller.clear();
    };

});


