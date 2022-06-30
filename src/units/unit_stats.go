package units

import (
	"configparser"
	"sort"
	"utils"
)

type UnitStat struct {
	RequestsCount        int
	AvgRequestsPerSecond float64
	MostRequestedPaths   []string
}

type UnitStats struct {
	Units           []Unit // List of units
	Config          configparser.Config
	MostAnnoyingIPs []string
	Stats           map[string]UnitStat // Statistics for every unit
}

func getUniqueIPs(u []Unit) []string {
	return utils.Map(
		utils.FilterUnique(
			u,
			func(u1 Unit, u2 Unit) bool {
				return u1.Ip == u2.Ip
			}),
		func(unit Unit) string {
			return unit.Ip
		})
}

func getUniquePaths(u []Unit) []string {
	return utils.Map(
		utils.FilterUnique(
			u,
			func(u1 Unit, u2 Unit) bool {
				return u1.Path == u2.Path
			}),
		func(unit Unit) string {
			return unit.Path
		})
}

func getRequestsCount(u []Unit) int {
	return len(u)
}

func setRequestsCount(u []Unit, stat *UnitStat) {
	stat.RequestsCount = getRequestsCount(u)
}

func setRequestsPerSecond(u []Unit, stat *UnitStat) {
	requestsCount := getRequestsCount(u)
	timeStarted := u[0].Time
	timeStopped := u[requestsCount-1].Time
	stat.AvgRequestsPerSecond = float64(requestsCount) / (timeStopped - timeStarted)
}

func setMostRequestedPath(u []Unit, stat *UnitStat) {
	uniquePaths := getUniquePaths(u)

	pathsToRequestsCount := make(map[string]int, len(uniquePaths))

	for _, path := range uniquePaths {
		pathsToRequestsCount[path] = len(utils.FilterArray(u, func(u Unit) bool { return u.Path == path }))
	}

	sort.SliceStable(uniquePaths, func(i, j int) bool {
		return pathsToRequestsCount[uniquePaths[i]] > pathsToRequestsCount[uniquePaths[j]]
	})
	stat.MostRequestedPaths = uniquePaths
}

func calculateStat(units []Unit, config configparser.Config) UnitStat {
	var stat UnitStat
	configKeyToFunction := map[string]func([]Unit, *UnitStat){
		"requests_count":      setRequestsCount,
		"requests_per_second": setRequestsPerSecond,
		"most_requested_path": setMostRequestedPath,
	}
	for key := range configKeyToFunction {
		if utils.IsElementInArray(config.Output.StatFields, key) {
			configKeyToFunction[key](units, &stat)
		}
	}
	return stat
}

func getMostAnnoyingIPs(units []Unit) []string {
	uniqueIPs := getUniqueIPs(units)
	ipsToRequestsCount := make(map[string]int, len(uniqueIPs))

	for _, unit := range units {
		_, ok := ipsToRequestsCount[unit.Ip]
		if ok {
			ipsToRequestsCount[unit.Ip] += 1
		} else {
			ipsToRequestsCount[unit.Ip] = 1
		}
	}

	sort.SliceStable(uniqueIPs, func(i, j int) bool {
		return ipsToRequestsCount[uniqueIPs[i]] > ipsToRequestsCount[uniqueIPs[j]]
	})
	return uniqueIPs
}

func (u *UnitStats) calculateStats(config configparser.Config) {
	uniqueIPs := getUniqueIPs(u.Units)

	u.Stats = make(map[string]UnitStat, len(uniqueIPs))
	for _, ip := range uniqueIPs {
		u.Stats[ip] = calculateStat(utils.FilterArray(
			u.Units,
			func(unit Unit) bool { return unit.Ip == ip }),
			config)
	}
}

func (u *UnitStats) CalculateStats(config configparser.Config) {
	u.calculateStats(config)

	u.MostAnnoyingIPs = getMostAnnoyingIPs(u.Units)
}

func (u *UnitStats) GetRequestsPerSecond() float64 {
	requestsCount := getRequestsCount(u.Units)
	firstRequestTime := u.Units[0].Time
	lastRequestTime := u.Units[len(u.Units)-1].Time
	return float64(requestsCount) / (lastRequestTime - firstRequestTime)
}
