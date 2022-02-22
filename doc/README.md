# document for think frame

Mainly used for config read, database read and write, and net routing.

Users generally use the `T/*B` struct make by `Maker`. Except for `T/*B` created with `Maker`, `Coder` cannot be set outside the package and will not work

## T

struct with Encode (`T`'s `V` to []byte) and Decode ([]byte to `T`'s `V`)

***NOTE:*** `T`'s `V` must be a pointer (for Decode)

test list (outside package):

- 非`Maker`方式创建的`T`无法正常使用
- `V`不是pointer情况下`Decode`无法工作，`Encode`不受影响
- 各种`Maker`创建的`T`检查
