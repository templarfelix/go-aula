package repository

import (
	"go.uber.org/fx"
	"microservice/domain/repository/category"
	tag "microservice/domain/repository/tag"
)

var Module = fx.Module("repository",
	fx.Provide(tag.ProvideTagRepository),
	fx.Provide(category.ProvideCategoryRepository),
)
