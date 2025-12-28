-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
-- Count unattributed transactions (beneficiary is null)
-- run before and after applying patterns, raw hints and manual fixups.
SELECT 
	month, count(*) as un_attributed 
FROM stemTransactions
where postedDate >= "2024-08" and beneficiary is null
group by month
order by month;
