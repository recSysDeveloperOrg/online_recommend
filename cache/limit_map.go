package cache

type limitArray struct {
	cur   int
	limit int
	array []interface{}
}

func makeLimitArray(limit int) *limitArray {
	return &limitArray{
		limit: limit,
	}
}

func (l *limitArray) add(e interface{}) {
	if l.cur == len(l.array) && l.cur < l.limit {
		array := make([]interface{}, maxInt(1, len(l.array)*2))
		copy(array[0:l.cur], l.array[0:l.cur])
		l.array = array
	}

	if l.cur == l.limit {
		copy(l.array[0:l.limit-1], l.array[1:l.limit])
	}
	l.array[minInt(l.cur, l.limit-1)] = e
	l.cur = minInt(l.cur+1, l.limit)
}

func (l *limitArray) find(e interface{}) bool {
	for _, v := range l.array {
		if v == e {
			return true
		}
	}

	return false
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// LimitMap 限制了每一个key对应的value数组的最大长度
type LimitMap struct {
	m         map[interface{}]*limitArray
	limitSize int
}

func NewLimitMap(limitSize int) *LimitMap {
	return &LimitMap{
		m:         make(map[interface{}]*limitArray),
		limitSize: limitSize,
	}
}

func (l *LimitMap) Add(k, v interface{}) {
	if _, ok := l.m[k]; !ok {
		l.m[k] = makeLimitArray(l.limitSize)
	}
	l.m[k].add(v)
}

func (l *LimitMap) CheckExist(k, v interface{}) bool {
	if _, ok := l.m[k]; !ok {
		return false
	}

	return l.m[k].find(v)
}
