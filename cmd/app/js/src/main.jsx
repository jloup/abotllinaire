var page                  = require('page');
//var Api                   = require('./api.js').Api;
//var parse                 = require('url-parse');
//var qs                    = require('querystringify');
//var classNames            = require('classnames');

page('/', function(context) {

  var root = (
              <div className='container'>
              </div>);

  ReactDOM.unmountComponentAtNode(document.getElementById('app'));
  ReactDOM.render(root, document.getElementById('app'));
});

$(document).ready(function() {
  page();
});
