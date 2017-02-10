package switch_point

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	"encoding/json"
	"fmt"
	"strconv"

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

func SwitchAlreadyPoint(pointType int64, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch 83 point", image.Md5))

	fineCount, fineres := getFineTuneRes(image)
	var finepRep *PointsRep
	if fineres != nil {
		var areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "leftEar", "rightEar", "mouth", "nouse", "face"}
		pRep := &PointsRep{}
		for _, area := range areas {
			if image.Results[strconv.Itoa(int(fineCount))][area] != nil {
				//
				switch area {
				case "leftEyebrow":
					pRep.LeftEyeBrow = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				case "rightEyebrow":
					pRep.RightEyeBrow = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				case "leftEye":
					pRep.LeftEye = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				case "rightEye":
					pRep.RightEye = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				case "leftEar":
					pRep.LeftEar = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				case "rightEar":
					pRep.RightEar = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				case "mouth":
					pRep.Mouth = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				case "nouse":
					pRep.Nouse = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				case "face":
					pRep.Face = image.Results[strconv.Itoa(int(fineCount))][area][0].Points
				}
			}
		}
		finepRep = pRep
	}

	//import point
	_, thrRes := getImportPoint(pointType, image)
	if thrRes == nil {
		//face++ point
		_, thrRes = faceResSwitchPoint(pointType, image)
	}

	// nil
	nilRep := SwitchNilPoint(pointType)

	if finepRep == nil && thrRes == nil {
		return nilRep
	}
	pRep := &PointsRep{}
	switch pointType {

	case 5:
		if finepRep == nil {
			if thrRes == nil {
				return nilRep
			} else {
				return thrRes
			}
		} else {
			if fineCount >= int(pointType) {
				rep := switchFullPoint(fineCount, pointType, finepRep)
				pRep = rep
			} else {
				if thrRes != nil {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, thrRes)
					pRep = rep
				} else {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, nilRep)
					pRep = rep
				}
			}
		}
		return pRep
	case 27:
		if finepRep == nil {
			if thrRes == nil {
				return nilRep
			} else {
				return thrRes
			}
		} else {
			if fineCount >= int(pointType) {
				rep := switchFullPoint(fineCount, pointType, finepRep)
				pRep = rep
			} else {
				if thrRes != nil {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, thrRes)
					pRep = rep
				} else {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, nilRep)
					pRep = rep
				}
			}
		}
		return pRep
	case 68:
		if finepRep == nil {
			if thrRes == nil {
				return nilRep
			} else {
				return thrRes
			}
		} else {
			if fineCount >= int(pointType) {
				rep := switchFullPoint(fineCount, pointType, finepRep)
				pRep = rep
			} else {
				if thrRes != nil {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, thrRes)
					pRep = rep
				} else {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, nilRep)
					pRep = rep
				}
			}
		}
		return pRep
	case 83:
		if finepRep == nil {
			if thrRes == nil {
				return nilRep
			} else {
				return thrRes
			}
		} else {
			if fineCount >= int(pointType) {
				rep := switchFullPoint(fineCount, pointType, finepRep)
				pRep = rep
			} else {
				if thrRes != nil {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, thrRes)
					pRep = rep
				} else {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, nilRep)
					pRep = rep
				}
			}
		}
		return pRep
	case 95:
		if finepRep == nil {
			if thrRes == nil {
				return nilRep
			} else {
				return thrRes
			}
		} else {
			if fineCount >= int(pointType) {
				rep := switchFullPoint(fineCount, pointType, finepRep)
				pRep = rep
			} else {
				if thrRes != nil {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, thrRes)
					pRep = rep
				} else {
					rep := switchNotFullPoint(fineCount, pointType, finepRep, nilRep)
					pRep = rep
				}
			}
		}
		return pRep
	}

	return nilRep
}

