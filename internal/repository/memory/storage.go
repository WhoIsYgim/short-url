package memory

import (
	"fmt"
	"short-link/internal/entities"
	"short-link/internal/utils"
	"short-link/pkg/errs"
	"sync"
	"time"
)

type LinkStorage struct {
	store map[string]*entities.Link
	mx    *sync.RWMutex
}

func NewLinkStorage() *LinkStorage {
	return &LinkStorage{
		store: map[string]*entities.Link{},
		mx:    new(sync.RWMutex),
	}
}

func (ls *LinkStorage) GetLink(token string) (*entities.Link, error) {
	ls.mx.RLock()
	defer ls.mx.RUnlock()

	link, ok := ls.store[token]
	if !ok {
		return nil, errs.NotFoundError()
	}
	return link, nil
}

func (ls *LinkStorage) GetLinkByOriginal(origLink string) (*entities.Link, error) {
	ls.mx.RLock()
	defer ls.mx.RUnlock()

	for _, link := range ls.store {
		if link.OriginalLink == origLink {
			return link, nil
		}
	}
	return nil, errs.NotFoundError()
}

func (ls *LinkStorage) StoreLink(link *entities.Link) error {
	ls.mx.Lock()
	ls.store[link.Token] = link
	ls.mx.Unlock()
	return nil
}

func (ls *LinkStorage) StartRecalculation(interval time.Duration, deleted chan []string) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			now := utils.CurrentTimeString()
			var deletedTokens []string
			ls.mx.Lock()
			for k, v := range ls.store {
				if v.Expired(now) {
					fmt.Println("Deleted: ", k)
					delete(ls.store, k)
					deletedTokens = append(deletedTokens, k)
				}
			}
			ls.mx.Unlock()
			if len(deleted) != 0 {
				deleted <- deletedTokens
			}
		}
	}()
	return
}

func (ls *LinkStorage) ShutDown() error { return nil }
