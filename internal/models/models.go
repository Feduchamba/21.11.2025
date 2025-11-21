package models

type LinksRequest struct {
	Links []string `json:"links"`
}

type LinksResponce struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}

type PastLinksRequest struct {
	LinksNum []int32 `json:"linksNum"`
}
