
select sum(amount) from Transactions;

select 
    strftime('%Y-%m', postedDate, 'start of month') as month,
    a.account,
    count(*) as num,
    sum(t.amount) as amount
from Transactions t
    join Accounts a on t.Account = a.Account
group by 
    date(month, 'start of month'),
    t.account
order by month, t.account;