package image

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

type AllPointsRep struct {
	FineResult []*imagemodel.Points `json:"fine_result"`
	ResultRep  *ResultRep           `json:"result"`
	PointType  int64                `json:"point_type"`
	Status     int64                `json:"status"`
}

type ResPointsRep struct {
	ResultRep *ResultRep `json:"result"`
	PointType int64      `json:"point_type"`
}

type ResultRep struct {
	LeftEyeBrow  []*imagemodel.Points `json:"leftEyebrow"`
	RightEyeBrow []*imagemodel.Points `json:"rightEyebrow"`
	LeftEye      []*imagemodel.Points `json:"leftEye"`
	RightEye     []*imagemodel.Points `json:"rightEye"`
	LeftEar      []*imagemodel.Points `json:"leftEar"`
	RightEar     []*imagemodel.Points `json:"rightEar"`
	Mouth        []*imagemodel.Points `json:"mouth"`
	Nouse        []*imagemodel.Points `json:"nouse"`
	Face         []*imagemodel.Points `json:"face"`
}

type Points struct {
	Points []float64 `json:"95"`
}

type ThrResRep struct {
	WidthScale  float64             `json:"width_scale"`
	HeightScale float64             `json:"height_scale"`
	Name        string              `json:"name"`
	Points      []*imagemodel.Point `json:"points"`
}

func SwitchPoint(pointType int64, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch point", image.Md5))
	switch pointType {
	case 5:
		return switchAllPoint(pointType, image)
	case 27:
		return switchAllPoint(pointType, image)
	case 68:
		return switchAllPoint(pointType, image)
	case 83:
		return SwitchEightPoint(pointType, image)
	case 95:
		return SwitchNinePoint(pointType, image)
	}
	return SwitchNilPoint(pointType)
}

func switchAllPoint(pointType int64, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch all point", image.Md5))
	if image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))] == nil {
		if image.ThrFaces["face++"] == nil {
			log.Info(fmt.Sprintf("image = %s thr point =nil", image.Md5))
			return SwitchNilPoint(pointType)
		}
		return faceResSwitchPoint(pointType, image)
	}

	res1B, _ := json.Marshal(image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))])
	//	fmt.Println(string(res1B)) //json
	imPoints := new(Points)
	if err := json.Unmarshal(res1B, &imPoints.Points); err != nil {
		fmt.Println("json unmarshal err=%s", err)
		return SwitchNilPoint(pointType)
	}
	fmt.Println(len(imPoints.Points))

	points := make([]*imagemodel.Point, 0, 0)
	points = append(points, &imagemodel.Point{})
	var p *imagemodel.Point
	for i := 0; i < len(imPoints.Points); i++ {
		if i%2 == 0 {
			x := imPoints.Points[i]
			p = &imagemodel.Point{
				X: x,
			}
		}
		if i%2 == 1 {
			y := imPoints.Points[i]
			p.Y = y
			points = append(points, p)
		}
	}
	fmt.Println(points)
	var pRep *PointsRep
	switch pointType {
	case 5:
		pRep = &PointsRep{
			LeftEye:  []*imagemodel.Point{points[1]},
			RightEye: []*imagemodel.Point{points[10]},
			Nouse:    []*imagemodel.Point{points[35]},
			Mouth:    []*imagemodel.Point{points[47], points[48]},
		}
		return pRep
	case 27:
		pRep = &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20]},
			RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28]},
			LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5]},
			RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14]},
			Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42]},
			Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60]},
		}
		return pRep
	case 68:
		pRep = &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}},
			RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28], &imagemodel.Point{X: (points[31].X + points[32].X) / 2, Y: (points[31].Y + points[32].Y) / 2}, &imagemodel.Point{X: (points[29].X + points[30].X) / 2, Y: (points[29].Y + points[30].Y) / 2}},
			LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5], points[6], points[8], points[9], points[7]},
			RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14], points[15], points[17], points[18], points[16]},
			Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42], &imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y + points[38].Y) / 2}, &imagemodel.Point{X: (points[37].X+points[38].X)/2 + 2*(points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + 2*(points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[43], points[44]},
			Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60], points[51], points[64], points[62], points[61], points[53], points[54], points[66], points[63]},
			Face:         []*imagemodel.Point{points[79], points[80], points[81], points[82], points[83], points[84], points[85], points[86], points[87], points[88], points[89], points[90], points[91], points[92], points[93], points[94], points[95]},
		}
		return pRep
	case 83:
		pRep = &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}, points[21], points[23], points[25], points[26], points[24]},
			RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28], &imagemodel.Point{X: (points[31].X + points[32].X) / 2, Y: (points[31].Y + points[32].Y) / 2}, &imagemodel.Point{X: (points[29].X + points[30].X) / 2, Y: (points[29].Y + points[30].Y) / 2}, points[29], points[31], points[33], points[34], points[32]},
			LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5], points[6], points[8], points[9], points[7]},
			RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14], points[15], points[17], points[18], points[16]},
			Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42], &imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y + points[38].Y) / 2}, &imagemodel.Point{X: (points[37].X+points[38].X)/2 + 2*(points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + 2*(points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[43], points[44], points[39], points[40]},
			Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60], points[51], points[52], points[65], points[64], points[57], points[58], points[62], points[61], points[53], points[54], &imagemodel.Point{X: (points[50].X + points[59].X) / 2, Y: (points[50].Y + points[59].Y) / 2}},
			Face:         []*imagemodel.Point{points[79], points[80], points[81], points[82], points[83], points[84], points[85], points[86], points[87], points[88], points[89], points[90], points[91], points[92], points[93], points[94], points[95]},
		}
		return pRep
	case 95:
		pRep = &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}, points[21], points[23], points[25], points[26], points[24]},
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

	return pRep
}

func SwitchNinePoint(pointType int64, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch 95 point", image.Md5))
	//	fmt.Println(image.ThrFaces["deepir_import"])
	if image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))] == nil {
		if image.ThrFaces["face++"] == nil {
			log.Info(fmt.Sprintf("image = %s thr point =nil", image.Md5))
			return SwitchNilPoint(pointType)
		}
		return faceResSwitchPoint(pointType, image)
	}

	res1B, _ := json.Marshal(image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))])
	//	fmt.Println(string(res1B)) //json
	imPoints := new(Points)
	if err := json.Unmarshal(res1B, &imPoints.Points); err != nil {
		fmt.Println("json unmarshal err=%s", err)
		return SwitchNilPoint(pointType)
	}
	fmt.Println(len(imPoints.Points))

	points := make([]*imagemodel.Point, 0, 0)
	points = append(points, &imagemodel.Point{})
	var p *imagemodel.Point
	for i := 0; i < len(imPoints.Points); i++ {
		if i%2 == 0 {
			x := imPoints.Points[i]
			p = &imagemodel.Point{
				X: x,
			}
		}
		if i%2 == 1 {
			y := imPoints.Points[i]
			p.Y = y
			points = append(points, p)
		}
	}
	fmt.Println(points)

	pRep := &PointsRep{
		LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}, points[21], points[23], points[25], points[26], points[24]},
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

func SwitchEightPoint(pointType int64, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch 83 point", image.Md5))

	if image.ThrFaces["deepir_import"] == nil {
		if image.ThrFaces["face++"] == nil {
			log.Info(fmt.Sprintf("image = %s thr point =nil", image.Md5))
			return SwitchNilPoint(pointType)
		}
		return faceResSwitchPoint(pointType, image)
	}
	var res1B []byte
	if image.ThrFaces["deepir_import"]["95"] != nil {
		res, err := json.Marshal(image.ThrFaces["deepir_import"]["95"])
		if err != nil {
			fmt.Println("json Marshal 95 err=%s", err)
			return SwitchNilPoint(pointType)
		}
		res1B = res
	} else {
		res, err := json.Marshal(image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))])
		if err != nil {
			fmt.Println("json Marshal deepir import err=%s", err)
			return SwitchNilPoint(pointType)
		}
		res1B = res
	}

	fmt.Println(string(res1B)) //json
	imPoints := new(Points)
	if err := json.Unmarshal(res1B, &imPoints.Points); err != nil {
		fmt.Println("json unmarshal err=%s", err)
		return SwitchNilPoint(pointType)
	}
	fmt.Println(len(imPoints.Points))
	points := make([]*imagemodel.Point, 0, 0)
	points = append(points, &imagemodel.Point{})
	var p *imagemodel.Point
	for i := 0; i < len(imPoints.Points); i++ {
		if i%2 == 0 {
			x := imPoints.Points[i]
			p = &imagemodel.Point{
				X: x,
			}
		}
		if i%2 == 1 {
			y := imPoints.Points[i]
			p.Y = y
			points = append(points, p)
		}
	}
	fmt.Println(points)
	if len(imPoints.Points) == 0 {
		return SwitchNilPoint(pointType)
	}
	pRep := &PointsRep{
		LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}, points[21], points[23], points[25], points[26], points[24]},
		RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28], &imagemodel.Point{X: (points[31].X + points[32].X) / 2, Y: (points[31].Y + points[32].Y) / 2}, &imagemodel.Point{X: (points[29].X + points[30].X) / 2, Y: (points[29].Y + points[30].Y) / 2}, points[29], points[31], points[33], points[34], points[32]},
		LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5], points[6], points[8], points[9], points[7]},
		RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14], points[15], points[17], points[18], points[16]},
		Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42], &imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y + points[38].Y) / 2}, &imagemodel.Point{X: (points[37].X+points[38].X)/2 + 2*(points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + 2*(points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[43], points[44], points[39], points[40]},
		Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60], points[51], points[52], points[65], points[64], points[57], points[58], points[62], points[61], points[53], points[54], &imagemodel.Point{X: (points[50].X + points[59].X) / 2, Y: (points[50].Y + points[59].Y) / 2}},
		Face:         []*imagemodel.Point{points[79], points[80], points[81], points[82], points[83], points[84], points[85], points[86], points[87], points[88], points[89], points[90], points[91], points[92], points[93], points[94], points[95]},
	}

	return pRep
}

