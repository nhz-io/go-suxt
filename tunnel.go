package suxt

import (
    "net"
    "golang.org/x/crypto/ssh"
    "time"
    "errors"
    "io/ioutil"
)

type Tunnel struct {
    User string
    Server string
    Port string
    KeyPath string
    Timeout time.Duration
    tunnel *ssh.Client
    socket net.Conn
}

func (this *Tunnel) getSigner() (signer ssh.Signer, err error) {
    buf, err := ioutil.ReadFile(this.KeyPath)
    if err != nil {
        return
    }

    signer, err = ssh.ParsePrivateKey(buf)
    return
}

func (this *Tunnel) getClientConfig () (conf *ssh.ClientConfig, err error) {
    signer, err := this.getSigner()
    if err != nil {
        return
    }

    conf = &ssh.ClientConfig{
        User: this.User,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout: this.Timeout,
        Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
    }

    return
}

func (this *Tunnel) Connect (addr string) (conn net.Conn, err error) {
    if this == nil {
        err = errors.New("missing tunnel")
        return
    }
    if this.socket != nil || this.tunnel != nil {
        err = errors.New("aready connected")
        return
    }

    conf, err := this.getClientConfig()
    if err != nil {
        return
    }

    tunnel, err := ssh.Dial("tcp", net.JoinHostPort(this.Server, this.Port), conf)
    if err != nil {
        return
    }

    conn, err = tunnel.Dial("unix", addr)
    if err != nil {
        _ = tunnel.Close()
        return
    }

    this.tunnel = tunnel
    this.socket = conn

    return
}

func (this *Tunnel) Disconnect () error {
    if this == nil {
        return errors.New("missing tunnel")
    }
    if this == nil || this.socket == nil || this.tunnel == nil {
        return errors.New("not connected")
    }

    serr := this.socket.Close()
    terr := this.tunnel.Close()

    this.socket = nil
    this.tunnel = nil

    if serr != nil {
        return serr
    }

    if terr != nil {
        return terr
    }

    return nil
}

func (this *Tunnel) Socket () net.Conn {
    return this.socket
}