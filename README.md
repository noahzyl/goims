GOIMS is a simple TCP instant messaging system with the basic framework provided by [@aceld](https://github.com/aceld) and implemented by GO.

## Structure

![goims-structure](goims-structure.png)

## Functions

- Basic TCP server and client
- User online and offline events
- Public chat
- Private chat
- Update(change) username
- Timeout-enforced offline

## Build

```
make
```

## Run

**server (use 127.0.0.1:5090)**

```
cd bin
./server
```

**client**

```
cd bin
./client -ip <ip of the server> -port <port of the server> 
```

## Clean

```
make clean
```

## Reference

[8小时转职Golang工程师(如果你想低成本学习Go语言)](https://www.bilibili.com/video/BV1gf4y1r79E/?spm_id_from=333.337.search-card.all.click&vd_source=28b15b3c60f368a0e9ca44b4ffcfdf19)