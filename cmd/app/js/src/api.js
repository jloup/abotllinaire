var Api = {};

Api.API_HOST = 'jlj.am';
Api.API_PORT = '';
Api.API_ROOT = 'https://' + Api.API_HOST + '/api/feeder';

var ApiEndpoints = {
  'FEED_BY_ID': {
            'type': 'GET',
            'parameters': [
                            {'name': 'id',     'mandatory': true,  'inquery': false},
                            {'name': 't',      'mandatory': false, 'inquery': true},
                            {'name': 'exclid', 'mandatory': false, 'inquery': true}
                          ],
            'buildUrl': function(params) {
              return '/feed/' + params.id;
            }
          },
};

Api.BuildRequest = function(endpointId, params) {
  var endpoint = ApiEndpoints[endpointId];
  var query    = {};
  var url      = '';
  var headers  = {};

  if (endpoint === undefined) {
    return {'err': 'cannot find endpoint ' + endpointId};
  }

  for (var i = 0; i < endpoint.parameters.length; i++) {
    var parameter = endpoint.parameters[i];
    if (parameter.mandatory === true && params[parameter.name] === undefined) {
      return {'err': 'parameter \'' + parameter.name + '\' is mandatory'};
    }

    if (parameter.inquery === true) {
      query[parameter.name] = params[parameter.name];
    }
  }

  url = endpoint.buildUrl(params);

  return {'err': null, 'url': Api.API_ROOT + url, 'headers': headers, 'params': query, 'type': endpoint.type};
};

Api.Call = function(endpointId, params, fn) {
  var request = Api.BuildRequest(endpointId, params);

  if (request.err !== null) {
    console.error('Api BuildRequest error', request.err);
    fn({'err': request.err}, -1, {});
    return null;
  }

  return Api.DoRequest(request.type, request.url, request.params, request.headers, fn);
};

Api.DoRequest = function(type, url, params, headers, callback) {
  return $.ajax({
    data:     params,
    headers:  headers,
    timeout:  60000,
    dataType: 'json',
    error: function(xhr, status, err) {
      if (status === 'abort') {
        return;
      }

      var apiStatus   = null;
      var apiResponse = null;

      if (xhr.responseJSON) {
        apiStatus   = xhr.responseJSON.status;
        apiResponse = xhr.responseJSON.res;
      }
      callback({'url': url, 'err': err, 'status': status, 'xhr': xhr}, apiStatus, apiResponse);
    },
    success: function(data) {
      var err = null;
      if (data.status !== 0) {
        err = {'url': url, 'status': status};
      }
      callback(err, data.status, data.res);
    },
    type: type,
    url: url
  });
};

module.exports.Api = Api;

