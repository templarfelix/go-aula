package service

import (
	"go.uber.org/fx"
	"microservice/domain/service/category"
	"microservice/domain/service/tag"
)

var Module = fx.Module("service",
	fx.Provide(tag.ProvideTagService),
	fx.Provide(category.ProvideCategoryService),
)
