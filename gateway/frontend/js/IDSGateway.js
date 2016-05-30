
var sampleApp = angular.module('gatewayApp', ["ngRoute","ngResource"]);

sampleApp.service('DataShare', function(){
  var service = {};
  service.data = false;
  service.setData = function(data) {
  this.data = data;
  }

  service.getData = function(){
  return this.data;}

   return service;
});

sampleApp .config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/ShowDetails', {
        templateUrl: '/html/brokerDetails.html',
        controller: 'ShowDetailsController'
      }).
      when('/ShowResults/:domain', {
        templateUrl: '/html/queryResults.html',
        controller: 'ShowResultsController'
      })
  }]);


 sampleApp.controller('ShowDetailsController', function($scope, $routeParams, DataShare) {
       console.log("BrokerDetailsController")
       $scope.domainInformation = DataShare.getData()

   });

  sampleApp.controller('ShowResultsController', function($scope, $resource, $routeParams, DataShare) {
    console.log("ShowResultsController")
    $scope.queryDomain = $routeParams.domain;
   $scope.onBrokerSelect = function(domainInformation) {
                    DataShare.setData(domainInformation);
                    location.href = "#ShowDetails";
               }


    var DomainInformation = $resource("rest/domainInformation/:domainName", {domainName: '@domainName'}, {search: {method:"GET", params: {domainName: "@domainName"}, isArray: true}});
    		$scope.get = function(domainName){
            			// Passing parameters to Book calls will become arguments if
            			// we haven't defined it as part of the path (we did with id)
                			DomainInformation.search({domainName:domainName}, function(data){
            				$scope.results = data;
            			});
            		};
            $scope.get($scope.queryDomain);
  });


  function init_map() {
          var var_location = new google.maps.LatLng(45.430817,12.331516);

          var var_mapoptions = {
            center: var_location,
            zoom: 1
          };

          var var_marker = new google.maps.Marker({
              position: var_location,
              map: var_map,
              title:"Venice"});

          var var_map = new google.maps.Map(document.getElementById("map-container"),
              var_mapoptions);

          var_marker.setMap(var_map);

        }


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

