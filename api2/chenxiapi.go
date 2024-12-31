package api2

import "github.com/chenxiio/chenxi/api"

type Cxapi interface {
	IOServerAPI
	api.RecipeApi
	api.JobApi
	ALMApi
}
