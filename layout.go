package bento

import (
	"fmt"

	"github.com/metafates/bento/casso"
)

const _floatPrecisionMultiplier float64 = 100.0

const (
	SpacerSizeEq     casso.Strength = casso.Required / 10.0
	MinSizeGTE       casso.Strength = casso.Strong * 100.0
	MaxSizeLTE       casso.Strength = casso.Strong * 100.0
	MinSizeLTE       casso.Strength = casso.Strong * 100.0
	LengthSizeEq     casso.Strength = casso.Strong * 10.0
	PercentageSizeEq casso.Strength = casso.Strong
	RatioSizeEq      casso.Strength = casso.Strong / 10.0
	MinSizeEq        casso.Strength = casso.Medium * 10.0
	MaxSizeEq        casso.Strength = casso.Medium * 10.0
	Grow             casso.Strength = 100.0
	FillGrow         casso.Strength = casso.Medium
	SpaceGrow        casso.Strength = casso.Weak * 10.0
	AllSegmentGrow   casso.Strength = casso.Weak
)

type Layout struct {
	Direction   Direction
	Constraints []Constraint
	Margin      Margin
	Flex        Flex
	Spacing     Spacing
}

func (l Layout) Split(area Rect) []Rect {
	segments, _, err := l.split(area)
	if err != nil {
		panic(err)
	}

	return segments
}

func (l Layout) split(area Rect) (segments []Rect, spacers []Rect, err error) {
	solver := casso.NewSolver()

	innerArea := area.Inner(l.Margin)

	var areaStart, areaEnd float64

	switch l.Direction {
	case DirectionHorizontal:
		areaStart = float64(innerArea.X) * _floatPrecisionMultiplier
		areaEnd = float64(innerArea.Right()) * _floatPrecisionMultiplier
	case DirectionVertical:
		areaStart = float64(innerArea.Y) * _floatPrecisionMultiplier
		areaEnd = float64(innerArea.Bottom()) * _floatPrecisionMultiplier
	}

	variableCount := len(l.Constraints)*2 + 2

	variables := make([]casso.Variable, variableCount)

	for i := 0; i < variableCount; i++ {
		variables[i] = casso.NewVariable()
	}

	spacerElements := newElements(variables)
	segmentElements := newElements(variables[1:])

	var spacing int

	switch s := l.Spacing.(type) {
	case SpacingSpace:
		spacing = int(s)
	case SpacingOverlap:
		spacing = -int(s)
	}

	areaSize := _Element{
		Start: variables[0],
		End:   variables[1],
	}

	if err := configureArea(&solver, areaSize, areaStart, areaEnd); err != nil {
		return nil, nil, fmt.Errorf("configure area: %w", err)
	}

	if err := configureVariableInAreaConstraints(&solver, variables, areaSize); err != nil {
		return nil, nil, fmt.Errorf("configure variable in area constraints: %w", err)
	}

	if err := configureVariableConstraints(&solver, variables); err != nil {
		return nil, nil, fmt.Errorf("configure variable constraints: %w", err)
	}

	if err := configureFlexConstraints(&solver, areaSize, spacerElements, l.Flex, spacing); err != nil {
		return nil, nil, fmt.Errorf("configure flex constraints: %w", err)
	}

	if err := configureConstraints(&solver, areaSize, segmentElements, l.Constraints, l.Flex); err != nil {
		return nil, nil, fmt.Errorf("configure constraints: %w", err)
	}

	for i := 0; i < len(segmentElements)-1; i++ {
		left := segmentElements[i]
		right := segmentElements[i+1]

		if err := solver.AddConstraint(left.hasSize(right.size(), AllSegmentGrow)); err != nil {
			return nil, nil, fmt.Errorf("add has size constraint: %w", err)
		}
	}

	changes := make(map[casso.Variable]float64)

	for _, c := range solver.FetchChanges() {
		changes[c.Variable] = c.Constant
	}

	segments = changesToRects(changes, segmentElements, innerArea, l.Direction)
	spacers = changesToRects(changes, spacerElements, innerArea, l.Direction)

	return segments, spacers, nil
}

func changesToRects(
	changes map[casso.Variable]float64,
	elements []_Element,
	area Rect,
	direction Direction,
) []Rect {
	var rects []Rect

	// TODO

	return rects
}

func configureFillConstraints(
	solver *casso.Solver,
	segments []_Element,
	constraints []Constraint,
	flex Flex,
) error {
	// TODO
	return nil
}

func configureConstraints(
	solver *casso.Solver,
	area _Element,
	segments []_Element,
	constraints []Constraint,
	flex Flex,
) error {
	for i := 0; i < min(len(constraints), len(segments)); i++ {
		constraint := constraints[i]
		segment := segments[i]

		switch constraint := constraint.(type) {
		case ConstraintMax:
			size := int(constraint)

			if err := solver.AddConstraint(segment.hasMaxSize(size, MaxSizeLTE)); err != nil {
				return fmt.Errorf("add has max size constraint: %w", err)
			}

			if err := solver.AddConstraint(segment.hasIntSize(size, MaxSizeEq)); err != nil {
				return fmt.Errorf("add has int size constraint: %w", err)
			}
		case ConstraintMin:
			size := int(constraint)

			if err := solver.AddConstraint(segment.hasMinSize(size, MinSizeGTE)); err != nil {
				return fmt.Errorf("add has min size constraint: %w", err)
			}

			if err := solver.AddConstraint(segment.hasSize(area.size(), FillGrow)); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}
		case ConstraintLength:
			length := int(constraint)

			if err := solver.AddConstraint(segment.hasIntSize(length, LengthSizeEq)); err != nil {
				return fmt.Errorf("add has int size constraint: %w", err)
			}
		case ConstraintPercentage:
			panic("unimplemented")
		case ConstraintRatio:
			panic("unimplemented")
		case ConstraintFill:
			panic("unimplemented")
		}
	}

	return nil
}