func switchFullPoint(fineCount int, pointType int64, finepRep *PointsRep) *PointsRep {
	pRep := &PointsRep{}
	switch pointType {
	case 5:
		pRep.LeftEye = append(pRep.LeftEye, finepRep.LeftEye[0])
		pRep.RightEye = append(pRep.RightEye, finepRep.RightEye[0])
		pRep.Nouse = append(pRep.Nouse, finepRep.Nouse[0])
		pRep.Mouth = finepRep.Mouth[0:2]
		return pRep
	case 27:

		pRep.LeftEyeBrow = finepRep.LeftEyeBrow[0:3]
		pRep.RightEyeBrow = finepRep.RightEyeBrow[0:3]
		pRep.LeftEye = finepRep.LeftEye[0:5]
		pRep.RightEye = finepRep.RightEye[0:5]
		pRep.Nouse = finepRep.Nouse[0:5]
		pRep.Mouth = finepRep.Mouth[0:6]

		return pRep
	case 68:
		pRep.LeftEyeBrow = finepRep.LeftEyeBrow[0:5]
		pRep.RightEyeBrow = finepRep.RightEyeBrow[0:5]
		pRep.LeftEye = finepRep.LeftEye[0:9]
		pRep.RightEye = finepRep.RightEye[0:9]
		pRep.Nouse = finepRep.Nouse[0:9]
		pRep.Mouth = finepRep.Mouth[0:14]
		pRep.Face = finepRep.Face[0:17]

		return pRep
	case 83:
		pRep.LeftEyeBrow = finepRep.LeftEyeBrow[0:10]
		pRep.RightEyeBrow = finepRep.RightEyeBrow[0:10]
		pRep.LeftEye = finepRep.LeftEye[0:9]
		pRep.RightEye = finepRep.RightEye[0:9]
		pRep.Nouse = finepRep.Nouse[0:11]
		pRep.Mouth = finepRep.Mouth[0:17]
		pRep.Face = finepRep.Face[0:17]

		return pRep
	case 95:
		pRep.LeftEyeBrow = finepRep.LeftEyeBrow[0:10]
		pRep.RightEyeBrow = finepRep.RightEyeBrow[0:10]
		pRep.LeftEye = finepRep.LeftEye[0:9]
		pRep.RightEye = finepRep.RightEye[0:9]
		pRep.LeftEar = finepRep.LeftEar[0:]
		pRep.RightEar = finepRep.RightEar[0:]
		pRep.Nouse = finepRep.Nouse[0:11]
		pRep.Mouth = finepRep.Mouth[0:19]
		pRep.Face = finepRep.Face[0:17]

		return pRep
	}

	return nil
}

