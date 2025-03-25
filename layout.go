package bento

import (
	"fmt"
	"math"

	"github.com/metafates/bento/internal/casso"
)

const _floatPrecisionMultiplier float64 = 100.0

const (
	_spacerSizeEq     casso.Strength = casso.Required / 10.0
	_minSizeGTE       casso.Strength = casso.Strong * 100.0
	_maxSizeLTE       casso.Strength = casso.Strong * 100.0
	_minSizeLTE       casso.Strength = casso.Strong * 100.0
	_lengthSizeEq     casso.Strength = casso.Strong * 10.0
	_percentageSizeEq casso.Strength = casso.Strong
	_ratioSizeEq      casso.Strength = casso.Strong / 10.0
	_minSizeEq        casso.Strength = casso.Medium * 10.0
	_maxSizeEq        casso.Strength = casso.Medium * 10.0
	_grow             casso.Strength = 100.0
	_fillGrow         casso.Strength = casso.Medium
	_spaceGrow        casso.Strength = casso.Weak * 10.0
	_allSegmentGrow   casso.Strength = casso.Weak
)

type Splitted []Rect

func (s Splitted) Assign(areas ...*Rect) {
	for i := range areas {
		if areas[i] != nil {
			*areas[i] = s[i]
		}
	}
}

func (s Splitted) Unwrap() Rect {
	return s[0]
}

func (s Splitted) Unwrap2() (a, b Rect) {
	return s[0], s[1]
}

func (s Splitted) Unwrap3() (a, b, c Rect) {
	return s[0], s[1], s[2]
}

func (s Splitted) Unwrap4() (a, b, c, d Rect) {
	return s[0], s[1], s[2], s[3]
}

type Layout struct {
	Direction   Direction
	Constraints []Constraint
	Padding     Padding
	Flex        Flex
	Spacing     Spacing
}

// NewLayout returns vertical layout with the given constraints
func NewLayout(constraints ...Constraint) Layout {
	return Layout{
		Direction:   DirectionVertical,
		Constraints: constraints,
		Padding:     NewPadding(),
		Flex:        FlexLegacy,
		Spacing:     SpacingSpace(0),
	}
}

func (l Layout) Vertical() Layout {
	return l.WithDirection(DirectionVertical)
}

func (l Layout) Horizontal() Layout {
	return l.WithDirection(DirectionHorizontal)
}

func (l Layout) WithDirection(direction Direction) Layout {
	l.Direction = direction
	return l
}

func (l Layout) WithPadding(padding Padding) Layout {
	l.Padding = padding
	return l
}

func (l Layout) WithFlex(flex Flex) Layout {
	l.Flex = flex
	return l
}

func (l Layout) WithSpacing(spacing Spacing) Layout {
	l.Spacing = spacing
	return l
}

func (l Layout) WithConstraints(constraints ...Constraint) Layout {
	l.Constraints = append(l.Constraints, constraints...)
	return l
}

func (l Layout) SplitWithSpacers(area Rect) (segments, spacers Splitted) {
	segments, spacers, err := l.split(area)
	if err != nil {
		panic(err)
	}

	return segments, spacers
}

func (l Layout) Split(area Rect) Splitted {
	segments, _ := l.SplitWithSpacers(area)

	return segments
}

