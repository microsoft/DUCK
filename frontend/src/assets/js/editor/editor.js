// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

var editorModule = angular.module("duck.editor");

editorModule.controller("EditorController", function (DocumentModel, TaxonomyService, EventBus, LocaleService, AssumptionSetService, DocumentExporter,
                                                      $stateParams, AbandonComponent, NotificationService, ObjectUtils, $scope, $rootScope,
                                                      $anchorScroll, $location) {

    var controller = this;

    controller.NotificationService = NotificationService;
    controller.locales = LocaleService.getLocales();
    controller.assumptionSets = AssumptionSetService.getAssumptionSets();

    controller.filterStatement = null; // statement groups to filter on; provided as part of the validation result
    controller.preview = false;

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
        controller.filterStatement = null;
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

    controller.selectOriginal = function () {
        DocumentModel.selectOriginal();
        $scope.document = DocumentModel.document;
    };

    controller.deleteStatement = DocumentModel.deleteStatement;

    controller.duplicateStatement = DocumentModel.duplicateStatement;

    controller.emptyStatement = DocumentModel.emptyStatement;

    controller.makePassive = DocumentModel.makePassive;

    controller.makeActive = DocumentModel.makeActive;

    controller.complianceCheck = function () {
        var scrolled = false;
        DocumentModel.complianceCheck().then(function () {
            DocumentModel.document.statements.forEach(function (statement) {
                // scroll document to first statement error
                if (scrolled) {
                    return;
                }
                if (statement.errors.useScope.active
                    || statement.errors.qualifier.active
                    || statement.errors.dataCategory.active
                    || statement.errors.sourceScope.active
                    || statement.errors.action.active
                    || statement.errors.resultScope.active

                ) {
                    scrolled = true;
                    $location.hash(statement.trackingId);
                    $anchorScroll();
                }
            });
        });
    };

    controller.getState = function () {
        return DocumentModel.state;
    };

    controller.dirty = function () {
        return DocumentModel.dirty;
    };

    controller.getLocalePrefix = function () {
        return "editor/" + DocumentModel.document.locale + "/" + DocumentModel.document.locale;
    };

    controller.addStatement = function () {
        var passive = false;
        if (DocumentModel.document.statements.length > 0) {
            passive = DocumentModel.document.statements[DocumentModel.document.statements.length - 1].passive;
        }
        DocumentModel.addStatement({
            useScope: null,
            qualifier: TaxonomyService.findTerm("qualifier", "identified_data", DocumentModel.document.locale, "identified"),
            dataCategory: null,
            sourceScope: null,
            action: null,
            resultScope: null,
            passive: passive,
            dataCategories: []
        });
    };

    controller.addNewCategory = function(statement, operator){
        if(controller.existingExceptOperator(statement)){
            if(operator == "and"){
                statement.dataCategories.push({
                    "operator": operator,
                    "qualifierCode": "",
                    "dataCategoryCode": ""
                });
            }
        }
        else{
            /*
            var categoryCode = TaxonomyService.findCategory("dataCategory", statement.dataCategories[statement.dataCategories.length-1].dataCategory, DocumentModel.document.locale, statement.dataCategories[statement.dataCategories.length-1].dataCategory);
            console.log("categoryCode " + categoryCode);
            */

            statement.dataCategories.push({
                    "operator": operator,
                    "qualifierCode": "",
                    "dataCategoryCode": ""
            });
        }

    };

    controller.existingExceptOperator = function(statement){
        var existing = false;
        statement.dataCategories.forEach(function(category){
            if(category.operator == 'except'){
                existing = true;
            }
        });
        return existing;
    };

    controller.deleteCategory = function(statement, index){
        statement.dataCategories.splice(index, 1);
    };

    controller.hasErrors = function (statement) {
        var errors = statement.errors;
        if (ObjectUtils.isNull(errors)) {
            return false;
        }

        return errors.useScope.active || errors.qualifier.active || errors.dataCategory.active || errors.sourceScope.active || errors.action.active || errors.resultScope.active;
    };

    controller.downloadDocument = function () {
        prepareDataCategories();
        DocumentExporter.export("text/plain", DocumentModel.document);
    };

    controller.exportDocument = function () {
        prepareDataCategories();
        DocumentExporter.export("text/plain", DocumentModel.document, "json");
    };

    var prepareDataCategories = function(){
        DocumentModel.document.statements.forEach(function(statement){
            statement.dataCategories.forEach(function(category){
                category.qualifierCode = TaxonomyService.findCode("qualifier", category.qualifier, DocumentModel.document.locale, category.qualifier);
                category.dataCategoryCode = TaxonomyService.findCode("dataCategory", category.dataCategory, DocumentModel.document.locale, category.dataCategory);
            });
        });
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
            var terms = TaxonomyService.lookup("scope", DocumentModel.document.locale, term, false, false, false);
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
            DocumentModel.currentFieldType = "scope";  // used to tunnel to new term modal; look to refactor
            DocumentModel.currentField = fieldName;
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
                    } else if (fieldName === "sourceScope") {
                        // if other scopes are empty, default them to the source scope
                        if (statement.useScope === null || statement.useScope.trim().length === 0) {
                            statement.useScope = newValue;
                        }
                        if (statement.resultScope === null || statement.resultScope.trim().length === 0) {
                            statement.resultScope = newValue;
                        }
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
                return TaxonomyService.lookup("qualifier", DocumentModel.document.locale, term, false, false, false)
            },
            on_detach: function (value) {
                DocumentModel.validateSyntax();
            }
        };

        $scope.dataCategoryCompletion = {
            suggest: function (term) {
                var terms = TaxonomyService.lookup("dataCategory", DocumentModel.document.locale, term, false, false, false);
                terms.push({value: "_new", label: "<span class='primary-text'>New term...</span>"});
                return terms
            },
            on_attach: function (value) {
                $scope.currentField = "dataCategory";
                $scope.currentFieldType = "dataCategory";
                DocumentModel.currentFieldType = "dataCategory";  // used to tunnel to new term modal; look to refactor
                DocumentModel.currentField = "dataCategory";

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
                DocumentModel.currentFieldType = "action";  // used to tunnel to new term modal; look to refactor
                DocumentModel.currentField = "action";

            },
            suggest: function (term) {
                return TaxonomyService.lookup("action", DocumentModel.document.locale, term, false, false, false)
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
        /*
        case_1: null,
        case_2: null,
        */
        category: null,
        categoryValue: null,
        dictionary: "document",
        location: null,
        locationValue: null
    };

    controller.clear = function () {
        controller.newTerm.value = null;
        /*
        controller.newTerm.case_1 = null;
        controller.newTerm.case_2 = null;
        */
        controller.newTerm.category = null;
        controller.newTerm.categoryValue = null;
        controller.newTerm.dictionary = "document";
        controller.newTerm.location = null;
        controller.newTerm.locationValue = null;
    };

    $scope.newCategoryCompletion = {
        suggest: function (term) {
            return TaxonomyService.lookup(DocumentModel.currentFieldType, DocumentModel.document.locale, term, true, false, false);
        },
        on_select: function (category) {
            controller.newTerm.category = category;
        },
        auto_select_first: false
    };

    $scope.newLocationCompletion = {
        suggest: function (term) {
            return TaxonomyService.lookup("location", DocumentModel.document.locale, term, false, false, true);
        },
        on_select: function (location) {
            controller.newTerm.location = location;
        },
        auto_select_first: false
    };

    controller.addTerm = function () {
        var dictionaryType = controller.newTerm.dictionary === "document" ? "document" : "global";
        var code = controller.newTerm.value.split(" ").join("").toLowerCase(); // replace blank lines and convert to lowercase
        /*
        var case_1 = controller.newTerm.case_1;
        var case_2 = controller.newTerm.case_2;
        */

        var categoryCode = TaxonomyService.findCategory(DocumentModel.currentFieldType, controller.newTerm.categoryValue, DocumentModel.document.locale, controller.newTerm.categoryValue);
        var locationCode = TaxonomyService.findLocation("location", controller.newTerm.locationValue, DocumentModel.document.locale, controller.newTerm.locationValue);
        //DocumentModel.addTerm(DocumentModel.currentFieldType, code, categoryCode, locationCode, controller.newTerm.value, dictionaryType, case_1, case_2);
        DocumentModel.addTerm(DocumentModel.currentFieldType, code, categoryCode, locationCode, controller.newTerm.value, dictionaryType);
        var statement = DocumentModel.getCurrentStatement();
        statement[DocumentModel.currentField] = controller.newTerm.value;
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


