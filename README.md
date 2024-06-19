# go-expense-tracker

Uses sqlite3 db to store expenses

## expense

```
$ expense

Track your finances with ease!
expense helps you add, list, and analyze your spending.
Run expense help for detailed commands.

Usage:
  expense [flags]
  expense [command]

Available Commands:
  add         Add new expenses to your tracker
  completion  Generate the autocompletion script for the specified shell
  delete      Delete an expenses from your tracker
  help        Help about any command
  set         Sets a value
  show        View your tracked expenses
  update      Update existing entries

Flags:
  -h, --help   help for expense

Use "expense [command] --help" for more information about a command.

```

## expense add

```
$ expense add -h

Add new expenses to your tracker!
Use expense add ITEM -amount AMOUNT -date DATE -category CATEGORY.
Replace item, amount, dates & categories with your details.
date and category is optional, defaults to current date and general category

Examples:

expense add lunch at mtr -a 500 -c food
expense add lunch at mtr -a 500 -d 12 -c food
expense add lunch at mtr --amount 500 --category food
expense add lunch at mtr --amount 500 --date 12 --category food
expense add lunch at mtr --amount 500 --date 12/06/2024 --category food

all the above commands will produce the same results

Usage:
  expense add ITEM -amount AMOUNT [-date DATE] [-category CATEGORY] [flags]

Flags:
  -a, --amount int        specify the amount of the expense item
  -c, --category string   specify a category of the expense item (default "general")
  -d, --date string       specify a date of the expense item (default "19/06/2024")
  -h, --help              help for add


```

## expense update

```
$ expense update -h

Update existing entries!
Use expense update ID [ITEM] [-amount AMOUNT] [-date DATE] [-category CATEGORY].
Replace ID with the entry's number, brackets indicate optional arguments.

Examples: assuming current date is 12/06/2024

expense update 123 -a 500 -c food
expense update 123 -a 500 -d 12 -c food
expense update 123 --amount 500 --category food
expense update 123 --amount 500 --date 12 -category food
expense update 123 --amount 500 --date 12/06/2024 -category food

all the above commands will produce the same results

Usage:
  expense update ID [ITEM] [-amount AMOUNT] [-date DATE] [-category CATEGORY] [flags]

Flags:
  -p, --amount int        specify the amount of the expense item
  -c, --category string   specify a category of the expense item
  -d, --date string       specify a date of the expense item
  -h, --help              help for update


```

## expense delete

```
$ expense delete -h

Permanently remove expenses!
Use expense delete ID.
Replace ID with the unique number of the expense you wish to remove.
Caution: This action cannot be undone.

expense delete 1 2 3 4

This command deletes expense id numbers 1, 2, 3 and 4

Usage:
  expense delete ID [flags]

Flags:
  -h, --help   help for delete

```

## expense show

```
$ expense show -h

View your tracked expenses!  Use expense show with optional filters:

--month MM or --month MM/YYYY for a specific month (e.g., -m 06/2024).
--group GROUP to categorize expenses by a group, can be any of date or category.

Examples:

expense show
expense show -m 06
expense show -m 06/2024
expense show -m 06/2024 -g date
expense show -m 06/2024 -g category
expense show --month 06
expense show --month 06/2024
expense show --month 06/2024 --group date
expense show --month 06/2024 --group category

Usage:
  expense show [--month MM/YYYY] [--group GROUP] [flags]
  expense show [command]

Available Commands:
  report      Gets the yearly report for a given year, defaults to current year
  today       View your todays tracked expenses
  week        View your current week's tracked expenses
  yesterday   View your yesterday's tracked expenses

Flags:
  -g, --group string   specify a group condition can be any one of [date, category]
  -h, --help           help for show
  -m, --month string   specify the month for list of expenses (default "06/2024")

Use "expense show [command] --help" for more information about a command.


```

## expense show today

```
$ expense show today -h

View your todays tracked expenses!  Use expense show today with optional filters:

--group GROUP to categorize expenses by a group, can be any of date or category.

Examples:

expense show today
expense show today -g date
expense show today -g category
expense show today --group date
expense show today --group category

Usage:
  expense show today [--group GROUP] [flags]

Flags:
  -g, --group string   specify a group condition can be any one of [date, category]
  -h, --help           help for today


```

## expense show yesterday

```
$ expense show yesterday -h

View your yesterday's tracked expenses!  Use expense show today with optional filters:

--group GROUP to categorize expenses by a group, can be any of date or category.

Examples:

expense show yesterday
expense show yesterday -g date
expense show yesterday -g category
expense show yesterday --group date
expense show yesterday --group category

Usage:
  expense show yesterday [--group GROUP] [flags]

Flags:
  -g, --group string   specify a group condition can be any one of [date, category]
  -h, --help           help for yesterday


```

## expense show report

```
$ expense show report -h

View your yearly expense report!
Use expense show report with optional filters:

Examples:

expense show report
expense show report -y 2024
expense show report --year 2024

Usage:
  expense show report [-y YYYY] [flags]

Flags:
  -h, --help          help for report
  -y, --year string   specify the year to generate report for (default "2024")


```

## expense set budget

```
$ expense set budget -h
Sets the budget for a given month(MM) or month/year(MM/YYYY) defaults to current month

Usage:
  expense set budget [-m MM/YYYY or MM] [flags]

Flags:
  -h, --help           help for budget
  -m, --month string   specify the month to set budget for (default "06/2024")

```
