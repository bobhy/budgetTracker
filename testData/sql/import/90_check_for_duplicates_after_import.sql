-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
-- Check for potential duplicates after import, where newly imported record duplicates one that has already beeen categorized.
-- run the report and scan for adjacent duplicate amounts.
SELECT 
	* 
FROM stemTransactions
where postedDate >= "2024-08"
order by postedDate, amount, description;