func faceResSwitchPoint(pointType int64, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch face++ point", image.Md5))
	if image.ThrFaces["face++"]["83"] == nil {
		log.Info(fmt.Sprintf("image = %s thr point =nil", image.Md5))
		return SwitchNilPoint(pointType)
	}

	res1B, _ := json.Marshal(image.ThrFaces["face++"]["83"])
	//	fmt.Println(string(res1B)) //json
	faceRes := &thrfacemodel.FaceModelV3{}
	if err := json.Unmarshal([]byte(string(res1B)), &faceRes); err != nil {
		fmt.Println("json unmarshal err=%s", err)
		return SwitchNilPoint(pointType)
	}

	if faceRes.Faces == nil || len(faceRes.Faces) == 0 {
		log.Info(fmt.Sprintf("image = %s face++ nil no face point", image.Md5))
		return SwitchNilPoint(pointType)
	}

	//	fmt.Println("---faceRes%s---", faceRes)
	points := make([]*imagemodel.Point, 0, 0)
	var p *imagemodel.Point
	for _, point := range faceRes.Faces[0].Landmark {
		if point != nil {
			p = &imagemodel.Point{
				X: point.X,
				Y: point.Y,
			}
			points = append(points, p)
		}
	}

	if len(points) < 83 {
		log.Info(fmt.Sprintf("face++ unmarshal len(points)=%d switch point<83", len(points)))
		return SwitchNilPoint(pointType)
	}

	landmark := faceRes.Faces[0].Landmark
	if landmark == nil {
		log.Info(fmt.Sprintf("face++ unmarshal landmark==nil"))
		return SwitchNilPoint(pointType)
	}
	fmt.Println("width=%d,/height=%d", faceRes.Faces[0].FaceRectangle.ImageWidth, faceRes.Faces[0].FaceRectangle.ImageHeight)
	pRep := &PointsRep{
		LeftEyeBrow: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["left_eyebrow_left_corner"].X, Y: landmark["left_eyebrow_left_corner"].Y},
			&imagemodel.Point{X: ((landmark["left_eyebrow_lower_middle"].X + landmark["left_eyebrow_upper_middle"].X) / 2), Y: (landmark["left_eyebrow_lower_middle"].Y + landmark["left_eyebrow_upper_middle"].Y) / 2},
			&imagemodel.Point{X: landmark["left_eyebrow_right_corner"].X, Y: landmark["left_eyebrow_right_corner"].Y},
			&imagemodel.Point{X: (landmark["left_eyebrow_upper_left_quarter"].X + landmark["left_eyebrow_lower_left_quarter"].X) / 2, Y: (landmark["left_eyebrow_upper_left_quarter"].Y + landmark["left_eyebrow_lower_left_quarter"].Y) / 2},
			&imagemodel.Point{X: (landmark["left_eyebrow_upper_right_quarter"].X + landmark["left_eyebrow_lower_right_quarter"].X) / 2, Y: (landmark["left_eyebrow_upper_right_quarter"].Y + landmark["left_eyebrow_lower_right_quarter"].Y) / 2},
			&imagemodel.Point{X: landmark["left_eyebrow_upper_middle"].X, Y: landmark["left_eyebrow_upper_middle"].Y},
			&imagemodel.Point{X: landmark["left_eyebrow_upper_left_quarter"].X, Y: landmark["left_eyebrow_upper_left_quarter"].Y},
			&imagemodel.Point{X: landmark["left_eyebrow_upper_right_quarter"].X, Y: landmark["left_eyebrow_upper_right_quarter"].Y},
			&imagemodel.Point{X: landmark["left_eyebrow_lower_right_quarter"].X, Y: landmark["left_eyebrow_lower_right_quarter"].Y},
			&imagemodel.Point{X: landmark["left_eyebrow_lower_left_quarter"].X, Y: landmark["left_eyebrow_lower_left_quarter"].Y},
		},
		RightEyeBrow: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["right_eyebrow_left_corner"].X, Y: landmark["right_eyebrow_left_corner"].Y},
			&imagemodel.Point{X: (landmark["right_eyebrow_lower_middle"].X + landmark["right_eyebrow_upper_middle"].X) / 2, Y: (landmark["right_eyebrow_lower_middle"].Y + landmark["right_eyebrow_upper_middle"].Y) / 2},
			&imagemodel.Point{X: landmark["right_eyebrow_right_corner"].X, Y: landmark["right_eyebrow_right_corner"].Y},
			&imagemodel.Point{X: (landmark["right_eyebrow_upper_left_quarter"].X + landmark["right_eyebrow_lower_left_quarter"].X) / 2, Y: (landmark["right_eyebrow_upper_left_quarter"].Y + landmark["right_eyebrow_lower_left_quarter"].Y) / 2},
			&imagemodel.Point{X: (landmark["right_eyebrow_upper_right_quarter"].X + landmark["right_eyebrow_lower_right_quarter"].X) / 2, Y: (landmark["right_eyebrow_upper_right_quarter"].Y + landmark["right_eyebrow_lower_right_quarter"].Y) / 2},
			&imagemodel.Point{X: landmark["right_eyebrow_upper_middle"].X, Y: landmark["right_eyebrow_upper_middle"].Y},
			&imagemodel.Point{X: landmark["right_eyebrow_upper_left_quarter"].X, Y: landmark["right_eyebrow_upper_left_quarter"].Y},
			&imagemodel.Point{X: landmark["right_eyebrow_upper_right_quarter"].X, Y: landmark["right_eyebrow_upper_right_quarter"].Y},
			&imagemodel.Point{X: landmark["right_eyebrow_lower_right_quarter"].X, Y: landmark["right_eyebrow_lower_right_quarter"].Y},
			&imagemodel.Point{X: landmark["right_eyebrow_lower_left_quarter"].X, Y: landmark["right_eyebrow_lower_left_quarter"].Y},
		},
		LeftEye: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["left_eye_pupil"].X, Y: landmark["left_eye_pupil"].Y},
			&imagemodel.Point{X: landmark["left_eye_left_corner"].X, Y: landmark["left_eye_left_corner"].Y},
			&imagemodel.Point{X: landmark["left_eye_top"].X, Y: landmark["left_eye_top"].Y},
			&imagemodel.Point{X: landmark["left_eye_right_corner"].X, Y: landmark["left_eye_right_corner"].Y},
			&imagemodel.Point{X: landmark["left_eye_bottom"].X, Y: landmark["left_eye_bottom"].Y},
			&imagemodel.Point{X: landmark["left_eye_upper_left_quarter"].X, Y: landmark["left_eye_upper_left_quarter"].Y},
			&imagemodel.Point{X: landmark["left_eye_upper_right_quarter"].X, Y: landmark["left_eye_upper_right_quarter"].Y},
			&imagemodel.Point{X: landmark["left_eye_lower_right_quarter"].X, Y: landmark["left_eye_lower_right_quarter"].Y},
			&imagemodel.Point{X: landmark["left_eye_lower_left_quarter"].X, Y: landmark["left_eye_lower_left_quarter"].Y},
		},
		RightEye: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["right_eye_pupil"].X, Y: landmark["right_eye_pupil"].Y},
			&imagemodel.Point{X: landmark["right_eye_left_corner"].X, Y: landmark["right_eye_left_corner"].Y},
			&imagemodel.Point{X: landmark["right_eye_top"].X, Y: landmark["right_eye_top"].Y},
			&imagemodel.Point{X: landmark["right_eye_right_corner"].X, Y: landmark["right_eye_right_corner"].Y},
			&imagemodel.Point{X: landmark["right_eye_bottom"].X, Y: landmark["right_eye_bottom"].Y},
			&imagemodel.Point{X: landmark["right_eye_upper_left_quarter"].X, Y: landmark["right_eye_upper_left_quarter"].Y},
			&imagemodel.Point{X: landmark["right_eye_upper_right_quarter"].X, Y: landmark["right_eye_upper_right_quarter"].Y},
			&imagemodel.Point{X: landmark["right_eye_lower_right_quarter"].X, Y: landmark["right_eye_lower_right_quarter"].Y},
			&imagemodel.Point{X: landmark["right_eye_lower_left_quarter"].X, Y: landmark["right_eye_lower_left_quarter"].Y},
		},
		Nouse: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["nose_tip"].X, Y: landmark["nose_tip"].Y},
			&imagemodel.Point{X: (landmark["nose_tip"].X + ((landmark["nose_contour_left1"].X+landmark["nose_contour_right1"].X)/2-landmark["nose_tip"].X)*2/3), Y: (landmark["nose_tip"].Y + ((landmark["nose_contour_left1"].Y+landmark["nose_contour_right1"].Y)/2-landmark["nose_tip"].Y)*2/3)},
			&imagemodel.Point{X: landmark["nose_left"].X, Y: landmark["nose_left"].Y},
			&imagemodel.Point{X: landmark["nose_contour_lower_middle"].X, Y: landmark["nose_contour_lower_middle"].Y},
			&imagemodel.Point{X: landmark["nose_right"].X, Y: landmark["nose_right"].Y},
			&imagemodel.Point{X: (landmark["nose_contour_left1"].X + landmark["nose_contour_right1"].X) / 2, Y: (landmark["nose_contour_left1"].Y + landmark["nose_contour_right1"].Y) / 2},
			&imagemodel.Point{X: (landmark["nose_tip"].X + ((landmark["nose_contour_left1"].X+landmark["nose_contour_right1"].X)/2-landmark["nose_tip"].X)/3), Y: (landmark["nose_tip"].Y + ((landmark["nose_contour_left1"].Y+landmark["nose_contour_right1"].Y)/2-landmark["nose_tip"].Y)/3)},
			&imagemodel.Point{X: landmark["nose_contour_left3"].X, Y: landmark["nose_contour_left3"].Y},
			&imagemodel.Point{X: landmark["nose_contour_right3"].X, Y: landmark["nose_contour_right3"].Y},
			&imagemodel.Point{X: landmark["nose_contour_left2"].X, Y: landmark["nose_contour_left2"].Y},
			&imagemodel.Point{X: landmark["nose_contour_right2"].X, Y: landmark["nose_contour_right2"].Y},
		},
		Mouth: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["mouth_left_corner"].X, Y: landmark["mouth_left_corner"].Y},
			&imagemodel.Point{X: landmark["mouth_right_corner"].X, Y: landmark["mouth_right_corner"].Y},
			&imagemodel.Point{X: landmark["mouth_upper_lip_top"].X, Y: landmark["mouth_upper_lip_top"].Y},
			&imagemodel.Point{X: landmark["mouth_upper_lip_bottom"].X, Y: landmark["mouth_upper_lip_bottom"].Y},
			&imagemodel.Point{X: landmark["mouth_lower_lip_top"].X, Y: landmark["mouth_lower_lip_top"].Y},
			&imagemodel.Point{X: landmark["mouth_lower_lip_bottom"].X, Y: landmark["mouth_lower_lip_bottom"].Y},
			&imagemodel.Point{X: landmark["mouth_upper_lip_left_contour1"].X, Y: landmark["mouth_upper_lip_left_contour1"].Y},
			&imagemodel.Point{X: landmark["mouth_upper_lip_right_contour1"].X, Y: landmark["mouth_upper_lip_right_contour1"].Y},
			&imagemodel.Point{X: landmark["mouth_lower_lip_right_contour3"].X, Y: landmark["mouth_lower_lip_right_contour3"].Y},
			&imagemodel.Point{X: landmark["mouth_lower_lip_left_contour3"].X, Y: landmark["mouth_lower_lip_left_contour3"].Y},
			&imagemodel.Point{X: landmark["mouth_upper_lip_left_contour3"].X, Y: landmark["mouth_upper_lip_left_contour3"].Y},
			&imagemodel.Point{X: landmark["mouth_upper_lip_right_contour3"].X, Y: landmark["mouth_upper_lip_right_contour3"].Y},
			&imagemodel.Point{X: landmark["mouth_lower_lip_right_contour1"].X, Y: landmark["mouth_lower_lip_right_contour1"].Y},
			&imagemodel.Point{X: landmark["mouth_lower_lip_left_contour1"].X, Y: landmark["mouth_lower_lip_left_contour1"].Y},
			&imagemodel.Point{X: landmark["mouth_upper_lip_left_contour2"].X, Y: landmark["mouth_upper_lip_left_contour2"].Y},
			&imagemodel.Point{X: landmark["mouth_upper_lip_right_contour2"].X, Y: landmark["mouth_upper_lip_right_contour2"].Y},
			&imagemodel.Point{X: (landmark["mouth_left_corner"].X + landmark["mouth_right_corner"].X) / 2, Y: (landmark["mouth_left_corner"].Y + landmark["mouth_right_corner"].Y) / 2},
			//			&imagemodel.Point{X: landmark["mouth_lower_lip_right_contour2"].X , Y: landmark["mouth_lower_lip_right_contour2"].Y },
			//			&imagemodel.Point{X: landmark["mouth_lower_lip_left_contour2"].X , Y: landmark["mouth_lower_lip_left_contour2"].Y },
		},
		Face: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["contour_chin"].X, Y: landmark["contour_chin"].Y},
			&imagemodel.Point{X: landmark["contour_left1"].X, Y: landmark["contour_left1"].Y},
			&imagemodel.Point{X: landmark["contour_left2"].X, Y: landmark["contour_left2"].Y},
			&imagemodel.Point{X: landmark["contour_left3"].X, Y: landmark["contour_left3"].Y},
			&imagemodel.Point{X: landmark["contour_left4"].X, Y: landmark["contour_left4"].Y},
			&imagemodel.Point{X: landmark["contour_left5"].X, Y: landmark["contour_left5"].Y},
			&imagemodel.Point{X: landmark["contour_left6"].X, Y: landmark["contour_left6"].Y},
			&imagemodel.Point{X: landmark["contour_left7"].X, Y: landmark["contour_left7"].Y},
			&imagemodel.Point{X: landmark["contour_left8"].X, Y: landmark["contour_left8"].Y},
			&imagemodel.Point{X: landmark["contour_right1"].X, Y: landmark["contour_right1"].Y},
			&imagemodel.Point{X: landmark["contour_right2"].X, Y: landmark["contour_right2"].Y},
			&imagemodel.Point{X: landmark["contour_right3"].X, Y: landmark["contour_right3"].Y},
			&imagemodel.Point{X: landmark["contour_right4"].X, Y: landmark["contour_right4"].Y},
			&imagemodel.Point{X: landmark["contour_right5"].X, Y: landmark["contour_right5"].Y},
			&imagemodel.Point{X: landmark["contour_right6"].X, Y: landmark["contour_right6"].Y},
			&imagemodel.Point{X: landmark["contour_right7"].X, Y: landmark["contour_right7"].Y},
			&imagemodel.Point{X: landmark["contour_right8"].X, Y: landmark["contour_right8"].Y},
		},
	}

	if pRep == nil {
		log.Info(fmt.Sprintf("image face++ switch point pRep = nil "))
		return SwitchNilPoint(pointType)
	}

	pRepre := &PointsRep{}
	switch pointType {
	case 5:
		pRepre.LeftEye = append(pRepre.LeftEye, pRep.LeftEye[0])
		pRepre.RightEye = append(pRepre.RightEye, pRep.RightEye[0])
		pRepre.Nouse = append(pRepre.Nouse, pRep.Nouse[0])
		pRepre.Mouth = pRep.Mouth[0:2]
		return pRepre
	case 27:
		pRepre.LeftEyeBrow = pRep.LeftEyeBrow[0:3]
		pRepre.RightEyeBrow = pRep.RightEyeBrow[0:3]
		pRepre.LeftEye = pRep.LeftEye[0:5]
		pRepre.RightEye = pRep.RightEye[0:5]
		pRepre.Nouse = pRep.Nouse[0:5]
		pRepre.Mouth = pRep.Mouth[0:6]
		return pRepre
	case 68:
		pRepre.LeftEyeBrow = pRep.LeftEyeBrow[0:5]
		pRepre.RightEyeBrow = pRep.RightEyeBrow[0:5]
		pRepre.LeftEye = pRep.LeftEye[0:9]
		pRepre.RightEye = pRep.RightEye[0:9]
		pRepre.Nouse = pRep.Nouse[0:9]
		pRepre.Mouth = pRep.Mouth[0:6]
		pRepre.Mouth = append(pRepre.Mouth, pRep.Mouth[6])
		pRepre.Mouth = append(pRepre.Mouth, pRep.Mouth[9])
		pRep.Mouth = append(pRep.Mouth, &imagemodel.Point{X: landmark["mouth_lower_lip_right_contour2"].X, Y: landmark["mouth_lower_lip_right_contour2"].Y})
		pRep.Mouth = append(pRep.Mouth, &imagemodel.Point{X: landmark["mouth_lower_lip_left_contour2"].X, Y: landmark["mouth_lower_lip_left_contour2"].Y})
		pRepre.Mouth = append(pRepre.Mouth, pRep.Mouth[12])
		pRepre.Mouth = append(pRepre.Mouth, pRep.Mouth[13])
		pRepre.Mouth = append(pRepre.Mouth, pRep.Mouth[14])
		pRepre.Mouth = append(pRepre.Mouth, pRep.Mouth[15])
		pRepre.Face = pRep.Face
		return pRepre
	case 83:
		pRepre = pRep
		return pRepre
	case 95:
		pRep.Mouth = append(pRep.Mouth, &imagemodel.Point{X: landmark["mouth_lower_lip_right_contour2"].X, Y: landmark["mouth_lower_lip_right_contour2"].Y})
		pRep.Mouth = append(pRep.Mouth, &imagemodel.Point{X: landmark["mouth_lower_lip_left_contour2"].X, Y: landmark["mouth_lower_lip_left_contour2"].Y})
		pp := SwitchNilPoint(pointType)
		pRep.LeftEar = pp.LeftEar
		pRep.RightEar = pp.RightEar
		pRepre = pRep
		return pRepre
	}

	return pRepre
}

