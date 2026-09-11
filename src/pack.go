package xtool

/*
主模式
-full 构建完整zip目录,输出全部文件的到meta,并输出o
-diff 基于full的输出和meta, 只将改变的文件写入out,不要修改meta
-sqld 先不做

-indir=输入目录/data/ui/web/
-C=切掉前缀如果是 `/data/ui` 则zip中只有 `/web/xxxx`
-o=输出的zip文件
-b=diff模式有效, 一个zip输入文件,以此为基础构建o
-meta=full模式的输出,diff模式的输入,csv格式:文件完整本地路径(-C=切掉前缀),文件大小byte,unix修改时间
--exclude-from=和tar的--exclude-from相同

main() 处理参数调用功能
diff() 对应diff模式
full() 对应full模式
sql_diff() 对应 sqld模式(空函数)
csv用开源库

代码2空格缩进, 函数内不要空行, 不要注释,
*/