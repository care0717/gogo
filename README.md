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
expr    = mul ("+" mul | "-" mul)*
mul     = primary ("*" primary | "/" primary)*
primary = num | "(" expr ")"
num     = ("0"|...|"9")+
```
