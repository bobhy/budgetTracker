-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
DROP TABLE IF EXISTS temp.xall;


CREATE virtual TABLE temp.xall
	USING csv(
	filename='/home/bobhy/db/budget/export/recat1.csv',
	KEY text,
	MONTH text,
	postedDate text,
	account text,
	amount REAL,
	description text,
	rawHint text,
	beneficiary text,
	category text,
	subCategory text,
	stemdesc text
);

UPDATE
	Transactions
SET 
	beneficiary = y.beneficiary,
	category = y.category,
	subcategory = y.subcategory
FROM
	(
	SELECT
		x.KEY,
		x.beneficiary,
		x.category,
		x.subcategory
	FROM
		temp.xall x 
	) AS y
WHERE
	transactions.key = y.key;

SELECT
	changes();

UPDATE
	Transactions
SET
	beneficiary = NULL
WHERE
	beneficiary = '';

UPDATE
	Transactions
SET
	category = NULL
WHERE
	category = '';

UPDATE
	Transactions
SET
	subcategory = NULL
WHERE
	subcategory = '';
