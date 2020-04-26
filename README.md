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
expr       = equality 
equality   = relational ("==" relational | "!=" relational)*
relational = add ("<" add | "<=" add | ">" add | ">=" add)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary)*
unary      = ("+" | "-")? primary
primary    = num | "(" expr ")"
num        = ("0"|...|"9")+
```