func configureFlexConstraints(
	solver *casso.Solver,
	area _Element,
	spacers []_Element,
	flex Flex,
	spacing int,
) error {
	var spacersExceptFirstAndLast []_Element

	if len(spacers) > 2 {
		spacersExceptFirstAndLast = spacers[1 : len(spacers)-1]
	}

	spacingF := float64(spacing) * _floatPrecisionMultiplier

	switch flex {
	case FlexSpaceAround:
		panic("not implemented")
	case FlexSpaceBetween:
		panic("not implemented")
	case FlexStart:
		for _, s := range spacersExceptFirstAndLast {
			if err := solver.AddConstraint(s.hasSize(casso.NewExpressionFromConstant(spacingF), SpacerSizeEq)); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}

			if len(spacers) >= 2 {
				first := spacers[0]
				last := spacers[len(spacers)-1]

				if err := solver.AddConstraint(first.isEmpty()); err != nil {
					return fmt.Errorf("add is empty constraint: %w", err)
				}

				if err := solver.AddConstraint(last.hasSize(area.size(), Grow)); err != nil {
					return fmt.Errorf("add has size constraint: %w", err)
				}
			}
		}

	case FlexCenter:
		panic("not implemented")
	case FlexEnd:
		panic("not implemented")
	}

	return nil
}

func configureVariableConstraints(
	solver *casso.Solver,
	variables []casso.Variable,
) error {
	variables = variables[1:]

	count := len(variables)

	for i := 0; i < count-count%2; i += 2 {
		left, right := variables[i], variables[i+1]

		constraint := casso.LessThanEqual(casso.Required).WithVariable(left).WithVariable(right)

		if err := solver.AddConstraint(constraint); err != nil {
			return fmt.Errorf("add constraint: %w", err)
		}
	}

	return nil
}

func configureVariableInAreaConstraints(
	solver *casso.Solver,
	variables []casso.Variable,
	area _Element,
) error {
	for _, v := range variables {
		start := casso.GreaterThanEqual(casso.Required).WithVariable(v).WithVariable(area.Start)
		end := casso.GreaterThanEqual(casso.Required).WithVariable(v).WithVariable(area.End)

		if err := solver.AddConstraint(start); err != nil {
			return fmt.Errorf("add start constraint: %w", err)
		}

		if err := solver.AddConstraint(end); err != nil {
			return fmt.Errorf("add end constraint: %w", err)
		}
	}

	return nil
}

func configureArea(
	solver *casso.Solver,
	area _Element,
	areaStart, areaEnd float64,
) error {
	startConstraint := casso.Equal(casso.Required).WithVariable(area.Start).WithConstant(areaStart)
	endConstraint := casso.Equal(casso.Required).WithVariable(area.End).WithConstant(areaEnd)

	if err := solver.AddConstraint(startConstraint); err != nil {
		return fmt.Errorf("add start constraint: %w", err)
	}

	if err := solver.AddConstraint(endConstraint); err != nil {
		return fmt.Errorf("add end constraint: %w", err)
	}

	return nil
}

func newElements(variables []casso.Variable) []_Element {
	var elements []_Element

	count := len(variables)

	for i := 0; i < count-count%2; i += 2 {
		start, end := variables[i], variables[i+1]

		elements = append(elements, _Element{Start: start, End: end})
	}

	return elements
}

type _Element struct {
	Start, End casso.Variable
}

func newElement() _Element {
	return _Element{
		Start: casso.NewVariable(),
		End:   casso.NewVariable(),
	}
}

func (e _Element) size() casso.Expression {
	return e.End.Sub(e.Start)
}

func (e _Element) isEmpty() casso.Constraint {
	return casso.
		Equal(casso.Required).
		WithExpression(e.size()).
		WithConstant(0)
}

func (e _Element) hasSize(
	size casso.Expression,
	strength casso.Strength,
) casso.Constraint {
	return casso.
		Equal(strength).
		WithExpression(e.size()).
		WithExpression(size)
}

func (e _Element) hasMaxSize(
	size int,
	strength casso.Strength,
) casso.Constraint {
	return casso.
		LessThanEqual(strength).
		WithExpression(e.size()).
		WithConstant(float64(size) * _floatPrecisionMultiplier)
}

func (e _Element) hasMinSize(
	size int,
	strength casso.Strength,
) casso.Constraint {
	return casso.
		GreaterThanEqual(strength).
		WithExpression(e.size()).
		WithConstant(float64(size) * _floatPrecisionMultiplier)
}

func (e _Element) hasIntSize(
	size int,
	strength casso.Strength,
) casso.Constraint {
	return casso.
		Equal(strength).
		WithExpression(e.size()).
		WithConstant(float64(size) * _floatPrecisionMultiplier)
}
