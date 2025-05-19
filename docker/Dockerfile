FROM golang:tip-bullseye AS builder
#它是golang1.19.2 + ubuntu20.04
#ARG SVR 在容器里面建立SVC的环境变量,并将宿主机的SVR环节变量的内容赋值给该变量
ARG SVC
ARG GOARCH
ARG GOARM
ARG VERSION
ARG COMMIT
ARG TIME
#ENV GOPROXY 在容器里面建立GOPROXY的环境变量,保证下载golang库的时候足够使用国内的源
ENV GOPROXY=https://goproxy.cn

#指定docker上的工作目录是/go/src/digger_assistant
WORKDIR /go/src/digger_assistant

#拷贝当前目录下的所有文件和子目录到docker的工作目录中(即/go/src/digger_assistant)
COPY . .
#下载相关的所有库
RUN go mod tidy -compat="1.17"


#编译将编译结果输出到 build/main
#参考 https://zhuanlan.zhihu.com/p/338891206
#RUN CGO_ENABLED=0 GOARCH=amd64 go build -o build/main -ldflags "-s -w" digger_assistant
#RUN CGO_ENABLED=1 GOARCH=amd64 CGO_LDFLAGS="-static" go build -o build/main -ldflags "-s -w" digger_assistant
RUN CGO_ENABLED=1 GOARCH=amd64 go build -o build/main -ldflags "-s -w" digger_assistant

#list dependencies of main
#RUN ldd build/main

#将 build/main 移动到 /main,至此编译完毕不需要在打包机上运行其他东西了
RUN mv build/main /main

# 从干净的系统开始重做docker镜像
#临时用这种方式解决动态库找不到的问题
#FROM scratch 
FROM busybox:glibc
#RUN apk add --no-cache tzdata
#ENV TZ Asia/Shanghai !!!不用写时区否则时间会出问题
#COPY --from=builder linux-vdso.so.1 .
COPY --from=builder /lib/x86_64-linux-gnu/libm.so.6 /lib/x86_64-linux-gnu/libm.so.6
COPY --from=builder /lib/x86_64-linux-gnu/libpthread.so.0 /lib/x86_64-linux-gnu/libpthread.so.0
COPY --from=builder /lib/x86_64-linux-gnu/libdl.so.2 /lib/x86_64-linux-gnu/libdl.so.2
COPY --from=builder /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so.6
COPY --from=builder /lib64/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2

# 新镜像工作目录在docker容器的根路径下
WORKDIR /
COPY ./conf/app.conf ./conf/
#COPY ./data/data.db ./data/

#拷贝打包机的/main文件到当前路径下,拷贝之后的文件也是 /main
COPY --from=builder /main .
#RUN chmod +x /main

#docker容器的启动程序为 /main
ENTRYPOINT ["/main"]
