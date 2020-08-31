## upsert 

upsert是一个mongo中的函数，指更新数据的时候没有匹配相应的doc则会插入一个新的doc。

在mysql中一般使用 `on duplicate key update` 实现

```sql
insert set ..=.. where index = 1 and index2 = 2 
on duplicate key update 
update set index=values(index)...
```

## laravel 

Laravel中的函数`updateOrCreate`是先search再判断返回结果执行insert/update。想要用 `on duplicate key update` 只好自己扩展Model了：

```php
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Query\Expression;
use Illuminate\Support\Arr;
use Illuminate\Support\Facades\DB;

trait ModelTrait
{
    public static function insertOrUpdate(array $values) {
        if (empty($values)) {
            return true;
        }

        if (! is_array(reset($values))) {
            $values = [$values];
        }

        /** @var Model $class */
        $class = new static;

        $sql = $class->getConnection()->getQueryGrammar()->compileInsert($class->getQuery(), $values);
        $bindings = array_values(array_filter(Arr::flatten($values, 1), function ($binding) {
            return ! $binding instanceof Expression;
        }));

        // 扩展 Illuminate\Database\Query\Builder::insert
        $columns = implode(',', array_map(function ($c){
            return "`$c`=values(`$c`)";
        },array_keys($values[0])));

        $sql .= " on duplicate key update $columns";

        return DB::insert($sql, $bindings);
    }

}
```

> 花了好些时间才理清 Laravel query builder