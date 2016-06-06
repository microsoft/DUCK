/**
 * Defines the gateway module. Gateway services manage client communications with the server.
 */
var gatewayModule = angular.module("duck.gateway");

/**
 * Signs the user into the system by creating a login.
 */
gatewayModule.service('SigninService', function (CurrentUser, $http, $q, $translate) {

    this.signin = function (username, password) {
        return $q(function (resolve, reject) {
            // FIXME workaround until backend fixed
            if (true) {
                CurrentUser.initializeWith({firstName: "Andy", lastName: "Author", id: "123", token: "124", locale: "en"});
                $translate.use(CurrentUser.locale).then(function () {
                    resolve();
                });
                return;
            }

            $http.post('/login', {username, password: password}).success(function (data) {
                CurrentUser.initializeWith(data);
                $translate.use(CurrentUser.locale).then(function () {
                    resolve();
                });
            }).error(function (data, status) {
                reject(status);
            });
        });
    }
});
gatewayModule.service('UserService', function (CurrentUser, $http, $q) {
    this.update = function () {
        return $q(function (resolve, reject) {
            // FIXME implement
            resolve();
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
                    config.headers["X-Authorization"] = CurrentUser.token;
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


