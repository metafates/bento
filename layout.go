package bento

import (
	"fmt"
	"log"
	"math"

	"github.com/metafates/bento/casso"
)

const _floatPrecisionMultiplier float64 = 100.0

const (
	SpacerSizeEq     casso.Priority = casso.Required / 10.0
	MinSizeGTE       casso.Priority = casso.Strong * 100.0
	MaxSizeLTE       casso.Priority = casso.Strong * 100.0
	MinSizeLTE       casso.Priority = casso.Strong * 100.0
	LengthSizeEq     casso.Priority = casso.Strong * 10.0
	PercentageSizeEq casso.Priority = casso.Strong
	RatioSizeEq      casso.Priority = casso.Strong / 10.0
	MinSizeEq        casso.Priority = casso.Medium * 10.0
	MaxSizeEq        casso.Priority = casso.Medium * 10.0
	Grow             casso.Priority = 100.0
	FillGrow         casso.Priority = casso.Medium
	SpaceGrow        casso.Priority = casso.Weak * 10.0
	AllSegmentGrow   casso.Priority = casso.Weak
)

type Layout struct {
	direction   Direction
	constraints []Constraint
	margin      Margin
	flex        Flex
	spacing     Spacing
}

func (l Layout) Split(area Rect) []Rect {
	segments, _, err := l.split(area)
	if err != nil {
		log.Panicf("failed to split: %w", err)
	}

	return segments
}

func (l Layout) split(area Rect) (segments []Rect, spacers []Rect, err error) {
	solver := casso.NewSolver()

	innerArea := area.Inner(l.margin)

	var areaStart, areaEnd float64

	switch l.direction {
	case DirectionHorizontal:
		areaStart = float64(innerArea.X) * _floatPrecisionMultiplier
		areaEnd = float64(innerArea.Right()) * _floatPrecisionMultiplier
	case DirectionVertical:
		areaStart = float64(innerArea.Y) * _floatPrecisionMultiplier
		areaEnd = float64(innerArea.Bottom()) * _floatPrecisionMultiplier
	}

	variableCount := len(l.constraints)*2 + 2

	variables := make([]casso.Symbol, variableCount)

	for i := 0; i < variableCount; i++ {
		variables[i] = casso.New()
	}

	spacerElements := newElements(variables)
	segmentElements := newElements(variables[1:])

	var spacing int

	switch s := l.spacing.(type) {
	case SpacingSpace:
		spacing = int(s)
	case SpacingOverlap:
		spacing = -int(s)
	}

	areaSize := _Element{
		Start: variables[0],
		End:   variables[1],
	}

	if err := configureArea(solver, areaSize, areaStart, areaEnd); err != nil {
		return nil, nil, fmt.Errorf("configure area: %w", err)
	}

	if err := configureVariableInAreaConstraints(solver, variables, areaSize); err != nil {
		return nil, nil, fmt.Errorf("configure variable in area constraints: %w", err)
	}

	if err := configureVariableConstraints(solver, variables); err != nil {
		return nil, nil, fmt.Errorf("configure variable constraints: %w", err)
	}

	if err := configureFlexConstraints(solver, areaSize, spacerElements, l.flex, spacing); err != nil {
		return nil, nil, fmt.Errorf("configure flex constraints: %w", err)
	}

	if err := configureConstraints(solver, areaSize, segmentElements, l.constraints, l.flex); err != nil {
		return nil, nil, fmt.Errorf("configure constraints: %w", err)
	}

	for i := 0; i < len(segmentElements)-1; i++ {
		left := segmentElements[i]
		right := segmentElements[i+1]

		if err := left.addHasSizeConstraint(solver, right.size(), AllSegmentGrow); err != nil {
			return nil, nil, fmt.Errorf("add has size constraint: %w", err)
		}
	}

	segments = valuesToRects(solver, segmentElements, innerArea, l.direction)
	spacers = valuesToRects(solver, spacerElements, innerArea, l.direction)

	return segments, spacers, nil
}

