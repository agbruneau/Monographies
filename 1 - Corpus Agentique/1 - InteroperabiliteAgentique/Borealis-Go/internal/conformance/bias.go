package conformance

import (
	"math"
	"sort"
)

// GroupDecision associe une décision (approuvée ou non) à un groupe protégé
// (segment, tranche d'âge…). Données synthétiques ; build-time.
type GroupDecision struct {
	Group    string
	Approved bool
}

// GroupStat est la statistique d'approbation d'un groupe.
type GroupStat struct {
	Group        string
	N            int
	ApprovalRate float64
}

// BiasReport agrège les taux d'approbation par groupe, la divergence maximale
// (impact disparate = max − min des taux) et le coefficient de Gini des taux.
type BiasReport struct {
	Groups        []GroupStat
	MaxDivergence float64
	Gini          float64
	Threshold     float64
	Alert         bool
}

// AnalyzeBias calcule les statistiques par groupe et **alerte** si la divergence
// maximale des taux d'approbation dépasse threshold (ex. 0,15 = 15 %) — §11.4.
func AnalyzeBias(decisions []GroupDecision, threshold float64) BiasReport {
	type acc struct{ n, approved int }
	byGroup := map[string]*acc{}
	for _, d := range decisions {
		a := byGroup[d.Group]
		if a == nil {
			a = &acc{}
			byGroup[d.Group] = a
		}
		a.n++
		if d.Approved {
			a.approved++
		}
	}

	groups := make([]GroupStat, 0, len(byGroup))
	for g, a := range byGroup {
		rate := 0.0
		if a.n > 0 {
			rate = float64(a.approved) / float64(a.n)
		}
		groups = append(groups, GroupStat{Group: g, N: a.n, ApprovalRate: rate})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Group < groups[j].Group })

	report := BiasReport{Groups: groups, Threshold: threshold}
	if len(groups) == 0 {
		return report
	}
	minR, maxR := groups[0].ApprovalRate, groups[0].ApprovalRate
	for _, g := range groups {
		minR = math.Min(minR, g.ApprovalRate)
		maxR = math.Max(maxR, g.ApprovalRate)
	}
	report.MaxDivergence = maxR - minR
	report.Gini = giniOf(groups)
	report.Alert = report.MaxDivergence > threshold
	return report
}

// giniOf calcule le coefficient de Gini des taux d'approbation (inégalité).
func giniOf(groups []GroupStat) float64 {
	n := len(groups)
	if n == 0 {
		return 0
	}
	sum := 0.0
	for _, g := range groups {
		sum += g.ApprovalRate
	}
	if sum == 0 {
		return 0
	}
	var absDiff float64
	for _, a := range groups {
		for _, b := range groups {
			absDiff += math.Abs(a.ApprovalRate - b.ApprovalRate)
		}
	}
	return absDiff / (2 * float64(n) * sum)
}
