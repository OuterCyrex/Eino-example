# 存储引擎

## 一.mysql体系结构

mysql体系结构分为四层：连接层、服务层、引擎层、存储层

![image-20250325105344441](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250325105344441.png)

**连接层**：主要管理tcp连接，管理连接池。

**服务层**：管理DML、DDL语句的封装，解析、对优化进行查询，对查询结果进行缓存。

**引擎层**：提供可插拔式的存储引擎

**存储层**：用于与硬盘进行交互，存储相关数据。

## 二.存储引擎特点

存储引擎，可以广义理解为表的类型。其是基于**表**的设置而非**数据库**。

此处主要探究 **InnoDB**、**MyISAM** 和 **Memory**。

### 1.InnoDB

**InnoDB**的逻辑存储结构：

![image-20250325123814686](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250325123814686.png)

其中，**page**是**磁盘操作**的**最小单元**。

且 **InnoDB** 是 mysql 的默认存储引擎，其**支持事务，支持外键，支持行锁**。

### 2.MyISAM

**MyISAM** 是 Mysql 早期的默认存储引擎。

其**不支持事务**，**不支持外键**，**支持表锁不支持行锁**，访问速度快。

### 3.Memory

**Memory引擎** 的表数据时存储在内存中的，由于受到硬件问题、或断电问题的影响，只能将这些表作为临时表或缓存使用。
其**基于内存存放**，且默认是 **hash索引** 。

### 4.总结

![image-20250325124420123](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250325124420123.png)

- **InnoDB**：是Mysql的默认存储引擎，支持事务、外键。如果应用对事务的完整性有比较高的要求，在并发条件下要求数据的一致

  性，数据操作除了插入和查询之外，还包含很多的更新、删除操作，那么InnoDB存储引擎是比较合适的选择。

- **MyISAM**：如果应用是以读操作和插入操作为主，只有很少的更新和删除操作，并且对事务的完整性、并发性要求不是很高，那

  么选择这个存储引擎是非常合适的。

- **MEMORY**：将所有数据保存在内存中，访问速度快，通常用于临时表及缓存。MEMORY的缺陷就是对表的大小有限制，太大的表

  无法缓存在内存中，而且无法保障数据的安全性。

# 索引

**索引 (index) **是帮助 MySQL 高效获取数据的**数据结构**（有序）。在数据之外，数据库系统还维护着满足特定查找算法的数据结构，这些数据结构以某种方式引用（指向）数据，这样就可以在这些数据结构上实现高级查找算法，这种数据结构就是索引；**索引**存储在**引擎层**内。

Mysql 的索引有四大类：**B+Tree索引**、**Hash索引**、**R-tree 索引**、**Full-text 索引 (全文索引)**。

![image-20250325130056884](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250325130056884.png)

其中，**Hash索引** 只支持 **Memory引擎**。 **R-tree索引** 只支持 **MyISAM引擎**。

**全文索引 **支持 **InnoDB** 与 **MyISAM** 引擎。

## 一. B+ Tree索引

### 1.B树

**N阶的B树** 指的是， **N - 1** 个 key，**N** 个指针。

![image-20250325130836978](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250325130836978.png)

**B树**的插入逻辑是，每次某个节点满足 **N - 1** 个 **key** 再插入时，会将这 N 个 key 排列，取中间的数值上提。

![image-20250325131425487](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250325131425487.png)

### 2.B+树

与B树不同的是，B+树的所有叶子节点之间形成了一个单向链表，且所有的元素都位于单向链表上 (包括上方的父节点)

![image-20250326164324673](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250326164324673.png)

每次满足**N - 1** 个 **key** 再插入时，会将中间元素上提，但仍然保留在叶子节点，且从该元素分裂，形成单向链表。

![image-20250326170049564](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250326170049564.png)

### 3.Mysql B+树索引

Mysql的 **B+树索引** 在 **B+树** 的基础上再优化，加上了双向的指针，在叶子节点之间形成了一个**双向链表**。

![image-20250326170552459](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250326170552459.png)

**B+树** 相较于其他搜索树，如**二叉搜索树**而言，

1. 层级更少，搜索效率高
2. 对于 **B树**，无论是叶子节点还是非叶子节点，都会保存数据，这样导致一
   页中存储的键值减少，指针跟着减少，要同样保存大量数据，只能增加树的
   高度，导致性能降低

当数据重复时，B+树会为重复的数据建立一个**链表**，将指针指向链表的头部。

### 4.B+树性能分析

通常而言，一**页**的大小为16KB。**InnoDB **的指针占用 **6个字节** 的空间，主键即使为bigint，占用字节数为8。
则当B+树的高度为2时：

```latex
n * 8 + (n + 1) * 6 = 16 * 1024
```

得到算出 n 约为1170，此时可以存储的行数为：

```latex
1171 * 16 = 18736
```

若高度为3，则可以存储：

```latx
1171 * 1171 * 16 = 21939856
```

## 二. Hash 索引

即通过一定的 **Hash算法** 将 key 映射到某个特点地址上。只支持 **Memory引擎**。

