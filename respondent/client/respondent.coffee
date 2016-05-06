App = React.createClass
    getInitialState: -> 
        loaded: false
    componentDidMount: ->
        $.ajax
            url: @props.pollURL
            dataType: "json"
            cache: false
            success: (poll) =>
                @setState
                    loaded: true
                    poll: poll     
            error: =>
                @setState
                    loaded: true
                    poll: null
    render: ->
        if not @state.loaded
            return `<NotLoaded />`
        if not @state.poll
            return `<NoPoll />`
        return `<Poll poll={this.state.poll}/>`

NotLoaded = React.createClass
    render: ->
        `<div>Загрузка...</div>`

NoPoll = React.createClass
    render: ->
        return `(
            <div>Нет опроса. Повторите запрос позднее.</div>
        )`
        
Poll = React.createClass
    render: ->
        poll = @props.poll
        questions = poll.questions.map (q, n)-> `<Question data={q} key={n} name={n}/>`
        return `(
        <section>
          <header>
            <h1>{poll.title}</h1>
            <p>{poll.caption}</p>
          </header>
          <div>
            <h2>Вопросы</h2>
            {questions}
          </div>
        </section>
        )`
        
Question = React.createClass
    select: (e) ->
        console.log e
    render: -> 
        q = @props.data
        name = @props.name
        select = @select
        options = q.options.map (o, n) -> 
            `<Option data={o} type={q.type} key={n} name={name} value={n} select={select}/>`
        return `(
          <div>
            <h3>{this.props.data.text}</h3>
            {options}
          </div>
        )`
    
Option = React.createClass
    render: -> 
        type = @props.type
        input = ""
        if type is 'single-option'
            input = "radio"
        else if type is 'multi-option'
            input = "checkbox"
        else
            console.log "Bad question type"
            return `<span>Bad question</span>`
        name = @props.name
        value = @props.value
        return `<li><input type={input} name={name} value={value} onChange={}/>{this.props.data}</li>`
    
ReactDOM.render `<App pollURL="/testapi/poll.json"/>`, $("#app")[0]