func ThrResults(pointType int64, image *imagemodel.ImageModel) []*ThrResRep {
	log.Info(fmt.Sprintf("image = %s switch face++ point", image.Md5))
	resultPoints := make([]*ThrResRep, 0, 0)
	if image.ThrFaces == nil {
		log.Info(fmt.Sprintf("image = %s thr point==nil ", image.Md5))
		return resultPoints
	}

	if image.ThrFaces["face++"]["83"] != nil {
		log.Info(fmt.Sprintf("image = %s thr point ", image.Md5))
		res1B, _ := json.Marshal(image.ThrFaces["face++"]["83"])
		//		fmt.Println(string(res1B)) //json

		faceRes := &thrfacemodel.FaceModelV3{}
		if err := json.Unmarshal([]byte(string(res1B)), &faceRes); err == nil {
			//			fmt.Println("---faceRes%s----", faceRes) //json
			if faceRes.Faces != nil && len(faceRes.Faces) != 0 {
				landmark := faceRes.Faces[0].Landmark
				if landmark != nil {

					points := make([]*imagemodel.Point, 0, 0)
					var p *imagemodel.Point
					for _, point := range faceRes.Faces[0].Landmark {
						if point != nil {
							p = &imagemodel.Point{
								X: point.X,
								Y: point.Y,
							}
							points = append(points, p)
						}
					}
					//				if len(points) >= 83 {
					thrResRep := &ThrResRep{
						HeightScale: 1,
						WidthScale:  1,
						Name:        "face++",
						Points:      points,
					}

					resultPoints = append(resultPoints, thrResRep)
				} else {
					log.Info(fmt.Sprintf("landmark==nil "))
				}

			} else {
				log.Info(fmt.Sprintf("image = %s face++ nil no face point", image.Md5))
			}

		} else {
			log.Info(fmt.Sprintf("face++json unmarshalerr "))
		}
		//
	}

	if image.ThrFaces["deepir_import"] != nil {
		var res1B []byte
		if image.ThrFaces["deepir_import"]["95"] != nil {
			res, err := json.Marshal(image.ThrFaces["deepir_import"]["95"])
			if err != nil {
				fmt.Println("json Marshal 95 err=%s", err)
				return resultPoints
			}
			res1B = res
		} else {
			res, err := json.Marshal(image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))])
			if err != nil {
				fmt.Println("json Marshal deepir import err=%s", err)
				return resultPoints
			}
			res1B = res
		}

		imPoints := new(Points)
		if err := json.Unmarshal(res1B, &imPoints.Points); err == nil {
			points := make([]*imagemodel.Point, 0, 0)
			var p *imagemodel.Point
			for i := 0; i < len(imPoints.Points); i++ {
				if i%2 == 0 {
					x := imPoints.Points[i]
					p = &imagemodel.Point{
						X: x,
					}
				}
				if i%2 == 1 {
					y := imPoints.Points[i]
					p.Y = y
					points = append(points, p)
				}
			}
			if len(points) > 0 {
				thrResRep := &ThrResRep{
					HeightScale: 1,
					WidthScale:  1,
					Name:        "import",
					Points:      points,
				}

				resultPoints = append(resultPoints, thrResRep)
			}

		} else {
			log.Info(fmt.Sprintf("import res json unmarshal err"))
		}

	}

	return resultPoints
}

