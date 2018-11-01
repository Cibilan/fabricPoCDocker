//SPDX-License-Identifier: Apache-2.0

var contarct = require('./controller.js');

module.exports = function(app){
  
  app.get('/get_all_con', function(req, res){
    contarct.get_all_con(req, res);
  });

  app.get('/get_con/:id', function(req, res){
    contarct.get_con(req, res);
  });

  app.get('/add_con/:newCon', function(req, res){
    contarct.add_con(req, res);
  });

  app.get('/add_party/:newParty', function(req, res){
    contarct.add_party(req, res);
  });

    app.get('/con_act/:conAct', function(req, res){
    contarct.con_act(req, res);
  });

  app.get('/con_sign/:conSign', function(req, res){
    contarct.con_sign(req, res);
  });
    
}
