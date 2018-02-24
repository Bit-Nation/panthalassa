package main

import (
  "gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p"
  "crypto/rand"
  "context"
  crypto "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
  "fmt"
  ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
  "log"
  host "gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
  "bufio"
  net "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
  dht "gx/ipfs/QmVSep2WwKcXxMonPASsAJ3nZVjfVMKgMcaSigxKnUWpJv/go-libp2p-kad-dht"
  ds "gx/ipfs/QmPpegoMqhAEqjncrzArm7KVWAkCm78rqL2DPuNjhPrshg/go-datastore"
)

func makei(ctx context.Context) host.Host {

  // Generate a key pair for this host. We will use it at least
  // to obtain a valid host ID.
  priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

  if err != nil {
    panic(err)
  }

  opts := []libp2p.Option{
    libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
    libp2p.Identity(priv),
  }

  opts = append(opts, libp2p.NoEncryption())

  basicHost, err := libp2p.New(ctx, opts...)

  if err != nil {
    panic(err)
  }

  hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

  // Now we can build a full multiaddress to reach this host
  // by encapsulating both addresses:
  addr := basicHost.Addrs()[0]
  fullAddr := addr.Encapsulate(hostAddr)

  log.Printf("I am %s\n", fullAddr)

  return basicHost

}

func main() {

  ctx := context.Background();

  h := makei(ctx)

  d := dht.NewDHT(ctx, h, ds.NewMapDatastore())

  d.Bootstrap(ctx)

  d.FindPeersConnectedToPeer();

}

// doEcho reads a line of data a stream and writes it back
func doEcho(s net.Stream) error {
  buf := bufio.NewReader(s)
  str, err := buf.ReadString('\n')
  if err != nil {
    return err
  }

  log.Printf("read: %s\n", str)
  _, err = s.Write([]byte(str))
  return err
}
