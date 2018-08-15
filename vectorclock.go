package vectorclock

type Msg struct {
  data interface{}
  vclock map[string]int
}

type Process struct {
  actor     string
  data      interface{}
  vclock    map[string]int
  recv      chan Msg
}


func New(actor string) *Process {
  m := map[string]int{}
  m[actor] = 0
  return &Process{
    vclock:m,
    actor: actor,
    recv: make(chan Msg),
  }
}

func (self *Process) incr() {
  if _, found := self.vclock[self.actor]; !found {
    self.vclock[self.actor] = 0
  }
  self.vclock[self.actor] += 1
}

func (self *Process) Set(data interface{}) {
    self.incr()
    self.data = data
}

func (self *Process) Get() interface{} {
  return self.data
}

func (self *Process) Vclock() interface{} {
  return self.vclock
}

func (self *Process) SendMsgFrom(p *Process) {
  msg := Msg{
    data: p.data,
    vclock: p.vclock,
  }
  self.sync(msg)
}

func (self *Process) Stop() {
  close(self.recv)
}

func (self *Process) sync(msg Msg) {
  self.incr()
  hasCorrectInfo := true
  for k, _ := range msg.vclock {
    _, ok := self.vclock[k]
    if ok {
      if msg.vclock[k] > self.vclock[k] {
        hasCorrectInfo = false
        self.vclock[k] = msg.vclock[k]
      }
    } else {
      hasCorrectInfo = false
      self.vclock[k] = msg.vclock[k]
    }
  }
  if !hasCorrectInfo {
    self.data = msg.data
  }
}
