package boomGroupInfo

import (
	. "sim_data_gen/utils"
)

func GetMaterialOfBoomGroupInfo(materialIndex int, boomGroupInfo BoomGroupInfo) float64 {
	switch materialIndex {
	case 1:
		return FloatPtrValue(boomGroupInfo.Material1, -1)
	case 2:
		return FloatPtrValue(boomGroupInfo.Material2, -1)
	case 3:
		return FloatPtrValue(boomGroupInfo.Material3, -1)
	case 4:
		return FloatPtrValue(boomGroupInfo.Material4, -1)
	case 5:
		return FloatPtrValue(boomGroupInfo.Material5, -1)
	case 6:
		return FloatPtrValue(boomGroupInfo.Material6, -1)
	case 7:
		return FloatPtrValue(boomGroupInfo.Material7, -1)
	case 8:
		return FloatPtrValue(boomGroupInfo.Material8, -1)
	case 9:
		return FloatPtrValue(boomGroupInfo.Material9, -1)
	case 10:
		return FloatPtrValue(boomGroupInfo.Material10, -1)
	case 11:
		return FloatPtrValue(boomGroupInfo.Material11, -1)
	case 12:
		return FloatPtrValue(boomGroupInfo.Material12, -1)
	case 13:
		return FloatPtrValue(boomGroupInfo.Material13, -1)
	case 14:
		return FloatPtrValue(boomGroupInfo.Material14, -1)
	case 15:
		return FloatPtrValue(boomGroupInfo.Material15, -1)
	case 16:
		return FloatPtrValue(boomGroupInfo.Material16, -1)
	case 17:
		return FloatPtrValue(boomGroupInfo.Material17, -1)
	case 18:
		return FloatPtrValue(boomGroupInfo.Material18, -1)
	case 19:
		return FloatPtrValue(boomGroupInfo.Material19, -1)
	case 20:
		return FloatPtrValue(boomGroupInfo.Material20, -1)
	}
	return -1
}
