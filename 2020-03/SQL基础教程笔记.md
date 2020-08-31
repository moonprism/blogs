## 数据库和SQL

### sql种类分类

* DDL（Data Definition Language，数据定义语言）
    * `CREATE`
    * `DROP`
    * `ALTER`
* DML（Data Manipulation Language，数据操纵语言）
    * `SELECT`
    * `INSERT`
    * `UPDATE`
    * `DELETE`
* DCL（Data Control Language，数据控制语言）
    * `COMMIT`
    * `ROLLBACK`
    * `GRANT`
    * `REVOKE`

### 数据类型的指定

### 约束的设置

* `NOT NULL`

### 表定义的更新（`ALTER TABLE`）

* 更新
```sql
ALTER TABLE table_name ADD column_name VARCHAR(100);
```

* 删除
```sql
ALTER TABLE table_name DROP column_name;
```

## 基础查询

> 不能对 null 使用比较运算符，应该使用 is null / is not null

## 聚合与排序

### 常用聚合函数

* `COUNT`
* `SUM`
* `MAX`
* `MIN`
* `AVG`

`GROUP BY`子句将表中数据分为多组进行管理，在 `GROUP BY` 子句中指定的列称为聚合键或者分组列

### 书写顺序

1. SELECT → 2. FROM → 3. WHERE → 4. GROUP BY

使用聚合时，SELECT 子句中只能存在以下三种元素。

1. 常数
2. 聚合函数
3. GROUP BY子句中指定的列名（也就是聚合键）

只有 SELECT 子句和 HAVING 子句（以及之后将要学到的 ORDER BY 子句）中能够使用 COUNT 等聚合函数

### `HAVING`子句

* HAVING子句要写在GROUP BY子句之后。
* WHERE子句用来指定数据行的条件，HAVING子句用来指定分组的条件。

```sql
SELECT product_type, COUNT(*)
  FROM Product
 GROUP BY product_type
HAVING COUNT(*) = 2;
```

## 数据更新

* insert
```sql
INSERT INTO ProductIns VALUES
 ('0002', '打孔器', '办公用品', 500, 320, '2009-09-11');
```
* delete
```sql
DELETE FROM Product
 WHERE sale_price >= 4000;
```
* update
```sql
UPDATE Product
   SET sale_price = sale_price + 1000
 WHERE product_name = 'T恤衫';
```

### 事务

```sql
START TRANSACTION
...
COMMIT / ROLLBACK
```

事务的ACID特性

* 原子性（Atomicity）
原子性是指在事务结束时，其中所包含的更新处理要么全部执行，要么完全不执行

* 一致性（Consistency）
一致性指的是事务中包含的处理要满足数据库提前设置的约束，如主键约束或者 NOT NULL 约束等

* 隔离性（Isolation）
隔离性指的是保证不同事务之间互不干扰的特性。该特性保证了事务之间不会互相嵌套。此外，在某个事务中进行的更改，在该事务结束之前，对其他事务而言是不可见的。因此，即使某个事务向表中添加了记录，在没有提交之前，其他事务也是看不到新添加的记录的。

* 持久性（Durability）
持久性也可以称为耐久性，指的是在事务（不论是提交还是回滚）结束后，DBMS 能够保证该时间点的数据状态会被保存的特性。即使由于系统故障导致数据丢失，数据库也一定能通过某种手段进行恢复。
如果不能保证持久性，即使是正常提交结束的事务，一旦发生了系统故障，也会导致数据丢失，一切都需要从头再来。
保证持久性的方法根据实现的不同而不同，其中最常见的就是将事务的执行记录保存到硬盘等存储介质中（该执行记录称为日志）。当发生故障时，可以通过日志恢复到故障发生前的状态。

## 复杂查询

### 视图

视图相当于是可重用的子查询

### 子查询

子查询相当于是一次性的视图，可以直接用于 FROM 后

```sql
SELECT product_type, cnt_product
  FROM (SELECT *
          FROM (SELECT product_type, COUNT(*) AS cnt_product
                  FROM Product
                 GROUP BY product_type) AS ProductSum -----①
         WHERE cnt_product = 4) AS ProductSum2; -----------②
```

### 标量子查询

必须而且只能返回 1 行 1 列的结果，可以在 SQL 中当作常量使用

### 关联子查询

![](https://kicoe-blog.oss-cn-shanghai.aliyuncs.com/VSqwreWlpzKLveThdwNF.jpg)

1. 关联子查询会在细分的组内进行比较时使用。
2. 关联子查询和GROUP BY子句一样，也可以对表中的数据进行切分。
3. 关联子查询的结合条件如果未出现在子查询之中就会发生错误。

## 函数、谓词、CASE 表达式

=、<、>、<> 等比较运算符，其正式的名称就是比较谓词。

* LIKE
* BETWEEN
* IS NULL、IS NOT NULL
* IN
* EXISTS

### CASE

```sql
SELECT SUM(CASE WHEN product_type = '衣服'
                THEN sale_price ELSE 0 END) AS sum_price_clothes,
       SUM(CASE WHEN product_type = '厨房用具'
                THEN sale_price ELSE 0 END) AS sum_price_kitchen,
       SUM(CASE WHEN product_type = '办公用品'
                THEN sale_price ELSE 0 END) AS sum_price_office
  FROM Product;
```

case 分组

```sql
SELECT 
CASE WHEN `pay_count` < 50 THEN 0
     WHEN `pay_count` >= 50 AND `pay_count` < 100 THEN 1
     WHEN `pay_count` >= 100 AND `pay_count` < 300 THEN 2
     WHEN `pay_count` >= 300 AND `pay_count` < 500 THEN 3
     WHEN `pay_count` >= 500 AND `pay_count` < 1000 THEN 4
     WHEN `pay_count` >= 1000 AND `pay_count` < 3000 THEN 5
     WHEN `pay_count` >= 3000 AND `pay_count` < 5000 THEN 6
     WHEN `pay_count` >= 5000 AND `pay_count` < 10000 THEN 7
     WHEN `pay_count` >= 10000 AND `pay_count` < 30000 THEN 8
     WHEN `pay_count` >= 30000 THEN 9
END AS lv, 
COUNT(*) as pay_user, 
SUM(`pay_count`) as pay_count 
FROM tf_report_user_pay 
WHERE 1 {$where} GROUP BY lv ORDER BY lv ASC
```
>[success] 上面这段是当初实习时候写的sql，现在还能在笔记里面翻到哈哈，可能那时写完被老大夸了高兴才特别记下的吧。

## 集合运算

主要是内联结与左右的两个外联结

```sql
SELECT SP.shop_id, SP.shop_name, SP.product_id, P.product_name, P.sale_price
  FROM ShopProduct AS SP INNER JOIN Product AS P ----①
    ON SP.product_id = P.product_id
 WHERE SP.shop_id = '000A';
```

外联结会关联所有数据，比如最近遇到的同步两张表的数据

```sql
update t1 inner join t2 on t1.id = t2.id set t1.c = t2.c;
```

如果以上sql用的是`left join`，那么就算没有在`on`中匹配到t1与t2有相等id，都会去更新t1.c = null。