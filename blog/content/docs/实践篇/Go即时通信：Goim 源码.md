# Go即时通信：Goim 源码

Tags: Go
Last edited time: June 13, 2024 10:24 AM
状态: 已发布

### **一.为什么想要看goim源码**

想了解go如何实现一个**长连接的网关**，并且如何通过服务端来下发配置。

**为什么选goim去看**

主要是看中goim的github介绍中的以下几点：

- 支持单个推送、多个推送和广播
- 支持多种协议：WebSocket，TCP，HTTP
- 支持一个key被多个用户订阅，并且可配置最大用户数

**还有一个比较重要的就是纯go实现。**

### **二.整体架构解析**

goim 的架构可参照github的架构图

![Untitled](../../../Go即时通信：Goim%20源码%2060a94c13395e4e0ea7f2390634b8670a/Untitled.png)

具体解析一下：

**2.1.客户端发起TCP连接与Comet 建立连接**

可以通过 `benchmarks/client/main.go` 查看到 Go 通过 TCP 发起连接

```go
// connnect to server
conn, err := net.Dial("tcp", address)
```

连接后需要带上参数来加入特定的房间和告诉 `comet` 自己想要接收的信息。

```go
// AuthToken auth token.
type AuthToken struct {
    Mid      int64   `json:"mid"`
    // 对应唯一标识的key
    // 这个也可以用来定位server 中channel 所在的bucket，server为comet的server
    // 定位逻辑是：
	  //   idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	  //   看名字是定位了同个城市的所有连接
    // 这个key可以定位到若没传则用uuid.new
    Key      string  `json:"key"` 
    
    // 要加入的房间id
    RoomID   string  `json:"room_id"`
    // 对应客户端平台，如ios，安卓
    Platform string  `json:"platform"`  
    // 需要接收的信息
    Accepts  []int32 `json:"accepts"` 
}
```

2.2.**comet 将链接建立起来并new 出一个Channel 结构体**，用于保存对应的单条链接。Channel 结构体如下

```go
// Channel used by message pusher send msg to write goroutine.
type Channel struct {
		// 所属的房间
    Room     *Room 
    // 这个给客户端发送的数据，结构是环形的
    CliProto Ring 
    // 读取要推送给客户端的chan
    signal   chan *grpc.Proto 
    // 接管了对应net.conn链接的写
    Writer   bufio.Writer 
    // 接管了对应net.conn链接的读
    Reader   bufio.Reader 
    Next     *Channel
    Prev     *Channel 

		// 客户端表示的mid
    Mid      int64 
    // 唯一标识链接的key
    Key      string 
    // 客户端的ip，赋值为net.SplitHostPort(conn.RemoteAddr().String())
    IP       string 
    // 需要监听的操作
    watchOps map[int32]struct{} 
    mutex    sync.RWMutex
}
```

2.3.comet 处理 TCP 连接

```go
func acceptTCP(server *Server, lis *net.TCPListener) {
	for {
		go serveTCP(server, conn, r)
	}
}
```

通过authTCP 方法，以grpc的方式调用了 logic 模块的 Connect

```go
// ServeTCP serve a tcp connection.
func (s *Server) ServeTCP(conn *net.TCPConn, rp, wp *bytes.Pool, tr *xtime.Timer) {
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Mid, ch.Key, rid, accepts, hb, err = s.authTCP(ctx, rr, wr, p); err == nil {
			ch.Watch(accepts...)
			b = s.Bucket(ch.Key)
		}
	}
}

```

2.4.将**对应的链接和key设置到redis中，即为在线链接。**

2.5.业务方将要推的消息发送给logic，push消息的协议如下

```go
message PushMsg {
    enum Type {
		    // 单条消息push
        PUSH = 0; 
        // 按房间push
        ROOM = 1; 
        // 广播
        BROADCAST = 2; 
    }
    Type type = 1;
    // 消息操作的类型
    int32 operation = 2; 
    int32 speed = 3;
    // 定位是在哪个comet server
    string server = 4; 
    // 房间号
    string room = 5; 
    // 要推的链接的key
    repeated string keys = 6;
    // 消息体
    bytes msg = 7; 
}
```

logic接收到消息后，将消息放入kafka 等待job消费。

2.6.job消费完之后，通过job包下的 `comet.process()` 以grpc的方式向comet进行消息的发送。

以下则是comet push单条消息的代码

```go
func (s *server) PushMsg(ctx context.Context, req *pb.PushMsgReq) (reply *pb.PushMsgReply, err error) {
    if len(req.Keys) == 0 || req.Proto == nil {
        return nil, errors.ErrPushMsgArg
    }
    for _, key := range req.Keys {
        if channel := s.srv.Bucket(key).Channel(key); channel != nil {
            // 检查客户端的accepts 参数，看这条消息是否需要推送
            if !channel.NeedPush(req.ProtoOp) {
                continue
            }
            // 调用channel push方法推送消息。
            // 这里其实是给ch.signal 发送了req.Proto数据，
            // 然后由dispatchTCP这个方法的携程去取出chan里的数据用于发送
            if err = channel.Push(req.Proto); err != nil {
                return
            }
        }
    }
    return &pb.PushMsgReply{}, nil
}
```

综上就是goim整个建立链接、存对应关系，再由业务放发送消息的整个流程。

另外发现Comet是用bilibli 服务发现项目discovery来做的，comet可以横向扩展.

其中chan.CliProto 使用的是Ring结构体。留下了todo ，所以将其抽取出来讨论一下。

