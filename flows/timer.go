package flows

// TimerID represents a single timer
type TimerID int

// TimerCallback is a function the gets called upon a timer event. This event receives the expiry time and the current time.
type TimerCallback func(EventContext, DateTimeNanoseconds)

var timerMaxID TimerID

// RegisterTimer registers a new timer and returns the new TimerID.
func RegisterTimer() TimerID {
	ret := timerMaxID
	timerMaxID++
	return ret
}

var (
	timerIdle   = RegisterTimer()
	timerActive = RegisterTimer()
)

type funcEntry struct {
	function TimerCallback
	context  EventContext
}

type funcEntries []funcEntry

func makeFuncEntries() funcEntries {
	return make(funcEntries, 2)
}

func (fe *funcEntries) expire(when DateTimeNanoseconds) DateTimeNanoseconds {
	var next DateTimeNanoseconds
	fep := *fe
	for i, v := range fep {
		if v.context.When != 0 {
			if v.context.When <= when {
				fep[i].function(v.context, when)
				fep[i].context.When = 0
			} else if next == 0 || v.context.When <= next {
				next = v.context.When
			}
		}
	}
	return next
}

func (fe *funcEntries) addTimer(id TimerID, f TimerCallback, context EventContext) {
	fep := *fe
	if int(id) >= len(fep) {
		fep = append(fep, make(funcEntries, len(fep)-int(id)+1)...)
		*fe = fep
	}
	fep[id].function = f
	fep[id].context = context
}

func (fe *funcEntries) hasTimer(id TimerID) bool {
	fep := *fe
	if int(id) >= len(fep) || id < 0 {
		return false
	}
	if fep[id].context.When == 0 {
		return false
	}
	return true
}
