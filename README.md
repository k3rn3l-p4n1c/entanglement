# Entanglement
Entanglement sync your data in cluster between node

#### How to use

| Node    | Ip       |
|---------|----------|
| Server1 | 10.0.0.1 |
| Server2 | 10.0.0.2 |

First run bootstrap

```go
conf := entanglement.DefaultConfig()

conf.RaftAddr = "10.0.0.1:12700"
conf.HttpAddr = "10.0.0.1:12701"

system := entanglement.Bootstrap(conf)
e := system.New("foo")

```

Now join other nodes

```go
conf := entanglement.DefaultConfig()

conf.RaftAddr = "10.0.0.2:12700"
conf.HttpAddr = "10.0.0.2:12701"
conf.JoinAddr = "10.0.0.1:12701"

system := entanglement.Bootstrap(conf)
e := system.New("foo")

```
Write data on Server1

```go
e.Set("baar")
```
Read data on the other node
```go
v, _ := e.Get("baar")
println(v)
```