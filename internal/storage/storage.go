package storage

import (
	"project/internal/models"
	"sync"
	"sync/atomic"
)

type Storage struct {
	mu    *sync.Mutex
	links map[int]models.LinksResponce
}

var linksNum int32

func New() *Storage {
	atomic.AddInt32(&linksNum, 1)
	return &Storage{
		mu:    &sync.Mutex{},
		links: make(map[int]models.LinksResponce),
	}
}

func (s *Storage) SaveLinks(links models.LinksResponce) {
	s.links[int(linksNum)] = links

	atomic.AddInt32(&linksNum, 1)
}

func (s *Storage) GetLinks(linksList []int32) []models.LinksResponce {
	slice := []models.LinksResponce{}
	for _, v := range linksList {
		slice = append(slice, s.links[int(v)])
	}
	return slice
}

func (s *Storage) GetLinksNum() int {
	return int(linksNum)
}
