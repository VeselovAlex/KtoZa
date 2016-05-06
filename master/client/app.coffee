app = angular.module('ktoza-master', ['ngWebSocket'])

class Poll
    constructor: ($http, $websocket) ->
        $http
            .get "api/poll",
                responseType:"json"
            .then (resp) =>
                @readPoll resp.data
        $http
            .get "api/stats",
                responseType: "json"
            .then (resp) =>
                    @readStatistics resp.data
             
        $websocket "ws://#{location.host}#{location.pathname}api/ws"
            .onMessage (msg) =>
                message = JSON.parse msg.data
                if message.event is "poll-update"
                    @readPoll message.data
                    @statistics = null
                else if message.event is "stats-update"
                    @readStatistics message.data
                   
        @checkAuth = =>
            $http
                .get "api/auth"
                .then =>
                        @auth = true
                    , =>
                        @auth = false
                        @badPwd = true
                            
        @authorize = =>
            $http
                .post "api/auth?p=#{@pwd}"
                .then =>
                        @checkAuth()
                        @authRequested = true
                        
        @checkAuth()
        
        @submit = =>
            if confirm "При обновлении опроса текущая статистика будет удалена. Продолжить?"
                poll = {
                    title: @poll.title
                    caption: @poll.caption
                    events: @poll.events
                }
                poll.questions = @angQuestions.map (q) ->
                    ret = {
                        type:  q.type
                        text:  q.text
                    }
                    ret.options = q.options.map (o) -> o.option
                    return ret
                $http
                    .put "api/poll", poll
                    .then =>
                            @poll = poll
                            console.log 'Submitted'
                        , =>
                            console.log 'Not submitted'
                            @badSubmission = true
                              
    createPoll: =>
        now = new Date()
        now.setMilliseconds(0)
        now.setSeconds(0)
        @poll = {
            events: {
                registration: now
                start: now
                end: now
            }
        }   
        @angQuestions = []
    
    readPoll: (from) =>
        @poll = from
        # Разбор дат
        @poll.events.registration = new Date(@poll.events.registration)
        @poll.events.registration.setMilliseconds(0)
        @poll.events.start = new Date(@poll.events.start)
        @poll.events.start.setMilliseconds(0)
        @poll.events.end = new Date(@poll.events.end)
        @poll.events.end.setMilliseconds(0)
        # Оборачиваем вопросы
        @angQuestions = @poll.questions.map (q) ->
            angQ = {
                text: q.text
                type: q.type
            }
            angQ.options = q.options.map (o) -> 
                option: o
            return angQ
    
    readStatistics: (from) =>
        @statistics = from
        @statistics.date = new Date(@statistics.date)
        @statistics.date.setMilliseconds(0)
    
    # Валидация
    hasPoll: => 
        @poll?
    hasStatistics: => 
        @statistics?
    hasTitle: =>
        @hasPoll() && @poll.title? && @poll.title isnt ""
    hasQuestions: =>
        @angQuestions? && @angQuestions.length > 0
    hasText: (q) ->
        q.text? && q.text isnt "" 
    hasType: (q) ->
        q.type? && q.type isnt "" 
    hasOptions: (q) ->
        q.options? && q.options.length > 1
    isValidOption: (o) ->
        o.option? && o.option isnt ""
    isValidQuestion: (q) ->
        ret = @hasText(q) && @hasType(q) && @hasOptions(q)
        # Проверка всех вариантов
        return ret && q.options.reduce (acc, opt) =>
                acc && @isValidOption(opt)
            , true
    isValidRegTime: =>
        # Больше 5 секунд до начала регистрации
        @hasPoll() && (@poll.events.registration - new Date()) > 5000 # ms
    isValidStartTime: =>
        # Хотя бы 1 минута на регистрацию
        @hasPoll() && (@poll.events.start - @poll.events.registration) >= 60000 # ms
    isValidEndTime: =>
        # Хотя бы 1 минута на опрос
        @hasPoll() && (@poll.events.end - @poll.events.start) >= 60000 # ms   
    isValidPoll: => 
        valid = @hasTitle()
        valid = valid && @isValidRegTime() && @isValidStartTime() && @isValidEndTime()
        valid = valid && @hasQuestions()
        return valid && @angQuestions.reduce (acc, q) =>
                acc && @isValidQuestion(q)
            , true
    
    # Управление авторизацией
    isAuthorized: => 
        @auth
    
    #Управление вопросами
    newQuestion: =>
        @angQuestions.push {
                options: []
            }
    deleteQuestion: (qNum) =>
        newQuestions = []
        @angQuestions.forEach (q, n) ->
            if n isnt qNum
                newQuestions.push q
        @angQuestions = newQuestions 
    newOption: (q) =>
        @angQuestions[q].options.push {option: ""}
    deleteOption: (qNum, optNum) =>
        question = @angQuestions[qNum]
        newOptions = []
        question.options.forEach (o, n) ->
            if n isnt optNum
                newOptions.push o
        question.options = newOptions
    
    # Управление текущим окном
    currentView: "intro"
    setView: (view) =>
        @currentView = view
    isCurrentView: (view) =>
        @currentView is view
          
    # Сворачивание статистики
    collapsed: []
    collapse: (q) =>
        @collapsed[q] = !@collapsed[q]
    isCollapsed: (q) =>  @collapsed[q]

    
app.controller("Poll", Poll)