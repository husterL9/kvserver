-- kvserver_test.lua

-- 加载必要的 LuaSocket 和 gRPC 库
local grpc = require("grpc")
local pb = require("protobuf")

-- 初始化 gRPC 客户端
local function init_grpc_client()
    local channel = grpc.channel("localhost:50051")  -- 替换为实际的 gRPC 服务地址
    local stub = grpc.stub(channel, "KvService")
    return stub
end

-- 测试写入操作
local function write_test(client)
    for i = 1, sysbench.opt.threads do
        local key = string.format("key-%d", i)
        local value = string.format("value-%d", math.random(1000000))
        client:Set({key = key, value = value})
    end
end

-- 测试读取操作
local function read_test(client)
    for i = 1, sysbench.opt.threads do
        local key = string.format("key-%d", math.random(sysbench.opt.threads))
        client:Get({key = key})
    end
end

function thread_init(thread_id)
    client = init_grpc_client()
end

function thread_done(thread_id)
    client = nil
end

function event(thread_id)
    if sysbench.opt.test_type == "write" then
        write_test(client)
    elseif sysbench.opt.test_type == "read" then
        read_test(client)
    end
end

sysbench.cmdline.options = {
    test_type = {"Type of test to run: write or read", "write"},
}