func GetThrResults(image *imagemodel.ImageModel) []*ThrResRep {
	log.Info(fmt.Sprintf("image = %s switch face++ point", image.Md5))
	resultPoints := make([]*ThrResRep, 0, 0)
	if image.ThrFaces == nil {
		log.Info(fmt.Sprintf("image = %s thr point==nil ", image.Md5))
		return resultPoints
	}

	if image.ThrFaces["face++"]["83"] != nil {
		log.Info(fmt.Sprintf("image = %s thr point ", image.Md5))
		res1B, _ := json.Marshal(image.ThrFaces["face++"]["83"])
		//		fmt.Println(string(res1B)) //json

		faceRes := &thrfacemodel.FaceModelV3{}
		if err := json.Unmarshal([]byte(string(res1B)), &faceRes); err == nil {
			//			fmt.Println("---faceRes%s----", faceRes) //json
			if faceRes.Faces != nil && len(faceRes.Faces) != 0 {
				landmark := faceRes.Faces[0].Landmark
				if landmark != nil {

					points := make([]*imagemodel.Point, 0, 0)
					var p *imagemodel.Point
					for _, point := range faceRes.Faces[0].Landmark {
						if point != nil {
							p = &imagemodel.Point{
								X: point.X,
								Y: point.Y,
							}
							points = append(points, p)
						}
					}
					//				if len(points) >= 83 {
					thrResRep := &ThrResRep{
						HeightScale: 1,
						WidthScale:  1,
						Name:        "face++",
						Points:      points,
					}

					resultPoints = append(resultPoints, thrResRep)
				} else {
					log.Info(fmt.Sprintf("landmark==nil "))
				}
			} else {
				log.Info(fmt.Sprintf("face==nil "))
			}

		} else {
			log.Info(fmt.Sprintf("face++json unmarshalerr "))
		}
		//
	}

	if image.ThrFaces["deepir_import"] != nil {
		var res1B []byte
		if image.ThrFaces["deepir_import"]["95"] != nil {
			res, err := json.Marshal(image.ThrFaces["deepir_import"]["95"])
			if err != nil {
				fmt.Println("json Marshal 95 err=%s", err)
				return resultPoints
			}
			res1B = res
		} else if image.ThrFaces["deepir_import"]["83"] != nil {
			res, err := json.Marshal(image.ThrFaces["deepir_import"]["83"])
			if err != nil {
				fmt.Println("json Marshal deepir import err=%s", err)
				return resultPoints
			}
			res1B = res
		} else if image.ThrFaces["deepir_import"]["68"] != nil {
			res, err := json.Marshal(image.ThrFaces["deepir_import"]["68"])
			if err != nil {
				fmt.Println("json Marshal deepir import err=%s", err)
				return resultPoints
			}
			res1B = res
		} else if image.ThrFaces["deepir_import"]["27"] != nil {
			res, err := json.Marshal(image.ThrFaces["deepir_import"]["27"])
			if err != nil {
				fmt.Println("json Marshal deepir import err=%s", err)
				return resultPoints
			}
			res1B = res
		} else if image.ThrFaces["deepir_import"]["5"] != nil {
			res, err := json.Marshal(image.ThrFaces["deepir_import"]["5"])
			if err != nil {
				fmt.Println("json Marshal deepir import err=%s", err)
				return resultPoints
			}
			res1B = res
		}

		imPoints := new(Points)
		if err := json.Unmarshal(res1B, &imPoints.Points); err == nil {
			points := make([]*imagemodel.Point, 0, 0)
			var p *imagemodel.Point
			for i := 0; i < len(imPoints.Points); i++ {
				if i%2 == 0 {
					x := imPoints.Points[i]
					p = &imagemodel.Point{
						X: x,
					}
				}
				if i%2 == 1 {
					y := imPoints.Points[i]
					p.Y = y
					points = append(points, p)
				}
			}
			if len(points) > 0 {
				thrResRep := &ThrResRep{
					HeightScale: 1,
					WidthScale:  1,
					Name:        "import",
					Points:      points,
				}

				resultPoints = append(resultPoints, thrResRep)
			}

		} else {
			log.Info(fmt.Sprintf("import res json unmarshal err"))
		}

	}

	return resultPoints
}

