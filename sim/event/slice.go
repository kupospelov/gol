package event

type Cells []Cell

func (events Cells) Len() int {
	return len(events)
}

func (events Cells) Swap(i, j int) {
	events[i], events[j] = events[j], events[i]
}

func (events Cells) Less(i, j int) bool {
	if events[i].X < events[j].X {
		return true
	} else {
		if events[i].X > events[j].X {
			return false
		}
	}

	if events[i].Y < events[j].Y {
		return true
	} else {
		if events[i].Y > events[j].Y {
			return false
		}
	}

	return !events[i].Alive
}
