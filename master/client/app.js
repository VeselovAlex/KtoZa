// Generated by CoffeeScript 1.10.0
(function() {
  var Poll, app,
    bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; };

  app = angular.module('ktoza-master', ['ngWebSocket']);

  Poll = (function() {
    function Poll($http, $websocket) {
      this.isCollapsed = bind(this.isCollapsed, this);
      this.collapse = bind(this.collapse, this);
      this.isCurrentView = bind(this.isCurrentView, this);
      this.setView = bind(this.setView, this);
      this.deleteOption = bind(this.deleteOption, this);
      this.newOption = bind(this.newOption, this);
      this.deleteQuestion = bind(this.deleteQuestion, this);
      this.newQuestion = bind(this.newQuestion, this);
      this.authorize = bind(this.authorize, this);
      this.isAuthorized = bind(this.isAuthorized, this);
      this.isValidPoll = bind(this.isValidPoll, this);
      this.isValidEndTime = bind(this.isValidEndTime, this);
      this.isValidStartTime = bind(this.isValidStartTime, this);
      this.isValidRegTime = bind(this.isValidRegTime, this);
      this.hasQuestions = bind(this.hasQuestions, this);
      this.hasTitle = bind(this.hasTitle, this);
      this.hasStatistics = bind(this.hasStatistics, this);
      this.hasPoll = bind(this.hasPoll, this);
      this.createPoll = bind(this.createPoll, this);
      $http.get("api/poll", {
        responseType: "json"
      }).then((function(_this) {
        return function(resp) {
          _this.poll = resp.data;
          _this.poll.events.registration = new Date(_this.poll.events.registration);
          _this.poll.events.registration.setMilliseconds(0);
          _this.poll.events.start = new Date(_this.poll.events.start);
          _this.poll.events.start.setMilliseconds(0);
          _this.poll.events.end = new Date(_this.poll.events.end);
          _this.poll.events.end.setMilliseconds(0);
          return _this.angQuestions = _this.poll.questions.map(function(q) {
            var angQ;
            angQ = {
              text: q.text,
              type: q.type
            };
            angQ.options = q.options.map(function(o) {
              return {
                option: o
              };
            });
            return angQ;
          });
        };
      })(this));
      $http.get("api/stats", {
        responseType: "json"
      }).then((function(_this) {
        return function(resp) {
          _this.statistics = resp.data;
          _this.statistics.date = new Date(_this.statistics.date);
          return _this.statistics.date.setMilliseconds(0);
        };
      })(this));
      $websocket("ws://" + location.host + location.pathname + "api/ws").onMessage((function(_this) {
        return function(msg) {
          var message;
          message = JSON.parse(msg.data);
          if (message.event === "poll-update") {
            _this.poll = message.data;
            return _this.statistics = null;
          } else if (message.event === "stats-update") {
            return _this.statistics = message.data;
          }
        };
      })(this));
      this.submit = (function(_this) {
        return function() {
          var poll;
          if (confirm("При обновлении опроса текущая статистика будет удалена. Продолжить?")) {
            poll = {
              title: _this.poll.title,
              caption: _this.poll.caption,
              events: _this.poll.events
            };
            poll.questions = _this.angQuestions.map(function(q) {
              var ret;
              ret = {
                type: q.type,
                text: q.text
              };
              ret.options = q.options.map(function(o) {
                return o.option;
              });
              return ret;
            });
            return $http.put("api/poll", poll).then(function() {
              _this.poll = poll;
              return console.log('Submitted');
            }, function() {
              console.log('Not submitted');
              return _this.badSubmission = true;
            });
          }
        };
      })(this);
    }

    Poll.prototype.createPoll = function() {
      var now;
      now = new Date();
      now.setMilliseconds(0);
      now.setSeconds(0);
      this.poll = {
        events: {
          registration: now,
          start: now,
          end: now
        }
      };
      return this.angQuestions = [];
    };

    Poll.prototype.hasPoll = function() {
      return this.poll != null;
    };

    Poll.prototype.hasStatistics = function() {
      return this.statistics != null;
    };

    Poll.prototype.hasTitle = function() {
      return this.hasPoll() && (this.poll.title != null) && this.poll.title !== "";
    };

    Poll.prototype.hasQuestions = function() {
      return (this.angQuestions != null) && this.angQuestions.length > 0;
    };

    Poll.prototype.hasText = function(q) {
      return (q.text != null) && q.text !== "";
    };

    Poll.prototype.hasType = function(q) {
      return (q.type != null) && q.type !== "";
    };

    Poll.prototype.hasOptions = function(q) {
      return (q.options != null) && q.options.length > 1;
    };

    Poll.prototype.isValidOption = function(o) {
      return (o.option != null) && o.option !== "";
    };

    Poll.prototype.isValidQuestion = function(q) {
      var ret;
      ret = this.hasText(q) && this.hasType(q) && this.hasOptions(q);
      return ret && q.options.reduce((function(_this) {
        return function(acc, opt) {
          return acc && _this.isValidOption(opt);
        };
      })(this), true);
    };

    Poll.prototype.isValidRegTime = function() {
      return this.hasPoll() && (this.poll.events.registration - new Date()) > 5000;
    };

    Poll.prototype.isValidStartTime = function() {
      return this.hasPoll() && (this.poll.events.start - this.poll.events.registration) >= 60000;
    };

    Poll.prototype.isValidEndTime = function() {
      return this.hasPoll() && (this.poll.events.end - this.poll.events.start) >= 60000;
    };

    Poll.prototype.isValidPoll = function() {
      var valid;
      valid = this.hasTitle();
      valid = valid && this.isValidRegTime() && this.isValidStartTime() && this.isValidEndTime();
      valid = valid && this.hasQuestions();
      return valid && this.angQuestions.reduce((function(_this) {
        return function(acc, q) {
          return acc && _this.isValidQuestion(q);
        };
      })(this), true);
    };

    Poll.prototype.isAuthorized = function() {
      if (this.auth != null) {
        return this.auth;
      }
      console.log('Checking auth');
      return this.auth = false;
    };

    Poll.prototype.authorize = function() {
      return this.auth = true;
    };

    Poll.prototype.newQuestion = function() {
      return this.angQuestions.push({
        options: []
      });
    };

    Poll.prototype.deleteQuestion = function(qNum) {
      var newQuestions;
      newQuestions = [];
      this.angQuestions.forEach(function(q, n) {
        if (n !== qNum) {
          return newQuestions.push(q);
        }
      });
      return this.angQuestions = newQuestions;
    };

    Poll.prototype.newOption = function(q) {
      return this.angQuestions[q].options.push({
        option: ""
      });
    };

    Poll.prototype.deleteOption = function(qNum, optNum) {
      var newOptions, question;
      question = this.angQuestions[qNum];
      newOptions = [];
      question.options.forEach(function(o, n) {
        if (n !== optNum) {
          return newOptions.push(o);
        }
      });
      return question.options = newOptions;
    };

    Poll.prototype.currentView = "intro";

    Poll.prototype.setView = function(view) {
      return this.currentView = view;
    };

    Poll.prototype.isCurrentView = function(view) {
      return this.currentView === view;
    };

    Poll.prototype.collapsed = [];

    Poll.prototype.collapse = function(q) {
      return this.collapsed[q] = !this.collapsed[q];
    };

    Poll.prototype.isCollapsed = function(q) {
      return this.collapsed[q];
    };

    return Poll;

  })();

  app.controller("Poll", Poll);

}).call(this);