func fineTuneSwitchPoint(pointType int64, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image fineTune = %s switch point", image.Md5))
	if image.Results[strconv.Itoa(int(pointType))] == nil {
		log.Info(fmt.Sprintf("image  = %s switch point", image.Md5))
		return SwitchPoint(pointType, image)
	}
	log.Info(fmt.Sprintf("image fineTune = %s switch point", image.Md5))
	var areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "leftEar", "rightEar", "mouth", "nouse", "face"}
	pRep := &PointsRep{}
	for _, area := range areas {
		if image.Results[strconv.Itoa(int(pointType))][area] != nil {
			//
			switch area {
			case "leftEyebrow":
				pRep.LeftEyeBrow = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			case "rightEyebrow":
				pRep.RightEyeBrow = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			case "leftEye":
				pRep.LeftEye = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			case "rightEye":
				pRep.RightEye = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			case "leftEar":
				pRep.LeftEar = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			case "rightEar":
				pRep.RightEar = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			case "mouth":
				pRep.Mouth = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			case "nouse":
				pRep.Nouse = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			case "face":
				pRep.Face = image.Results[strconv.Itoa(int(pointType))][area][0].Points
			}
		}
	}
	fmt.Println(pRep)
	return pRep
}

func GetAllResults(image *imagemodel.ImageModel) *AllPointsRep {
	log.Info(fmt.Sprintf("image all res = %s switch point", image.Md5))
	apres := &AllPointsRep{}
	resultRep := &ResultRep{
		LeftEyeBrow:  make([]*imagemodel.Points, 0, 0),
		RightEyeBrow: make([]*imagemodel.Points, 0, 0),
		LeftEye:      make([]*imagemodel.Points, 0, 0),
		RightEye:     make([]*imagemodel.Points, 0, 0),
		LeftEar:      make([]*imagemodel.Points, 0, 0),
		RightEar:     make([]*imagemodel.Points, 0, 0),
		Mouth:        make([]*imagemodel.Points, 0, 0),
		Nouse:        make([]*imagemodel.Points, 0, 0),
		Face:         make([]*imagemodel.Points, 0, 0),
	}
	if image.FineResults["95"] != nil {
		apres.FineResult = fineResToPoints(image.FineResults["95"])
		apres.ResultRep = resultRep
		apres.Status = 1
		apres.PointType = 95
		return apres
	}

	if image.FineResults["83"] != nil {
		apres.FineResult = fineResToPoints(image.FineResults["83"])
		apres.ResultRep = resultRep
		apres.Status = 1
		apres.PointType = 83
		return apres
	}

	if image.FineResults["68"] != nil {
		apres.FineResult = fineResToPoints(image.FineResults["68"])
		apres.ResultRep = resultRep
		apres.Status = 1
		apres.PointType = 68
		return apres
	}

	if image.FineResults["27"] != nil {
		apres.FineResult = fineResToPoints(image.FineResults["27"])
		apres.ResultRep = resultRep
		apres.Status = 1
		apres.PointType = 27
		return apres
	}

	if image.FineResults["5"] != nil {
		apres.FineResult = fineResToPoints(image.FineResults["5"])
		apres.ResultRep = resultRep
		apres.Status = 1
		apres.PointType = 5
		return apres
	}

	if image.Results["95"] != nil {
		apres.ResultRep = resToPoints(image.Results["95"])
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = 95
		apres.Status = 0
		return apres
	}

	if image.Results["83"] != nil {
		apres.ResultRep = resToPoints(image.Results["83"])
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = 83
		apres.Status = 0
		return apres
	}

	if image.Results["68"] != nil {
		apres.ResultRep = resToPoints(image.Results["68"])
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = 68
		apres.Status = 0
		return apres
	}

	if image.Results["27"] != nil {
		apres.ResultRep = resToPoints(image.Results["27"])
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = 27
		apres.Status = 0
		return apres
	}

	if image.Results["5"] != nil {
		apres.ResultRep = resToPoints(image.Results["5"])
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = 5
		apres.Status = 0
		return apres
	}

	fmt.Println(apres)
	return apres
}

func GetPointTypeResults(image *imagemodel.ImageModel, pointType int64) *AllPointsRep {
	log.Info(fmt.Sprintf("image all res = %s switch point", image.Md5))
	apres := &AllPointsRep{}

	if image.FineResults[strconv.Itoa(int(pointType))] != nil {
		apres.FineResult = fineResToPoints(image.FineResults[strconv.Itoa(int(pointType))])
		apres.ResultRep = &ResultRep{}
		apres.Status = 1
		apres.PointType = pointType
		return apres
	}

	if image.Results[strconv.Itoa(int(pointType))] != nil {
		apres.ResultRep = resToPoints(image.Results[strconv.Itoa(int(pointType))])
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = pointType
		apres.Status = 0
		return apres
	}

	fmt.Println(apres)
	return apres
}

func GetPointTypeNotFineResults(image *imagemodel.ImageModel, pointType int64) *AllPointsRep {
	log.Info(fmt.Sprintf("image all res = %s switch point", image.Md5))
	apres := &AllPointsRep{}

	if image.Results[strconv.Itoa(int(pointType))] != nil {
		apres.ResultRep = resToPoints(image.Results[strconv.Itoa(int(pointType))])
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = pointType
		apres.Status = 0
		return apres
	}

	fmt.Println(apres)
	return apres
}

func GetTaskNotFineResults(image *imagemodel.ImageModel, task *taskmodel.TaskModel) *AllPointsRep {
	log.Info(fmt.Sprintf("get task image all not fine res = %s switch point", image.Md5))
	apres := &AllPointsRep{}

	if image.Results[strconv.Itoa(int(task.PointType))] != nil {
		apres.ResultRep = taskResToPoints(image.Results[strconv.Itoa(int(task.PointType))], task.TaskId)
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = task.PointType
		apres.Status = 0
		return apres
	}

	fmt.Println(apres)
	return apres
}

func GetSmallTaskNotFineResults(image *imagemodel.ImageModel, task *smalltaskmodel.SmallTaskModel) *AllPointsRep {
	log.Info(fmt.Sprintf("get task image all not fine res = %s switch point", image.Md5))
	apres := &AllPointsRep{}

	if image.Results[strconv.Itoa(int(task.PointType))] != nil {
		apres.ResultRep = smallTaskResToPoints(image.Results[strconv.Itoa(int(task.PointType))], task.SmallTaskId, task.Areas)
		apres.FineResult = make([]*imagemodel.Points, 0, 0)
		apres.PointType = task.PointType
		apres.Status = 0
		return apres
	}

	fmt.Println(apres)
	return apres
}

func smallTaskResToPoints(result map[string][]*imagemodel.Points, taskId string, area string) *ResultRep {

	resultRep := &ResultRep{
		LeftEyeBrow:  make([]*imagemodel.Points, 0, 0),
		RightEyeBrow: make([]*imagemodel.Points, 0, 0),
		LeftEye:      make([]*imagemodel.Points, 0, 0),
		RightEye:     make([]*imagemodel.Points, 0, 0),
		LeftEar:      make([]*imagemodel.Points, 0, 0),
		RightEar:     make([]*imagemodel.Points, 0, 0),
		Mouth:        make([]*imagemodel.Points, 0, 0),
		Nouse:        make([]*imagemodel.Points, 0, 0),
		Face:         make([]*imagemodel.Points, 0, 0),
	}

	if result[area] != nil {
		res := make([]*imagemodel.Points, 0, 0)
		for _, point := range result[area] {
			if strings.EqualFold(point.SmallTaskId, taskId) {
				res = append(res, point)
			}
		}
		switch area {
		case "leftEyebrow":
			resultRep.LeftEyeBrow = res
		case "rightEyebrow":
			resultRep.RightEyeBrow = res
		case "leftEye":
			resultRep.LeftEye = res
		case "rightEye":
			resultRep.RightEye = res
		case "leftEar":
			resultRep.LeftEar = res
		case "rightEar":
			resultRep.RightEar = res
		case "mouth":
			resultRep.Mouth = res
		case "nouse":
			resultRep.Nouse = res
		case "face":
			resultRep.Face = res
		}
	}

	return resultRep
}

