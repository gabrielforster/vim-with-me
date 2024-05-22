package chat

type ChatMsg struct {
	Name string
	Msg  string
	Bits int
}

type FilterCB func(msg string) bool
type MapCB func(msg string) string
type Occurrence struct {
	Count int
	Msg   string
}

type ChatAggregator struct {
	mapFn       MapCB
	filter      FilterCB
	occurrences []*Occurrence
	max         *Occurrence
}

var defaultMax = Occurrence{Count: 0, Msg: ""}
var identity = func(x string) string { return x }

func NewChatAggregator() ChatAggregator {
	return ChatAggregator{
		filter:      nil,
		occurrences: make([]*Occurrence, 0),
		max:         &defaultMax,
		mapFn:       identity,
	}
}

func (c ChatAggregator) WithFilter(filterFn FilterCB) ChatAggregator {
	c.filter = filterFn
	return c
}


func (c ChatAggregator) WithMap(mapFn MapCB) ChatAggregator {
	c.mapFn = mapFn
	return c
}

func (c *ChatAggregator) Add(msg string) {
	msg = c.mapFn(msg)
	if !c.filter(msg) {
		return
	}

	for _, occurrence := range c.occurrences {
		if occurrence.Msg == msg {
			occurrence.Count++
			if occurrence.Count > c.max.Count {
				c.max = occurrence
			}
			return
		}
	}

	c.occurrences = append(c.occurrences, &Occurrence{Count: 1, Msg: msg})
}

func (c *ChatAggregator) Reset() Occurrence {
	max := c.max
	c.max = &defaultMax
	c.occurrences = make([]*Occurrence, 0)

	return *max
}

func (c *ChatAggregator) Pipe(ch chan ChatMsg) {
	for msg := range ch {
		c.Add(msg.Msg)
	}
}
