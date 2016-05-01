// Generated by CoffeeScript 1.10.0
(function() {
  var Poll, StateView;

  Poll = React.createClass({
    getInitialState: function() {
      return {
        loaded: false
      };
    },
    getCurrentView: function() {
      return this.state.view;
    },
    componentDidMount: function() {
      return $.ajax({
        url: this.props.fromUrl,
        dataType: 'json',
        cache: false,
        success: (function(_this) {
          return function(poll) {
            poll.loaded = true;
            poll.registrationAt = new Date(poll.events.registration);
            poll.startAt = new Date(poll.events.start);
            poll.endAt = new Date(poll.events.end);
            _this.setState(poll);
            _this.setState({
              view: 'intro'
            });
            return setTimeout(function() {
              return _this.setState({
                view: 'registration'
              });
            }, 3000);
          };
        })(this),
        error: (function(_this) {
          return function(xhr, statu, err) {
            return console.log("Bad request", status, err.toString());
          };
        })(this)
      });
    },
    render: function() {
      if (this.state.loaded) {
        return (
              <div className="poll">
                <header>
                  <h1>{this.state.title}</h1>
                  <p>{this.state.caption}</p>
                  <div className="event-timings">
                    <div className="timing">
                      Начало регистрации {this.state.registrationAt.toLocaleDateString()} 
                      {' '} в {this.state.registrationAt.toLocaleTimeString()}
                    </div>
                    <div className="timing">
                      Начало опроса {this.state.startAt.toLocaleDateString()} 
                      {' '} в {this.state.startAt.toLocaleTimeString()}
                    </div>
                    <div className="timing">
                      Окончание опроса {this.state.endAt.toLocaleDateString()} 
                      {' '} в {this.state.endAt.toLocaleTimeString()}
                    </div>
                  </div>  
                </header>
                <nav>
                    <a href="/#intro">Intro</a><br />
                    <a href="/#registration">Registration</a><br />
                    <a href="/#answers">Answers</a><br />
                    <a href="/#finish">Finish</a><br />
                    <a href="/#stats">Stats</a>
                </nav>
                <StateView getView={this.getCurrentView}/>
              </div>
            );
      } else {
        return <div class="loading">Loading...</div>;
      }
    }
  });

  StateView = React.createClass({
    render: function() {
      return (
            <section className="view">
              {this.props.getView()}
            </section>
        );
    }
  });


  /*
   */

  ReactDOM.render(<Poll fromUrl="/testapi/poll.json"/>, $("#poll")[0]);

}).call(this);