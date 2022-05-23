package history

func NewQueue(capacity int) *Queue {
	return &Queue{
		changes:     make([]Change, capacity),
		changeIndex: 0,
	}
}

type Queue struct {
	changes     []Change
	changeIndex int
}

func (q *Queue) LastChange() Change {
	return q.changes[q.changeIndex]
}

func (q *Queue) Push(change Change) error {
	if err := change.Apply(); err != nil {
		return err
	}

	if q.changeIndex > 0 {
		// prepend new change
		q.changeIndex--
		q.changes[q.changeIndex] = change
		// erase unpop history
		for i := q.changeIndex - 1; i >= 0; i-- {
			q.changes[i] = nil
		}
	} else {
		// shift changes right
		copy(q.changes[1:], q.changes)
		// prepend new change
		q.changes[0] = change
	}
	return nil
}

func (q *Queue) Pop() error {
	if err := q.changes[q.changeIndex].Revert(); err != nil {
		return err
	}
	// move to previous change
	q.changeIndex++
	return nil
}

func (q *Queue) Unpop() error {
	if err := q.changes[q.changeIndex-1].Apply(); err != nil {
		return err
	}
	// move to next change
	q.changeIndex--
	return nil
}

func (q *Queue) CanUnpop() bool {
	nextChange := q.changeIndex - 1
	return (nextChange >= 0) && (q.changes[nextChange] != nil)
}

func (q *Queue) CanPop() bool {
	previousIndex := q.changeIndex + 1
	return (q.changes[q.changeIndex] != nil) && (previousIndex < len(q.changes))
}
