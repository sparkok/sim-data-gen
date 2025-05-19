package location

func (this *Location) IsValid() bool {
	if this.Status == nil {
		return false
	}
	if this.X == nil || this.Y == nil || *this.X == 0 || *this.Y == 0 {
		return false
	}
	if *this.Status&1 != 0 {
		//gps
		return true
	}
	if *this.Status&2 != 0 {
		//cell
		return true
	}
	if *this.Status&4 != 0 {
		//calculate
		return false
	}
	return false

}
