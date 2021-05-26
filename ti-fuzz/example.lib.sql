select * from c where c.a = c.b;
select x, y from a order by z;
select x from a where x < y order by x;
select * from a where y != 1 and z > 100 group by x;
select * from (select * from a) b join (select * from c) d;
select count(a) from t where a = 1 and b + a > 100;
select t1.a, t1.b, t2.c from t1 inner join t2 on t1.a = t2.b + t2.c;
select avg(x) from t where x >= y;
select * from t where x - y <= z and y > (select count(*) from b);
