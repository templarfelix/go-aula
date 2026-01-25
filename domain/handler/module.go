package handler

import (
	"go.uber.org/fx"
	"microservice/domain/handler/category"
	"microservice/domain/handler/health"
	"microservice/domain/handler/tag"
)

var Module = fx.Module("handler",
	fx.Provide(tag.ProvideTagHandler),
	fx.Provide(category.ProvideCategoryHandler),
	fx.Provide(health.NewHealthHandler),
	fx.Invoke(tag.RegisterHooks),
	fx.Invoke(category.RegisterHooks),
	fx.Invoke(health.RegisterRoutes),
)
