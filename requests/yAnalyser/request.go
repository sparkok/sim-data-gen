package yAnalyser
import (
	"bytes"
	"sim_data_gen/requests/common"
	yAnalyserModel "sim_data_gen/models/yAnalyser"
)
import "time"


/**
* YAnalyser请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"analyserNum"`
	AnalyserNum *string  
	//`json:"createdAt"`
	CreatedAt *string  
	//`json:"crushingPlant"`
	CrushingPlant *string  
	//`json:"flux"`
	Flux *string  
	//`json:"load"`
	Load *string  
	//`json:"mat1"`
	Mat1 *string  
	//`json:"mat10"`
	Mat10 *string  
	//`json:"mat11"`
	Mat11 *string  
	//`json:"mat12"`
	Mat12 *string  
	//`json:"mat13"`
	Mat13 *string  
	//`json:"mat14"`
	Mat14 *string  
	//`json:"mat15"`
	Mat15 *string  
	//`json:"mat16"`
	Mat16 *string  
	//`json:"mat17"`
	Mat17 *string  
	//`json:"mat18"`
	Mat18 *string  
	//`json:"mat19"`
	Mat19 *string  
	//`json:"mat2"`
	Mat2 *string  
	//`json:"mat20"`
	Mat20 *string  
	//`json:"mat3"`
	Mat3 *string  
	//`json:"mat4"`
	Mat4 *string  
	//`json:"mat5"`
	Mat5 *string  
	//`json:"mat6"`
	Mat6 *string  
	//`json:"mat7"`
	Mat7 *string  
	//`json:"mat8"`
	Mat8 *string  
	//`json:"mat9"`
	Mat9 *string  
	//`json:"speed"`
	Speed *string  
	//`json:"status"`
	Status *string  
	//`json:"testAt"`
	TestAt *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.AnalyserNum != nil {
      condition.WriteString("and (y_analyser.analyser_num like '%" + *this.AnalyserNum + "%')")
    }
    if this.CreatedAt != nil {
      condition.WriteString("and (y_analyser.created_at = '" + *this.CreatedAt + "')")
    }
    if this.CrushingPlant != nil {
      condition.WriteString("and (y_analyser.crushing_plant like '%" + *this.CrushingPlant + "%')")
    }
    if this.Flux != nil {
     condition.WriteString("and (y_analyser.flux = " + *this.Flux + ")")
    }
    if this.Load != nil {
     condition.WriteString("and (y_analyser.load = " + *this.Load + ")")
    }
    if this.Mat1 != nil {
     condition.WriteString("and (y_analyser.mat1 = " + *this.Mat1 + ")")
    }
    if this.Mat10 != nil {
     condition.WriteString("and (y_analyser.mat10 = " + *this.Mat10 + ")")
    }
    if this.Mat11 != nil {
     condition.WriteString("and (y_analyser.mat11 = " + *this.Mat11 + ")")
    }
    if this.Mat12 != nil {
     condition.WriteString("and (y_analyser.mat12 = " + *this.Mat12 + ")")
    }
    if this.Mat13 != nil {
     condition.WriteString("and (y_analyser.mat13 = " + *this.Mat13 + ")")
    }
    if this.Mat14 != nil {
     condition.WriteString("and (y_analyser.mat14 = " + *this.Mat14 + ")")
    }
    if this.Mat15 != nil {
     condition.WriteString("and (y_analyser.mat15 = " + *this.Mat15 + ")")
    }
    if this.Mat16 != nil {
     condition.WriteString("and (y_analyser.mat16 = " + *this.Mat16 + ")")
    }
    if this.Mat17 != nil {
     condition.WriteString("and (y_analyser.mat17 = " + *this.Mat17 + ")")
    }
    if this.Mat18 != nil {
     condition.WriteString("and (y_analyser.mat18 = " + *this.Mat18 + ")")
    }
    if this.Mat19 != nil {
     condition.WriteString("and (y_analyser.mat19 = " + *this.Mat19 + ")")
    }
    if this.Mat2 != nil {
     condition.WriteString("and (y_analyser.mat2 = " + *this.Mat2 + ")")
    }
    if this.Mat20 != nil {
     condition.WriteString("and (y_analyser.mat20 = " + *this.Mat20 + ")")
    }
    if this.Mat3 != nil {
     condition.WriteString("and (y_analyser.mat3 = " + *this.Mat3 + ")")
    }
    if this.Mat4 != nil {
     condition.WriteString("and (y_analyser.mat4 = " + *this.Mat4 + ")")
    }
    if this.Mat5 != nil {
     condition.WriteString("and (y_analyser.mat5 = " + *this.Mat5 + ")")
    }
    if this.Mat6 != nil {
     condition.WriteString("and (y_analyser.mat6 = " + *this.Mat6 + ")")
    }
    if this.Mat7 != nil {
     condition.WriteString("and (y_analyser.mat7 = " + *this.Mat7 + ")")
    }
    if this.Mat8 != nil {
     condition.WriteString("and (y_analyser.mat8 = " + *this.Mat8 + ")")
    }
    if this.Mat9 != nil {
     condition.WriteString("and (y_analyser.mat9 = " + *this.Mat9 + ")")
    }
    if this.Speed != nil {
     condition.WriteString("and (y_analyser.speed = " + *this.Speed + ")")
    }
    if this.Status != nil {
     condition.WriteString("and (y_analyser.status = " + *this.Status + ")")
    }
    if this.TestAt != nil {
     condition.WriteString("and (y_analyser.test_at = " + *this.TestAt + ")")
    }
    if this.Token != nil {
      condition.WriteString("and (y_analyser.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"analyserNum"`
	AnalyserNum *string 
	//`json:"createdAt"`
	CreatedAt *time.Time 
	//`json:"crushingPlant"`
	CrushingPlant *string 
	//`json:"flux"`
	Flux *float64 
	//`json:"load"`
	Load *float64 
	//`json:"mat1"`
	Mat1 *float64 
	//`json:"mat10"`
	Mat10 *float64 
	//`json:"mat11"`
	Mat11 *float64 
	//`json:"mat12"`
	Mat12 *float64 
	//`json:"mat13"`
	Mat13 *float64 
	//`json:"mat14"`
	Mat14 *float64 
	//`json:"mat15"`
	Mat15 *float64 
	//`json:"mat16"`
	Mat16 *float64 
	//`json:"mat17"`
	Mat17 *float64 
	//`json:"mat18"`
	Mat18 *float64 
	//`json:"mat19"`
	Mat19 *float64 
	//`json:"mat2"`
	Mat2 *float64 
	//`json:"mat20"`
	Mat20 *float64 
	//`json:"mat3"`
	Mat3 *float64 
	//`json:"mat4"`
	Mat4 *float64 
	//`json:"mat5"`
	Mat5 *float64 
	//`json:"mat6"`
	Mat6 *float64 
	//`json:"mat7"`
	Mat7 *float64 
	//`json:"mat8"`
	Mat8 *float64 
	//`json:"mat9"`
	Mat9 *float64 
	//`json:"speed"`
	Speed *float64 
	//`json:"status"`
	Status *int 
	//`json:"testAt"`
	TestAt *int 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2YAnalyser() yAnalyserModel.YAnalyser  {
	var yAnalyser = yAnalyserModel.YAnalyser{}

	yAnalyser.AnalyserNum = this.AnalyserNum
	yAnalyser.CreatedAt = this.CreatedAt
	yAnalyser.CrushingPlant = this.CrushingPlant
	yAnalyser.Flux = this.Flux
	yAnalyser.Load = this.Load
	yAnalyser.Mat1 = this.Mat1
	yAnalyser.Mat10 = this.Mat10
	yAnalyser.Mat11 = this.Mat11
	yAnalyser.Mat12 = this.Mat12
	yAnalyser.Mat13 = this.Mat13
	yAnalyser.Mat14 = this.Mat14
	yAnalyser.Mat15 = this.Mat15
	yAnalyser.Mat16 = this.Mat16
	yAnalyser.Mat17 = this.Mat17
	yAnalyser.Mat18 = this.Mat18
	yAnalyser.Mat19 = this.Mat19
	yAnalyser.Mat2 = this.Mat2
	yAnalyser.Mat20 = this.Mat20
	yAnalyser.Mat3 = this.Mat3
	yAnalyser.Mat4 = this.Mat4
	yAnalyser.Mat5 = this.Mat5
	yAnalyser.Mat6 = this.Mat6
	yAnalyser.Mat7 = this.Mat7
	yAnalyser.Mat8 = this.Mat8
	yAnalyser.Mat9 = this.Mat9
	yAnalyser.Speed = this.Speed
	yAnalyser.Status = this.Status
	yAnalyser.TestAt = this.TestAt
	yAnalyser.Token = this.Token
	return yAnalyser
}

type UpdateObj struct {

	//`json:"analyserNum"`
	AnalyserNum *string 
	//`json:"createdAt"`
	CreatedAt *time.Time 
	//`json:"crushingPlant"`
	CrushingPlant *string 
	//`json:"flux"`
	Flux *float64 
	//`json:"load"`
	Load *float64 
	//`json:"mat1"`
	Mat1 *float64 
	//`json:"mat10"`
	Mat10 *float64 
	//`json:"mat11"`
	Mat11 *float64 
	//`json:"mat12"`
	Mat12 *float64 
	//`json:"mat13"`
	Mat13 *float64 
	//`json:"mat14"`
	Mat14 *float64 
	//`json:"mat15"`
	Mat15 *float64 
	//`json:"mat16"`
	Mat16 *float64 
	//`json:"mat17"`
	Mat17 *float64 
	//`json:"mat18"`
	Mat18 *float64 
	//`json:"mat19"`
	Mat19 *float64 
	//`json:"mat2"`
	Mat2 *float64 
	//`json:"mat20"`
	Mat20 *float64 
	//`json:"mat3"`
	Mat3 *float64 
	//`json:"mat4"`
	Mat4 *float64 
	//`json:"mat5"`
	Mat5 *float64 
	//`json:"mat6"`
	Mat6 *float64 
	//`json:"mat7"`
	Mat7 *float64 
	//`json:"mat8"`
	Mat8 *float64 
	//`json:"mat9"`
	Mat9 *float64 
	//`json:"speed"`
	Speed *float64 
	//`json:"status"`
	Status *int 
	//`json:"testAt"`
	TestAt *int 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2YAnalyser() yAnalyserModel.YAnalyser  {
	var yAnalyser = yAnalyserModel.YAnalyser{}

	yAnalyser.AnalyserNum = this.AnalyserNum
	yAnalyser.CreatedAt = this.CreatedAt
	yAnalyser.CrushingPlant = this.CrushingPlant
	yAnalyser.Flux = this.Flux
	yAnalyser.Load = this.Load
	yAnalyser.Mat1 = this.Mat1
	yAnalyser.Mat10 = this.Mat10
	yAnalyser.Mat11 = this.Mat11
	yAnalyser.Mat12 = this.Mat12
	yAnalyser.Mat13 = this.Mat13
	yAnalyser.Mat14 = this.Mat14
	yAnalyser.Mat15 = this.Mat15
	yAnalyser.Mat16 = this.Mat16
	yAnalyser.Mat17 = this.Mat17
	yAnalyser.Mat18 = this.Mat18
	yAnalyser.Mat19 = this.Mat19
	yAnalyser.Mat2 = this.Mat2
	yAnalyser.Mat20 = this.Mat20
	yAnalyser.Mat3 = this.Mat3
	yAnalyser.Mat4 = this.Mat4
	yAnalyser.Mat5 = this.Mat5
	yAnalyser.Mat6 = this.Mat6
	yAnalyser.Mat7 = this.Mat7
	yAnalyser.Mat8 = this.Mat8
	yAnalyser.Mat9 = this.Mat9
	yAnalyser.Speed = this.Speed
	yAnalyser.Status = this.Status
	yAnalyser.TestAt = this.TestAt
	yAnalyser.Token = this.Token
	return yAnalyser
}
