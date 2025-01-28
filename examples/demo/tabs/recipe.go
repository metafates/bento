package tabs

import (
	"fmt"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/clearwidget"
	"github.com/metafates/bento/examples/demo/theme"
	"github.com/metafates/bento/paragraphwidget"
	"github.com/metafates/bento/scrollwidget"
	"github.com/metafates/bento/textwidget"
)

var _ bento.StatefulWidget[RecipeTabState] = (*RecipeTab)(nil)

type RecipeIngredient struct {
	name     string
	quantity string
}

type Recipe struct {
	steps       []string
	ingredients []RecipeIngredient
}

// https://www.realsimple.com/food-recipes/browse-all-recipes/salmon-nigiri
var SalmonNigiriRecipe = Recipe{
	steps: []string{
		"Stir vinegar, sugar, and salt in a small bowl until dissolved.",
		"Thoroughly rinse rice with running water and then cook according to package directions.",
		"Spread rice in a large glass baking dish, drizzle with vinegar mixture, and then gently fold in to incorporate.",
		"Form rice mixture into 16 (3-by-1-inch) pieces. Spread wasabi (to taste) on the underside of each piece of salmon. Lay the salmon on the rice. Serve with soy sauce or tamari for drizzling on the salmon.",
	},
	ingredients: []RecipeIngredient{
		{
			name:     "unseasoned rice vinegar",
			quantity: "¼ cup",
		},
		{
			name:     "granulated sugar",
			quantity: "1 tablespoon",
		},
		{
			name:     "kosher salt",
			quantity: "1 teaspoon",
		},
		{
			name:     "uncooked sushi rice",
			quantity: "1 cup",
		},
		{
			name:     "sushi-grade salmon, thinly sliced",
			quantity: "¾ pounds",
		},
		{
			name:     "Wasabi paste",
			quantity: "",
		},
		{
			name:     "Soy sauce or tamari",
			quantity: "",
		},
	},
}

type RecipeTabState struct {
	rowIndex int
}

func NewRecipeTabState() RecipeTabState {
	return RecipeTabState{}
}

type RecipeTab struct {
	recipe Recipe
}

func NewRecipeTab(recipe Recipe) RecipeTab {
	return RecipeTab{
		recipe: recipe,
	}
}

func (r *RecipeTab) RenderStateful(area bento.Rect, buffer *bento.Buffer, state RecipeTabState) {
	area = area.Inner(bento.NewPadding(1, 2))
	clearwidget.New().Render(area, buffer)

	title := blockwidget.NewTitle(
		textwidget.
			NewLineStr("Salmon Nigiri Recipe").
			WithStyle(bento.NewStyle().Bold().White()),
	)

	blockwidget.
		New().
		WithTitle(title).
		WithTitlesAlignment(bento.AlignmentCenter).
		WithStyle(theme.Global.Content).
		WithPadding(bento.NewPadding(1, 1, 2, 1)).
		Render(area, buffer)

	scrollbarArea := bento.Rect{
		X:      area.X,
		Y:      area.Y + 2,
		Width:  area.Width,
		Height: area.Height - 3,
	}

	r.renderScrollbar(state.rowIndex, scrollbarArea, buffer)

	area = area.Inner(bento.NewPadding(1, 2))

	var recipe, ingredients bento.Rect

	bento.
		NewLayout(
			bento.ConstraintLen(44),
			bento.ConstraintMin(0),
		).
		Horizontal().
		Split(area).
		Assign(&recipe, &ingredients)

	r.renderRecipe(recipe, buffer)
	r.renderIngredients(ingredients, buffer)
}

func (r *RecipeTab) renderRecipe(area bento.Rect, buffer *bento.Buffer) {
	lines := make([]textwidget.Line, 0, len(r.recipe.steps))

	for step, text := range r.recipe.steps {
		stepSpan := textwidget.
			NewSpan(fmt.Sprintf("Step %d: ", step+1)).
			WithStyle(bento.NewStyle().White().Bold())

		textSpan := textwidget.NewSpan(text).WithStyle(bento.NewStyle().White().Italic())

		lines = append(lines, textwidget.NewLine(stepSpan, textSpan))
	}

	paragraphwidget.
		New(textwidget.NewText(lines...)).
		WithWrap(paragraphwidget.NewWrap().WithTrim(true)).
		WithBlock(blockwidget.New().WithPadding(bento.NewPadding(0, 1, 0, 0))).
		Render(area, buffer)
}

func (r *RecipeTab) renderIngredients(area bento.Rect, buffer *bento.Buffer) {
	textwidget.NewTextStr("TODO: Ingredients here").Render(area, buffer)
}

func (r *RecipeTab) renderScrollbar(position int, area bento.Rect, buffer *bento.Buffer) {
	state := scrollwidget.NewState(len(r.recipe.ingredients))
	state.SetViewportContentLen(6)
	state.SetPosition(position)

	scrollwidget.New(scrollwidget.OrientationVerticalRight).RenderStateful(area, buffer, state)
}