在**Hash冲突**时，则在对应地址上创建一个**单向链表**。

特点：

1. Hash索引只能用于**对等比较**(=，in)，不支持范围查询（(between,>，<，...）
2. 无法利用索引完成排序操作
3. 查询效率高，通常只需要一次检索就可以了，效率通常要高于B+tree索引

**InnoDB引擎** 中具有自适应hash功能，hash索引是存储引擎根据B+Tree索引在指定条件下自动构建的。

即：如果 **InnoDB** 发现某些查询模式（例如频繁的等值查询）可以通过构建Hash索引来加速查询，它会自动在内存中为这些频繁访问的 **B+Tree 索引**节点创建 **Hash索引**。

## 三.索引分类

根据索引类型常常分为：**主键索引、唯一索引、常规索引、全文索引**。

![image-20250326180200220](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250326180200220.png)

根据索引的存储方式，又可以分为：**聚集索引、二级索引**

![image-20250326180305368](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250326180305368.png)

**聚集索引**选取规则：

- 如果存在主键，主键索引就是聚集索引。
- 如果不存在主键，将使用第一个 **唯一(UNIQUE)** 索引作为聚集索引。
- 如果表没有主键，或没有合适的唯一索引，则 **InnoDB** 会自动生成一个 **row_id** 作为隐藏的聚集索引。

![image-20250326195230783](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250326195230783.png)

如上图所示，上方依靠 `ID` 进行**聚簇索引**的结构，下方为依据 `name` 字段进行**二级索引**的结构。

通过二级索引获取聚簇索引的信息，并再次通过聚簇索引获取该行全部信息的过程，称为 **回表查询**。

# 优化

## 一.性能分析

### 1. 慢查询日志

**慢查询日志 **记录了所有执行时间超过指定参数 (long_query_time，单位：秒，默认10秒）的所有SQL语句的日志。
MySQL 的 **慢查询日志** 默认没有开启，需要在MySQL的配置文件（/etc/my.cnf）中配置。

### 2. profile

**show profiles** 能够在做 **SQL优化** 时帮助我们了解时间都耗费到哪里去了。

通过 **have_profiling** 参数，能够看到当前MySQL是否支持profile操作：
```mysql
SELECT @@have_profiling;
```

效果：

![image-20250326204446172](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20250326204446172.png)

### 3. explain

explain用于阐述一次查询的具体步骤，其返回结果有多个字段：

1. **ID**：SELECT 查询的序列号，表示查询中执行 SELECT 子句或者是操作表的顺序 (id相同，执行顺序从上到下；id不同，值越大，越先执行)。
2. **Select_type**：表示 SELECT 的类型，常见的取值有**SIMPLE**（简单表，即不使用表连接或者子查询）、**PRIMARY**（主查询，即外层的查询）、**UNION**（UNION中的第二个或者后面的查询语句）、**SUBQUERY**（SELECT / WHERE之后包含了子查询）等
3. **type**：表示连接类型，性能由好到差的连接类型为NULL、system、const、eq_ref、ref、range、index、all。
4. **possible_key**：可能应用在这张表上的索引。
5. **Key**：实际使用的索引
6. **Key_len**：表示索引中使用的字节数，该值为索引字段最大可能长度，并非实际使用长度，在不损失精确性的前提下，长度越短越好。如：`SELECT * FROM users WHERE username LIKE 'abc%';`，则key_len为3，因为只需要匹配前三个。
7. **rows**：MySQL认为必须要执行查询的行数，在 **innoDB** 引擎的表中，是一个估计值，可能并不总是准确的。
8. **filtered**：表示返回结果的行数占需读取行数的百分比，filtered 的值越大越好。

| type类型 | 说明                                                         | 示例                                                      |
| -------- | ------------------------------------------------------------ | --------------------------------------------------------- |
| NULL     | 没有查询任何表                                               | `SELECT 'A';`                                             |
| system   | 表中只有一行数据，是 `const` 类型的特例。                    | `SELECT * FROM User WHERE ID = 1 LIMIT 1;`                |
| const    | 查询的条件为常数，MySQL 只需要读取一次即可找到匹配的行。     | `SELECT * FROM User WHERE ID = 1;`                        |
| eq_ref   | 使用索引列进行等值连接，通常用于主键或唯一索引的等值查询。   | `SELECT * FROM User u JOIN Profile p ON u.ID = p.UserID;` |
| ref      | 使用索引列进行等值查询，但索引列不是主键或唯一索引。         | `SELECT * FROM User WHERE Status = 'active';`             |
| range    | 使用索引列进行范围查询，通常用于范围条件（如 `BETWEEN` 或 `>`）。 | `SELECT * FROM User WHERE Age BETWEEN 18 AND 30;`         |
| index    | 全索引扫描，扫描整个索引，但不扫描整个表。                   | `SELECT * FROM User ORDER BY Age;`                        |
| all      | 全表扫描，扫描整个表，性能最差。                             | `SELECT * FROM User;`                                     |