### **三.发现的Ring 的 TODO padding**

Ring 结构体代码如下

> 关于完整的代码可以看：[https://github.com/Terry-Mao/goim/blob/master/internal/comet/ring.go](https://link.zhihu.com/?target=https%3A//github.com/Terry-Mao/goim/blob/master/internal/comet/ring.go)
> 

```go
// Ring ring proto buffer.
type Ring struct {
    // read
    rp   uint64 // 读取的位置
    num  uint64mask uint64// TODO split cacheline, many cpu cache line size is 64
    // pad [40]byte
    // write
    wp   uint64 // 写的位置
    data []grpc.Proto // grpc调用的数据
}
```

不难理解，这是一个用于存放 `grpc.Proto` 的环形队列。但是这个存在伪共享的问题

**什么是伪共享？**

当CPU执行运算的时候，它先去L1查找所需的数据、再去L2、然后是L3，如果最后这些缓存中都没有，所需的数据就要去主内存拿。

**走得越远，运算耗费的时间就越长。**

所以如果在做一些很频繁的事，要尽量确保数据在L1缓存中。

另外，线程之间共享一份数据的时候，需要一个线程把数据写回主存，而另一个线程访问主存中相应的数据。

Cache是由很多个cache line组成的。每个cache line通常是64字节，并且它有效地引用主内存中的一块儿地址。CPU每次从主存中拉取数据时，会把相邻的数据也存入同一个cache line。

但其实这里rp 跟wp 其实是**不需要每次都同时进入缓存行中的，**只有在队满的时候或者队空的时候才需要两个同时判断，中间读取数据`wp--` 的时候，其实rp并不受影响，但是因为每次要读入整个缓存行中，所以rp和wp又是被同时读入的，刷会主存的时候其实rp没有任何变化，但还是需要刷回主存。

这里留下了TODO，**为什么不能像java一样通过空间换时间的padding 来凑齐缓存行的64个字节**

**java 的 解决方式可以通过padding 的方式来防止rb 和 wb被同时读到同一个CPU缓存行内，防止需要频繁的刷新主存。**

而go 的推荐做法其实是想通过通信来共享内存，而不是内存来通信。所以也并不提供内存可见的关键字。

如果真的需要解决伪共享的问题，可以通过 `sync.atomic` 包或者 `sync.Mutex` 来解决。

## 四.拓展思考

Goim这里实线了在线用户房间的单播、多播和广播，像我们平时用的最多的群聊是会加入到一个群组中，然后就算我们不在线，消息依然会在我们上线后推送到我们的消息列表中，那这个该如何实现呢？

最简单的做法是通过数据库保存下来群聊的消息，如果有用户加入了群聊并且有用户发送了消息，这个时候就将消息写入到数据库中，等到用户上线的时候，再去查询群聊里的所有新消息进行展示。

这种方式属于写扩散，一旦群聊人员变多之后，数据量会有一个激增。

还有另一种方式，每次有新消息不立即写入每个用户的消息列表中，等到用户上线再去拉取自己相关群聊的消息，这种则属于读扩散。

它们各自有优缺点。

**写扩散优点：** 数据独立，每个用户维护好自己的消息队列，某个用户群聊的消息丢了并不影响其他的用户。读取数据简单，用户只需要拉取自己的消息列表既可以拉取到所有数据，读取时无需额外计算

**写扩散缺点：** 要等到所有数据扩散完才能看到，实时性差。浪费存储，如果有的用户已经不登录了也不需要写入这条消息

**读扩散优点：** 数据实时性实时高，用户收到群聊的消息后可以立马去群聊的地方拉取到消息。写入逻辑简单。当读少写多时可以省下扩散成本

**读扩散缺点：** 读取数据逻辑复杂，如果用户消息的数据做了分区，那么需要从多个地方读取数据。并且需要聚合排序等。读取的数据会造成热点问题，有些群聊特别活跃的话会频繁的读取，如果刚好都分在同一个分区内则会造成资源倾斜的情况

如果每条消息都要发送给群聊所有用户，资源消耗会比较大，所以只在用户触发到群聊或者页面上浏览到的时候进行消息推送，平时则只拉取是否有 **最新的消息以及当前状态的最后一条消息即可。**

我们将读扩散和写扩散进行结合：

如果群聊用户是 **在线的状态，** 没有屏蔽群聊，则及时进行消息推送。

用户屏蔽了群聊， **则写入群聊未读信息统计及最后一条消息的预览，为了提升用户在浏览列表的体验。** 这里需要维护一个用户群聊消息未读和最后一条消息的表，来实现用户群聊信息的预览

用户是 **离线状态，** 用户未设置群聊消息免打扰，则推入消息中心。用户设置消息免打扰，则不推入消息中心中。

![Untitled](../../../Go即时通信：Goim%20源码%2060a94c13395e4e0ea7f2390634b8670a/Untitled1.png)

可以通过在线状态的拉取来看是否需要给用户发送， **用户在线状态则可以通过定时发送心跳来维护， 并且用** ETCD 实现状态变更和监听来进行同步。

如果用户在线的话则需要通过 `kv-store` 来找到存储 `websocket` 连接的服务器，转发给对应的 `chat-server` 进行消息的推送。

如果不在线则通过 api-server 拉取查看是否有屏蔽消息，没有的话直接推入 Notification Server 中。

如果用户关闭了通知，则等待用户刷新消息列表的时候再进行消息拉取。