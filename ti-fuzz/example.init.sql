create table t1 (i1 integer, c1 char);
insert into t1 values (2, 'a'), (1, 'b'), (3, 'c'), (0, null);
create table t2 (i2 integer, c2 char, f2 float);
insert into t2 values (0, 'c', null), (1, null, 0.1), (3, 'b', 0.01), (2, 'q', 0.12), (null, 'a', -0.1), (null, null, null);
select * from t1 where i1 = 1;
