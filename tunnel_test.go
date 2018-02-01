package suxt

import (
    "testing"
    "github.com/nu7hatch/gouuid"
    "net"
    "time"
    "os/user"
    "bytes"
    "io"
    "strings"
)

func TestTunnel_ConnectDisconnect(t *testing.T) {
    id, err := uuid.NewV4()
    if err != nil {
        t.Error(err)
        t.FailNow()
    }
    path := "/tmp/" + id.String()
    currentUser, err := user.Current()
    username := currentUser.Username
    if err != nil {
        t.Error(err)
        t.FailNow()
    }
    tunnel := &Tunnel{
        User: username,
        Server: "localhost",
        Port: "22",
        KeyPath: "/home/" + username + "/.ssh/id_rsa",
        Timeout: time.Minute,
    }

    go func() {
        server, err := net.Listen("unix", path)
        if err != nil {
            t.Error(err)
            t.FailNow()
        }
        client, err := server.Accept()
        if err != nil {
            t.Error(err)
            t.FailNow()
        }
        _, err = io.Copy(client, strings.NewReader("pass"))
        if err != nil {
            t.Error(err)
            t.FailNow()
        }
        client.Close()
        server.Close()
    }()

    conn, err := tunnel.Connect(path)
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    buf := new(bytes.Buffer)
    _, err = io.Copy(buf, conn)
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    if buf.String() != "pass" {
        t.Error("invalid response:", buf.String())
    }
}

