-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
SELECT 
strftime('%Y-%m', postedDate, 'start of month') month
, beneficiary
, count(*) num
, sum(amount) amount
from stemTransactions st 
where category not in ('income', 'transfers')
group by month, beneficiary
order by month, beneficiary