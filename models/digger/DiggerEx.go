package digger


//计算挖机的产量
func (this *Digger)CalcYield(nowAsUtc int64) float64 {
	if(*this.Utc > 0) {
		spendSeconds := nowAsUtc - (int64)(*this.Utc)
		var hour float64 = float64(spendSeconds) / 3600.0
		totalMine := (*this.Speed) * hour
		return totalMine
	}
	return 0
}
