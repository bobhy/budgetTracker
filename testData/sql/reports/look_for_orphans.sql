-- sqlite
select *
from Transactions t
where beneficiary is null
order by description,
    beneficiary,
    category,
    subcategory,
    postedDate
;
