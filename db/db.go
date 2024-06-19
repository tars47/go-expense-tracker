package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
)

type Expense struct {
	Id       int
	Item     string
	Amount   int
	Date     time.Time
	Category string
}

type Budget struct {
	Budget int
	Month  int
	Year   int
}

type Report struct {
	Month      time.Month
	Year       int
	Spent      int
	Budget     int
	Saved      int
	Percentage string
}

type ExpenseDb struct {
	*sql.DB
	FilePath string
}

func OpenDB() (*ExpenseDb, error) {

	etdir, err := initExpenseDbDir()
	if err != nil {
		return nil, err
	}

	etpath := filepath.Join(etdir, "expense.db")

	db, err := sql.Open("sqlite3", etpath)
	if err != nil {
		return nil, err
	}

	etdb := ExpenseDb{db, etpath}
	if err := etdb.createTablesIfNotExists(); err != nil {
		return nil, err
	}
	return &etdb, nil
}

func initExpenseDbDir() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	etdir := filepath.Join(homeDir, ".expense")

	if _, err := os.Stat(etdir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(etdir, 0o770); err != nil {
				return "", err
			}
			return etdir, nil
		}
		return "", err
	}
	return etdir, nil
}

func (e *ExpenseDb) createTablesIfNotExists() error {
	_, err := e.Exec(`
	   
	        CREATE TABLE IF NOT EXISTS expenses 
			       ( 
		              id INTEGER, 
					  item TEXT NOT NULL, 
					  amount INTEGER NOT NULL, 
					  category TEXT NOT NULL, 
					  date DATETIME NOT NULL, 
													
					  PRIMARY KEY(id AUTOINCREMENT)
					)
		   `)
	if err != nil {
		return err
	}

	_, err = e.Exec(`
	  
	        CREATE TABLE IF NOT EXISTS  budget
			       (
			          budget INTEGER NOT NULL, 
					  month INTEGER NOT NULL, 
					  year INTEGER NOT NULL, 
					
					  PRIMARY KEY (month,year)
					)
			`)
	if err != nil {
		return err
	}
	return nil

}

func (e *ExpenseDb) Insert(expense Expense) error {

	_, err := e.Exec(
		`
		    INSERT INTO expenses 
		           (item,amount,category,date) 
				   VALUES( ?, ?, ?, date(?))
		`,
		expense.Item,
		expense.Amount,
		expense.Category,
		expense.Date.Format("2006-01-02"))

	return err
}

func (e *ExpenseDb) Delete(ids []int) (int, error) {

	q := "DELETE FROM expenses WHERE id = ? "

	idany := make([]any, 0, len(ids))
	idany = append(idany, ids[0])

	for i := 1; i < len(ids); i++ {
		q += "OR id = ? "
		idany = append(idany, ids[i])
	}

	res, err := e.Exec(q, idany...)
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()

	return int(rows), err
}

func (e *ExpenseDb) ListRange(first, last, group string) ([]Expense, error) {

	var q string
	var es []Expense

	switch {
	case group == "date":
		{

			q = `
			          SELECT date,sum(amount) AS amount 
					  FROM expenses 
					  WHERE date BETWEEN date(?) AND date(?) 
					  GROUP BY date
					  ORDER BY date DESC, amount DESC
				`
			rows, err := e.Query(q, first, last)

			if err != nil {
				return nil, err
			}
			for rows.Next() {
				var ex Expense
				err = rows.Scan(
					&ex.Date,
					&ex.Amount,
				)
				if err != nil {
					return es, err
				}
				es = append(es, ex)
			}
			return es, err
		}

	case group == "category":
		{
			q = `
			         SELECT category,sum(amount) AS amount  
					 FROM expenses 
					 WHERE date BETWEEN date(?) AND date(?) 
					 GROUP BY category
					 ORDER BY date DESC, amount DESC
				`
			rows, err := e.Query(q, first, last)

			if err != nil {
				return nil, err
			}
			for rows.Next() {
				var ex Expense
				err = rows.Scan(
					&ex.Category,
					&ex.Amount,
				)
				if err != nil {
					return es, err
				}
				es = append(es, ex)
			}
			return es, err
		}

	default:
		{
			q = `
			         SELECT * FROM expenses 
					 WHERE date BETWEEN date(?) AND date(?) 
					 ORDER BY date DESC, amount DESC
				`
			rows, err := e.Query(q, first, last)

			if err != nil {
				return nil, err
			}

			for rows.Next() {
				var ex Expense
				err = rows.Scan(
					&ex.Id,
					&ex.Item,
					&ex.Amount,
					&ex.Category,
					&ex.Date,
				)
				if err != nil {
					return es, err
				}
				es = append(es, ex)
			}
			return es, err
		}

	}
}

func (e *ExpenseDb) Update(exp Expense) error {

	q := `UPDATE expenses SET `
	qParts := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)

	if exp.Item != "" {
		qParts = append(qParts, `item = ?`)
		args = append(args, exp.Item)
	}

	if exp.Amount > 0 {
		qParts = append(qParts, `amount = ?`)
		args = append(args, exp.Amount)
	}

	if !exp.Date.IsZero() {
		qParts = append(qParts, `date = date(?)`)
		args = append(args, exp.Date.Format("2006-01-02"))
	}

	if exp.Category != "" {
		qParts = append(qParts, `category = ?`)
		args = append(args, exp.Category)
	}

	q += strings.Join(qParts, ",") + ` WHERE id = ?`
	args = append(args, exp.Id)

	fmt.Println(q)
	_, err := e.Exec(q, args...)

	return err
}

func (e *ExpenseDb) UpsertBudget(b Budget) error {

	_, err := e.Exec(
		`
		   INSERT INTO budget 
		          (budget, month, year) 
				  VALUES( ?, ?, ?)
		   ON CONFLICT(month, year)
		   DO UPDATE SET budget = ?;
		`,
		b.Budget,
		b.Month,
		b.Year,
		b.Budget)

	return err
}

func (e *ExpenseDb) GetBudget(m, y int) (Budget, error) {

	var b Budget
	row := e.QueryRow(` 
	                     SELECT budget, month, year FROM budget 
		                 WHERE month = ? AND year = ?
		              `,

		m,
		y,
	)

	err := row.Scan(&b.Budget, &b.Month, &b.Year)

	return b, err
}

func (e *ExpenseDb) GetReport(first, last string) ([]Report, error) {

	var rs []Report
	rows, err := e.Query(` 
	                     SELECT E.year,
						        E.month,
								E.spent,
								B.budget-E.spent AS saved,
								B.budget
						 FROM 
						     (
                                
						        SELECT  strftime('%m', date) AS month,
								        strftime('%Y', date) AS year, 
										sum(amount) AS spent  
								FROM expenses E
                                WHERE date >= date(?) 
								      AND 
									  date <= date(?)
                                GROUP BY month
                                ORDER BY month ASC
                             ) E
                         JOIN budget B
                         ON  (
						        B.month = E.month 
							    AND 
							    B.year = E.year
							 )
`,

		first,
		last,
	)

	if err != nil {
		return rs, err
	}
	for rows.Next() {
		var r Report
		err = rows.Scan(
			&r.Year,
			&r.Month,
			&r.Spent,
			&r.Saved,
			&r.Budget,
		)
		if err != nil {
			return rs, err
		}
		rs = append(rs, r)
	}
	return rs, err
}