func (l Layout) split(area Rect) (segments []Rect, spacers []Rect, err error) {
	solver := casso.NewSolver()

	innerArea := area.Inner(l.Padding)

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
		End:   variables[len(variables)-1],
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

	if err := configureFillConstraints(&solver, segmentElements, l.Constraints, l.Flex); err != nil {
		return nil, nil, fmt.Errorf("configure fill constraints: %w", err)
	}

	if l.Flex != FlexLegacy {
		for i := 0; i < len(segmentElements)-1; i++ {
			left := segmentElements[i]
			right := segmentElements[i+1]

			if err := solver.AddConstraint(left.hasSize(right.size(), _allSegmentGrow)); err != nil {
				return nil, nil, fmt.Errorf("add has size constraint: %w", err)
			}
		}
	}

	fetched := solver.FetchChanges()

	changes := make(map[casso.Variable]float64, len(fetched))
	for _, c := range fetched {
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

	for _, e := range elements {
		start := changes[e.Start]
		end := changes[e.End]

		startRounded := int(math.Round(math.Round(start) / _floatPrecisionMultiplier))
		endRounded := int(math.Round(math.Round(end) / _floatPrecisionMultiplier))

		size := max(0, endRounded-startRounded)

		switch direction {
		case DirectionHorizontal:
			rect := Rect{
				Position: Position{
					X: startRounded,
					Y: area.Y,
				},
				Width:  size,
				Height: area.Height,
			}

			rects = append(rects, rect)
		case DirectionVertical:
			rect := Rect{
				Position: Position{
					X: area.X,
					Y: startRounded,
				},
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
	var (
		validConstraints []Constraint
		validSegments    []_Element
	)

	for i := 0; i < min(len(constraints), len(segments)); i++ {
		c := constraints[i]
		s := segments[i]

		switch c.(type) {
		case ConstraintFill, ConstraintMin:
			if _, ok := c.(ConstraintMin); ok && flex == FlexLegacy {
				continue
			}

			validConstraints = append(validConstraints, c)
			validSegments = append(validSegments, s)
		}
	}

	if len(validConstraints) < 2 {
		return nil
	}

	for _, indices := range combinations(len(validConstraints), 2) {
		i, j := indices[0], indices[1]

		leftConstraint := validConstraints[i]
		leftSegment := validSegments[i]

		rightConstraint := validConstraints[j]
		rightSegment := validSegments[j]

		getScalingFactor := func(c Constraint) float64 {
			var scalingFactor float64

			switch c := c.(type) {
			case ConstraintFill:
				scale := float64(c)

				scalingFactor = 1e-6
				scalingFactor = max(scalingFactor, scale)

			case ConstraintMin:
				scalingFactor = 1
			}

			return scalingFactor
		}

		leftScalingFactor := getScalingFactor(leftConstraint)
		rightScalingFactor := getScalingFactor(rightConstraint)

		lhs := leftSegment.size().MulConstant(rightScalingFactor)
		rhs := rightSegment.size().MulConstant(leftScalingFactor)

		constraint := casso.Equal(_grow).ExpressionLHS(lhs).ExpressionRHS(rhs)
		if err := solver.AddConstraint(constraint); err != nil {
			return fmt.Errorf("add constraint: %w", err)
		}
	}

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

			err := solver.AddConstraints(
				segment.hasMaxSize(size, _maxSizeLTE),
				segment.hasIntSize(size, _maxSizeEq),
			)
			if err != nil {
				return fmt.Errorf("add constraints: %w", err)
			}

		case ConstraintMin:
			size := int(constraint)

			if err := solver.AddConstraint(segment.hasMinSize(size, _minSizeGTE)); err != nil {
				return fmt.Errorf("add has min size constraint: %w", err)
			}

			if flex == FlexLegacy {
				if err := solver.AddConstraint(segment.hasIntSize(size, _minSizeEq)); err != nil {
					return fmt.Errorf("add has size constraint: %w", err)
				}
			} else {
				if err := solver.AddConstraint(segment.hasSize(area.size(), _fillGrow)); err != nil {
					return fmt.Errorf("add has size constraint: %w", err)
				}
			}

		case ConstraintLen:
			length := int(constraint)

			if err := solver.AddConstraint(segment.hasIntSize(length, _lengthSizeEq)); err != nil {
				return fmt.Errorf("add has int size constraint: %w", err)
			}

		case ConstraintPercentage:
			size := area.size().MulConstant(float64(constraint)).DivConstant(100)

			if err := solver.AddConstraint(segment.hasSize(size, _percentageSizeEq)); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}

		case ConstraintRatio:
			size := area.size().MulConstant(float64(constraint.Num)).DivConstant(float64(max(1, constraint.Den)))

			if err := solver.AddConstraint(segment.hasSize(size, _ratioSizeEq)); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}

		case ConstraintFill:
			if err := solver.AddConstraint(segment.hasSize(area.size(), _fillGrow)); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}
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
	case FlexLegacy:
		for _, s := range spacersExceptFirstAndLast {
			if err := solver.AddConstraint(s.hasSize(casso.NewExpressionFromConstant(spacingF), _spacerSizeEq)); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}
		}

		if len(spacers) >= 2 {
			first, last := spacers[0], spacers[len(spacers)-1]

			err := solver.AddConstraints(first.isEmpty(), last.isEmpty())
			if err != nil {
				return fmt.Errorf("add constraints: %w", err)
			}
		}
	case FlexSpaceAround:
		if len(spacersExceptFirstAndLast) >= 2 {
			for _, indices := range combinations(len(spacersExceptFirstAndLast), 2) {
				i, j := indices[0], indices[1]

				left, right := spacersExceptFirstAndLast[i], spacersExceptFirstAndLast[j]

				if err := solver.AddConstraint(left.hasSize(right.size(), _spacerSizeEq)); err != nil {
					return fmt.Errorf("add has size constraint: %w", err)
				}
			}
		}

		for _, s := range spacersExceptFirstAndLast {
			err := solver.AddConstraints(
				s.hasMinSize(spacing, _spacerSizeEq),
				s.hasSize(area.size(), _spaceGrow),
			)
			if err != nil {
				return fmt.Errorf("add constraints: %w", err)
			}

		}

	case FlexSpaceBetween:
		if len(spacersExceptFirstAndLast) >= 2 {
			for _, indices := range combinations(len(spacersExceptFirstAndLast), 2) {
				i, j := indices[0], indices[1]

				left, right := spacersExceptFirstAndLast[i], spacersExceptFirstAndLast[j]

				if err := solver.AddConstraint(left.hasSize(right.size(), _spacerSizeEq)); err != nil {
					return fmt.Errorf("add has size constraint: %w", err)
				}
			}
		}

		for _, s := range spacersExceptFirstAndLast {
			err := solver.AddConstraints(
				s.hasMinSize(spacing, _spacerSizeEq),
				s.hasSize(area.size(), _spaceGrow),
			)
			if err != nil {
				return fmt.Errorf("add constraints: %w", err)
			}
		}

		if len(spacers) >= 2 {
			first, last := spacers[0], spacers[len(spacers)-1]

			err := solver.AddConstraints(first.isEmpty(), last.isEmpty())
			if err != nil {
				return fmt.Errorf("add constraints: %w", err)
			}

		}
	case FlexStart:
		for _, s := range spacersExceptFirstAndLast {
			if err := solver.AddConstraint(s.hasSize(casso.NewExpressionFromConstant(spacingF), _spacerSizeEq)); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}
		}

		if len(spacers) >= 2 {
			first := spacers[0]
			last := spacers[len(spacers)-1]

			err := solver.AddConstraints(
				first.isEmpty(),
				last.hasSize(area.size(), _grow),
			)
			if err != nil {
				return fmt.Errorf("add constraints: %w", err)
			}
		}

	case FlexCenter:
		for _, s := range spacersExceptFirstAndLast {
			constraint := s.hasSize(casso.NewExpressionFromConstant(spacingF), _spacerSizeEq)

			if err := solver.AddConstraint(constraint); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}
		}

		if len(spacers) >= 2 {
			first, last := spacers[0], spacers[len(spacers)-1]

			err := solver.AddConstraints(
				first.hasSize(area.size(), _grow),
				last.hasSize(area.size(), _grow),
				first.hasSize(last.size(), _spacerSizeEq),
			)
			if err != nil {
				return fmt.Errorf("add constraints: %w", err)
			}
		}
	case FlexEnd:
		for _, s := range spacersExceptFirstAndLast {
			if err := solver.AddConstraint(s.hasSize(casso.NewExpressionFromConstant(spacingF), _spacerSizeEq)); err != nil {
				return fmt.Errorf("add has size constraint: %w", err)
			}
		}

		if len(spacers) >= 2 {
			first := spacers[0]
			last := spacers[len(spacers)-1]

			err := solver.AddConstraints(
				last.isEmpty(),
				first.hasSize(area.size(), _grow),
			)
			if err != nil {
				return fmt.Errorf("add constraints: %w", err)
			}
		}
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

		constraint := casso.LessThanEqual(casso.Required).VariableLHS(left).VariableRHS(right)

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
		start := casso.GreaterThanEqual(casso.Required).VariableLHS(v).VariableRHS(area.Start)
		end := casso.LessThanEqual(casso.Required).VariableLHS(v).VariableRHS(area.End)

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
	startConstraint := casso.Equal(casso.Required).VariableLHS(area.Start).ConstantRHS(areaStart)
	endConstraint := casso.Equal(casso.Required).VariableLHS(area.End).ConstantRHS(areaEnd)

	if err := solver.AddConstraint(startConstraint); err != nil {
		return fmt.Errorf("add start constraint: %w", err)
	}

	if err := solver.AddConstraint(endConstraint); err != nil {
		return fmt.Errorf("add end constraint: %w", err)
	}

	return nil
}

