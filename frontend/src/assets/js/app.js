/**
 * This module bootstraps the application, including defining URLs and pages.
 *
 * @type {angular.Module}
 */

var app = angular.module("duck.application", ["duck.main", "duck.user", "duck.signin", "ui.router"]);

app.factory("AppInfo", function () {

    return {
        name: "DUCK Application",
        version: "1.0.0"
    };
});

app.config(["$urlRouterProvider", "$locationProvider", "$stateProvider", "$logProvider", "$provide",
    function ($urlRouterProvider, $locationProvider, $stateProvider, $logProvider, $provide) {

        // set the debug log level
        $logProvider.debugEnabled(true);

        // decorate the logger to prepend DUCK prefix
        $provide.decorator('$log', ['$delegate', function ($delegate) {
            // pointer to original log methods
            var origDebug = $delegate.debug;
            var origInfo = $delegate.info;

            // override method to prepend DUCK prefix
            $delegate.debug = function () {
                var args = [].slice.call(arguments);
                args[0] = ["DUCK", ': ', args[0]].join('');

                //invoke original method
                origDebug.apply(null, args)
            };

            $delegate.info = function () {
                var args = [].slice.call(arguments);
                args[0] = ["DUCK", ': ', args[0]].join('');

                //invoke original method
                origInfo.apply(null, args)
            };

            return $delegate;
        }]);

        // setup URL structure
        $urlRouterProvider.otherwise("/");

        $locationProvider.html5Mode({
            enabled: false,
            requireBase: false
        });

        $locationProvider.hashPrefix("!");

        // define the application states
        $stateProvider
            .state("main", {  // the top-level state for protected (signed in) areas of the application
                url: "/",
                templateUrl: "../../main.html"
            })

            .state("signin", {   // signin and registration process
                url: "/signin",
                templateUrl: "../../signin.html"
            })
    }]);


app.controller("AppController", function (CurrentUser, AppInfo, $log) {
    $log.info("Initializing version " + AppInfo.version);

    CurrentUser.initialize();

});

app.run(function ($rootScope) {
    // load Foundation after the main controller has been initialized (as determined by the target scope)
    $rootScope.$on('$viewContentLoaded', function (event) {
        if (event.targetScope.initFoundation) {
            $(document).foundation();
        }
    });

});

/**
 * Monkey patches the Foundation reflow() method to not abort a reload if an existing plugin of the same type is found on a DOM element. The original
 * implementation aborted the reflow operation if more than one plugin was detected per DOM element without checking if the plugins were the same. This
 * caused an error with Angular UI router as Foundation needs to be reloaded when a view section changes.
 *
 * See the PATCH section inline.
 */

Foundation.reflow = function (elem, plugins) {

    // If plugins is undefined, just grab everything
    if (typeof plugins === 'undefined') {
        plugins = Object.keys(this._plugins);
    }
    // If plugins is a string, convert it to an array with one item
    else if (typeof plugins === 'string') {
        plugins = [plugins];
    }

    var _this = this;

    // Iterate through each plugin
    $.each(plugins, function (i, name) {
        // Get the current plugin
        var plugin = _this._plugins[name];

        // Localize the search to all elements inside elem, as well as elem itself, unless elem === document
        var $elem = $(elem).find('[data-' + name + ']').addBack('[data-' + name + ']');

        // For each plugin found, initialize it
        $elem.each(function () {
            var $el = $(this),
                opts = {};
            // Don't double-dip on plugins

            // PATCH: start replaced code
            var data = $el.data('zfPlugin');
            if (data && data.constructor.name !== plugin.name) {
                console.warn("Tried to initialize " + name + " on an element that already has a Foundation plugin.");
                return;
            }
            // PATCH: end replaced code

            if ($el.attr('data-options')) {
                $el.attr('data-options').split(';').forEach(function (e) {
                    var opt = e.split(':').map(function (el) {
                        return el.trim();
                    });
                    if (opt[0]) opts[opt[0]] = parseValue(opt[1]);
                });
            }
            try {
                $el.data('zfPlugin', new plugin($(this), opts));
            } catch (er) {
                console.error(er);
            }
        });
    });
};
