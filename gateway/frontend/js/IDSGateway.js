
var sampleApp = angular.module('gatewayApp', ["ui.router","ngResource",'uiGmapgoogle-maps',"google.places"]);

 sampleApp.config(function($stateProvider,$urlRouterProvider) {
    $urlRouterProvider.otherwise("/")



    $stateProvider

        .state("results", {
                url: "/results/:domainName",
                templateUrl: "/html/overview.html",
                controller: "ResultController",
                params : {
                    location: null
                }
                })

        .state("details", {
                url: "/details/:brokerId",
                templateUrl: "/html/details.html",
                controller: "ResultDetailsController"})
      });


  sampleApp.controller('MainController', ["$scope",'$state', function($scope, $state){
          $scope.queryInformation = function() {
            if ($scope.query.location != null) {
            var parsedLocation = parseLocation($scope.query.location);
            }
            $state.go("results",{"domainName": $scope.query.domain, "location":parsedLocation,"name": $scope.query.name});
          }


          function parseLocation(location) {
                var parsedLocation = {};
                var address_components = location.address_components;
                console.log(location);

              for (index = 0; index < address_components.length; ++index) {
                     var component = address_components[index];
                     console.log(component);
                     var types = component.types;
                     var type = types[0];
                     switch (type) {
                     case "locality":
                        parsedLocation.city = component.long_name;
                        break;
                     case "administrative_area_level_1":
                        parsedLocation.region = component.long_name;
                        break;
                     case "country":
                        parsedLocation.country = component.long_name;
                         break;
                     default:
                         break;
                     }
                }
                return parsedLocation;
          }
    }]);

  sampleApp.controller('ResultController', ["$scope",'$state','$stateParams','$resource' ,function($scope, $state,$stateParams,$resource){
        console.log("ResultController loaded")
        $scope.queryDomain = $stateParams.domainName;
        $scope.location = $stateParams.location;
        console.log($scope.location)
        console.log( $scope.queryDomain)
        $scope.onBrokerSelect = function(broker) {
                    console.log(broker)
                    $state.go("details",{"brokerId": broker.id})
                    }

    var Brokers = $resource("rest/brokers/:domainName", {domainName: '@domainName'}, {search: {method:"GET", params: {domainName: "@domainName", country: "country",region: "region",city:"city"}, isArray: true}});
     $scope.get = function(domainName,location){
                        console.lo
            			// Passing parameters to Book calls will become arguments if
            			// we haven't defined it as part of the path (we did with id)
                			Brokers.search({domainName:domainName,country:location.country,region:location.region, city: location.city}, function(data){
            				$scope.results = data;
            				console.log(data);
            			});
            		};
     $scope.get($scope.queryDomain,$scope.location);
  }]);

  sampleApp.controller('ResultDetailsController', ["$scope",'$state','$stateParams','$resource' ,function($scope, $state,$stateParams,$resource){
            var selectedBrokerId = $stateParams.brokerId;



            var DomainInformation = $resource("rest/brokers/:brokerId/domainInformation", {brokerId: '@brokerId'}, {search: {method:"GET", params: {brokerId: "@brokerId"}, isArray: false}});
                 $scope.getDomainInformation = function(brokerId){
                        			// Passing parameters to Book calls will become arguments if
                        			// we haven't defined it as part of the path (we did with id)
                            			DomainInformation.search({brokerId:brokerId}, function(data){
                        				$scope.details = data;
                        				$scope.broker = data.broker;
                        				$scope.getDomainController($scope.broker.realWorldDomain.name);
                        				console.log(data);
                                        showMapForBrokerLocation(data.broker.geolocation);
                        			});
                        		};
                 $scope.getDomainInformation(selectedBrokerId);

                 var DomainController = $resource("rest/domainControllers/:domainName", {domainName: '@domainName'}, {search: {method:"GET", params: {domainName: "@domainName"}, isArray: false}});
                                  $scope.getDomainController = function(domainName){
                                         			// Passing parameters to Book calls will become arguments if
                                         			// we haven't defined it as part of the path (we did with id)
                                             			DomainController.search({domainName:domainName}, function(data){
                                         				$scope.domainController = data;
                                         				console.log(data);
                                                        subscribeBrokerInformation();
                                         			});
                                         		};


                 function showMapForBrokerLocation(geolocation) {
                               console.log(geolocation)
                               var brokerLongitude = geolocation.longitude;
                                var brokerLatitude = geolocation.latitude;
                                var brokerLocation = new google.maps.LatLng(brokerLatitude,brokerLongitude);
                                 $scope.marker = {"id": 1,"location": brokerLocation};
                                 $scope.map = { center: brokerLocation, zoom: 10 };
                  };

                 function subscribeBrokerInformation() {
                         // Create a client instance

                   var address = $scope.domainController.brokerAddress.Host.split(":");
                   var host = address[0];
                   var port = address[1];
                   client = new Paho.MQTT.Client(host, parseInt(port), "clientId");

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
                       console.log("Subscribing to topic: " + $scope.broker.id );
                       client.subscribe($scope.broker.id);
                     }

                     // called when the client loses its connection
                     function onConnectionLost(responseObject) {
                       if (responseObject.errorCode !== 0) {
                         console.log("onConnectionLost:"+responseObject.errorMessage);
                       }
                     }
                      // called when a message arrives
                     function onMessageArrived(message) {
                       console.log("Received new Domain Information");
                       var json = JSON.parse(message.payloadString);
                       $scope.details = json;
                        $scope.$apply();
                     }
    }]);






