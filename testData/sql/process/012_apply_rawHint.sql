-- CONNECTION: database=/home/bobhy/db/budget/budget1.db
-- categorize based on 'rawHint', which is provided only by Capitalone currently
-- don't overwrite categories if transaction already categorized (i.e, beneficiary is non null)
-- if the hint doesn't specify beneficiary (usually the case), then beneficiary is based on account of transaction.

update
	Transactions
set
	beneficiary = stptb.beneficiary
	, category = stptb.mainCat
	, subCategory = stptb.subCat
	, tag = stptb.tag
from
	(
	select
			st.key
		, ptb.rawHintPat as tag
		, coalesce(ptb.beneficiary, a.beneficiary) beneficiary
		, ptb.mainCat
		, ptb.subCat
	from
			stemTransactions st
	join PatToBudcat ptb 
	on
			ptb.rawHintPat is not null
		and st.rawHint like ptb.rawHintPat
	join Accounts a 
			on
		st.account = a.account
	where st.beneficiary is null
	) stptb
where
	transactions.key = stptb.key
	and transactions.beneficiary is null
;

select changes();