func valuesToRects(
	solver *casso.Solver,
	elements []_Element,
	area Rect,
	direction Direction,
) []Rect {
	rects := make([]Rect, 0, len(elements))

	for _, e := range elements {
		startRaw := solver.Val(e.Start)
		endRaw := solver.Val(e.End)

		process := func(value float64) int {
			rounded := math.Round(math.Round(value) / _floatPrecisionMultiplier)

			return int(math.Trunc(rounded))
		}

		start := process(startRaw)
		end := process(endRaw)

		size := end - start

		switch direction {
		case DirectionHorizontal:
			rect := Rect{
				X:      start,
				Y:      area.Y,
				Width:  size,
				Height: area.Height,
			}

			rects = append(rects, rect)
		case DirectionVertical:
			rect := Rect{
				X:      area.X,
				Y:      start,
				Width:  area.Width,
				Height: size,
			}

			rects = append(rects, rect)
		}
	}

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

			if err := segment.addHasMaxSizeConstraint(solver, size, MaxSizeLTE); err != nil {
				return fmt.Errorf("add has max size constraint: %w", err)
			}

			if err := segment.addHasIntSizeConstraint(solver, size, MaxSizeEq); err != nil {
				return fmt.Errorf("add has int size constraint: %w", err)
			}
		case ConstraintMin:
			size := int(constraint)

			if err := segment.addHasMinSizeConstraint(solver, size, MinSizeGTE); err != nil {
				return fmt.Errorf("add has min size constraint: %w", err)
			}

			if err := segment.addHasSizeConstraint(solver, area.size(), FillGrow); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}
		case ConstraintLength:
			length := int(constraint)

			if err := segment.addHasIntSizeConstraint(solver, length, LengthSizeEq); err != nil {
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
			if err := s.addHasSizeConstraint(solver, casso.NewExpr(spacingF), SpacerSizeEq); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}

			if len(spacers) >= 2 {
				first := spacers[0]
				last := spacers[len(spacers)-1]

				if err := first.addIsEmptyConstraint(solver); err != nil {
					return fmt.Errorf("add is empty constraint: %w", err)
				}

				if err := last.addHasSizeConstraint(solver, area.size(), Grow); err != nil {
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
	variables []casso.Symbol,
) error {
	variables = variables[1:]

	count := len(variables)

	for i := 0; i < count-count%2; i += 2 {
		left, right := variables[i], variables[i+1]

		if _, err := solver.AddConstraint(casso.NewConstraint(
			casso.LTE,
			0,
			left.T(1),
			right.T(1),
		)); err != nil {
			return fmt.Errorf("add constraint: %w", err)
		}
	}

	return nil
}

func configureVariableInAreaConstraints(
	solver *casso.Solver,
	variables []casso.Symbol,
	area _Element,
) error {
	for _, v := range variables {
		startConstraint := casso.NewConstraint(casso.GTE, 0, v.T(1), area.Start.T(1))
		endConstraint := casso.NewConstraint(casso.LTE, 0, v.T(1), area.End.T(1))

		if _, err := solver.AddConstraint(startConstraint); err != nil {
			return fmt.Errorf("add start constraint: %w", err)
		}

		if _, err := solver.AddConstraint(endConstraint); err != nil {
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
	if _, err := solver.AddConstraint(area.Start.EQ(areaStart)); err != nil {
		return fmt.Errorf("add start constraint: %w", err)
	}

	if _, err := solver.AddConstraint(area.End.EQ(areaEnd)); err != nil {
		return fmt.Errorf("add end constraint: %w", err)
	}

	return nil
}

func newElements(variables []casso.Symbol) []_Element {
	var elements []_Element

	count := len(variables)

	for i := 0; i < count-count%2; i += 2 {
		start, end := variables[i], variables[i+1]

		elements = append(elements, _Element{Start: start, End: end})
	}

	return elements
}

type _Element struct {
	Start, End casso.Symbol
}

func (e _Element) size() casso.Expr {
	return casso.NewExpr(0, e.End.T(1), e.Start.T(-1))
}

func (e _Element) addIsEmptyConstraint(solver *casso.Solver) error {
	_, err := solver.AddConstraintWithPriority(
		casso.Required-1,
		casso.NewConstraint2(
			casso.EQ,
			e.size(),
			casso.NewExpr(0),
		),
	)

	return err
}

func (e _Element) addHasSizeConstraint(
	solver *casso.Solver,
	size casso.Expr,
	priority casso.Priority,
) error {
	_, err := solver.AddConstraintWithPriority(
		priority,
		casso.NewConstraint2(casso.EQ, e.size(), size),
	)

	return err
}

func (e _Element) addHasMaxSizeConstraint(
	solver *casso.Solver,
	size int,
	priority casso.Priority,
) error {
	_, err := solver.AddConstraintWithPriority(
		priority,
		casso.NewConstraint2(
			casso.LTE,
			e.size(),
			casso.NewExpr(float64(size)*_floatPrecisionMultiplier),
		),
	)

	return err
}

func (e _Element) addHasMinSizeConstraint(
	solver *casso.Solver,
	size int,
	priority casso.Priority,
) error {
	_, err := solver.AddConstraintWithPriority(
		priority,
		casso.NewConstraint2(
			casso.GTE,
			e.size(),
			casso.NewExpr(float64(size)*_floatPrecisionMultiplier),
		),
	)

	return err
}

func (e _Element) addHasIntSizeConstraint(
	solver *casso.Solver,
	size int,
	priority casso.Priority,
) error {
	_, err := solver.AddConstraintWithPriority(
		priority,
		casso.NewConstraint2(
			casso.EQ,
			e.size(),
			casso.NewExpr(float64(size)*_floatPrecisionMultiplier),
		),
	)

	return err
}
