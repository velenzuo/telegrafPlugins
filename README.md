telegraf output 插件，将数据push 到open-falcon transfer    
将outputs/openfalcon 目录放到到 telegraf/plugins/outputs/    
编辑 telegraf/plugins/outputs/all/all.go 最后一行添加    
```
    _ "github.com/influxdata/telegraf/plugins/outputs/openfalcon"
```
重新编译telegraf

telegraf.conf    
```
[[outputs.openfalcon]]
  addr = "127.0.0.1:8433"
```