func switchNotFullPoint(fineCount int, pointType int64, finepRep *PointsRep, thrRes *PointsRep) *PointsRep {
	pRep := &PointsRep{}
	switch pointType {
	case 5:
		pRep.LeftEye = append(pRep.LeftEye, finepRep.LeftEye[0])
		pRep.RightEye = append(pRep.RightEye, finepRep.RightEye[0])
		pRep.Nouse = append(pRep.Nouse, finepRep.Nouse[0])
		pRep.Mouth = finepRep.Mouth[0:2]
		return pRep
	case 27:

		pRep.LeftEyeBrow = thrRes.LeftEyeBrow
		pRep.RightEyeBrow = thrRes.RightEyeBrow
		pRep.LeftEye = finepRep.LeftEye
		pRep.LeftEye = append(pRep.LeftEye, thrRes.LeftEye[1], thrRes.LeftEye[2], thrRes.LeftEye[3], thrRes.LeftEye[4])
		pRep.RightEye = finepRep.RightEye
		pRep.RightEye = append(pRep.RightEye, thrRes.RightEye[1], thrRes.RightEye[2], thrRes.RightEye[3], thrRes.RightEye[4])
		pRep.Nouse = finepRep.Nouse
		pRep.Nouse = append(pRep.Nouse, thrRes.Nouse[1], thrRes.Nouse[2], thrRes.Nouse[3], thrRes.Nouse[4])
		pRep.Mouth = finepRep.Mouth[0:2]
		pRep.Mouth = append(pRep.Mouth, thrRes.Mouth[2], thrRes.Mouth[3], thrRes.Mouth[4], thrRes.Mouth[5])

		return pRep
	case 68:
		rep := &PointsRep{
			LeftEyeBrow:  make([]*imagemodel.Point, 5),
			RightEyeBrow: make([]*imagemodel.Point, 5),
			LeftEye:      make([]*imagemodel.Point, 9),
			RightEye:     make([]*imagemodel.Point, 9),
			Mouth:        make([]*imagemodel.Point, 14),
			Nouse:        make([]*imagemodel.Point, 9),
			Face:         make([]*imagemodel.Point, 17),
		}
		if fineCount == 5 {
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)

			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])

			copy(rep.LeftEyeBrow, thrRes.LeftEyeBrow)
			copy(rep.RightEyeBrow, thrRes.RightEyeBrow)
			copy(rep.Face, thrRes.Face)

		} else if fineCount == 27 {
			copy(rep.LeftEyeBrow, finepRep.LeftEyeBrow)
			copy(rep.RightEyeBrow, finepRep.RightEyeBrow)
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)

			copy(rep.LeftEyeBrow[len(finepRep.LeftEyeBrow):], thrRes.LeftEyeBrow[len(finepRep.LeftEyeBrow):])
			copy(rep.RightEyeBrow[len(finepRep.RightEyeBrow):], thrRes.RightEyeBrow[len(finepRep.RightEyeBrow):])
			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])

			copy(rep.Face, thrRes.Face)
		}

		return rep
	case 83:
		rep := &PointsRep{
			LeftEyeBrow:  make([]*imagemodel.Point, 10),
			RightEyeBrow: make([]*imagemodel.Point, 10),
			LeftEye:      make([]*imagemodel.Point, 9),
			RightEye:     make([]*imagemodel.Point, 9),
			Mouth:        make([]*imagemodel.Point, 17),
			Nouse:        make([]*imagemodel.Point, 11),
			Face:         make([]*imagemodel.Point, 17),
		}

		if fineCount == 5 {
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)

			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])

			copy(rep.LeftEyeBrow, thrRes.LeftEyeBrow)
			copy(rep.RightEyeBrow, thrRes.RightEyeBrow)
			copy(rep.Face, thrRes.Face)

		} else if fineCount == 27 {
			copy(rep.LeftEyeBrow, finepRep.LeftEyeBrow)
			copy(rep.RightEyeBrow, finepRep.RightEyeBrow)
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)

			copy(rep.LeftEyeBrow[len(finepRep.LeftEyeBrow):], thrRes.LeftEyeBrow[len(finepRep.LeftEyeBrow):])
			copy(rep.RightEyeBrow[len(finepRep.RightEyeBrow):], thrRes.RightEyeBrow[len(finepRep.RightEyeBrow):])
			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])

			copy(rep.Face, thrRes.Face)
		} else if fineCount == 68 {
			copy(rep.LeftEyeBrow, finepRep.LeftEyeBrow)
			copy(rep.RightEyeBrow, finepRep.RightEyeBrow)
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)
			copy(rep.Face, finepRep.Face)

			copy(rep.LeftEyeBrow[len(finepRep.LeftEyeBrow):], thrRes.LeftEyeBrow[len(finepRep.LeftEyeBrow):])
			copy(rep.RightEyeBrow[len(finepRep.RightEyeBrow):], thrRes.RightEyeBrow[len(finepRep.RightEyeBrow):])
			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])
		}

		return rep
	case 95:
		rep := &PointsRep{
			LeftEyeBrow:  make([]*imagemodel.Point, 10),
			RightEyeBrow: make([]*imagemodel.Point, 10),
			LeftEye:      make([]*imagemodel.Point, 9),
			RightEye:     make([]*imagemodel.Point, 9),
			LeftEar:      make([]*imagemodel.Point, 5),
			RightEar:     make([]*imagemodel.Point, 5),
			Mouth:        make([]*imagemodel.Point, 19),
			Nouse:        make([]*imagemodel.Point, 11),
			Face:         make([]*imagemodel.Point, 17),
		}

		if fineCount == 5 {
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)

			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])

			copy(rep.LeftEyeBrow, thrRes.LeftEyeBrow)
			copy(rep.RightEyeBrow, thrRes.RightEyeBrow)
			copy(rep.LeftEar, thrRes.LeftEar)
			copy(rep.RightEar, thrRes.RightEar)
			copy(rep.Face, thrRes.Face)

		} else if fineCount == 27 {
			copy(rep.LeftEyeBrow, finepRep.LeftEyeBrow)
			copy(rep.RightEyeBrow, finepRep.RightEyeBrow)
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)

			copy(rep.LeftEyeBrow[len(finepRep.LeftEyeBrow):], thrRes.LeftEyeBrow[len(finepRep.LeftEyeBrow):])
			copy(rep.RightEyeBrow[len(finepRep.RightEyeBrow):], thrRes.RightEyeBrow[len(finepRep.RightEyeBrow):])
			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])

			copy(rep.LeftEar, thrRes.LeftEar)
			copy(rep.RightEar, thrRes.RightEar)
			copy(rep.Face, thrRes.Face)

		} else if fineCount == 68 {
			copy(rep.LeftEyeBrow, finepRep.LeftEyeBrow)
			copy(rep.RightEyeBrow, finepRep.RightEyeBrow)
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)
			copy(rep.Face, finepRep.Face)
			copy(rep.LeftEar, thrRes.LeftEar)
			copy(rep.RightEar, thrRes.RightEar)

			copy(rep.LeftEyeBrow[len(finepRep.LeftEyeBrow):], thrRes.LeftEyeBrow[len(finepRep.LeftEyeBrow):])
			copy(rep.RightEyeBrow[len(finepRep.RightEyeBrow):], thrRes.RightEyeBrow[len(finepRep.RightEyeBrow):])
			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])
		} else if fineCount == 83 {
			copy(rep.LeftEyeBrow, finepRep.LeftEyeBrow)
			copy(rep.RightEyeBrow, finepRep.RightEyeBrow)
			copy(rep.LeftEye, finepRep.LeftEye)
			copy(rep.RightEye, finepRep.RightEye)
			copy(rep.Nouse, finepRep.Nouse)
			copy(rep.Mouth, finepRep.Mouth)
			copy(rep.Face, finepRep.Face)
			copy(rep.LeftEar, thrRes.LeftEar)
			copy(rep.RightEar, thrRes.RightEar)

			copy(rep.LeftEyeBrow[len(finepRep.LeftEyeBrow):], thrRes.LeftEyeBrow[len(finepRep.LeftEyeBrow):])
			copy(rep.RightEyeBrow[len(finepRep.RightEyeBrow):], thrRes.RightEyeBrow[len(finepRep.RightEyeBrow):])
			copy(rep.LeftEye[len(finepRep.LeftEye):], thrRes.LeftEye[len(finepRep.LeftEye):])
			copy(rep.RightEye[len(finepRep.RightEye):], thrRes.RightEye[len(finepRep.RightEye):])
			copy(rep.Nouse[len(finepRep.Nouse):], thrRes.Nouse[len(finepRep.Nouse):])
			copy(rep.Mouth[len(finepRep.Mouth):], thrRes.Mouth[len(finepRep.Mouth):])
		}

		return rep
	}

	return nil
}

