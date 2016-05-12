###
# Александр Веселов <veselov143@gmail.com>
# СПбГУ, Математико-механический факультет, гр. 442
# Май, 2016 г.
###
app = angular.module('ktoza-respondent', [])

# Временное хранилище данных опроса и статистики для запрета отрисовки данных
# в несоответствующем состоянии
dataBuf = {}

class Poll
    constructor: ($http, $timeout, $interval) ->
        # Получение опроса
        $http
            .get "api/poll",
                responseType: "json"
            .then (resp) =>
                ###
                # Успешное завершение запроса
                ###
                poll = resp.data
                if not poll?
                    @loaded = false
                    return
                @header.title = poll.title
                @header.caption = poll.caption
                # Инициализация списка вопросов
                questions = poll.questions
                @collapsed = questions.map (q) -> false
                @questionTypes = questions.map (q) -> q.type
                @questionsAnswered = questions.map (q) -> false
                @applies = questions.map (q) ->
                    q.options.map (o) -> 0
                dataBuf.questions = questions
                # Разбор строк дат
                @events = {}
                @events.registrationAt = new Date(poll.events.registration)
                @events.startAt = new Date(poll.events.start)
                @events.endAt = new Date(poll.events.end)
                @loaded = true
                @doChangePollState($timeout)
                cd = $interval =>
                        @countdown.update @countdown.raw-1
                        if @countdown.raw < -1
                            $interval.cancel cd
                    , 1000
                
                    
        # Получение статистики
        @loadStat = =>
            $http
                .get "api/stats",
                    responseType: "json"
                .then (resp) =>
                    ###
                    # Успешное завершение запроса
                    ###
                    statistics = resp.data
                    statistics.date = new Date(statistics.date)
                    dataBuf.statistics = statistics
        @checkReg = =>
            $http
                .get "api/register",
                    responseType: "json"
                .then (resp) =>
                        @reg = resp.data
        @register = =>
            $http
                .post "api/register"
                .then =>
                        @reg = true
                    , =>
                        @reg = false
                
        @checkReg()
        
        # Отправка ответа      
        @submit = =>    
            toSubmit = @applies.map (a) ->
                ans = []
                a.forEach (opt, num) ->
                    if opt is 1 
                        ans.push(num)
                return ans
            $http
                .post "api/submit", toSubmit
                .then =>
                        # Успех
                        console.log 'Submitted'
                        @reg = false
                        @setView('intro')
                    , =>
                        # Отказ
                        console.log 'Not submitted'
                        @badSubmission = true
        
        # Начальное состояние -- "intro"
        @setView('intro')
                
    #Управление состоянием опроса
    pollState: "before"
    doChangePollState: ($timeout) =>
        now = new Date()
        reg = @events.registrationAt - now
        start = @events.startAt - now
        end = @events.endAt - now
        
        if end <= 0 
            @pollState = "ended"
            @countdown.update(-1)
            @loadStat()
        else if start <= 0
            @pollState = "started"
            @countdown.update(Math.ceil(end / 1000))
            $timeout => 
                    @pollState = "ended"
                    @countdown.update(-1)
                    @loadStat()
                , end
        else if reg <= 0
            @countdown.update(Math.ceil(start / 1000))
            @pollState = "registration"
            $timeout => 
                    @pollState = "ended"
                    @countdown.update(-1)
                    @loadStat()
                , end
            $timeout => 
                    @pollState = "started"
                    @countdown.update Math.ceil((@events.endAt - new Date()) / 1000)
                , start
        else
            @pollState = "before"
            @countdown.update(Math.ceil(reg / 1000))
            $timeout => 
                    @pollState = "ended"
                    @countdown.update(-1)
                    @loadStat()
                , end   
            $timeout => 
                    @pollState = "started"
                    @countdown.update Math.ceil((@events.endAt - new Date()) / 1000)
                , start
            $timeout => 
                    @pollState = "registration"
                    @countdown.update Math.ceil((@events.startAt - new Date()) / 1000)
                , reg
        return
            
    # Управление навигацией
    canShowQuestions: =>
        @loaded && @isRegistered() && (@pollState is "started")        
    canShowStatistics: =>
        @loaded && (@pollState is "ended")        
    
    # Управление регистрацией
    isRegistered: => @reg
    
    # Управление обратным отсчетом
    countdown: 
        raw: 0
        days: 0
        minutes: 0
        seconds: 0
        update: (raw) ->
            @raw = raw
            minutes = (raw - (@seconds = raw % 60)) / 60    
            hours =  (minutes - (@minutes = minutes % 60)) / 60    
            @days = (hours - (@hours = hours % 24)) / 24    
            @timeString = "#{@days} д. #{@hours} ч. #{@minutes} м. #{@seconds} с."
    # Управление текущим окном
    currentView: "intro"
    setView: (view) =>
        if @currentView is view
            return
        @currentView = view
        switch view
            when "intro"
                @questions = []
                @statistics = {}
            when "questions"
                @questions = dataBuf.questions
            when "stats"
                @questions = dataBuf.questions
                @statistics = dataBuf.statistics
            else
                @currentView = "intro"
    isCurrentView: (view) =>
        @currentView is view
    # Функции для сбора ответов
    apply: (q, o) =>
        type = @questionTypes[q]
        if type is "single-option"
            @applySingleOption(q, o)
        else if type is "multi-option"
            @applyMultiOption(q, o)
    
    applyMultiOption: (q, o) =>
        @applies[q][o] = (@applies[q][o] + 1) % 2
        @questionsAnswered[q] = @applies[q].indexOf(1) isnt -1
        
    applySingleOption: (q, o) =>
        last = @applies[q].indexOf 1
        if last isnt -1
            # Сброс последнего ответа
            @applies[q][last] = 0
            @questionsAnswered[q] = false
        if o isnt last
            # Выбран иной вариант ответа
            @applies[q][o] = 1
            @questionsAnswered[q] = true
    
    isApplied: (q, o) =>
        return @applies[q][o] is 1
    
    questionsAnswered: []
    numAnswered: =>
        @questionsAnswered.filter (x)->
                x
            .length
    
    isAnswered: (q) =>
        return @questionsAnswered[q]
    
    # Сворачивание ответов
    collapsed: []
    collapse: (q) =>
        @collapsed[q] = !@collapsed[q]
    isCollapsed: (q) =>  @collapsed[q]
    header: {}
    
app.controller("Poll", Poll)
