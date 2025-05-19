package auto

import (
	. "sim_data_gen/models/boomGroup"
	. "sim_data_gen/models/boomGroupInfo"
	. "sim_data_gen/models/boomPile"
	. "sim_data_gen/models/composition"
	. "sim_data_gen/models/contentPercent"
	. "sim_data_gen/models/digger"
	. "sim_data_gen/models/diggerProductBinding"
	. "sim_data_gen/models/diggerSwitchBoomGroupLog"
	. "sim_data_gen/models/lorryDiggerBindingLog"
	. "sim_data_gen/models/lorryNearbyTargetSpan"
	. "sim_data_gen/models/mineChangeStatus"
	. "sim_data_gen/models/mineProduct"
	. "sim_data_gen/models/productAndBoomGroup"
	. "sim_data_gen/models/setOfBoomGroups"
	. "sim_data_gen/models/weighLogger"
	. "sim_data_gen/models/yAnalyser"
)

var AutoMigrateDbList = []interface{}{
	&BoomGroup{},
	&BoomGroupInfo{},
	&BoomPile{},
	&Composition{},
	&ContentPercent{},
	&Digger{},

	&DiggerProductBinding{},
	&DiggerSwitchBoomGroupLog{},
	&LorryDiggerBindingLog{},
	&LorryNearbyTargetSpan{},
	&MineChangeStatus{},
	&MineProduct{},
	&ProductAndBoomGroup{},
	&SetOfBoomGroups{},
	&WeighLogger{},
	&YAnalyser{},
}