func getImportPoint(pointType int64, image *imagemodel.ImageModel) (int, *PointsRep) {
	log.Info(fmt.Sprintf("image = %s import point", image.Md5))

	if image.ThrFaces["deepir_import"] == nil {
		return 0, nil
	}

	var (
		res1B []byte
		count int = 0
	)
	if image.ThrFaces["deepir_import"]["95"] == nil {
		if image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))] == nil {
			return 0, nil
		} else {
			res, err := json.Marshal(image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))])
			if err != nil {
				fmt.Println("json Marshal deepir import err=%s", err)
				return 0, nil
			}
			res1B = res
			count = int(pointType)
		}
	} else {
		res, err := json.Marshal(image.ThrFaces["deepir_import"]["95"])
		if err != nil {
			fmt.Println("json Marshal 95 err=%s", err)
			return 0, nil
		}
		res1B = res
		count = 95
	}

	fmt.Println(string(res1B)) //json
	imPoints := new(Points)
	if err := json.Unmarshal(res1B, &imPoints.Points); err != nil {
		fmt.Println("json unmarshal err=%s", err)
		return 0, nil
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
		return 0, nil
	}
	var pRep *PointsRep
	switch pointType {
	case 5:
		pRep = &PointsRep{
			LeftEye:  []*imagemodel.Point{points[1]},
			RightEye: []*imagemodel.Point{points[10]},
			Nouse:    []*imagemodel.Point{points[35]},
			Mouth:    []*imagemodel.Point{points[47], points[48]},
		}
		return int(pointType), pRep
	case 27:
		pRep = &PointsRep{
			LeftEyeBrow:  []*imagemodel.Point{points[19], &imagemodel.Point{X: (points[21].X + points[22].X + points[23].X + points[24].X) / 4, Y: (points[21].Y + points[22].Y + points[23].Y + points[24].Y) / 4}, points[20]},
			RightEyeBrow: []*imagemodel.Point{points[27], &imagemodel.Point{X: (points[31].X + points[32].X + points[29].X + points[30].X) / 4, Y: (points[31].Y + points[32].Y + points[29].Y + points[30].Y) / 4}, points[28]},
			LeftEye:      []*imagemodel.Point{points[1], points[2], points[4], points[3], points[5]},
			RightEye:     []*imagemodel.Point{points[10], points[11], points[13], points[12], points[14]},
			Nouse:        []*imagemodel.Point{points[35], &imagemodel.Point{X: (points[37].X+points[38].X)/2 + (points[35].X-((points[37].X+points[38].X)/2))/3, Y: (points[37].Y+points[38].Y)/2 + (points[35].Y-((points[37].Y+points[38].Y)/2))/3}, points[41], points[36], points[42]},
			Mouth:        []*imagemodel.Point{points[47], points[48], points[49], points[50], points[59], points[60]},
		}
		return int(pointType), pRep
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
		return int(pointType), pRep
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
		return int(pointType), pRep
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
		return int(pointType), pRep
	}

	return count, pRep
}

