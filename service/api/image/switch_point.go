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
	log.Info(fmt.Sprintf("image = %s switch 95 point", image.Md5))
	//	fmt.Println(image.ThrFaces["deepir_import"])
	if image.ThrFaces["deepir_import"] == nil {
		if image.ThrFaces["face++"] == nil {
			log.Info(fmt.Sprintf("image = %s thr point =nil", image.Md5))
			return &PointsRep{}
		}
	}
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

	pRep := &PointsRep{
		LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}, points[23], points[21], points[25], points[26], points[24]},
		RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28], &imagemodel.Point{X: (points[31].X + points[32].X) / 2, Y: (points[31].Y + points[32].Y) / 2}, &imagemodel.Point{X: (points[29].X + points[30].X) / 2, Y: (points[29].Y + points[30].Y) / 2}, points[29], points[31], points[33], points[34], points[32]},
		LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5], points[6], points[8], points[9], points[7]},
		RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14], points[15], points[17], points[18], points[16]},
		LeftEar:      []*imagemodel.Point{points[76], points[72], points[73], points[74], points[78]},
		RightEar:     []*imagemodel.Point{points[75], points[69], points[70], points[71], points[77]},
		Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42], &imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y + points[38].Y) / 2}, &imagemodel.Point{X: (points[37].X+points[38].X)/2 + 2*(points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + 2*(points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[43], points[44], points[39], points[40]},
		Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60], points[51], points[52], points[65], points[64], points[57], points[58], points[62], points[61], points[53], points[54], &imagemodel.Point{X: (points[50].X + points[59].X) / 2, Y: (points[50].Y + points[59].Y) / 2}, points[66], points[63]},
		Face:         []*imagemodel.Point{points[79], points[80], points[81], points[82], points[83], points[84], points[85], points[86], points[87], points[88], points[89], points[90], points[91], points[92], points[93], points[94], points[95]},
	}

	return pRep
}

func SwitchEightPoint(image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch 83 point", image.Md5))
	//	fmt.Println(image.ThrFaces["deepir_import"])
	if image.ThrFaces["deepir_import"] == nil {
		if image.ThrFaces["face++"] == nil {
			log.Info(fmt.Sprintf("image = %s thr point =nil", image.Md5))
			return &PointsRep{}
		}
	}
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

	pRep := &PointsRep{
		LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}, points[23], points[21], points[25], points[26], points[24]},
		RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28], &imagemodel.Point{X: (points[31].X + points[32].X) / 2, Y: (points[31].Y + points[32].Y) / 2}, &imagemodel.Point{X: (points[29].X + points[30].X) / 2, Y: (points[29].Y + points[30].Y) / 2}, points[29], points[31], points[33], points[34], points[32]},
		LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5], points[6], points[8], points[9], points[7]},
		RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14], points[15], points[17], points[18], points[16]},
		//		LeftEar:      []*imagemodel.Point{points[76], points[72], points[73], points[74], points[78]},
		//		RightEar:     []*imagemodel.Point{points[75], points[69], points[70], points[71], points[77]},
		Nouse: []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42], &imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y + points[38].Y) / 2}, &imagemodel.Point{X: (points[37].X+points[38].X)/2 + 2*(points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + 2*(points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[43], points[44], points[39], points[40]},
		Mouth: []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60], points[51], points[52], points[65], points[64], points[57], points[58], points[62], points[61], points[53], points[54], &imagemodel.Point{X: (points[50].X + points[59].X) / 2, Y: (points[50].Y + points[59].Y) / 2}},
		Face:  []*imagemodel.Point{points[79], points[80], points[81], points[82], points[83], points[84], points[85], points[86], points[87], points[88], points[89], points[90], points[91], points[92], points[93], points[94], points[95]},
	}

	return pRep
}

func fineTuneSwitchPoint(pointType string, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image fineTune = %s switch point", image.Md5))
	if image.Results[pointType] == nil {
		log.Info(fmt.Sprintf("image  = %s switch point", image.Md5))
		return SwitchPoint(image)
	}
	log.Info(fmt.Sprintf("image fineTune = %s switch point", image.Md5))
	var areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "leftEar", "rightEar", "mouth", "nouse", "face"}
	pRep := &PointsRep{}
	for _, area := range areas {
		if image.Results[pointType][area] != nil {
			//
			switch area {
			case "leftEyebrow":
				pRep.LeftEyeBrow = image.Results[pointType][area][0].Points
			case "rightEyebrow":
				pRep.RightEyeBrow = image.Results[pointType][area][0].Points
			case "leftEye":
				pRep.LeftEye = image.Results[pointType][area][0].Points
			case "rightEye":
				pRep.RightEye = image.Results[pointType][area][0].Points
			case "leftEar":
				pRep.LeftEar = image.Results[pointType][area][0].Points
			case "rightEar":
				pRep.RightEar = image.Results[pointType][area][0].Points
			case "mouth":
				pRep.Mouth = image.Results[pointType][area][0].Points
			case "nouse":
				pRep.Nouse = image.Results[pointType][area][0].Points
			case "face":
				pRep.Face = image.Results[pointType][area][0].Points
			}
		}
	}
	return pRep
}
