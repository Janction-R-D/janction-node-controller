package utils

type Page struct {
	Offset  uint64
	Limit   uint64
	PageCnt int
	Total   int
	Pre     interface{}
	Next    interface{}
}

func (Page) GetPaginator(pageIndex, pageSize, total int) (p Page) {
	if pageIndex <= 0 {
		pageIndex = 1
	}

	if pageSize <= 0 || pageSize > 100000 {
		pageSize = 100000
	}

	p.Total = total
	p.PageCnt = total / pageSize
	if total%pageSize != 0 {
		p.PageCnt += 1
	}

	if pageIndex <= 0 {
		p.Offset = 0
		p.Limit = uint64(pageSize)

	} else {
		p.Offset = uint64(pageSize * (pageIndex - 1))
		p.Limit = uint64(pageSize)
	}

	if pageIndex > 1 {
		p.Pre = pageIndex - 1
	} else {
		p.Pre = nil
	}
	if pageIndex < p.PageCnt {
		p.Next = pageIndex + 1
	} else {
		p.Next = nil
	}

	return p

}

func (Page) GetRangeFromPage(p Page) (start, end uint64) {
	start = p.Offset
	end = p.Offset + p.Limit
	if start > uint64(p.Total) {
		start = uint64(p.Total)
	}
	if end > uint64(p.Total) {
		end = uint64(p.Total)
	}
	return start, end
}

func (p Page) PageData(key string, value interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"count":    p.Total,
		"next":     p.Next,
		"previous": p.Pre,
		key:        value,
	}
	return data
}

func (p Page) PageDataWithExtra(key string, value interface{}, extraValue map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"count":    p.Total,
		"next":     p.Next,
		"previous": p.Pre,
		key:        value,
	}
	for k, v := range extraValue {
		data[k] = v
	}
	return data
}
