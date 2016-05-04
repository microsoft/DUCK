/**
 * Defines the gateway module. Gateway services manage client communications with the server.
 */
var gatewayModule = angular.module("duck.gateway", []);

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
                } else {
                    // FIXME
                    alert("Server Error");
                }
                return $q.reject(response);
            }
        };
    }]);

}]);

