-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
-- Apply categories based on patterns in (stemmed) description of transaction
-- only apply to transactions that haven't already been categorized (i.e beneficiary is non-null)

update
	Transactions
set
	beneficiary = stptb.beneficiary,
	category = stptb.mainCat,
	subCategory = stptb.subCat,
	tag = stptb.pat
from
	(
		select
			st.key,
			ptb.pat,
			ptb.beneficiary,
			ptb.mainCat,
			ptb.subCat
		from
			stemTransactions st
		join PatToBudcat ptb 
	on
			st.account = ptb.account
			and ptb.pat is not null
			and st.stemDesc like ptb.pat
	) stptb
where
	transactions.key = stptb.key
	and transactions.beneficiary is null
;

select changes();
