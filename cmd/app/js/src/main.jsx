var page                  = require('page');
var Api                   = require('./api.js').Api;
//var parse                 = require('url-parse');
//var qs                    = require('querystringify');
//var classNames            = require('classnames');

var VERSES_COUNT = 5;

var VerseUI = React.createClass({
  render: function() {

    var verses = this.props.verses.split('\n').map(function(verse) {
      return <p>{verse}</p>;
    });

    var temperature = ((this.props.temperature * 100) | 0) / 100.0;
    return (<div className='poem'>
              <div>{verses}</div>
              <div>(seed: {this.props.seed} - temp: {temperature})</div>
            </div>);
  }
});

var VersesUI = React.createClass({
  render: function() {
    var els = this.props.elements.map(function(element) {
      element.key = element.id;
      return React.createElement(VerseUI, element);
    });

    return (
    <div>{els}</div>
    );
  }
});

var VersesController = React.createClass({
  getInitialState: function() {
    return {'fromId': null, 'verses': []};
  },
  buildQuery: function() {
    if (this.state.fromId === null) {
      return {'count': VERSES_COUNT};
    } else {
      return {'count': VERSES_COUNT, 'fromId': this.state.fromId};
    }
  },
  fetch: function() {
    Api.Call('FB_VERSES', this.buildQuery(), function(err, status, response) {

      if (err !== null) {
        console.error(err);
        return;
      }

      this.setState({'verses': this.state.verses.concat(response), 'fromId': response[response.length - 1].id});
    }.bind(this));
  },
  componentDidMount: function() {
    this.fetch();
  },
  onNext: function() {
    this.fetch();
  },
  render: function() {
    var next = null;

    if (this.state.verses.length !== 0) {
      next = <button type='button' onClick={this.onNext}>next</button>;
    }

    return (<div>
              <VersesUI elements={this.state.verses}/>
              {next}
            </div>
           );
  }
});

page('/abotllinaire', function(context) {
  var root = (
            <div className='container'>
              <VersesController />
            </div>);

  ReactDOM.unmountComponentAtNode(document.getElementById('app'));
  ReactDOM.render(root, document.getElementById('app'));
});

$(document).ready(function() {
  page();
});
