-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
-- all transactions by stem, then categories; used to manually update some orphans
SELECT * FROM stemTransactions st 
ORDER BY beneficiary, category, subcategory, stemdesc, account, postedDate;