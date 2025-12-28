-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
-- import raw transactions (from capitalone.)
drop table if exists temp.raw;
CREATE virtual TABLE temp.raw USING csv(
	--filename='/home/bobhy/db/budget/import/c1_2024-06-30_2024-04-01_transaction_download.csv',
	--filename='/home/bobhy/db/budget/import/c1.csv',
filename='/home/bobhy/db/budget/import/c1_2024-11-05_transaction_download.csv',
	txDate text,
	postedDate text,
	cardNo text,
	description text,
	category text,
	debit REAL,
	credit real	
);

--select cardno, count(*), sum(debit), sum(credit) from temp.raw group by cardno;
--select * from temp.raw where cardno = '7223';

with all_key as (
	select 
	r.postedDate
	, r.description
	, r.category
	, (debit - credit) as r_amount
	, a.account as account
	, r.category as rawHint
	, t.key as tran_key
	--, t.PostedDate as tran_postedDate
	--, t.Description as tran_description
	--, t.amount as tran_amount
	--, t.account as tran_account
	from temp.raw r
	left outer join Accounts a
	on r.cardNo = a.payment_method
	left outer join Transactions t
	on r.postedDate = t.postedDate
	and a.account = t.account
	and r.description = t.description
	and r_amount = t.amount
)	
--select * from all_key
--limit 100;
insert into Transactions (postedDate, account, amount, description, rawHint)
select a.postedDate, a.account, r_amount, a.description, a.rawHint
from all_key a where tran_key is null;

select changes();