func getFineTuneRes(image *imagemodel.ImageModel) (int, *imagemodel.FineResult) {
	log.Info(fmt.Sprintf("image = %s fineTune point", image.Md5))
	var fineres *imagemodel.FineResult
	if image.FineResults != nil {
		if len(image.FineResults["95"]) != 0 {
			if image.FineResults["95"][0].Result != nil {
				fineres = image.FineResults["95"][0]
				return 95, fineres
			}
		}
	}

	if image.FineResults != nil {
		if len(image.FineResults["83"]) != 0 {
			if image.FineResults["83"][0].Result != nil {
				fineres = image.FineResults["83"][0]
				return 83, fineres
			}
		}
	}

	if image.FineResults != nil {
		if len(image.FineResults["68"]) != 0 {
			if image.FineResults["68"][0].Result != nil {
				fineres = image.FineResults["68"][0]
				return 68, fineres
			}
		}
	}

	if image.FineResults != nil {
		if len(image.FineResults["27"]) != 0 {
			if image.FineResults["27"][0].Result != nil {
				fineres = image.FineResults["27"][0]
				return 27, fineres
			}
		}
	}

	if image.FineResults != nil {
		if len(image.FineResults["5"]) != 0 {
			if image.FineResults["5"][0].Result != nil {
				fineres = image.FineResults["5"][0]
				return 5, fineres
			}
		}
	}
	return 0, nil
}

