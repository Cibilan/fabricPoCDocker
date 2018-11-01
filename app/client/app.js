// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', ['ngMaterial','md-steppers']);

// Angular Controller
app.controller('appController', function($scope, appFactory, $mdDialog){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();

			$scope.stepper = {};
			$scope.stage1 = [];
			$scope.stage2 = [];
			$scope.stage3 = [];
			$scope.stage4 = [];
			$scope.stage5 = [];
			$scope.stage6 = [];

	$scope.user = "user1";		
	
	$scope.queryAllCon = function(){

		appFactory.queryAllCon(function(data){
		
			var array = [];
			for (var i = 0; i < data.length; i++){
	
				data[i].Record.Key = data[i].Key;
				array.push(data[i].Record);
			}

			$scope.allcon = true;

			console.log(array);
			$scope.all_con = array;
		});
	}

	$scope.queryCon = function(){

		var id = $scope.con_id;		

		$scope.stepper = {
    		step1Completed : false,
    		step2Completed : false,
    		step3Completed : false,
    		step4Completed : false,
    		step5Completed : false,
    		step6Completed : false,
    		disable : false,
    		selected : 0
    		};
    	$scope.conHome = true;	  		


		appFactory.queryCon(id, function(data){
			delete $scope.con;	
    		delete $scope.stage1;
    		delete $scope.stage2;
    		delete $scope.stage3;
    		delete $scope.stage4;
    		delete $scope.all_party;
    		delete $scope.new_Party_Success;
    		delete $scope.newParty;
    		delete $scope.act;
    		delete $scope.new_conAct_Success;
    		  

    		console.log(data);

			$scope.con = data;
			$scope.con.Key = $scope.con_id;
			$scope.stage1 = [];
			$scope.stage2 = [];
			$scope.stage3 = [];
			$scope.stage4 = [];
			$scope.all_party = [];

			console.log($scope.con);

			$scope.conDetail = true;

    		angular.forEach($scope.con.historylist, function(list){
				if(list.stage == "Contract Created"){
					$scope.stage1.push(list);	
				}
				if(list.stage == "Contract Activation"){
					$scope.stage2.push(list);	
				}
				if(list.stage == "Contract Signing"){
					$scope.stage3.push(list);	
				}
				if(list.stage == "Contract Validation"){
					$scope.stage4.push(list);	
				}
			})

				if($scope.con.stage == "Contract Created"){
					
					$scope.showSign = false;
					$scope.showParty = true;
				}
				if($scope.con.stage == "Contract Activation"){
					$scope.showParty = true;
					
					$scope.showSign = false;
				}
				if($scope.con.stage == "Contract Signing"){
					$scope.showSign = true;
					$scope.showParty = false;
					
				}
				if($scope.con.stage == "Contract Validation"){
					$scope.showSign = true;
					$scope.showParty = false;
					
				}
			$scope.all_party = $scope.con.partylist;	

			if ($scope.query_con == "Could not locate Contract"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	// $scope.addCon = function(){

	// 	$scope.newCon.user = $scope.user;
	// 	console.log($scope.newCon);
	// 	appFactory.addCon($scope.newCon, function(data){
	// 		$scope.new_Con_Success = data;
	// 	});
	// }

	$scope.addParty = function(){
		$scope.newParty.user = $scope.user;
		$scope.newParty.key = $scope.con_id;
		console.log($scope.newParty);
		appFactory.addParty($scope.newParty, function(data){
			$scope.new_Party_Success = data;
		});
	}

	$scope.conAct = function(){

		$scope.act.fileName = angular.element('#fileName')[0].value;
		$scope.act.fileHash = angular.element('#fileHash')[0].value;
		$scope.act.user = $scope.user;
		$scope.act.key = $scope.con_id;
		console.log($scope.act);
		appFactory.conAct($scope.act, function(data){
			$scope.new_conAct_Success = data;
		});

	}

	$scope.conSign = function(){
		$scope.conSign.user = $scope.user;
		$scope.conSign.key = $scope.con_id;
		console.log($scope.conSign);
		appFactory.conSign($scope.conSign, function(data){
			$scope.partySign_Success = data;
		});
	}


	$scope.showPrompt = function(ev) {
    $mdDialog.show({
      controller: DialogController,
      templateUrl: 'views/createContract.html',
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose:true,
    })
  	};

   function DialogController($scope, $mdDialog) {
    $scope.addCon1 = function(){
		$scope.newCon.user = "user1";
		console.log($scope.newCon);
		appFactory.addCon($scope.newCon, function(data){
			$scope.new_Con_Success = data;
		});
		}
  	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllCon = function(callback){

    	$http.get('/get_all_con/').success(function(output){
			callback(output)
		});
	}

	factory.queryCon = function(id, callback){
    	$http.get('/get_con/'+id).success(function(output){
			callback(output)
		});
	}

	factory.addCon = function(data, callback){

		var con = data.id + "-" + data.detail + "-" + data.user;

    	$http.get('/add_con/'+con).success(function(output){
			callback(output)
		});
	}

	factory.addParty = function(data, callback){

		var newParty = data.key + "-" + data.partyName + "-" + data.mandatory  + "-" + data.user + "-" + data.partyWeight;

		console.log(newParty);

    	$http.get('/add_party/'+newParty).success(function(output){
			callback(output)
		});
	}

	factory.conAct = function(data, callback){

		var conAct = data.key + "-" + data.fileName + "-" + data.fileHash + "-" + data.condition + "-" + data.user;

		console.log(conAct);

    	$http.get('/con_act/'+conAct).success(function(output){
			callback(output)
		});
	}

	factory.conSign = function(data, callback){

		var conSign = data.key + "-" + data.user;

		console.log(conSign);

    	$http.get('/con_sign/'+conSign).success(function(output){
			callback(output)
		});
	}

	return factory;
});


