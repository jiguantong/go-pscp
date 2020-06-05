package main

type Host struct {
	User string
	Ip   string
	Dir  string
	Pwd  string
}

func (h Host) String() string {
	return "ip: " + h.Ip + "\npwd: " + h.Pwd + "\ndir: " + h.Dir
}
