# File Parsers

Parse texts (Json, Yaml, Excel) to code's objects. Excel use for game's configuration table

* `Json:` Use `IParser.LoadFile` or `IParser.LoadJsonBytes`
* `Yaml:` Use `IParser.LoadFile` or `IParser.LoadYamlBytes`
* `Excel:` Use `UnmarshalExcel`

## Excel

Support `xlsx` `xlsm` `xltm` `xltx` format file.
Use `excel:"tag"` to specified field.
the value only support base type

*NOTE:* not support `xls` format. please convert to xlsx before
