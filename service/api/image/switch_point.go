package image

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	"fmt"

	log "github.com/inconshreveable/log15"
)

type PointsRep struct {
	LeftEyeBrow  []*imagemodel.Point `json:"leftEyebrow"`
	RightEyeBrow []*imagemodel.Point `json:"rightEyebrow"`
	LeftEye      []*imagemodel.Point `json:"leftEye"`
	RightEye     []*imagemodel.Point `json:"rightEye"`
	LeftEar      []*imagemodel.Point `json:"leftEar"`
	RightEar     []*imagemodel.Point `json:"rightEar"`
	Mouth        []*imagemodel.Point `json:"mouth"`
	Nouse        []*imagemodel.Point `json:"nouse"`
	Face         []*imagemodel.Point `json:"face"`
}

func SwitchPoint(image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch point", image.Md5))
	//	fmt.Println(image.ThrFaces["deepir_import"])
	points := make([]*imagemodel.Point, 0, 0)
	points = append(points, &imagemodel.Point{})
	var p *imagemodel.Point
	for i := 0; i < len(image.ThrFaces["deepir_import"]); i++ {

		if i%2 == 0 {
			x, _ := image.ThrFaces["deepir_import"][i].(float64)
			p = &imagemodel.Point{
				X: x,
			}
		}
		if i%2 == 1 {
			y, _ := image.ThrFaces["deepir_import"][i].(float64)
			p.Y = y
			points = append(points, p)
		}
	}
	//	fmt.Println(points)

	pRep := &PointsRep{
		LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[19].X + points[20].X) / 2, Y: (points[19].Y + points[20].Y) / 2}, points[20], points[24], points[26], points[23], points[21], points[25], &imagemodel.Point{X: (points[19].X + 3*points[20].X) / 4, Y: (points[19].Y + 3*points[20].Y) / 4}, &imagemodel.Point{X: (3*points[19].X + points[20].X) / 4, Y: (3*points[19].Y + points[20].Y) / 4}},
		RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[27].X + points[28].X) / 2, Y: (points[27].Y + points[28].Y) / 2}, points[28], &imagemodel.Point{X: (3*points[27].X + points[28].X) / 4, Y: (3*points[27].Y + points[28].Y) / 4}, &imagemodel.Point{X: (points[27].X + 3*points[28].X) / 4, Y: (points[27].Y + 3*points[28].Y) / 4}, points[31], points[29], points[33], points[34], points[32]},
		LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5], points[6], points[8], points[9], points[7]},
		RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14], points[15], points[17], points[18], points[16]},
		LeftEar:      []*imagemodel.Point{points[76], points[72], points[73], points[74], points[78]},
		RightEar:     []*imagemodel.Point{points[75], points[69], points[70], points[71], points[77]},
		Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42], &imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y + points[38].Y) / 2}, &imagemodel.Point{X: (points[37].X+points[38].X)/2 + 2*(points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + 2*(points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[43], points[44], points[39], points[40]},
		Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60], points[53], points[54], points[66], points[63], points[57], points[58], points[62], points[61], points[51], points[52], &imagemodel.Point{X: (points[50].X + points[59].X) / 2, Y: (points[50].Y + points[59].Y) / 2}, points[64], points[65]},
		Face:         []*imagemodel.Point{points[79], points[87], points[86], points[85], points[84], points[83], points[82], points[81], points[80], points[95], points[94], points[93], points[92], points[91], points[90], points[89], points[88]},
	}

	//	fmt.Println(pRep)
	return pRep
}

//&imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y+points[38].Y)/2 + 12}
