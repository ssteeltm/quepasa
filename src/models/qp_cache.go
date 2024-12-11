package models

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

type QpCache struct {
	counter  atomic.Uint64
	cacheMap sync.Map
}

func (source *QpCache) Count() uint64 {
	return source.counter.Load()
}

func (source *QpCache) SetAny(key string, value interface{}, expiration time.Duration) {
	item := QpCacheItem{key, value, time.Now().Add(expiration)}
	source.SetCacheItem(item, "any")
}

// returns if it is a valid object, testing for now, it will not be necessary after debug
func (source *QpCache) SetCacheItem(item QpCacheItem, from string) bool {
	previous, loaded := source.cacheMap.Swap(item.Key, item)
	if loaded {
		// debugging messages in cache
		if strings.HasPrefix(from, "message") {
			prevItem := previous.(QpCacheItem)
			log.Warnf("[%s][%s] updating cache item ...", item.Key, from)
			log.Warnf("[%s][%s] old type: %s, %v", item.Key, from, reflect.TypeOf(prevItem.Value), prevItem.Value)
			log.Warnf("[%s][%s] new type: %s, %v", item.Key, from, reflect.TypeOf(item.Value), item.Value)
			log.Warnf("[%s][%s] equals: %v, deep equals: %v", item.Key, from, item.Value == prevItem.Value, reflect.DeepEqual(item.Value, prevItem.Value))

			var prevContent interface{}
			if prevWaMsg, ok := prevItem.Value.(*whatsapp.WhatsappMessage); ok {
				prevContent = prevWaMsg.Content

				if nee, ok := prevContent.(*waE2E.Message); ok {
					prevContent = nee.ProtoReflect()
				}
			}

			var newContent interface{}
			if newWaMsg, ok := item.Value.(*whatsapp.WhatsappMessage); ok {
				newContent = newWaMsg.Content

				if nee, ok := newContent.(*waE2E.Message); ok {
					newContent = nee.ProtoReflect()
				}
			}

			if prevContent != nil && newContent != nil {
				prevContentType := reflect.TypeOf(prevContent)
				newContentType := reflect.TypeOf(newContent)
				log.Warnf("[%s][%s] prev content type: %s, new content type: %s", item.Key, from, prevContentType, newContentType)
				if prevContentType == newContentType {

					b, err := json.Marshal(prevContent)
					if err == nil {
						log.Warnf("[%s][%s] old content as json: %s", item.Key, from, b)
					}

					b, err = json.Marshal(newContent)
					if err == nil {
						log.Warnf("[%s][%s] new content as json: %s", item.Key, from, b)
					}

					log.Warnf("[%s][%s] content equals: %v, content deep equals: %v", item.Key, from, prevContent == newContent, reflect.DeepEqual(prevContent, newContent))

					// if equals, deny triggers
					return !reflect.DeepEqual(prevContent, newContent)
				}
			}
		}
	} else {
		source.counter.Add(1)
	}

	return true
}

func (source *QpCache) GetAny(key string) (interface{}, bool) {
	if val, ok := source.cacheMap.Load(key); ok {
		item := val.(QpCacheItem)
		if time.Now().Before(item.Expiration) {
			return item.Value, true
		} else {
			source.DeleteByKey(key)
		}
	}
	return nil, false
}

func (source *QpCache) Delete(item QpCacheItem) {
	source.DeleteByKey(item.Key)
}

func (source *QpCache) DeleteByKey(key string) {
	_, loaded := source.cacheMap.LoadAndDelete(key)
	if loaded {
		source.counter.Add(^uint64(0))
	}
}

// gets a copy as array of cached items
func (source *QpCache) GetSliceOfCachedItems() (items []QpCacheItem) {
	source.cacheMap.Range(func(key, value any) bool {
		item := value.(QpCacheItem)
		items = append(items, item)
		return true
	})
	return items
}

// get a copy as array of cached items, ordered by expiration
func (source *QpCache) GetOrdered() (items []QpCacheItem) {

	// filling array
	items = source.GetSliceOfCachedItems()

	// ordering
	sort.Sort(QpCacheOrdering(items))
	return
}

// remove old ones, by timestamp, until a maximum length
func (source *QpCache) CleanUp(max uint64) {
	if max > 0 {
		length := source.counter.Load()
		amount := length - max
		if amount > 0 {
			items := source.GetOrdered()
			for i := 0; i < int(amount); i++ {
				source.DeleteByKey(items[i].Key)
			}
		}
	}
}
