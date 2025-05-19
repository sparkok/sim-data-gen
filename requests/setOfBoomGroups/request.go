package setOfBoomGroups

import (
	"bytes"
	setOfBoomGroupsModel "sim_data_gen/models/setOfBoomGroups"
	"sim_data_gen/requests/common"
)
import "time"

/**
* SetOfBoomGroups请求类
 */
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {

	//`json:"boomGroupIds"`
	BoomGroupIds *string
	//`json:"createdAt"`
	CreatedAt *string
	//`json:"dateFlag"`
	DateFlag *string
	//`json:"diggers"`
	Diggers *string
	//`json:"matContents"`
	MatContents *string
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string
	//`json:"mineProductDesp"`
	MineProductDesp *string
	//`json:"name"`
	Name *string
	//`json:"nt"`
	Nt *string

	//`json:"status"`
	Status *string
	//`json:"token"`
	Token *string
	//`json:"updateAt"`
	UpdateAt *string
}

func (this *SearchInfo) GetConditions() string {
	var condition bytes.Buffer
	if this.BoomGroupIds != nil {
		condition.WriteString("and (set_of_boom_groups.boom_group_ids like '%" + *this.BoomGroupIds + "%')")
	}
	if this.CreatedAt != nil {
		condition.WriteString("and (set_of_boom_groups.created_at = '" + *this.CreatedAt + "')")
	}
	if this.DateFlag != nil {
		condition.WriteString("and (set_of_boom_groups.date_flag like '%" + *this.DateFlag + "%')")
	}
	if this.Diggers != nil {
		condition.WriteString("and (set_of_boom_groups.diggers like '%" + *this.Diggers + "%')")
	}
	if this.MatContents != nil {
		condition.WriteString("and (set_of_boom_groups.mat_contents like '%" + *this.MatContents + "%')")
	}
	if this.TokenOfMineProduct != nil {
		condition.WriteString("and (set_of_boom_groups.mine_product_id = '" + *this.TokenOfMineProduct + "')")
	}
	if this.MineProductDesp != nil {
		condition.WriteString("and (mine_product0.name = '" + *this.MineProductDesp + "')")
	}
	if this.Name != nil {
		condition.WriteString("and (set_of_boom_groups.name like '%" + *this.Name + "%')")
	}
	if this.Nt != nil {
		condition.WriteString("and (set_of_boom_groups.nt like '%" + *this.Nt + "%')")
	}
	if this.Status != nil {
		condition.WriteString("and (set_of_boom_groups.status = " + *this.Status + ")")
	}
	if this.Token != nil {
		condition.WriteString("and (set_of_boom_groups.token like '%" + *this.Token + "%')")
	}
	if this.UpdateAt != nil {
		condition.WriteString("and (set_of_boom_groups.update_at = '" + *this.UpdateAt + "')")
	}

	if condition.Len() > 4 {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"boomGroupIds"`
	BoomGroupIds *string
	//`json:"createdAt"`
	CreatedAt *time.Time
	//`json:"dateFlag"`
	DateFlag *string
	//`json:"diggers"`
	Diggers *string
	//`json:"matContents"`
	MatContents *string
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string
	//`json:"name"`
	Name *string
	//`json:"nt"`
	Nt *string

	//`json:"status"`
	Status *int
	//`json:"token"`
	Token *string
	//`json:"updateAt"`
	UpdateAt *time.Time
}

func (this *CreateObj) Convert2SetOfBoomGroups() setOfBoomGroupsModel.SetOfBoomGroups {
	var setOfBoomGroups = setOfBoomGroupsModel.SetOfBoomGroups{}

	setOfBoomGroups.BoomGroupIds = this.BoomGroupIds
	setOfBoomGroups.CreatedAt = this.CreatedAt
	setOfBoomGroups.DateFlag = this.DateFlag
	setOfBoomGroups.Diggers = this.Diggers
	setOfBoomGroups.MatContents = this.MatContents
	setOfBoomGroups.TokenOfMineProduct = this.TokenOfMineProduct
	setOfBoomGroups.Name = this.Name
	setOfBoomGroups.Nt = this.Nt

	setOfBoomGroups.Status = this.Status
	setOfBoomGroups.Token = this.Token
	setOfBoomGroups.UpdateAt = this.UpdateAt
	return setOfBoomGroups
}

type UpdateObj struct {

	//`json:"boomGroupIds"`
	BoomGroupIds *string
	//`json:"createdAt"`
	CreatedAt *time.Time
	//`json:"dateFlag"`
	DateFlag *string
	//`json:"diggers"`
	Diggers *string
	//`json:"matContents"`
	MatContents *string
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string
	//`json:"name"`
	Name *string
	//`json:"nt"`
	Nt *string

	//`json:"status"`
	Status *int
	//`json:"token"`
	Token *string
	//`json:"updateAt"`
	UpdateAt *time.Time
}

func (this *UpdateObj) Convert2SetOfBoomGroups() setOfBoomGroupsModel.SetOfBoomGroups {
	var setOfBoomGroups = setOfBoomGroupsModel.SetOfBoomGroups{}

	setOfBoomGroups.BoomGroupIds = this.BoomGroupIds
	setOfBoomGroups.CreatedAt = this.CreatedAt
	setOfBoomGroups.DateFlag = this.DateFlag
	setOfBoomGroups.Diggers = this.Diggers
	setOfBoomGroups.MatContents = this.MatContents
	setOfBoomGroups.TokenOfMineProduct = this.TokenOfMineProduct
	setOfBoomGroups.Name = this.Name
	setOfBoomGroups.Nt = this.Nt

	setOfBoomGroups.Status = this.Status
	setOfBoomGroups.Token = this.Token
	setOfBoomGroups.UpdateAt = this.UpdateAt
	return setOfBoomGroups
}
