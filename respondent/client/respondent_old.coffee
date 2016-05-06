Poll = React.createClass
    getInitialState: -> {loaded: false}
    
    getCurrentView: -> @state.view
    
    componentDidMount: ->
        $.ajax
            url: @props.fromUrl
            dataType: 'json'
            cache: false
            success: (poll) =>
                if not poll
                    window.location += "/nopoll"
                    return
                poll.loaded = true
                poll.registrationAt= new Date poll.events.registration
                poll.startAt = new Date poll.events.start
                poll.endAt = new Date poll.events.end
                
                @setState poll
                @setState
                    view: 'intro'
                    
                setTimeout =>
                        @setState
                            view: 'registration'
                    , 3000
            error: (xhr, statu, err) =>
                console.log "Bad request", status, err.toString()
    render: -> 
        if @state.loaded
            return `(
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
            )`
        else return `<div class="loading">Loading...</div>`

# "intro"        -- начало опроса
# "registration" -- регистрация
# "answers"      -- прием ответов
# "finish"       -- окончание опроса
# "stats"        -- отображение статистики
StateView = React.createClass
    render: -> 
        return `(
            <section className="view">
              {this.props.getView()}
            </section>
        )`

ReactDOM.render `<Poll fromUrl="testapi/poll"/>`, $("#poll")[0]