package bento

type Constraint interface{}

type (
	ConstraintMin        int
	ConstraintMax        int
	ConstraintLength     int
	ConstraintPercentage int
	ConstraintRatio      struct{ A, B int }
	ConstraintFill       int
)
