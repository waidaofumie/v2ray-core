package router

import (
	"context"

	core "github.com/v2fly/v2ray-core/v4"
	"github.com/v2fly/v2ray-core/v4/app/observatory"
	"github.com/v2fly/v2ray-core/v4/common"
	"github.com/v2fly/v2ray-core/v4/features/extension"
)

type LeastPingStrategy struct {
	ctx         context.Context
	observatory extension.Observatory
}

func (l *LeastPingStrategy) InjectContext(ctx context.Context) {
	common.Must(core.RequireFeatures(ctx, func(observatory extension.Observatory) error {
		l.observatory = observatory
		return nil
	}))
	l.ctx = ctx
}

func (l *LeastPingStrategy) PickOutbound(strings []string) string {
	observeReport, err := l.observatory.GetObservation(l.ctx)
	if err != nil {
		newError("cannot get observe report").Base(err).WriteToLog()
		return ""
	}
	outboundsList := outboundList(strings)
	if result, ok := observeReport.(*observatory.ObservationResult); ok {
		status := result.Status
		leastPing := int64(99999999)
		selectedOutboundName := ""
		for _, v := range status {
			if outboundsList.contains(v.OutboundTag) && v.Alive && v.Delay < leastPing {
				selectedOutboundName = v.OutboundTag
			}
		}
		return selectedOutboundName
	}

	//No way to understand observeReport
	return ""
}

type outboundList []string

func (o outboundList) contains(name string) bool {
	for _, v := range o {
		if v == name {
			return true
		}
	}
	return false
}
