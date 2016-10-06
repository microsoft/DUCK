var editorModule = angular.module("duck.editor");

editorModule.controller("EditorController", function (DocumentModel, TaxonomyService, EventBus, LocaleService, AssumptionSetService, DocumentExporter,
                                                      $stateParams, AbandonComponent, NotificationService, ObjectUtils, $scope, $rootScope) {

    var controller = this;

    controller.NotificationService = NotificationService;
    controller.locales = LocaleService.getLocales();
    controller.assumptionSets = AssumptionSetService.getAssumptionSets();

    controller.filterStatement = null; // statement groups to filter on; provided as part of the validation result

    controller.filterOnStatement = function (statement) {
        // This will trigger the 'compatibleFilter' filter to remove statements determined (from validation) not compatible with the
        // given statement
        controller.filterStatement = statement;
    };

    controller.showAllStatements = function () {
        controller.filterStatement = null;
    };

    controller.isComplianceChecked = function () {
        return DocumentModel.state !== "UNKNOWN";
    };

    controller.filtering = function () {
        return controller.filterStatement != null;
    };

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
        controller.filterStatement = null
        DocumentModel.release();
    });

    initializeCompletions();

    controller.setDocumentLocale = DocumentModel.setDocumentLocale;

    controller.getLocale = function () {
        return DocumentModel.document ? DocumentModel.document.locale : null;
    };

    controller.setAssumptionSet = DocumentModel.setAssumptionSet;

    controller.getAssumptionSet = function () {
        return DocumentModel.document ? DocumentModel.document.assumptionSet : null;
    };

    controller.isSelected = function (document) {
        return DocumentModel.document === document;
    };

    controller.save = DocumentModel.save;

    controller.isEditable = DocumentModel.isEditable;

    controller.selectOriginal = function () {
        DocumentModel.selectOriginal();
        $scope.document = DocumentModel.document;
    };

    controller.revert = DocumentModel.revert;

    controller.deleteStatement = DocumentModel.deleteStatement;

    controller.emptyStatement = DocumentModel.emptyStatement;

    controller.makePassive = DocumentModel.makePassive;

    controller.makeActive = DocumentModel.makeActive;

    controller.complianceCheck = DocumentModel.complianceCheck;

    controller.getState = function () {
        return DocumentModel.state;
    };

    controller.toggleEdit = function (statement) {
        DocumentModel.toggleEdit(statement);
    };

    controller.dirty = function () {
        return DocumentModel.dirty;
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

    controller.downloadDocument = function () {
        DocumentExporter.export("text/plain", DocumentModel.document);
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
                console.log("attached:" + value);
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
        var dictionaryType = controller.newTerm.dictionary === "document" ? "document" : "global";
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

editorModule.filter('compatibleFilter', function (ObjectUtils) {
    return function (statements, filterStatement) {
        if (ObjectUtils.isNull(statements)) {
            return;
        }
        if (ObjectUtils.isNull(filterStatement) || ObjectUtils.isNull(filterStatement.$$statementExplanation)) {
            return statements; // no compliance check is active, do not filter
        }
        var set = new HashSet();
        set.add(filterStatement.trackingId);
        set.addAll(filterStatement.$$statementExplanation.compatiblePurpose);

        var filtered = [];
        statements.forEach(function (statement) {
            if (set.contains(statement.trackingId)) {
                filtered.push(statement);
            }
        });
        return filtered;
    }
});


