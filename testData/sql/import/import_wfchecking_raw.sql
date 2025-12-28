-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
-- import raw transactions (from wf checking.)

DROP TABLE IF EXISTS temp.raw;
CREATE virtual TABLE temp.raw USING csv(
	--filename='/home/bobhy/db/budget/import/wfchecking_2024-06-30.csv',
	--filename='/home/bobhy/db/budget/import/Checking1.csv',
filename='/home/bobhy/db/budget/import/wf_20241105_Checking1.csv',
	header=FALSE,
	postedDate text,
	credit REAL,
	cleared text,
	rawHint text,
	description text,	
);

--SELECT * FROM temp.raw LIMIT 10;
--select distinct cardno from temp.c1raw limit 10;

with all_key as (
	select 
	 date(concat_ws('-'
	 	, substr(r.postedDate, 7,4)
	 	, substr(r.postedDate, 1,2)
	 	, substr(r.postedDate, 4,2) )) fixedDate
	, r.description
	, (0 - credit) as amount
	, 'WF Checking' as account	-- hardcoded account
	, r.rawHint
	, t.key
	from temp.raw r
	left outer join Transactions t
	on fixedDate = t.postedDate
	and account = t.account
	and r.description = t.description
	and amount = t.amount
)
--select * from all_key  limit 100;
insert into Transactions (postedDate, account, amount, description, rawHint)
select a.fixedDate, a.account, a.amount, a.description, a.rawHint
from all_key a where a.key is null;

select changes();
