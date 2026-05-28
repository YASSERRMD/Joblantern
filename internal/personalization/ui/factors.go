// Package ui surfaces the personalized factors in the verdict page.
package ui

// FactorRow is one factor as rendered to the user.
type FactorRow struct {
	Label  string
	Value  string
	Impact string // "raises risk", "lowers risk", "neutral"
}

// Render produces a stable, presentation-ready slice of rows from a
// minimal personalization payload. Empty values are filtered.
func Render(roleFit int, salaryBump int, locationAnomaly bool, yearsGap int) []FactorRow {
	var rows []FactorRow
	if roleFit > 0 {
		impact := "neutral"
		if roleFit >= 70 {
			impact = "lowers risk"
		} else if roleFit <= 30 {
			impact = "raises risk"
		}
		rows = append(rows, FactorRow{"Role fit", itoa(roleFit) + "/100", impact})
	}
	if salaryBump != 0 {
		impact := "raises risk"
		if salaryBump < 0 {
			impact = "lowers risk"
		}
		rows = append(rows, FactorRow{"Salary vs your level", signed(salaryBump), impact})
	}
	if locationAnomaly {
		rows = append(rows, FactorRow{"Destination uncommon for your background", "yes", "raises risk"})
	}
	if yearsGap != 0 {
		rows = append(rows, FactorRow{"Years-of-experience gap", signed(yearsGap), "raises risk"})
	}
	return rows
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func signed(n int) string {
	if n > 0 {
		return "+" + itoa(n)
	}
	return itoa(n)
}
