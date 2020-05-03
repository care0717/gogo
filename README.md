# gogo 
gogo is toy golang compiler.

## How to

### compile
```shell script
make
```
### test
```shell script
make test
```

## Syntax
```ABNF
program    = stmt*
stmt       = expr ";" | "return" expr ";"
expr       = assign
assign     = equality ("=" assign)? 
equality   = relational ("==" relational | "!=" relational)*
relational = add ("<" add | "<=" add | ">" add | ">=" add)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary)*
unary      = ("+" | "-")? primary
primary    = num | ident | "(" expr ")"
num        = ("0"|...|"9")+
ident      = ("a"|...|"z")
```
