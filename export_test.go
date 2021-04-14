package think

var (
	Export_calcWeight    = calcWeight
	Export_modulesSorted = modulesSorted
	Export_makeWeightMap = func() map[Module]*weight {
		return make(map[Module]*weight)
	}
	Export_weightValue = func(weight *weight) int {
		return weight.toInt()
	}
)
