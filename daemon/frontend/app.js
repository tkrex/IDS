// app.js
// create our angular app and inject ngAnimate and ui-router
// =============================================================================
angular.module('setupApp', ['ngAnimate', 'ui.router'])

// configuring our routes
// =============================================================================
.config(function($stateProvider, $urlRouterProvider) {

    $stateProvider

        // route to show our basic form (/form)
        .state('setup', {
            url: '/setup',
            templateUrl: 'setup.html',
            controller: 'setupController'
        })

        // nested states
        // each of these sections will have their own view
        // url will be nested (/form/profile)
        .state('setup.address', {
            url: '/profile',
            templateUrl: 'form-profile.html'
        })

        // url will be /form/interests
        .state('form.interests', {
            url: '/interests',
            templateUrl: 'form-interests.html'
        })

        // url will be /form/payment
        .state('form.payment', {
            url: '/payment',
            templateUrl: 'form-payment.html'
        });

    // catch all route
    // send users to the form page
    $urlRouterProvider.otherwise('/form/profile');
})

// our controller for the form
// =============================================================================
.controller('formController', function($scope) {

    // we will store all of our form data in this object
    $scope.formData = {};

    // function to process the form
    $scope.processForm = function() {
        alert('awesome!');
    };

});