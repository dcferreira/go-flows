package packet

import (
	"github.com/CN-TU/go-flows/flows"
	"github.com/google/gopacket/layers"
)

type tcpFlow struct {
	flows.BaseFlow
	srcFIN, dstFIN, dstACK, srcACK bool
}

type uniFlow struct {
	flows.BaseFlow
}

// NewFlow creates a new flow based on a given event, table, key, context, and flow-id
//
// Depending on the event this will either be a tcp flow, or a standard flow
func NewFlow(event flows.Event, table *flows.FlowTable, key string, lowToHigh bool, context *flows.EventContext, id uint64) flows.Flow {
	if table.FiveTuple() {
		tp := event.(Buffer).TransportLayer()
		if tp != nil && tp.LayerType() == layers.LayerTypeTCP {
			ret := new(tcpFlow)
			ret.Init(table, key, lowToHigh, context, id)
			return ret
		}
	}
	ret := new(uniFlow)
	ret.Init(table, key, lowToHigh, context, id)
	return ret
}

func (flow *tcpFlow) Event(event flows.Event, context *flows.EventContext) {
	flow.BaseFlow.Event(event, context)
	if !flow.Active() {
		return
	}
	buffer := event.(Buffer)
	tcp := buffer.TransportLayer().(*layers.TCP)
	if tcp.RST {
		flow.Export(flows.FlowEndReasonEnd, context, context.When())
		return
	}
	if context.Forward() {
		if tcp.FIN {
			flow.srcFIN = true
		}
		if flow.dstFIN && tcp.ACK {
			flow.dstACK = true
		}
	} else {
		if tcp.FIN {
			flow.dstFIN = true
		}
		if flow.srcFIN && tcp.ACK {
			flow.srcACK = true
		}
	}

	if flow.srcFIN && flow.srcACK && flow.dstFIN && flow.dstACK {
		flow.Export(flows.FlowEndReasonEnd, context, context.When())
	}
}