func taskResToPoints(result map[string][]*imagemodel.Points, taskId string) *ResultRep {
	var areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "leftEar", "rightEar", "mouth", "nouse", "face"}

	resultRep := &ResultRep{
		LeftEyeBrow:  make([]*imagemodel.Points, 0, 0),
		RightEyeBrow: make([]*imagemodel.Points, 0, 0),
		LeftEye:      make([]*imagemodel.Points, 0, 0),
		RightEye:     make([]*imagemodel.Points, 0, 0),
		LeftEar:      make([]*imagemodel.Points, 0, 0),
		RightEar:     make([]*imagemodel.Points, 0, 0),
		Mouth:        make([]*imagemodel.Points, 0, 0),
		Nouse:        make([]*imagemodel.Points, 0, 0),
		Face:         make([]*imagemodel.Points, 0, 0),
	}
	for _, area := range areas {
		if result[area] != nil {
			points := getSmallTaskRes(result[area], taskId, area)
			switch area {
			case "leftEyebrow":
				resultRep.LeftEyeBrow = points
			case "rightEyebrow":
				resultRep.RightEyeBrow = points
			case "leftEye":
				resultRep.LeftEye = points
			case "rightEye":
				resultRep.RightEye = points
			case "leftEar":
				resultRep.LeftEar = points
			case "rightEar":
				resultRep.RightEar = points
			case "mouth":
				resultRep.Mouth = points
			case "nouse":
				resultRep.Nouse = points
			case "face":
				resultRep.Face = points
			}
		}
	}
	return resultRep
}

func getSmallTaskRes(points []*imagemodel.Points, taskId string, area string) []*imagemodel.Points {
	res := make([]*imagemodel.Points, 0, 0)
	smallTask, err := smalltaskmodel.QueryfineTuneTask(taskId, area)
	if err != nil {
		log.Error(fmt.Sprintf("query fineTune small task not found err %s", err))
		return res
	}

	for _, point := range points {
		if strings.EqualFold(point.SmallTaskId, smallTask.SmallTaskId) {
			res = append(res, point)
		}
	}
	return res
}

func resToPoints(result map[string][]*imagemodel.Points) *ResultRep {
	var areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "leftEar", "rightEar", "mouth", "nouse", "face"}
	resultRep := &ResultRep{
		LeftEyeBrow:  make([]*imagemodel.Points, 0, 0),
		RightEyeBrow: make([]*imagemodel.Points, 0, 0),
		LeftEye:      make([]*imagemodel.Points, 0, 0),
		RightEye:     make([]*imagemodel.Points, 0, 0),
		LeftEar:      make([]*imagemodel.Points, 0, 0),
		RightEar:     make([]*imagemodel.Points, 0, 0),
		Mouth:        make([]*imagemodel.Points, 0, 0),
		Nouse:        make([]*imagemodel.Points, 0, 0),
		Face:         make([]*imagemodel.Points, 0, 0),
	}
	for _, area := range areas {
		if result[area] != nil {
			//
			switch area {
			case "leftEyebrow":
				resultRep.LeftEyeBrow = result[area]
			case "rightEyebrow":
				resultRep.RightEyeBrow = result[area]
			case "leftEye":
				resultRep.LeftEye = result[area]
			case "rightEye":
				resultRep.RightEye = result[area]
			case "leftEar":
				resultRep.LeftEar = result[area]
			case "rightEar":
				resultRep.RightEar = result[area]
			case "mouth":
				resultRep.Mouth = result[area]
			case "nouse":
				resultRep.Nouse = result[area]
			case "face":
				resultRep.Face = result[area]
			}
		}
	}
	return resultRep
}

func fineResToPoints(fineResults []*imagemodel.FineResult) []*imagemodel.Points {
	fineRes := make([]*imagemodel.Points, 0, 0)

	for _, fine := range fineResults {
		points := make([]*imagemodel.Point, 0, 0)
		//		var p *imagemodel.Point
		for _, point := range fine.Result {
			if point != nil {
				for _, res := range point {
					points = append(points, res)
				}
			}
		}

		pointsRes := &imagemodel.Points{
			SmallTaskId: fine.SmallTaskId,
			User:        fine.User,
			Points:      points,
			Sys:         fine.Sys,
			CreatedAt:   fine.CreatedAt,
			FinishedAt:  fine.FinishedAt,
		}

		fineRes = append(fineRes, pointsRes)
	}

	return fineRes
}

func SwitchNilPoint(pointType int64) *PointsRep {
	log.Info(fmt.Sprintf("image switch nil point"))
	var (
		nilPoint = []float64{156.8617258703446, 274.1566045936671, 128.2278233297365, 278.321535872301, 178.98792328808722, 282.22615894602035, 157.38234228017384, 269.7313651101186, 156.34110946051538, 280.14369330670337, 142.54477460004057, 270.7725979297771, 142.28446639512595, 278.58184407721564, 171.6992935504779, 272.0741389543502, 169.61682791116095, 280.6643097165326, 271.397336032777, 270.25198151994783, 247.44898118063202, 281.7055425361911, 296.3869237045804, 271.2932143396063, 270.87671962294775, 265.3061208601845, 272.95918526226467, 275.97876202806947, 257.6010011723022, 269.47105690520397, 260.20408322144834, 277.01999484772796, 283.6318216637641, 266.6076666511432, 284.67305448342256, 273.63598818383787, 111.30779001028628, 254.8937974299853, 185.4956284109527, 254.11287281524147, 148.27155510816218, 238.49438052036433, 148.27155510816218, 249.68763333169295, 127.18659051007805, 238.75468872527895, 129.00874794448038, 250.2082497415222, 169.35651970624633, 242.91962000391285, 168.05497868167322, 253.59225640541223, 234.43357887887706, 253.59225640541223, 309.92295830411655, 245.78301025797367, 273.74011782098455, 230.94544257784037, 275.56227525538685, 242.65931179899823, 252.91546142781502, 237.9737641105351, 254.99792706713197, 247.605167692376, 294.56477421415406, 230.68513437292575, 294.8250824190687, 242.138695389169, 212.04706531224383, 335.32904863655455, 212.56768172207308, 356.15370502972405, 201.8950453205737, 281.4452343312765, 223.7609345334017, 279.6230768968741, 189.13994327975738, 325.6976291667617, 240.68096785285192, 324.91670455201785, 171.43898534556328, 342.3573701692493, 256.03915194281444, 341.3161373495908, 189.400251484672, 351.7284655461755, 237.2969611889619, 352.2490819560048, 193.30487455839128, 347.82384247245625, 232.61141350049874, 348.6047670872001, 162.3281981735516, 384.78760757033217, 260.7246996312776, 386.0891485949053, 209.70429146801226, 378.27990244746667, 209.70429146801226, 383.7463747506737, 198.77134686159826, 375.4165121934059, 218.5547704351093, 376.45774501306437, 180.28946431266033, 380.102059881869, 243.0237416970835, 382.184525521186, 185.23532020603807, 378.8005188572959, 237.81757759879113, 380.62267629169827, 186.27655302569656, 382.9654501359298, 237.5572693938765, 384.52729936541755, 213.34860633681694, 400.92671627503853, 213.34860633681694, 417.32613318465957, 185.23532020603807, 395.4602439718316, 238.59850221353497, 396.2411685865754, 173.52145098488023, 398.8442506357216, 187.57809405026964, 411.3390444716233, 237.03665298404727, 411.0787362667087, 249.7917550248636, 399.8854834553801, 179.50853969791646, 406.1328803733309, 244.06497451674198, 405.6122639635017, 50.306121670469956, 242.52551331812026, 57.95918289495975, 297.62755413444677, 84.48979513985772, 338.9540847466917, 374.6683665684292, 233.3418398487325, 371.1020283601722, 290.99999688596137, 345.0765298337353, 330.2806153589366, 81.30101962965364, 255.89286648497293, 343.55100795200894, 242.632656175263, 83.34183595618424, 280.3826624033403, 342.525509425572, 269.15817260742193, 213.78571101597376, 457.7142987932478, 86.64285387311664, 306.28572736467635, 90.92856815883091, 328.42858450753346, 92.35713958740234, 352.7142987932478, 100.2142824445452, 380.5714416503906, 116.64285387311664, 404.1428702218192, 136.64285387311662, 419.8571559361049, 158.07142530168807, 434.8571559361049, 182.35713958740234, 450.5714416503906, 339.4999967302595, 298.42858450753346, 338.07142530168807, 327.7142987932478, 333.7857110159738, 354.8571559361049, 322.35713958740234, 382.7142987932478, 306.6428538731166, 403.42858450753346, 288.07142530168807, 420.5714416503906, 268.7857110159738, 437.00001307896207, 245.92856815883093, 451.28572736467635}
	)
	points := make([]*imagemodel.Point, 0, 0)
	points = append(points, &imagemodel.Point{})
	var p *imagemodel.Point
	for i := 0; i < len(nilPoint); i++ {

		if i%2 == 0 {
			x := nilPoint[i]
			p = &imagemodel.Point{
				X: x,
			}
		}
		if i%2 == 1 {
			y := nilPoint[i]
			p.Y = y
			points = append(points, p)
		}
	}

	switch pointType {
	case 5:
		pRep := &PointsRep{
			LeftEye:  []*imagemodel.Point{points[1]},
			RightEye: []*imagemodel.Point{points[10]},
			Nouse:    []*imagemodel.Point{points[35]},
			Mouth:    []*imagemodel.Point{points[47], points[48]},
		}
		return pRep
	case 27:
		pRep := &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20]},
			RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28]},
			LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5]},
			RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14]},
			Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42]},
			Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60]},
		}
		return pRep
	case 68:
		pRep := &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}},
			RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28], &imagemodel.Point{X: (points[31].X + points[32].X) / 2, Y: (points[31].Y + points[32].Y) / 2}, &imagemodel.Point{X: (points[29].X + points[30].X) / 2, Y: (points[29].Y + points[30].Y) / 2}},
			LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5], points[6], points[8], points[9], points[7]},
			RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14], points[15], points[17], points[18], points[16]},
			Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42], &imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y + points[38].Y) / 2}, &imagemodel.Point{X: (points[37].X+points[38].X)/2 + 2*(points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + 2*(points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[43], points[44]},
			Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60], points[51], points[64], points[62], points[61], points[53], points[54], points[66], points[63]},
			Face:         []*imagemodel.Point{points[79], points[80], points[81], points[82], points[83], points[84], points[85], points[86], points[87], points[88], points[89], points[90], points[91], points[92], points[93], points[94], points[95]},
		}
		return pRep
	case 83:
		pRep := &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}, points[21], points[23], points[25], points[26], points[24]},
			RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28], &imagemodel.Point{X: (points[31].X + points[32].X) / 2, Y: (points[31].Y + points[32].Y) / 2}, &imagemodel.Point{X: (points[29].X + points[30].X) / 2, Y: (points[29].Y + points[30].Y) / 2}, points[29], points[31], points[33], points[34], points[32]},
			LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5], points[6], points[8], points[9], points[7]},
			RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14], points[15], points[17], points[18], points[16]},
			Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42], &imagemodel.Point{X: (points[37].X + points[38].X) / 2, Y: (points[37].Y + points[38].Y) / 2}, &imagemodel.Point{X: (points[37].X+points[38].X)/2 + 2*(points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + 2*(points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[43], points[44], points[39], points[40]},
			Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60], points[51], points[52], points[65], points[64], points[57], points[58], points[62], points[61], points[53], points[54], &imagemodel.Point{X: (points[50].X + points[59].X) / 2, Y: (points[50].Y + points[59].Y) / 2}},
			Face:         []*imagemodel.Point{points[79], points[80], points[81], points[82], points[83], points[84], points[85], points[86], points[87], points[88], points[89], points[90], points[91], points[92], points[93], points[94], points[95]},
		}
		return pRep
	case 95:
		pRep := &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20], &imagemodel.Point{X: (points[23].X + points[24].X) / 2, Y: (points[23].Y + points[24].Y) / 2}, &imagemodel.Point{X: (points[21].X + points[22].X) / 2, Y: (points[21].Y + points[22].Y) / 2}, points[21], points[23], points[25], points[26], points[24]},
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
	return &PointsRep{}

}