func newElements(variables []casso.Variable) []_Element {
	count := len(variables)

	elements := make([]_Element, 0, count/2+1)

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
		Equal(casso.Required - 1).
		ExpressionLHS(e.size()).
		ConstantRHS(0)
}

func (e _Element) hasSize(
	size casso.Expression,
	strength casso.Strength,
) casso.Constraint {
	return casso.
		Equal(strength).
		ExpressionLHS(e.size()).
		ExpressionRHS(size)
}

func (e _Element) hasMaxSize(
	size int,
	strength casso.Strength,
) casso.Constraint {
	return casso.
		LessThanEqual(strength).
		ExpressionLHS(e.size()).
		ConstantRHS(float64(size) * _floatPrecisionMultiplier)
}

func (e _Element) hasMinSize(
	size int,
	strength casso.Strength,
) casso.Constraint {
	return casso.
		GreaterThanEqual(strength).
		ExpressionLHS(e.size()).
		ConstantRHS(float64(size) * _floatPrecisionMultiplier)
}

func (e _Element) hasIntSize(
	size int,
	strength casso.Strength,
) casso.Constraint {
	return casso.
		Equal(strength).
		ExpressionLHS(e.size()).
		ConstantRHS(float64(size) * _floatPrecisionMultiplier)
}

func combinations(n, k int) [][]int {
	combins := binomial(n, k)
	data := make([][]int, combins)
	if len(data) == 0 {
		return data
	}

	data[0] = make([]int, k)
	for i := range data[0] {
		data[0][i] = i
	}

	for i := 1; i < combins; i++ {
		next := make([]int, k)
		copy(next, data[i-1])
		nextCombination(next, n, k)
		data[i] = next
	}

	return data
}

func nextCombination(s []int, n, k int) {
	for j := k - 1; j >= 0; j-- {
		if s[j] == n+j-k {
			continue
		}
		s[j]++
		for l := j + 1; l < k; l++ {
			s[l] = s[j] + l - j
		}
		break
	}
}

func binomial(n, k int) int {
	if n < 0 || k < 0 {
		panic("negative input")
	}
	if n < k {
		panic("bad set size")
	}
	// (n,k) = (n, n-k)
	if k > n/2 {
		k = n - k
	}
	b := 1
	for i := 1; i <= k; i++ {
		b = (n - k + i) * b / i
	}
	return b
}
