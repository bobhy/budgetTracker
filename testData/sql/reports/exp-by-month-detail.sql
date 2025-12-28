-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
select
	month
	, beneficiary
	, category
	, subCategory
	, account
	, count(*) num
	, sum(amount) amount
from
	stemTransactions st 
group by 
	month
	, beneficiary
	, category
	, subCategory
	, account
order by
	month
	, beneficiary
	, category
	, subCategory
	, account
