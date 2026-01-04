package process

import (
	"context"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func ShowLoading(ctx context.Context) {
	wailsRuntime.EventsEmit(ctx, "setLoadingStatus", true)
}

func HideLoading(ctx context.Context) {
	wailsRuntime.EventsEmit(ctx, "setLoadingStatus", false)
}
