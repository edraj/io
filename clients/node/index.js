var express = require('express');
var app = express();
var bodyParser = require('body-parser');
var request = require('request');

app.use(bodyParser.json());

app.listen(8080, function() { console.log('Server started at http://localhost:8080'); });

app.get('/', function(req, res) { res.send('Hello world!'); });

