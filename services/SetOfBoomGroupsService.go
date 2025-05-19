package services

import (
	"bytes"
	setOfBoomGroupsDao "sim_data_gen/daos/setOfBoomGroups"
	setOfBoomGroupsModel "sim_data_gen/models/setOfBoomGroups"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

func ChangeSetOfBoomGroups(request map[string]interface{}) (err error) {
	err = GetDb().Transaction(func(tx *gorm.DB) (err error) {
		var (
			condition           bytes.Buffer
			setOfBoomGroupsList []setOfBoomGroupsModel.SetOfBoomGroupsFully
		)
		if request["DateFlag"] != nil {
			condition.WriteString(" where (set_of_boom_groups.date_flag like '%" + request["DateFlag"].(string) + "%')")
		}
		if request["Diggers"] != nil {
			condition.WriteString("and (set_of_boom_groups.diggers like '%" + request["Diggers"].(string) + "%')")
		}
		if request["MineProductDesp"] != nil {
			condition.WriteString("and (mine_product0.name = '" + request["MineProductDesp"].(string) + "')")
		}
		if setOfBoomGroupsList, err = setOfBoomGroupsDao.PageObj(condition.String(), "", 1, 1000, tx); err != nil {
			return err
		}
		for _, fully := range setOfBoomGroupsList {
			setOfBoomGroups := fully.SetOfBoomGroups
			if *setOfBoomGroups.Token == request["Token"].(string) {
				setOfBoomGroups.Status = RefInt(0)
			} else {
				setOfBoomGroups.Status = RefInt(1)
			}
			if _, err = setOfBoomGroupsDao.SaveObj(&setOfBoomGroups, tx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
