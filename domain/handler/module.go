package handler

import (
	"go.uber.org/fx"
	"microservice/domain/handler/category"
	"microservice/domain/handler/tag"
)

var Module = fx.Module("handler",
	fx.Provide(tag.ProvideTagHandler),
	fx.Provide(category.ProvideCategoryHandler),
	fx.Invoke(tag.RegisterHooks),
	fx.Invoke(category.RegisterHooks),
)