/*
func v2() {

	pRep := &PointsRep{
		LeftEyeBrow: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["left_eyebrow_left_corner"].X * width, Y: landmark["left_eyebrow_left_corner"].Y * height},
			&imagemodel.Point{X: ((landmark["left_eyebrow_lower_middle"].X + landmark["left_eyebrow_upper_middle"].X) / 2) * width, Y: (landmark["left_eyebrow_lower_middle"].Y + landmark["left_eyebrow_upper_middle"].Y) / 2 * height},
			&imagemodel.Point{X: landmark["left_eyebrow_right_corner"].X * width, Y: landmark["left_eyebrow_right_corner"].Y * height},
			&imagemodel.Point{X: (landmark["left_eyebrow_upper_left_quarter"].X + landmark["left_eyebrow_lower_left_quarter"].X) / 2 * width, Y: (landmark["left_eyebrow_upper_left_quarter"].Y + landmark["left_eyebrow_lower_left_quarter"].Y) / 2 * height},
			&imagemodel.Point{X: (landmark["left_eyebrow_upper_right_quarter"].X + landmark["left_eyebrow_lower_right_quarter"].X) / 2 * width, Y: (landmark["left_eyebrow_upper_right_quarter"].Y + landmark["left_eyebrow_lower_right_quarter"].Y) / 2 * height},
			&imagemodel.Point{X: landmark["left_eyebrow_upper_middle"].X * width, Y: landmark["left_eyebrow_upper_middle"].Y * height},
			&imagemodel.Point{X: landmark["left_eyebrow_upper_left_quarter"].X * width, Y: landmark["left_eyebrow_upper_left_quarter"].Y * height},
			&imagemodel.Point{X: landmark["left_eyebrow_upper_right_quarter"].X * width, Y: landmark["left_eyebrow_upper_right_quarter"].Y * height},
			&imagemodel.Point{X: landmark["left_eyebrow_lower_right_quarter"].X * width, Y: landmark["left_eyebrow_lower_right_quarter"].Y * height},
			&imagemodel.Point{X: landmark["left_eyebrow_lower_left_quarter"].X * width, Y: landmark["left_eyebrow_lower_left_quarter"].Y * height},
		},
		RightEyeBrow: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["right_eyebrow_left_corner"].X * width, Y: landmark["right_eyebrow_left_corner"].Y * height},
			&imagemodel.Point{X: (landmark["right_eyebrow_lower_middle"].X + landmark["right_eyebrow_upper_middle"].X) / 2 * width, Y: (landmark["right_eyebrow_lower_middle"].Y + landmark["right_eyebrow_upper_middle"].Y) / 2 * height},
			&imagemodel.Point{X: landmark["right_eyebrow_right_corner"].X * width, Y: landmark["right_eyebrow_right_corner"].Y * height},
			&imagemodel.Point{X: (landmark["right_eyebrow_upper_left_quarter"].X + landmark["right_eyebrow_lower_left_quarter"].X) / 2 * width, Y: (landmark["right_eyebrow_upper_left_quarter"].Y + landmark["right_eyebrow_lower_left_quarter"].Y) / 2 * height},
			&imagemodel.Point{X: (landmark["right_eyebrow_upper_right_quarter"].X + landmark["right_eyebrow_lower_right_quarter"].X) / 2 * width, Y: (landmark["right_eyebrow_upper_right_quarter"].Y + landmark["right_eyebrow_lower_right_quarter"].Y) / 2 * height},
			&imagemodel.Point{X: landmark["right_eyebrow_upper_middle"].X * width, Y: landmark["right_eyebrow_upper_middle"].Y * height},
			&imagemodel.Point{X: landmark["right_eyebrow_upper_left_quarter"].X * width, Y: landmark["right_eyebrow_upper_left_quarter"].Y * height},
			&imagemodel.Point{X: landmark["right_eyebrow_upper_right_quarter"].X * width, Y: landmark["right_eyebrow_upper_right_quarter"].Y * height},
			&imagemodel.Point{X: landmark["right_eyebrow_lower_right_quarter"].X * width, Y: landmark["right_eyebrow_lower_right_quarter"].Y * height},
			&imagemodel.Point{X: landmark["right_eyebrow_lower_left_quarter"].X * width, Y: landmark["right_eyebrow_lower_left_quarter"].Y * height},
		},
		LeftEye: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["left_eye_pupil"].X * width, Y: landmark["left_eye_pupil"].Y * height},
			&imagemodel.Point{X: landmark["left_eye_left_corner"].X * width, Y: landmark["left_eye_left_corner"].Y * height},
			&imagemodel.Point{X: landmark["left_eye_top"].X * width, Y: landmark["left_eye_top"].Y * height},
			&imagemodel.Point{X: landmark["left_eye_right_corner"].X * width, Y: landmark["left_eye_right_corner"].Y * height},
			&imagemodel.Point{X: landmark["left_eye_bottom"].X * width, Y: landmark["left_eye_bottom"].Y * height},
			&imagemodel.Point{X: landmark["left_eye_upper_left_quarter"].X * width, Y: landmark["left_eye_upper_left_quarter"].Y * height},
			&imagemodel.Point{X: landmark["left_eye_upper_right_quarter"].X * width, Y: landmark["left_eye_upper_right_quarter"].Y * height},
			&imagemodel.Point{X: landmark["left_eye_lower_right_quarter"].X * width, Y: landmark["left_eye_lower_right_quarter"].Y * height},
			&imagemodel.Point{X: landmark["left_eye_lower_left_quarter"].X * width, Y: landmark["left_eye_lower_left_quarter"].Y * height},
		},
		RightEye: []*imagemodel.Point{

			&imagemodel.Point{X: landmark["right_eye_center"].X * width, Y: landmark["right_eye_center"].Y * height},
			&imagemodel.Point{X: landmark["right_eye_left_corner"].X * width, Y: landmark["right_eye_left_corner"].Y * height},
			&imagemodel.Point{X: landmark["right_eye_top"].X * width, Y: landmark["right_eye_top"].Y * height},
			&imagemodel.Point{X: landmark["right_eye_right_corner"].X * width, Y: landmark["right_eye_right_corner"].Y * height},
			&imagemodel.Point{X: landmark["right_eye_bottom"].X * width, Y: landmark["right_eye_bottom"].Y * height},
			&imagemodel.Point{X: landmark["right_eye_upper_left_quarter"].X * width, Y: landmark["right_eye_upper_left_quarter"].Y * height},
			&imagemodel.Point{X: landmark["right_eye_upper_right_quarter"].X * width, Y: landmark["right_eye_upper_right_quarter"].Y * height},
			&imagemodel.Point{X: landmark["right_eye_lower_right_quarter"].X * width, Y: landmark["right_eye_lower_right_quarter"].Y * height},
			&imagemodel.Point{X: landmark["right_eye_lower_left_quarter"].X * width, Y: landmark["right_eye_lower_left_quarter"].Y * height},
		},
		Nouse: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["nose_tip"].X * width, Y: landmark["nose_tip"].Y * height},
			&imagemodel.Point{X: (landmark["nose_tip"].X + ((landmark["nose_contour_left1"].X+landmark["nose_contour_right1"].X)/2-landmark["nose_tip"].X)*2/3) * width, Y: (landmark["nose_tip"].Y + ((landmark["nose_contour_left1"].Y+landmark["nose_contour_right1"].Y)/2-landmark["nose_tip"].Y)*2/3) * height},
			&imagemodel.Point{X: landmark["nose_left"].X * width, Y: landmark["nose_left"].Y * height},
			&imagemodel.Point{X: landmark["nose_contour_lower_middle"].X * width, Y: landmark["nose_contour_lower_middle"].Y * height},
			&imagemodel.Point{X: landmark["nose_right"].X * width, Y: landmark["nose_right"].Y * height},
			&imagemodel.Point{X: (landmark["nose_contour_left1"].X + landmark["nose_contour_right1"].X) / 2 * width, Y: (landmark["nose_contour_left1"].Y + landmark["nose_contour_right1"].Y) / 2 * height},
			&imagemodel.Point{X: (landmark["nose_tip"].X + ((landmark["nose_contour_left1"].X+landmark["nose_contour_right1"].X)/2-landmark["nose_tip"].X)/3) * width, Y: (landmark["nose_tip"].Y + ((landmark["nose_contour_left1"].Y+landmark["nose_contour_right1"].Y)/2-landmark["nose_tip"].Y)/3) * height},
			&imagemodel.Point{X: landmark["nose_contour_left3"].X * width, Y: landmark["nose_contour_left3"].Y * height},
			&imagemodel.Point{X: landmark["nose_contour_right3"].X * width, Y: landmark["nose_contour_right3"].Y * height},
			&imagemodel.Point{X: landmark["nose_contour_left2"].X * width, Y: landmark["nose_contour_left2"].Y * height},
			&imagemodel.Point{X: landmark["nose_contour_right2"].X * width, Y: landmark["nose_contour_right2"].Y * height},
		},
		Mouth: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["mouth_left_corner"].X * width, Y: landmark["mouth_left_corner"].Y * height},
			&imagemodel.Point{X: landmark["mouth_right_corner"].X * width, Y: landmark["mouth_right_corner"].Y * height},
			&imagemodel.Point{X: landmark["mouth_upper_lip_top"].X * width, Y: landmark["mouth_upper_lip_top"].Y * height},
			&imagemodel.Point{X: landmark["mouth_upper_lip_bottom"].X * width, Y: landmark["mouth_upper_lip_bottom"].Y * height},
			&imagemodel.Point{X: landmark["mouth_lower_lip_top"].X * width, Y: landmark["mouth_lower_lip_top"].Y * height},
			&imagemodel.Point{X: landmark["mouth_lower_lip_bottom"].X * width, Y: landmark["mouth_lower_lip_bottom"].Y * height},
			&imagemodel.Point{X: landmark["mouth_upper_lip_left_contour1"].X * width, Y: landmark["mouth_upper_lip_left_contour1"].Y * height},
			&imagemodel.Point{X: landmark["mouth_upper_lip_right_contour1"].X * width, Y: landmark["mouth_upper_lip_right_contour1"].Y * height},
			&imagemodel.Point{X: landmark["mouth_lower_lip_right_contour3"].X * width, Y: landmark["mouth_lower_lip_right_contour3"].Y * height},
			&imagemodel.Point{X: landmark["mouth_lower_lip_left_contour3"].X * width, Y: landmark["mouth_lower_lip_left_contour3"].Y * height},
			&imagemodel.Point{X: landmark["mouth_upper_lip_left_contour3"].X * width, Y: landmark["mouth_upper_lip_left_contour3"].Y * height},
			&imagemodel.Point{X: landmark["mouth_upper_lip_right_contour3"].X * width, Y: landmark["mouth_upper_lip_right_contour3"].Y * height},
			&imagemodel.Point{X: landmark["mouth_lower_lip_right_contour1"].X * width, Y: landmark["mouth_lower_lip_right_contour1"].Y * height},
			&imagemodel.Point{X: landmark["mouth_lower_lip_left_contour1"].X * width, Y: landmark["mouth_lower_lip_left_contour1"].Y * height},
			&imagemodel.Point{X: landmark["mouth_upper_lip_left_contour2"].X * width, Y: landmark["mouth_upper_lip_left_contour2"].Y * height},
			&imagemodel.Point{X: landmark["mouth_upper_lip_right_contour2"].X * width, Y: landmark["mouth_upper_lip_right_contour2"].Y * height},
			&imagemodel.Point{X: (landmark["mouth_left_corner"].X + landmark["mouth_right_corner"].X) / 2 * width, Y: (landmark["mouth_left_corner"].Y + landmark["mouth_right_corner"].Y) / 2 * height},
			//			&imagemodel.Point{X: landmark["mouth_lower_lip_right_contour2"].X * width, Y: landmark["mouth_lower_lip_right_contour2"].Y * height},
			//			&imagemodel.Point{X: landmark["mouth_lower_lip_left_contour2"].X * width, Y: landmark["mouth_lower_lip_left_contour2"].Y * height},
		},
		Face: []*imagemodel.Point{
			&imagemodel.Point{X: landmark["contour_chin"].X * width, Y: landmark["contour_chin"].Y * height},
			&imagemodel.Point{X: landmark["contour_left1"].X * width, Y: landmark["contour_left1"].Y * height},
			&imagemodel.Point{X: landmark["contour_left2"].X * width, Y: landmark["contour_left2"].Y * height},
			&imagemodel.Point{X: landmark["contour_left3"].X * width, Y: landmark["contour_left3"].Y * height},
			&imagemodel.Point{X: landmark["contour_left4"].X * width, Y: landmark["contour_left4"].Y * height},
			&imagemodel.Point{X: landmark["contour_left5"].X * width, Y: landmark["contour_left5"].Y * height},
			&imagemodel.Point{X: landmark["contour_left6"].X * width, Y: landmark["contour_left6"].Y * height},
			&imagemodel.Point{X: landmark["contour_left7"].X * width, Y: landmark["contour_left7"].Y * height},
			&imagemodel.Point{X: landmark["contour_left8"].X * width, Y: landmark["contour_left8"].Y * height},
			&imagemodel.Point{X: landmark["contour_right1"].X * width, Y: landmark["contour_right1"].Y * height},
			&imagemodel.Point{X: landmark["contour_right2"].X * width, Y: landmark["contour_right2"].Y * height},
			&imagemodel.Point{X: landmark["contour_right3"].X * width, Y: landmark["contour_right3"].Y * height},
			&imagemodel.Point{X: landmark["contour_right4"].X * width, Y: landmark["contour_right4"].Y * height},
			&imagemodel.Point{X: landmark["contour_right5"].X * width, Y: landmark["contour_right5"].Y * height},
			&imagemodel.Point{X: landmark["contour_right6"].X * width, Y: landmark["contour_right6"].Y * height},
			&imagemodel.Point{X: landmark["contour_right7"].X * width, Y: landmark["contour_right7"].Y * height},
			&imagemodel.Point{X: landmark["contour_right8"].X * width, Y: landmark["contour_right8"].Y * height},
		},
	}
}
*/
