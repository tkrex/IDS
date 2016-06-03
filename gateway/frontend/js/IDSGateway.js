
var sampleApp = angular.module('gatewayApp', ["ui.router","ngResource",'uiGmapgoogle-maps']);

 sampleApp.config(function($stateProvider,$urlRouterProvider) {
    $urlRouterProvider.otherwise("/")



    $stateProvider

        .state("results", {
                url: "/results/:domainName",
                templateUrl: "/html/results-overview.html",
                controller: "ResultController"
                })

        .state("details", {
                url: "/details/:brokerId",
                templateUrl: "/html/results-details.html",
                controller: "ResultDetailsController"})
      });




  sampleApp.controller('MainController', ["$scope",'$state', function($scope, $state){
          $scope.queryInformation = function(domainName) {
            $state.go("results",{"domainName": domainName});
          }
    }]);

  sampleApp.controller('ResultController', ["$scope",'$state','$stateParams','$resource' ,function($scope, $state,$stateParams,$resource){
        console.log("ResultController loaded")
        $scope.queryDomain = $stateParams.domainName;
        $scope.onBrokerSelect = function(broker) {
                    console.log(broker)
                    $state.go("details",{"brokerId": broker.id})
                    }

    var Brokers = $resource("rest/brokers/:domainName", {domainName: '@domainName'}, {search: {method:"GET", params: {domainName: "@domainName"}, isArray: true}});
     $scope.get = function(domainName){
            			// Passing parameters to Book calls will become arguments if
            			// we haven't defined it as part of the path (we did with id)
                			Brokers.search({domainName:domainName}, function(data){
            				$scope.results = data;
            				console.log(data);
            			});
            		};
     $scope.get($scope.queryDomain);
  }]);

  sampleApp.controller('ResultDetailsController', ["$scope",'$state','$stateParams','$resource' ,function($scope, $state,$stateParams,$resource){
            var selectedBrokerId = $stateParams.brokerId;



            var DomainInformation = $resource("rest/brokers/:brokerId/domainInformation", {brokerId: '@brokerId'}, {search: {method:"GET", params: {brokerId: "@brokerId"}, isArray: false}});
                 $scope.get = function(brokerId){
                        			// Passing parameters to Book calls will become arguments if
                        			// we haven't defined it as part of the path (we did with id)
                            			DomainInformation.search({brokerId:brokerId}, function(data){
                        				$scope.details = data;
                        				console.log(data);
                                        showMapForBrokerLocation(data.broker.geolocation);
                        			});
                        		};
                 $scope.get(selectedBrokerId);

                 function showMapForBrokerLocation(geolocation) {
                               console.log(geolocation)
                               var brokerLongitude = geolocation.longitude;
                                var brokerLatitude = geolocation.latitude;
                                var brokerLocation = new google.maps.LatLng(brokerLatitude,brokerLongitude);
                                 $scope.marker = {"id": 1,"location": brokerLocation};
                                 $scope.map = { center: brokerLocation, zoom: 10 };
                                                                     };


    }]);





  function subscribeBrokerInformation() {
        // Create a client instance
  client = new Paho.MQTT.Client("localhost", 1883, "clientId");

  // set callback handlers
  client.onConnectionLost = onConnectionLost;
  client.onMessageArrived = onMessageArrived;

  // connect the client
  client.connect({onSuccess:onConnect});
  }

  // called when the client connects
    function onConnect() {
      // Once a connection has been made, make a subscription and send a message.
      console.log("onConnect");
      client.subscribe("#");
    }

    // called when the client loses its connection
    function onConnectionLost(responseObject) {
      if (responseObject.errorCode !== 0) {
        console.log("onConnectionLost:"+responseObject.errorMessage);
      }
    }
     // called when a message arrives
    function onMessageArrived(message) {
      console.log("onMessageArrived:"+message.payloadString);
    }

