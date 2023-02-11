## Disclaimer

This implementation uses sequentional read. It means, that parts(chunks) placed on remote nodes are combined one-by-one. For small files it's pretty enough. But for larger ones parallel approach might be better option. Sequentional one has been choosen because of it doesn't block client on GET request, but starts transmission instantly.  

## Intro

This is a mono-repositary consisted from two parts - proxy and node.  

Node is a storage server. Any number of instances are OK. The only thing to do is to change proxy config correspodently by adding new nodes. 

Proxy is in charge of splitting and combining files regarding request. On PUT it splits file and uploads it's parts onto nodes, on GET it requests file parts from corresponding nodes in particular order and return it to the client.

Both parts support the same requests PUT/{ID} and GET/{ID} the only difference is that nodes work with parts(chunks) instead of proxy which works with full payload.

## Install and Run

```sh
go get github.com/3d0c/storage

# build proxy and node
make build
```

## Running and testing

Please take a look at `configs/dev` to be aware of pathes these programs use. If it's not suitted for you, change them. 
By default configuration there should be a 10 nodes, which proxy are expected to be.

```sh
# start nodes. (better to run it in separated terminal)
./start_nodes.sh
```

Now let's start the proxy

```sh
./proxy_main run --config=configs/dev/proxy.yaml
```

It should be up and running. All logs will be writtin to the STDOUT.

## Testing

A very simple test is by using `curl` just like that:

```sh
# Upload a file
# 127.0.0.1:8443 is defined in proxy.yaml
curl -H "Content-Type:application/octet-stream" --data-binary @anyfile -XPUT http://127.0.0.1:8443/file/1234656
# As a result you should get a 200OK

# Download the file
curl -O anyfile.new http://127.0.0.1:8443/file/1234656

# Test it works
md5 anyfile anyfile.new
```
