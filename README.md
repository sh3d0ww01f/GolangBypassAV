# GolangBypassAV
研究利用golang来bypassAV

## 说明
免杀这块本来就不是web狗擅长的，而且作为一个web狗也没必要花太多时间来折腾这个，达到能用就行，不要追求全部免杀，能免杀目标就行。


## 思路

## 命令

```bash

效果一般
go build -ldflags="-s -w" -o main1.exe -race main.go

效果还可以
go build -ldflags="-s -w" -o main1.exe

效果还可以
go build -ldflags="-s -w -H=windowsgui" -o main2.exe

```