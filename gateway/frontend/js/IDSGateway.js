
var sampleApp = angular.module('gatewayApp', ["ui.router","ngResource",'uiGmapgoogle-maps',"google.places"]);

 sampleApp.config(function($stateProvider,$urlRouterProvider) {
    $urlRouterProvider.otherwise("/")



    $stateProvider

        .state("results", {
                url: "/results/:domainName?country&city",
                templateUrl: "/html/overview.html",
                controller: "ResultController",
                params : {
                    location: null
                }
                })

        .state("details", {
                url: "/details/:brokerId/:domain?country&city",
                templateUrl: "/html/details.html",
                controller: "ResultDetailsController"})
      });


  sampleApp.controller('MainController', ["$scope",'$state','$resource', function($scope, $state, $resource){

                 var Domains = $resource("rest/domains", {}, {search: {method:"GET", params: {}, isArray: true}});
                     $scope.getDomains = function(){
                            			// Passing parameters to Book calls will become arguments if
                            			// we haven't defined it as part of the path (we did with id)
                                			Domains.search({}, function(data){
                            				$scope.domains = data;
                            				console.log(data);
                            			});
                            		};
                       $scope.getDomains();
          $scope.queryInformation = function() {
            if ($scope.query.location != null) {
            var parsedLocation = parseLocation($scope.query.location);
            }
            $state.go("results",{"domainName": $scope.query.domain, "country":parsedLocation.country,"city":parsedLocation.city,"name": $scope.query.name});
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
        var location = {}
        location.country = $stateParams.country;
        location.city = $stateParams.city;
        $scope.location = location;
        $scope.name = $stateParams.name;

        $scope.onBrokerSelect = function(broker) {
                    console.log(broker)
                    $state.go("details",{"brokerId": broker.id,"domain":$scope.queryDomain,"country":$scope.location.country,"city":$scope.location.city,"name":$scope.name})
                    }

    var Brokers = $resource("rest/brokers/:domainName", {domainName: '@domainName'}, {search: {method:"GET", params: {domainName: "@domainName", country: "country"}, isArray: false}});
     $scope.getBrokers = function(domainName,location,name){
            			// Passing parameters to Book calls will become arguments if
            			// we haven't defined it as part of the path (we did with id)
            			    var encodedDomain = encodeURIComponent(domainName)
            			    console.log(encodedDomain)
                			Brokers.search({domainName:encodedDomain,location:location,name:name}, function(data){
            				$scope.results = data;
            				console.log(data);
            			});
            		};


     $scope.getBrokers($scope.queryDomain,$scope.location,$scope.name);



  }]);

  sampleApp.controller('ResultDetailsController', ["$scope",'$state','$stateParams','$resource' ,function($scope, $state,$stateParams,$resource){
initializeMap();
    var selectedBrokerId = $stateParams.brokerId;
     var domain = $stateParams.domain;
     var name = $stateParams.name;
           var location = {}
             location.country = $stateParams.country;
             location.city = $stateParams.city;



            var DomainInformation = $resource("rest/brokers/:brokerId/:domain", {brokerId: '@brokerId', domain: "@domain"}, {search: {method:"GET", params: {brokerId: "@brokerId",domain: "domain"}, isArray: false}});
                 $scope.getDomainInformation = function(brokerId,domain,location,name){
                        			// Passing parameters to Book calls will become arguments if
                        			// we haven't defined it as part of the path (we did with id)
                            			DomainInformation.search({brokerId:brokerId,domain:domain,location:location,name:name}, function(data){
                        				$scope.details = data;
                        				$scope.broker = data.broker;
                        				$scope.getDomainController(domain);
                                       showMapForBrokerLocation(data.broker.geolocation);
                        			});
                        		};
                 $scope.getDomainInformation(selectedBrokerId,domain,location,name);

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
                               var brokerLongitude = geolocation.longitude;
                                var brokerLatitude = geolocation.latitude;
                                var brokerLocation = new google.maps.LatLng(brokerLatitude,brokerLongitude);
                                console.log("Broker Location: "+brokerLocation)
                                 $scope.marker = {idkey: 1,coords: brokerLocation};
                                 $scope.map = {center: {latitude: brokerLatitude, longitude: brokerLongitude}, zoom: 10, refresh: true };
                  };

                  function initializeMap() {
                      $scope.map = {center: {latitude: -33.865143, longitude: 151.209900}, zoom: 10, refresh: "true"};
                      $scope.marker = {idkey: 1,coords:{latitude: -33.865143, longitude: 151.209900}};
                  }

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






