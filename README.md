# autobuild

通过webhook实现自动拉取代码，然后执行docker-compose 重新构建运行命令

该程序运行在dokcer容器中，通过ssh连接到宿主机，在宿主机执行build命令
