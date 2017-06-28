// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

/**
 * Defines the gateway module. Gateway services manage client communications with the server.
 */
var gatewayModule = angular.module("duck.gateway");

/**
 * Signs the user into the system by creating a login.
 */
gatewayModule.service('SigninService', function (CurrentUser, $http, $q, $translate) {
    var context = this;
    this.signin = function (email, password) {
        return $q(function (resolve, reject) {
            $http.post('login', {email, password: password}).success(function (data) {
                CurrentUser.initializeWith(data);
                $translate.use(CurrentUser.locale).then(function () {
                    resolve();
                });
            }).error(function (data, status) {
                reject(status);
            });
        });
    };

    this.register = function (user) {
        return $q(function (resolve, reject) {
            // create the user and then log them in
            $http.post('/v1/users', user).success(function (data) {
                context.signin(user.email, user.password).then(function () {
                        resolve();
                    }, function (code) {
                        // FIXME
                        alert('Login Failed: ' + code);
                    }
                );

            }).error(function (data, status) {
                reject(status);
            });
        });
    }

});
gatewayModule.service('UserService', function (CurrentUser, $http, $q) {
    this.update = function () {
        CurrentUser.save();
        return $q(function (resolve, reject) {
            // FIXME implement
            resolve();
        });
    }
});

gatewayModule.service('RulebaseService', function (CurrentUser, $http, $q) {
    this.rulebases = null;
    var context = this;

    this.getRulebases = function () {
        return context.rulebases;
    };

    this.initialize = function () {
        return $q(function (resolve, reject) {
            // make sure the user is signed in
            if (!CurrentUser.loggedIn) {
                $state.go('signin');
            }
            var url = "/v1/rulebases";

            $http.get(url).success(function (data) {
                context.rulebases = angular.fromJson(data);
                resolve(context.rulebases);
            }).error(function (data, status) {
                reject(status);
            });
        });

    }
});

/**
 * Interceptor that adds an authorization token to the outbound request and handles errors reported from the server.
 */
gatewayModule.config(["$httpProvider", function ($httpProvider, $injector) {
    $httpProvider.interceptors.push(["$q", "$injector", function ($q, $injector) {
        return {
            "request": function (config) {
                if (config.url.lastIndexOf("https://", 0) === 0 || config.url.lastIndexOf("http://", 0) === 0) {
                    // don"t set header for requests out of the domain
                    return config;
                }
                var CurrentUser = $injector.get("CurrentUser");
                if (CurrentUser.token != null) {
                    config.headers["Authorization"] = "Bearer " + CurrentUser.token;
                }
                return config;
            },
            "responseError": function (response) {
                if (response.config.customErrorHandling) {
                    return $q.reject(response);
                }
                if (response.status === 401) {
                    $injector.get("$state").go("signin");
                }
                return $q.reject(response);
            }
        };
    }]);

}]);


gatewayModule.run(function ($rootScope, $state, CurrentUser) {

    $rootScope.$on('$stateChangeStart', function (event, toState) {
        if (!toState.requireSignin || CurrentUser.loggedIn) {
            // first case: no login required for this state, transition
            // second case: user logged in, transition
            return;
        }

        // user not logged in and attempting to access a restricted state, transition to signin instead
        event.preventDefault();
        $state.go('signin');
    });
});


