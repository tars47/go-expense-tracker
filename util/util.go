package util

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/tars47/go-expense-tracker/db"
	"github.com/tars47/go-expense-tracker/printf"
)

func ValidateDateFlag(ds string) (time.Time, error) {

	if len(strings.Split(ds, "/")) == 1 {
		ss := strings.Split(time.Now().Format("02/01/2006"), "/")
		ds = ds + "/" + strings.Join(ss[1:], "/")
	}
	d, err := time.Parse("02/01/2006", ds)
	if err != nil {
		printf.Red("Error: DATE value is invalid, it should be of format DD or DD/MM/YY")
		return time.Time{}, err
	}
	return d, nil
}

func ValidateMonthFlag(ms string) (time.Time, error) {
	var ds string
	t := time.Now()
	if len(strings.Split(ms, "/")) == 1 {
		ds = "01/" + ms + "/" + fmt.Sprint(t.Year())
	} else {
		ds = "01/" + ms
	}
	d, err := time.Parse("02/01/2006", ds)
	if err != nil {
		printf.Red("Error: MONTH value is invalid, it should be of format MM or MM/YY")
		return time.Time{}, err
	}

	return d, nil
}

func PrintTable(g string, es []db.Expense) int {
	switch g {
	case "date":
		{
			h := []string{"Date", "Amount"}
			var r [][]string
			var t int
			for _, e := range es {
				r = append(r, []string{
					e.Date.Format("02/01/2006"),
					FmtNum(e.Amount),
				})
				t += e.Amount
			}
			fmt.Println(Table(h, r))
			return t
		}
	case "category":
		{
			h := []string{"Category", "Amount"}
			var r [][]string
			var t int
			for _, e := range es {
				r = append(r, []string{
					e.Category,
					FmtNum(e.Amount),
				})
				t += e.Amount
			}
			fmt.Println(Table(h, r))
			return t
		}
	default:
		{
			h := []string{"ID", "Date", "Item", "Amount", "Category"}
			var r [][]string
			var t int
			for _, e := range es {
				r = append(r, []string{
					fmt.Sprintf("%d", e.Id),
					e.Date.Format("02/01/2006"),
					e.Item,
					FmtNum(e.Amount),
					e.Category,
				})
				t += e.Amount
			}
			fmt.Println(Table(h, r))
			return t
		}
	}
}

func Table(h []string, r [][]string) *table.Table {

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers(h...).
		Rows(r...).
		Width(72).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("#04B575")).
					Border(lipgloss.NormalBorder()).
					BorderTop(false).
					BorderLeft(false).
					BorderRight(false).
					BorderBottom(true).
					Bold(true)
			}
			if row%2 == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
			}
			return lipgloss.NewStyle()
		})
	return t
}

func PrintSummary(s int, d time.Time, e *db.ExpenseDb) {

	printf.BlueS(" Spent\t\t", FmtNum(s))

	b, err := e.GetBudget(int(d.Month()), d.Year())
	if err != nil {
		printf.Red("Error: There was an error querying budget from the database, %v", err)
		return
	}

	printf.BlueS(" Budget\t\t", FmtNum(b.Budget))

	bl := b.Budget - s
	if bl > 0 {
		printf.BlueS(" Balance\t", FmtNum(bl))
		printf.BlueS(" Summary\t", fmt.Sprintf("spent %.1f%% of the budget\n", (float64(s)/float64(b.Budget))*100))
	} else {
		printf.BlueS(" Balance\t", 0)
		printf.BlueS(" Summary\t", fmt.Sprintf("exceeded the budget by %.1f%%\n", math.Abs((float64(s)/float64(b.Budget))*100)))
	}
}

func FmtNum(n int) string {
	ret := make([]string, 0, 10)
	s := strings.Split(fmt.Sprint(n), "")
	slices.Reverse(s)
	for i, v := range s {
		if i > 2 && i%2 == 1 {
			ret = append(ret, ",")
		}
		ret = append(ret, v)
	}
	slices.Reverse(ret)
	return strings.Join(ret, "")
}