func switchAllPoint(pointType int64, image *imagemodel.ImageModel) *PointsRep {
	log.Info(fmt.Sprintf("image = %s switch all point", image.Md5))
	if image.ThrFaces["deepir_import"][strconv.Itoa(int(pointType))] == nil {
		if image.ThrFaces["face++"] == nil {
			log.Info(fmt.Sprintf("image = %s thr point =nil", image.Md5))
			return SwitchNilPoint(pointType)
		}
		_, faceRes := faceResSwitchPoint(pointType, image)
		return faceRes
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

func faceResSwitchPoint(pointType int64, image *imagemodel.ImageModel) (int, *PointsRep) {
	log.Info(fmt.Sprintf("image = %s switch face++ point", image.Md5))
	if image.ThrFaces["face++"]["83"] == nil {
		log.Info(fmt.Sprintf("image = %s thr point =nil", image.Md5))
		return 0, nil
	}

	res1B, _ := json.Marshal(image.ThrFaces["face++"]["83"])
	//	fmt.Println(string(res1B)) //json
	faceRes := &thrfacemodel.FaceModelV3{}
	if err := json.Unmarshal([]byte(string(res1B)), &faceRes); err != nil {
		fmt.Println("json unmarshal err=%s", err)
		return 0, nil
	}

	if faceRes.Faces == nil || len(faceRes.Faces) == 0 {
		log.Info(fmt.Sprintf("image = %s face++ nil no face point", image.Md5))
		return 0, nil
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
		return 0, nil
	}

	landmark := faceRes.Faces[0].Landmark
	if landmark == nil {
		log.Info(fmt.Sprintf("face++ unmarshal landmark==nil"))
		return 0, nil
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
		return 0, nil
	}

	pRepre := &PointsRep{}
	switch pointType {
	case 5:
		pRepre.LeftEye = append(pRepre.LeftEye, pRep.LeftEye[0])
		pRepre.RightEye = append(pRepre.RightEye, pRep.RightEye[0])
		pRepre.Nouse = append(pRepre.Nouse, pRep.Nouse[0])
		pRepre.Mouth = pRep.Mouth[0:2]
		return int(pointType), pRepre
	case 27:
		pRepre.LeftEyeBrow = pRep.LeftEyeBrow[0:3]
		pRepre.RightEyeBrow = pRep.RightEyeBrow[0:3]
		pRepre.LeftEye = pRep.LeftEye[0:5]
		pRepre.RightEye = pRep.RightEye[0:5]
		pRepre.Nouse = pRep.Nouse[0:5]
		pRepre.Mouth = pRep.Mouth[0:6]
		return int(pointType), pRepre
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
		return int(pointType), pRepre
	case 83:
		pRepre = pRep
		return int(pointType), pRepre
	case 95:
		pRep.Mouth = append(pRep.Mouth, &imagemodel.Point{X: landmark["mouth_lower_lip_right_contour2"].X, Y: landmark["mouth_lower_lip_right_contour2"].Y})
		pRep.Mouth = append(pRep.Mouth, &imagemodel.Point{X: landmark["mouth_lower_lip_left_contour2"].X, Y: landmark["mouth_lower_lip_left_contour2"].Y})
		pp := SwitchNilPoint(pointType)
		pRep.LeftEar = pp.LeftEar
		pRep.RightEar = pp.RightEar
		pRepre = pRep
		return int(pointType), pRepre
	}

	return int(pointType), pRepre
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
