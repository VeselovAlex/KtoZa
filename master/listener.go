package main

import (
	"encoding/json"
	"log"

	"github.com/VeselovAlex/KtoZa/model"
)

func ListenPubSub() {
	for {
		msg, ok := App.PubSub.Await().(eventMessage)
		if !ok {
			// Пропускаем
			log.Println("LISTENER :: Not event message received")
		}

		switch msg.Event {
		case EventNewAnswerCache:
			cache := &model.Statistics{}
			err := json.Unmarshal(msg.Data, cache)
			if err != nil {
				log.Println("LISTENER ::Bad message decoding:", err)
				break
			}
			// Обновляем статистику
			upd := App.StatisticsStorage.CreateOrJoinWith(cache)
			if upd != nil {
				// Статистика обновлена
				func() {
					upd.Lock.RLock()
					defer upd.Lock.Unlock()
					App.PubSub.NotifyAll(About.UpdatedStatistics(upd))
				}()
			}

		default:
			// Пропускаем
		}
	}
}
