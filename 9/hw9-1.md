### 几种socket粘包的解包方式

1. Fixed-Length

根据指定长度进行包分解，如abcdef，如果给定指定长度为2，则分解为ab|cd|ef。

例: netty [FixedLengthFrameDecoder](https://netty.io/4.0/api/io/netty/handler/codec/FixedLengthFrameDecoder.html)

2. Delimiter-Based

根据指定分隔符进行分解，如`\n`

例: netty [DelimiterBasedFrameDecoder](https://netty.io/4.0/api/io/netty/handler/codec/DelimiterBasedFrameDecoder.html)

3. Length-Field-Based

根据消息头中指定的信息长度来确定消息长度。

例: netty [LengthFieldBasedFrameDecoder](https://netty.io/4.0/api/io/netty/handler/codec/LengthFieldBasedFrameDecoder.html)


