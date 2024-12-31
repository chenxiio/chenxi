package api

import (
	"context"

	"github.com/chenxiio/chenxi/cfg"
)

type RecipeApi interface {
	ReadUnitRecipeList(ctx context.Context, t string) ([]string, error)              //perm:none
	ReadProcessRecipeList(ctx context.Context) ([]string, error)                     //perm:none
	ReadUnitRecipe(ctx context.Context, t, name string) (cfg.UnitRecipe, error)      //perm:none
	ReadProcessRecipe(ctx context.Context, name string) (cfg.ProcessRecipe, error)   //perm:none
	SaveUnitRecipe(ctx context.Context, t, name string, rcp cfg.UnitRecipe) error    //perm:none
	SaveProcessRecipe(ctx context.Context, name string, rcp cfg.ProcessRecipe) error //perm:none
}
