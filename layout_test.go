package bento

// func letters(t *testing.T, flex Flex, constraints []Constraint, width int, expected string) {
// 	area := Rect{Width: width, Height: 1}
//
// 	layout := Layout{
// 		direction:   DirectionHorizontal,
// 		constraints: constraints,
// 		flex:        flex,
// 		spacing:     SpacingSpace(0),
// 	}.Split(area)
//
// 	buffer := NewBufferEmpty(area)
//
// 	latin := []rune("abcdefghijklmnopqrstuvwxyz")
//
// 	for i := 0; i < min(len(constraints), len(layout)); i++ {
// 		c := latin[i]
// 		area := layout[i]
//
// 		s := strings.Repeat(string(c), area.Width)
// 	}
//
// 	_ = layout
// 	_ = area
// }
