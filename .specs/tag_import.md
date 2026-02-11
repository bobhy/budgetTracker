# Tag Imports

In Import view, we want to apply budget links to recognized imported transactions.
Use the raw sql query below, which works like this:
1. use other raw transaction fields to create a tag value.  At the moment, the tag is only based on description field.
2. map tags to budgets through the tags table

So:
1. add a "tag" button to the import view, which runs the sql query.  The button should appear just to the left of the "finalize" button.
2. When the user clicks the button and the query runs, report the number of records changed by the last update.  If the query fails, show the error in an error field in the view.

```sql
-- database: ./budget.db

-- remove fluff from description to create a good tag for mapping to a budget.
-- reviews *all* the raw_transactions, so this is the first filter to run.BEGIN IMMEDIATE
with 
t1 AS (
	-- undo WF "helpful" annotations
	select
		*,
		case
			when description like 'money transfer authorized on __/__ %' then substring(description, 35 + 1)
			when description like 'purchase authorized on __/__ %' then substring(description, 29 + 1)
			when description like 'purchase intl authorized on __/__ %' then substring(description, 34 + 1)
			
			-- purchase with cash back $$ authorized on dd
			else description
		end as stem1
	from
		raw_transactions
),
t2 as (
	-- remove merchant prefixes (other than the ones we want to use for tags)
	select
		*
		,
		case
			when stem1 like '___*%' then substring(stem1, 4 + 1)
			when stem1 like 'cash app*%' then substring(stem1, 9 + 1)
			when stem1 like 'zelle to %' then substring(stem1, 9 + 1)
			when stem1 like 'paypal *%' then substring(stem1, 8 + 1)
			else stem1
		end as stem
	from
		t1
)
update raw_transactions 
set tag = t2.stem
from t2 
where raw_transactions.id = t2.id
;

 -- apply budgets to known tag:budget mappings.
update raw_transactions
set budget = t.budget
from tags t where substr(raw_transactions.tag, 1, length(t.name)) = t.name;
```
