package imageUtils

import "mjpeg_multiplexer/src/utils"

func fraction(a, b int) float64 {
	return float64(a) / float64(b)
}

var Slots3 = PanelLayout{
	FirstWidth:     fraction(2, 3),
	FirstHeight:    1,
	ChildrenWidth:  fraction(1, 3),
	ChildrenHeight: fraction(1, 2),
	ChildrenPositions: []utils.FloatPoint{
		{X: fraction(2, 3), Y: 0},
		{X: fraction(2, 3), Y: fraction(1, 2)},
	},

	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: fraction(2, 3), Y: 0}, T2: utils.FloatPoint{X: fraction(2, 3), Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: fraction(2, 3), Y: 0.5}, T2: utils.FloatPoint{X: 1, Y: 0.5}},
	},
}

var Slots4 = PanelLayout{
	FirstWidth:     fraction(3, 4),
	FirstHeight:    1,
	ChildrenWidth:  fraction(1, 4),
	ChildrenHeight: fraction(1, 3),
	ChildrenPositions: []utils.FloatPoint{
		{X: fraction(3, 4), Y: 0},
		{X: fraction(3, 4), Y: fraction(1, 3)},
		{X: fraction(3, 4), Y: fraction(2, 3)},
	},
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: fraction(3, 4), Y: 0}, T2: utils.FloatPoint{X: fraction(3, 4), Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: fraction(3, 4), Y: fraction(1, 3)}, T2: utils.FloatPoint{X: 1, Y: fraction(1, 3)}},
		{T1: utils.FloatPoint{X: fraction(3, 4), Y: fraction(2, 3)}, T2: utils.FloatPoint{X: 1, Y: fraction(2, 3)}},
	},
}

var Slots6 = PanelLayout{
	FirstWidth:     fraction(2, 3),
	FirstHeight:    fraction(2, 3),
	ChildrenWidth:  fraction(1, 3),
	ChildrenHeight: fraction(1, 3),
	ChildrenPositions: []utils.FloatPoint{
		{X: 0, Y: fraction(2, 3)},
		{X: fraction(1, 3), Y: fraction(2, 3)},
		{X: fraction(2, 3), Y: fraction(2, 3)},
		{X: fraction(2, 3), Y: fraction(1, 3)},
		{X: fraction(2, 3), Y: 0},
	},
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: fraction(1, 3), Y: fraction(2, 3)}, T2: utils.FloatPoint{X: fraction(1, 3), Y: 1}},
		{T1: utils.FloatPoint{X: fraction(2, 3), Y: 0}, T2: utils.FloatPoint{X: fraction(2, 3), Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: fraction(2, 3), Y: fraction(1, 3)}, T2: utils.FloatPoint{X: 1, Y: fraction(1, 3)}},
		{T1: utils.FloatPoint{X: 0, Y: fraction(2, 3)}, T2: utils.FloatPoint{X: 1, Y: fraction(2, 3)}},
	},
}

var Slots8 = PanelLayout{
	FirstWidth:     fraction(3, 4),
	FirstHeight:    fraction(3, 4),
	ChildrenWidth:  fraction(1, 4),
	ChildrenHeight: fraction(1, 4),
	ChildrenPositions: []utils.FloatPoint{
		{X: 0, Y: fraction(3, 4)},
		{X: fraction(1, 4), Y: fraction(3, 4)},
		{X: fraction(2, 4), Y: fraction(3, 4)},
		{X: fraction(3, 4), Y: fraction(3, 4)},
		{X: fraction(3, 4), Y: fraction(2, 4)},
		{X: fraction(3, 4), Y: fraction(1, 4)},
		{X: fraction(3, 4), Y: 0},
	},
	VerticalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: fraction(1, 4), Y: fraction(3, 4)}, T2: utils.FloatPoint{X: fraction(1, 4), Y: 1}},
		{T1: utils.FloatPoint{X: fraction(2, 4), Y: fraction(3, 4)}, T2: utils.FloatPoint{X: fraction(2, 4), Y: 1}},
		{T1: utils.FloatPoint{X: fraction(3, 4), Y: 0}, T2: utils.FloatPoint{X: fraction(3, 4), Y: 1}},
	},
	HorizontalBorderPoints: []utils.Tuple[utils.FloatPoint]{
		{T1: utils.FloatPoint{X: fraction(3, 4), Y: fraction(1, 4)}, T2: utils.FloatPoint{X: 1, Y: fraction(1, 4)}},
		{T1: utils.FloatPoint{X: fraction(3, 4), Y: fraction(2, 4)}, T2: utils.FloatPoint{X: 1, Y: fraction(2, 4)}},
		{T1: utils.FloatPoint{X: 0, Y: fraction(3, 4)}, T2: utils.FloatPoint{X: 1, Y: fraction(3, 4)}},
	},
}
