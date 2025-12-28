-- reset all existing categorization

--drop table if exists backupTransactions;
--create table backupTransactions as select * from Transactions;

UPDATE
	transactions
SET
	beneficiary = NULL
	, category = NULL
	, subcategory = NULL
	, tag = NULL;