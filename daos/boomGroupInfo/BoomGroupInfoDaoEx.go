package boomGroupInfo

import (
	boomGroupInfoModel "sim_data_gen/models/boomGroupInfo"
	. "sim_data_gen/utils"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

// 根据ID获取采矿扩展
func GetObjByBoomGroupToken(boomGroupToken *string, tx ...*gorm.DB) (boomGroupInfoModel.BoomGroupInfo, error) {
	boomGroupInfo := boomGroupInfoModel.BoomGroupInfo{TokenOfBoomGroup: boomGroupToken}
	result := boomGroupInfoModel.BoomGroupInfo{}
	db := GetDb(tx...).Where(&boomGroupInfo).Take(&result)
	return result, db.Error
}

// 列出 采矿单元
func ListLastContentsByIds(ids []string, tx ...*gorm.DB) ([]boomGroupInfoModel.BoomGroupInfo, error) {
	builder := strings.Builder{}
	builder.WriteString("SELECT last_boom_group_info.boom_group_id,last_boom_group_info.created_at,boom_group_info.*,boom_group.* FROM (select boom_group_id,max(created_at) as created_at from boom_group_info group by boom_group_id) AS last_boom_group_info\n")
	builder.WriteString("left join boom_group_info ON (last_boom_group_info.boom_group_id = boom_group_info.boom_group_id AND last_boom_group_info.created_at = boom_group_info.created_at)\n")
	builder.WriteString("left join boom_group on (boom_group.token = last_boom_group_info.boom_group_id)\n")
	builder.WriteString(fmt.Sprintf("where last_boom_group_info.boom_group_id in ('%s')", strings.Join(ids, "','")))
	sql := builder.String()
	switch CurrentDbType {
	case "sqlite":
		list := []boomGroupInfoModel.BoomGroupInfoEx{}
		db := GetDb(tx...).Raw(sql).Find(&list)
		return convert2BoomGroupInfo(list), db.Error
	case "postgres":
		list := []boomGroupInfoModel.BoomGroupInfo{}
		db := GetDb(tx...).Raw(sql).Find(&list)
		return list, db.Error
	default:
		return nil, errors.Errorf("unsupported driver")
	}
}

func convert2BoomGroupInfo(list []boomGroupInfoModel.BoomGroupInfoEx) (ret []boomGroupInfoModel.BoomGroupInfo) {
	ret = make([]boomGroupInfoModel.BoomGroupInfo, len(list))
	for idx, item := range list {
		ret[idx] = *item.Convert2Obj()
	}
	return ret
}
