package cache

import (
	"container/list"
	"sync"
)

type LRUCache interface {
	Add(key string, value string) bool
	Get(key string) (string, bool)
	Delete(key string) bool
}

type item struct {
	key   string
	value string
}

type lruCache struct {
	capacity int
	cache    map[string]*list.Element
	queue    *list.List
	mutex    sync.Mutex
}

func NewLRUCache(capacity int) LRUCache {
	return &lruCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

// Add Добавляет новое значение с ключом в кеш (с наивысшим приоритетом), возвращает true, если все прошло успешно
// В случае дублирования ключа вернуть false
// В случае превышения размера - вытесняется наименее приоритетный элемент
func (c *lruCache) Add(key string, value string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, ok := c.cache[key]; ok {
		c.queue.MoveToFront(element)
		element.Value.(*item).value = value
		return false
	}

	elementNew := c.queue.PushFront(&item{
		key:   key,
		value: value,
	})
	c.cache[key] = elementNew

	if c.capacity != 0 && c.capacity < c.queue.Len() {
		c.removeElementLast()
	}

	return true
}

// Get Возвращает значение под ключом и флаг его наличия в кеше
// В случае наличия в кеше элемента повышает его приоритет
func (c *lruCache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, ok := c.cache[key]; ok {
		c.queue.MoveToFront(element)
		return element.Value.(*item).value, true
	}
	return "", false
}

// Delete Удаляет элемент из кеша, в случае успеха возврашает true, в случае отсутствия элемента - false
func (c *lruCache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, ok := c.cache[key]; ok {
		c.removeElement(element)
		return true
	}
	return false
}

func (c *lruCache) removeElement(element *list.Element) {
	c.queue.Remove(element)
	kv := element.Value.(*item)
	delete(c.cache, kv.key)
}

func (c *lruCache) removeElementLast() {
	elementLast := c.queue.Back()
	if elementLast != nil {
		c.removeElement(elementLast)
	}
}
